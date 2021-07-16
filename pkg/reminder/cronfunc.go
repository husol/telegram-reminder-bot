package reminder

//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/husol/telegram-reminder-bot/pkg/chatpreference"
	"github.com/husol/telegram-reminder-bot/pkg/cron"
	"github.com/husol/telegram-reminder-bot/pkg/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
)

type CronFuncServicer interface {
	Complete(r *Reminder) error
	UpdateReminderWithNextRun(rem *Reminder) error
	UpdateReminderWithRepeatSchedule(rem *Reminder) error
}

type CronFuncService struct {
	b                   telegram.TBWrapBot
	scheduler           cron.Scheduler
	reminderStore       Storer
	chatPreferenceStore chatpreference.Storer
}

func NewCronFuncService(
	b telegram.TBWrapBot,
	scheduler cron.Scheduler,
	reminderStore Storer,
	chatPreferenceStore chatpreference.Storer,
) *CronFuncService {
	return &CronFuncService{
		b:                   b,
		scheduler:           scheduler,
		reminderStore:       reminderStore,
		chatPreferenceStore: chatPreferenceStore,
	}
}

func (s *CronFuncService) Complete(r *Reminder) error {
	r.Status = cron.Completed

	timeNow := time.Now().In(time.UTC)
	r.CompletedAt = &timeNow

	err := s.reminderStore.UpdateReminder(r)
	if err != nil {
		return err
	}

	s.scheduler.Remove(r.CronID)
	return nil
}

// NewCronFunc creates a function which is called when a reminder is due
// Note: repeatable jobs can be of two kinds:
// - Reminders set as "remind me every 31 april at 13:52" will have a cron job like "52 13 31 April *"
//   These are occurrences which occur once a year
//   In this case RunOnlyOnce = false as we want to keep the schedule
// - Reminders set as "remind me every 3 minutes" will have a cron job set on a very specific date like "46 15 4 4 *".
//   These reminders are set with RunOnlyOnce = true as they should only run once.
//   They will have a RepeatSchedule which will reschedule the job for the following occurrence (e.g. in 3 minutes from now)
func NewCronFunc(s CronFuncServicer, b telegram.TBWrapBot, r *Reminder) func() {
	return func() {
		buttons := NewButtons()
		var inlineKeys [][]tb.InlineButton
		var inlineButtons []tb.InlineButton

		snoozeBtn := *buttons[SnoozeBtn]
		snoozeBtn.Data = strconv.Itoa(r.ID)
		inlineButtons = append(
			inlineButtons,
			snoozeBtn,
		)

		// if repeatable job add button to complete it
		if !r.Job.RunOnlyOnce || (r.Job.RunOnlyOnce && r.Job.RepeatSchedule != nil) {
			completeBtn := *buttons[CompleteBtn]
			completeBtn.Data = strconv.Itoa(r.ID)
			inlineButtons = append(inlineButtons, completeBtn)
		}
		inlineKeys = append(inlineKeys, inlineButtons)

		messageWithIcon := fmt.Sprintf("🗓 %s", r.Data.Message)
		_, err := b.Send(&tb.Chat{ID: int64(r.Data.RecipientID)}, messageWithIcon, &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
		if err != nil {
			log.Printf("NewReminderCronFunc err: %q", err)
			return
		}

		timeNow := time.Now().In(time.UTC)
		r.LastRunAt = &timeNow

		if !r.Job.RunOnlyOnce {
			// update the next run at field of the reminder if it is a recurring reminder
			err = s.UpdateReminderWithNextRun(r)
			if err != nil {
				log.Printf("NewReminderCronFunc UpdateReminderWithNextRun err: %q", err)
				return
			}
			return
		}

		if r.Job.RepeatSchedule != nil {
			// if the reminder has a RepeatSchedule then we don't want to Complete() it
			// but instead calculate the next time it should run and reschedule it
			updateErr := s.UpdateReminderWithRepeatSchedule(r)
			if updateErr != nil {
				log.Printf("NewReminderCronFunc UpdateReminderWithRepeatSchedule err: %q", updateErr)
				return
			}
			return
		}

		err = s.Complete(r)
		if err != nil {
			log.Printf("NewReminderCronFunc complete err: %q", err)
			return
		}
	}
}

// UpdateReminderWithRepeatSchedule updates the reminder setting the schedule
// date to be in the future according to the definition of RepeatSchedule.
// The current reminder on the scheduler gets removed and a new one is created
// with the newly calculated schedule
func (s *CronFuncService) UpdateReminderWithRepeatSchedule(rem *Reminder) error {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(rem.Job.ChatID)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return err
	}

	addedTime := time.Now().In(loc).Add(
		time.Duration(rem.RepeatSchedule.Days)*24*time.Hour +
			time.Duration(rem.RepeatSchedule.Hours)*time.Hour +
			time.Duration(rem.RepeatSchedule.Minutes)*time.Minute,
	)

	schedule := fmt.Sprintf("%d %d %d %d *",
		addedTime.Minute(),
		addedTime.Hour(),
		addedTime.Day(),
		addedTime.Month(),
	)
	rem.Job.Schedule = schedule

	// remove previous cron job before scheduling new one
	s.scheduler.Remove(rem.CronID)

	scheduleWithTZ := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, schedule)
	reminderCronID, err := s.scheduler.Add(scheduleWithTZ, NewCronFunc(s, s.b, rem))
	if err != nil {
		return err
	}
	rem.CronID = reminderCronID

	cronEntry := s.scheduler.GetEntryByID(reminderCronID)
	rem.NextRunAt = &cronEntry.Next

	err = s.reminderStore.UpdateReminder(rem)
	if err != nil {
		return err
	}

	return nil
}

// UpdateReminderWithNextRun updates the reminder NextRunAt field
// with the newly calculated schedule
func (s *CronFuncService) UpdateReminderWithNextRun(rem *Reminder) error {
	cronEntry := s.scheduler.GetEntryByID(rem.CronID)
	rem.NextRunAt = &cronEntry.Next

	err := s.reminderStore.UpdateReminder(rem)
	if err != nil {
		return err
	}

	return nil
}
