package command

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/husol/telegram-reminder-bot/pkg/date"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
)

type MessageRemindDayOfWeek struct {
	Day     string `regexpGroup:"day"`
	When    string `regexpGroup:"when"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePatternRemindDayOfWeek = `/remind me ?(on)? (?P<day>(((M|m)(on)|(T|t)(ues)|(W|w)(ednes)|(T|t)(hurs)|(F|f)(ri)|(S|s)(atur)|(S|s)(un))(day))) ?(?P<when>morning|afternoon|evening|night)? ?(at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)?)? (?P<message>.*)`

func HandleRemindDayOfWeek(service reminder.ServiceReminder) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageRemindDayOfWeek)
		if err := c.Bind(message); err != nil {
			return err
		}

		dateTime := mapMessageRemindDayOfWeekToReminderDateTime(message)
		nextSchedule, err := service.AddReminderOnDateTime(int(c.ChatID()), c.Text(), dateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))
		return err
	}
}

// nolint:goconst
func mapMessageRemindDayOfWeekToReminderDateTime(m *MessageRemindDayOfWeek) reminder.DateTime {
	dt := reminder.DateTime{
		DayOfWeek: strconv.Itoa(date.ToNumericDayOfWeek(m.Day)),
		Hour:      9,
		Minute:    0,
	}

	switch m.When {
	case "morning":
		dt.Hour = 9
		dt.Minute = 0

	case "afternoon":
		dt.Hour = 15
		dt.Minute = 0

	case "evening", "night":
		dt.Hour = 20
		dt.Minute = 0

	default:
	}

	if m.Hour != nil {
		hour, minute := date.ConvertTo24H(*m.Hour, m.Minute, m.AMPM)

		dt.Hour = hour
		dt.Minute = minute
	}

	return dt
}
