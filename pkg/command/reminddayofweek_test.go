package command_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/husol/telegram-reminder-bot/pkg/command"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
	"github.com/husol/telegram-reminder-bot/pkg/reminder/mocks"
	fakeBot "github.com/husol/telegram-reminder-bot/pkg/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

type TestCaseHandleRemindDayOfWeek struct {
	Text             string
	ExpectedDateTime reminder.DateTime
}

func TestHandleRemindDayOfWeek_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindDayOfWeek)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}
	testCases := newTestHandleRemindDayOfWeekTestCases()

	for name := range testCases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			bot := fakeBot.NewTBWrapBot()
			c := tbwrap.NewContext(bot, &tb.Message{Text: testCases[name].Text, Chat: chat}, nil, handlerPattern)
			mockReminderService := mocks.NewMockServicer(mockCtrl)
			mockReminderService.
				EXPECT().
				AddReminderOnDateTime(
					1,
					testCases[name].Text,
					testCases[name].ExpectedDateTime,
					"update weekly report").
				Return(reminder.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

			err := command.HandleRemindDayOfWeek(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func newTestHandleRemindDayOfWeekTestCases() map[string]TestCaseHandleRemindDayOfWeek {
	return map[string]TestCaseHandleRemindDayOfWeek{
		"without hours and minutes": {
			Text: "/remind me on tuesday update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      9,
				Minute:    0,
			},
		},
		"with hours and minutes": {
			Text: "/remind me on tuesday at 23:34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    34,
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me on tuesday at 23.34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    34,
			},
		},
		"with only hour": {
			Text: "/remind me on tuesday at 23 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    0,
			},
		},
		"with only hour pm": {
			Text: "/remind me on tuesday at 8pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    0,
			},
		},
		"with hour minute pm": {
			Text: "/remind me on tuesday at 8:30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    30,
			},
		},
		"with hour minute pm dot separator": {
			Text: "/remind me on tuesday at 8.30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    30,
			},
		},
		"with time of day and without hours and minutes": {
			Text: "/remind me on tuesday evening update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    0,
			},
		},
		"with time of day and hours and minutes": {
			Text: "/remind me on tuesday evening at 23:34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    34,
			},
		},
		"with time of day and hours and minutes dot separator": {
			Text: "/remind me on tuesday night at 23.34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    34,
			},
		},
		"with time of day and only hour": {
			Text: "/remind me on tuesday evening at 23 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    0,
			},
		},
		"with time of day and only hour pm": {
			Text: "/remind me on tuesday evening at 8pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    0,
			},
		},
		"with time of day and hour minute pm": {
			Text: "/remind me on tuesday night at 8:30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    30,
			},
		},
		"with time of day and hour minute pm dot separator": {
			Text: "/remind me on tuesday evening at 8.30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    30,
			},
		},
		"with time of day only": {
			Text: "/remind me tuesday evening update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    0,
			},
		},
	}
}

func TestHandleRemindDayOfWeek_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindDayOfWeek)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me on tuesday update weekly report"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	bot := fakeBot.NewTBWrapBot()
	c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
	mockReminderService := mocks.NewMockServicer(mockCtrl)
	mockReminderService.
		EXPECT().
		AddReminderOnDateTime(
			1,
			text,
			reminder.DateTime{
				DayOfWeek: "2",
				Hour:      9,
				Minute:    0,
			},
			"update weekly report").
		Return(reminder.NextScheduleChatTime{}, errors.New("error"))

	err = command.HandleRemindDayOfWeek(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
