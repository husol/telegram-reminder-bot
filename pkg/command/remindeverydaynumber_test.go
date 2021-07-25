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

func TestHandleRemindEveryDayNumber_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindEveryDayNumber)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	type TestCase struct {
		Text                       string
		ExpectedRepeatableDateTime *reminder.RepeatableDateTime
	}

	testCases := map[string]TestCase{
		"without hours and minutes": {
			Text: "/remind me every 4th of the month update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "9",
				Minute:     "0",
			},
		},
		"with hours and minutes": {
			Text: "/remind me every 4th of the month at 23:34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "23",
				Minute:     "34",
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me every 4th of the month at 23.34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "23",
				Minute:     "34",
			},
		},
		"with only hour": {
			Text: "/remind me every 4th of the month at 23 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "23",
				Minute:     "0",
			},
		},
		"with only hour pm": {
			Text: "/remind me every 4th of the month at 8pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "20",
				Minute:     "0",
			},
		},
		"with hour minute pm": {
			Text: "/remind me every 4th of the month at 8:30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "20",
				Minute:     "30",
			},
		},
		"with hour minute pm dot separator": {
			Text: "/remind me every 4th of the month at 8.30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "20",
				Minute:     "30",
			},
		},
	}

	for name := range testCases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			bot := fakeBot.NewTBWrapBot()
			c := tbwrap.NewContext(bot, &tb.Message{Text: testCases[name].Text, Chat: chat}, nil, handlerPattern)
			mockReminderService := mocks.NewMockServicer(mockCtrl)
			mockReminderService.
				EXPECT().
				AddRepeatableReminderOnDateTime(
					1,
					testCases[name].Text,
					testCases[name].ExpectedRepeatableDateTime,
					"update weekly report").
				Return(reminder.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

			err := command.HandleRemindEveryDayNumber(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func TestHandleRemindEveryDayNumber_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindEveryDayNumber)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me every 4th of the month update weekly report"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	bot := fakeBot.NewTBWrapBot()
	c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
	mockReminderService := mocks.NewMockServicer(mockCtrl)
	mockReminderService.
		EXPECT().
		AddRepeatableReminderOnDateTime(
			1,
			text,
			&reminder.RepeatableDateTime{
				DayOfMonth: "4",
				Month:      "*",
				Hour:       "9",
				Minute:     "0",
			},
			"update weekly report").
		Return(reminder.NextScheduleChatTime{}, errors.New("error"))

	err = command.HandleRemindEveryDayNumber(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
