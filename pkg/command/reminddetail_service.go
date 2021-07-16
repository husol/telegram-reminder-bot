package command

//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks

import (
	"errors"
	"time"

	"github.com/husol/telegram-reminder-bot/pkg/chatpreference"
	"github.com/husol/telegram-reminder-bot/pkg/cron"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
)

type RemindDetailServicer interface {
	GetReminder(chatID, reminderID int) (*ReminderDetail, error)
	DeleteReminder(chatID, ID int) error
}

type RemindDetailService struct {
	reminderStore       reminder.Storer
	scheduler           cron.Scheduler
	chatPreferenceStore chatpreference.Storer
}

func NewRemindDetailService(
	reminderStore reminder.Storer,
	scheduler cron.Scheduler,
	chatPreferenceStore chatpreference.Storer,
) *RemindDetailService {
	return &RemindDetailService{
		reminderStore:       reminderStore,
		scheduler:           scheduler,
		chatPreferenceStore: chatPreferenceStore,
	}
}

func (s *RemindDetailService) GetReminder(chatID, reminderID int) (*ReminderDetail, error) {
	rem, err := s.reminderStore.GetReminder(chatID, reminderID)
	if err != nil {
		return nil, err
	}

	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return nil, err
	}

	reminderDetail := &ReminderDetail{Reminder: *rem}
	if rem.Status == cron.Active {
		cronEntry := s.scheduler.GetEntryByID(rem.CronID)
		nextScheduleInChatTimezone := cronEntry.Next.In(loc)
		reminderDetail.NextSchedule = &nextScheduleInChatTimezone
	}
	if rem.Status == cron.Completed && rem.CompletedAt != nil {
		completedAtChatTimezone := rem.CompletedAt.In(loc)
		reminderDetail.CompletedAt = &completedAtChatTimezone
	}

	return reminderDetail, nil
}

func (s *RemindDetailService) DeleteReminder(chatID, id int) error {
	rem, err := s.reminderStore.GetReminder(chatID, id)
	if err != nil {
		return err
	}

	if chatID != rem.ChatID {
		return errors.New("unauthorised to delete reminder")
	}

	s.scheduler.Remove(rem.CronID)

	return s.reminderStore.DeleteReminder(rem.ChatID, id)
}
