package command_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/golang/mock/gomock"
	"github.com/husol/telegram-reminder-bot/pkg/command"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
	"github.com/husol/telegram-reminder-bot/pkg/reminder/mocks"
	fakeBot "github.com/husol/telegram-reminder-bot/pkg/telegram/fakes"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

type TestCaseHandleRemindDayMonth struct {
	Text             string
	ExpectedDateTime reminder.DateTime
}

func TestHandleRemindDayMonth_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindDayMonth)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}
	testCases := newTestHandleRemindDayMonthTestCases()

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

			err := command.HandleRemindDayMonth(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func newTestHandleRemindDayMonthTestCases() map[string]TestCaseHandleRemindDayMonth {
	return map[string]TestCaseHandleRemindDayMonth{
		"without hours and minutes": {
			Text: "/remind me on the 4th of march update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       9,
				Minute:     0,
			},
		},
		"with hours and minutes": {
			Text: "/remind me on the 4th of March at 23:34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       23,
				Minute:     34,
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me on the 4th of march at 23.34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       23,
				Minute:     34,
			},
		},
		"with only hour": {
			Text: "/remind me on the 4th of march at 23 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       23,
				Minute:     0,
			},
		},
		"with only hour pm": {
			Text: "/remind me on the 4th of march at 8pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       20,
				Minute:     0,
			},
		},
		"with hour minute pm": {
			Text: "/remind me on the 4th of march at 8:30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       20,
				Minute:     30,
			},
		},
		"with hour minute pm dot separator": {
			Text: "/remind me on the 4th of march at 8.30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       20,
				Minute:     30,
			},
		},
		"with only day and without hours and minutes": {
			Text: "/remind me on the 4th update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      0,
				Hour:       9,
				Minute:     0,
			},
		},
		"with only day and hours and minutes": {
			Text: "/remind me on the 4th at 23:34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      0,
				Hour:       23,
				Minute:     34,
			},
		},
		"with only day and hours and minutes dot separator": {
			Text: "/remind me on the 4th at 23.34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      0,
				Hour:       23,
				Minute:     34,
			},
		},
		"with only day and only hour": {
			Text: "/remind me on the 4th at 23 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      0,
				Hour:       23,
				Minute:     0,
			},
		},
		"with only day and only hour pm": {
			Text: "/remind me on the 4th at 8pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      0,
				Hour:       20,
				Minute:     0,
			},
		},
		"with only day and hour minute pm": {
			Text: "/remind me on the 4th at 8:30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      0,
				Hour:       20,
				Minute:     30,
			},
		},
		"with only day and hour minute pm dot separator": {
			Text: "/remind me on the 4th at 8.30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      0,
				Hour:       20,
				Minute:     30,
			},
		},
	}
}

func TestHandleRemindDayMonth_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindDayMonth)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me on the 4th of march update weekly report"
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
				DayOfMonth: 4,
				Month:      3,
				Hour:       9,
				Minute:     0,
			},
			"update weekly report").
		Return(reminder.NextScheduleChatTime{}, errors.New("error"))

	err = command.HandleRemindDayMonth(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
