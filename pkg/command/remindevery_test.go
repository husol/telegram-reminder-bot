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

// nolint:dupl
func TestHandleRemindEveryPattern1(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindEvery)
	require.NoError(t, err)
	text := "/remind me every 2 minutes update weekly report"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderEvery(
				1,
				text,
				reminder.AmountDateTime{
					Days:    0,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminder.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

		err := command.HandleRemindEvery(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderEvery(
				1,
				text,
				reminder.AmountDateTime{
					Days:    0,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminder.NextScheduleChatTime{}, errors.New("error"))

		err := command.HandleRemindEvery(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}

// nolint:dupl
func TestHandleRemindEveryPattern2(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindEvery)
	require.NoError(t, err)
	text := "/remind me every 2 minutes, 3 days update weekly report"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderEvery(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminder.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

		err := command.HandleRemindEvery(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderEvery(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminder.NextScheduleChatTime{}, errors.New("error"))

		err := command.HandleRemindEvery(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}

// nolint:dupl
func TestHandleRemindEveryPattern3(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindEvery)
	require.NoError(t, err)
	text := "/remind me every 2 minutes, 1 hour, 3 days update weekly report"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderEvery(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   1,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminder.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

		err := command.HandleRemindEvery(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderEvery(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   1,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminder.NextScheduleChatTime{}, errors.New("error"))

		err := command.HandleRemindEvery(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
