package command

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/husol/telegram-reminder-bot/pkg/date"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
)

type MessageRemindWhen struct {
	When    string `regexpGroup:"when"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePatternRemindWhen = `/remind me (?P<when>this afternoon|this evening|tonight|tomorrow morning|tomorrow afternoon|tomorrow evening|tomorrow) ?(at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)?)? (?P<message>.*)`

func HandleRemindWhen(service reminder.ServiceReminder) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageRemindWhen)
		if err := c.Bind(message); err != nil {
			return err
		}

		dateTime, err := mapMessageRemindWhenToReminderWordDateTime(message)
		if err != nil {
			return err
		}

		nextSchedule, err := service.AddReminderOnWordDateTime(int(c.ChatID()), c.Text(), dateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))
		return err
	}
}

func mapMessageRemindWhenToReminderWordDateTime(m *MessageRemindWhen) (reminder.WordDateTime, error) {
	var wdt reminder.WordDateTime

	switch m.When {
	case "this afternoon":
		wdt = reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   15,
			Minute: 0,
		}
	case "this evening", "tonight":
		wdt = reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   20,
			Minute: 0,
		}
	case "tomorrow", "tomorrow morning":
		wdt = reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   9,
			Minute: 0,
		}
	case "tomorrow afternoon":
		wdt = reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   15,
			Minute: 0,
		}
	case "tomorrow evening":
		wdt = reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   20,
			Minute: 0,
		}
	default:
		return wdt, fmt.Errorf("time not recognised: %s", m.When)
	}

	if m.Hour != nil {
		hour, minute := date.ConvertTo24H(*m.Hour, m.Minute, m.AMPM)

		wdt.Hour = hour
		wdt.Minute = minute
	}

	return wdt, nil
}
