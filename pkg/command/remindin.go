package command

import (
	"log"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
)

type MessageRemindIn struct {
	Who      string `regexpGroup:"who"`
	Amount1  int    `regexpGroup:"amount1"`
	Measure1 string `regexpGroup:"measure1"`
	Amount2  int    `regexpGroup:"amount2"`
	Measure2 string `regexpGroup:"measure2"`
	Amount3  int    `regexpGroup:"amount3"`
	Measure3 string `regexpGroup:"measure3"`
	Message  string `regexpGroup:"message"`
}

// nolint:lll
const HandlePatternRemindIn = `/remind me in (?P<amount1>\d{1,2}) (?P<measure1>minute|minutes|hour|hours|day|days)?(, (?P<amount2>\d{1,2}) (?P<measure2>minute|minutes|hour|hours|day|days)?(, (?P<amount3>\d{1,2}) (?P<measure3>minute|minutes|hour|hours|day|days))?)? (?P<message>.*)`

func HandleRemindIn(service reminder.ServiceReminder) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageRemindIn)
		if err := c.Bind(message); err != nil {
			log.Println("bind error")
			return err
		}

		amountDateTime := mapMessageRemindInToAmountDateTime(message)
		nextSchedule, err := service.AddReminderIn(int(c.ChatID()), c.Text(), amountDateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))

		return err
	}
}

func mapMessageRemindInToAmountDateTime(m *MessageRemindIn) reminder.AmountDateTime {
	amountDateTime := reminder.AmountDateTime{}
	amounts := []int{m.Amount1, m.Amount2, m.Amount3}
	measures := []string{m.Measure1, m.Measure2, m.Measure3}

	for i := 0; i < 3; i++ {
		switch measures[i] {
		case "minute", "minutes":
			amountDateTime.Minutes = amounts[i]
		case "hour", "hours":
			amountDateTime.Hours = amounts[i]
		case "day", "days":
			amountDateTime.Days = amounts[i]
		}
	}

	return amountDateTime
}
