package command

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/husol/telegram-reminder-bot/pkg/date"
	"github.com/husol/telegram-reminder-bot/pkg/reminder"
)

type MessageRemindEveryDayNumberMonth struct {
	Day     int    `regexpGroup:"day"`
	Month   string `regexpGroup:"month"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePatternRemindEveryDayNumberMonth = `/remind me every (?P<day>\d{1,2})(?:(st|nd|rd|th))? of (?P<month>(J|j)anuary|(F|f)ebruary|(M|m)arch|(A|a)pril|(M|m)ay|(J|j)une|(J|j)uly|(A|a)ugust|(S|s)eptember|(O|o)ctober|(N|n)ovember|(D|d)ecember) ?(at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)?)? (?P<message>.*)`

func HandleRemindEveryDayNumberMonth(service reminder.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageRemindEveryDayNumberMonth)
		if err := c.Bind(message); err != nil {
			return err
		}

		repeatDateTime := mapMessageRemindEveryDayNumberMonthToReminderDateTime(message)
		nextSchedule, err := service.AddRepeatableReminderOnDateTime(int(c.ChatID()), c.Text(), &repeatDateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))

		return err
	}
}

func mapMessageRemindEveryDayNumberMonthToReminderDateTime(m *MessageRemindEveryDayNumberMonth) reminder.RepeatableDateTime {
	rdt := reminder.RepeatableDateTime{
		DayOfMonth: strconv.Itoa(m.Day),
		Month:      strconv.Itoa(date.ToNumericMonth(m.Month)),
		Hour:       "9",
		Minute:     "0",
	}

	if m.Hour != nil {
		hour, minute := date.ConvertTo24H(*m.Hour, m.Minute, m.AMPM)

		rdt.Hour = strconv.Itoa(hour)
		rdt.Minute = strconv.Itoa(minute)
	}

	return rdt
}
