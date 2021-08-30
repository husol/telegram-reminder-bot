package command_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/golang/mock/gomock"
	"github.com/husol/telegram-reminder-bot/pkg/chatpreference"
	"github.com/husol/telegram-reminder-bot/pkg/chatpreference/mocks"
	"github.com/husol/telegram-reminder-bot/pkg/command"
	fakeBot "github.com/husol/telegram-reminder-bot/pkg/telegram/fakes"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleGetTimezone(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternGetTimezone)
	require.NoError(t, err)
	text := "/gettimezone"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockChatPreferenceStore := mocks.NewMockStorer(mockCtrl)
		mockChatPreferenceStore.
			EXPECT().
			GetChatPreference(1).
			Return(&chatpreference.ChatPreference{ChatID: 1, TimeZone: "Asia/Ho_Chi_Minh"}, nil)

		err := command.HandleGetTimezone(mockChatPreferenceStore)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockChatPreferenceStore := mocks.NewMockStorer(mockCtrl)
		mockChatPreferenceStore.
			EXPECT().
			GetChatPreference(1).
			Return(nil, errors.New("error"))

		err := command.HandleGetTimezone(mockChatPreferenceStore)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
