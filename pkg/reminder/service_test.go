package reminder_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/husol/telegram-reminder-bot/pkg/chatpreference"
	chatpreferenceMocks "github.com/husol/telegram-reminder-bot/pkg/chatpreference/mocks"
	"github.com/husol/telegram-reminder-bot/pkg/cron"
	"github.com/husol/telegram-reminder-bot/pkg/date"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
	reminderMocks "github.com/husol/telegram-reminder-bot/pkg/reminder/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Mocks struct {
	ReminderStore       *reminderMocks.MockStorer
	Scheduler           *reminderMocks.MockScheduler
	ChatPreferenceStore *chatpreferenceMocks.MockStorer
}

const (
	message    = "message"
	command    = "command"
	timezone   = "Asia/Ho_Chi_Minh"
	chatID     = 1
	cronID     = 2
	reminderID = 3
)

var (
	stubNextScheduleTime = timeNow()
	stubCreatedAt        = timeNow()
)

func TestService_AddReminderOnDateTime(t *testing.T) {
	loc, err := time.LoadLocation(timezone)
	require.NoError(t, err)

	t.Run("success with day of month and month", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				CreatedAt:   stubCreatedAt,
				NextRunAt:   &stubNextScheduleTime,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)

		service := reminder.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderOnDateTime(chatID, command, reminder.DateTime{
			DayOfMonth: 1,
			Month:      date.ToNumericMonth(time.April.String()),
			Hour:       13,
			Minute:     52,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, reminder.NextScheduleChatTime{Time: timeNow(), Location: loc}, nextScheduleTime)
	})

	t.Run("success with day of month without month", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 1 * *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 1 * *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				CreatedAt:   stubCreatedAt,
				NextRunAt:   &stubNextScheduleTime,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)

		service := reminder.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderOnDateTime(chatID, command, reminder.DateTime{
			DayOfMonth: 1,
			Month:      0,
			Hour:       13,
			Minute:     52,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, reminder.NextScheduleChatTime{Time: timeNow(), Location: loc}, nextScheduleTime)
	})

	t.Run("success with day of week", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 * * 1",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 * * 1",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				CreatedAt:   stubCreatedAt,
				NextRunAt:   &stubNextScheduleTime,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)

		service := reminder.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderOnDateTime(chatID, command, reminder.DateTime{
			DayOfWeek: "1",
			Hour:      13,
			Minute:    52,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, reminder.NextScheduleChatTime{Time: timeNow(), Location: loc}, nextScheduleTime)
	})
}

func TestService_AddReminderOnWordDateTime(t *testing.T) {
	loc, err := time.LoadLocation(timezone)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				CreatedAt:   stubCreatedAt,
				NextRunAt:   &stubNextScheduleTime,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)

		service := reminder.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderOnWordDateTime(chatID, command, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   13,
			Minute: 52,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, reminder.NextScheduleChatTime{Time: timeNow(), Location: loc}, nextScheduleTime)
	})
}

func TestService_AddRepeatableReminderOnDateTime(t *testing.T) {
	loc, err := time.LoadLocation(timezone)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 31 April *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: false,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 31 April *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: false,
				CreatedAt:   stubCreatedAt,
				NextRunAt:   &stubNextScheduleTime,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)

		service := reminder.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddRepeatableReminderOnDateTime(chatID, command, &reminder.RepeatableDateTime{
			DayOfMonth: "31",
			Month:      time.April.String(),
			Hour:       "13",
			Minute:     "52",
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, reminder.NextScheduleChatTime{Time: timeNow(), Location: loc}, nextScheduleTime)
	})
}

func TestService_AddReminderIn(t *testing.T) {
	loc, err := time.LoadLocation(timezone)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				CreatedAt:   stubCreatedAt,
				NextRunAt:   &stubNextScheduleTime,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)

		service := reminder.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderIn(chatID, command, reminder.AmountDateTime{
			Minutes: 1,
			Hours:   2,
			Days:    3,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, reminder.NextScheduleChatTime{Time: timeNow(), Location: loc}, nextScheduleTime)
	})
}

func TestService_AddReminderEvery(t *testing.T) {
	loc, err := time.LoadLocation(timezone)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				RepeatSchedule: &cron.JobRepeatSchedule{
					Minutes: 1,
					Hours:   2,
					Days:    3,
				},
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				RepeatSchedule: &cron.JobRepeatSchedule{
					Minutes: 1,
					Hours:   2,
					Days:    3,
				},
				CreatedAt: stubNextScheduleTime,
				NextRunAt: &stubNextScheduleTime,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: timezone,
		}, nil)

		service := reminder.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderEvery(chatID, command, reminder.AmountDateTime{
			Minutes: 1,
			Hours:   2,
			Days:    3,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, reminder.NextScheduleChatTime{Time: timeNow(), Location: loc}, nextScheduleTime)
	})
}

func createMocks(mockCtrl *gomock.Controller) Mocks {
	return Mocks{
		ReminderStore:       reminderMocks.NewMockStorer(mockCtrl),
		Scheduler:           reminderMocks.NewMockScheduler(mockCtrl),
		ChatPreferenceStore: chatpreferenceMocks.NewMockStorer(mockCtrl),
	}
}

func timeNow() time.Time {
	timeLoc, _ := time.LoadLocation(timezone)
	return time.Date(2020, time.April, 1, 13, 45, 0, 0, timeLoc).In(time.UTC)
}
