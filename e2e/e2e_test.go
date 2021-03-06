package e2e_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/husol/telegram-reminder-bot/pkg/bot"
	"github.com/husol/telegram-reminder-bot/pkg/db"
	"github.com/husol/telegram-reminder-bot/pkg/telegram/fakes"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

const chatID = 123456

func TestE2E(t *testing.T) {
	checkSkip(t)

	dbFile := mustGetEnv("TEST_E2E_DB_FILE")
	allowedChats := []int{chatID}

	telebot, database, err := setup(dbFile, allowedChats)
	defer database.Close()
	require.NoError(t, err)

	// RemindAt
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me tomorrow at 20:45 MSG1_")
	require.Contains(t, telebot.OutboundSendMessages[0], `Reminder "MSG1_" has been added`)

	// RemindDayMonth
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me on the 1st of december at 8:23 MSG2_")
	require.Contains(t, telebot.OutboundSendMessages[1], `Reminder "MSG2_" has been added`)

	// RemindEvery
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 2 minutes MSG3_")
	require.Contains(t, telebot.OutboundSendMessages[2], `Reminder "MSG3_" has been added`)

	// RemindEveryDayNumber
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 1st of the month at 8:23 MSG4_")
	require.Contains(t, telebot.OutboundSendMessages[3], `Reminder "MSG4_" has been added`)

	// RemindEveryDayNumberMonth
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 1st of december at 8:23 MSG5_")
	require.Contains(t, telebot.OutboundSendMessages[4], `Reminder "MSG5_" has been added`)

	// RemindEveryDayOfWeek
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every tuesday at 8:23 MSG6_")
	require.Contains(t, telebot.OutboundSendMessages[5], `Reminder "MSG6_" has been added`)

	// RemindIn
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me in 5 minutes MSG7_")
	require.Contains(t, telebot.OutboundSendMessages[6], `Reminder "MSG7_" has been added`)

	// RemindWhen
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me tomorrow MSG8_")
	require.Contains(t, telebot.OutboundSendMessages[7], `Reminder "MSG8_" has been added`)

	// RemindDayOfWeek
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me on tuesday MSG9_")
	require.Contains(t, telebot.OutboundSendMessages[8], `Reminder "MSG9_" has been added`)

	// RemindEveryDay
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every day MSG10_")
	require.Contains(t, telebot.OutboundSendMessages[9], `Reminder "MSG10_" has been added`)

	// RemindList
	telebot.SimulateIncomingMessageToChat(chatID, "/remindlist")
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG1_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG2_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG3_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG4_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG5_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG6_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG7_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG8_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG9_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG10_`)

	telebot.SimulateIncomingMessageToChat(chatID, "/gettimezone")
	require.Contains(t, telebot.OutboundSendMessages[11], `Asia/Ho_Chi_Minh`)

	telebot.SimulateIncomingMessageToChat(chatID, "/settimezone Europe/Rome")
	require.Contains(t, telebot.OutboundSendMessages[12], `Europe/Rome`)

	telebot.SimulateIncomingMessageToChat(chatID, "/gettimezone")
	require.Contains(t, telebot.OutboundSendMessages[13], `Europe/Rome`)

	// Delete a reminder
	telebot.SimulateIncomingMessageToChat(chatID, "/remindlist")
	reminderID, err := getReminderIDForMessageFromRemindList(telebot.OutboundSendMessages[14], "MSG4_")
	require.NoError(t, err)

	telebot.SimulateIncomingMessageToChat(chatID, fmt.Sprintf("/reminddelete %s", reminderID))
	require.Contains(t, telebot.OutboundSendMessages[15], fmt.Sprintf("Reminder %s has been deleted", reminderID))

	telebot.SimulateIncomingMessageToChat(chatID, "/remindlist")
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG1_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG2_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG3_`)
	require.NotContains(t, telebot.OutboundSendMessages[16], `MSG4_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG5_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG6_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG7_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG8_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG9_`)
	require.Contains(t, telebot.OutboundSendMessages[16], `MSG10_`)
}

func setup(dbFile string, allowedChats []int) (*fakes.TeleBot, *bolt.DB, error) {
	database, err := db.SetupDB(dbFile, allowedChats)
	if err != nil {
		return nil, nil, err
	}

	teleBot := fakes.NewTeleBot()
	botConfig := tbwrap.Config{
		AllowedChats: allowedChats,
		TBot:         teleBot,
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		return nil, nil, err
	}

	appBot := bot.New(allowedChats, database, telegramBot)
	appBot.Start()

	return teleBot, database, nil
}

func mustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}

	return value
}

func checkSkip(t *testing.T) {
	testDBFile := os.Getenv("TEST_E2E_DB_FILE")
	if testDBFile == "" {
		t.Skip()
	}
}

func getReminderIDForMessageFromRemindList(text, msg string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`.*%s.*\[\[\/r_([0-9]+)\]\]`, msg))
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return "", errors.New("error getting reminder id")
	}

	return match[1], nil
}
