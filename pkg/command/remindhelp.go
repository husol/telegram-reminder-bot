package command

import (
	"github.com/enrico5b1b4/tbwrap"
)

const HandlePatternHelp = "/remindhelp"

func HandleRemindHelp() func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		_, err := c.Send(remindHelpText)

		return err
	}
}

const remindHelpText = `
*Available commands*

_list reminders_
/remindlist

_get details of a reminder_
[/r_ID]

_delete a reminder_
[/reminddelete_ID]

_set a reminder_
/remind me on the 1st of december Update your report
/remind me on the 1st of december at 8:23 Update your report
/remind me tonight/this evening/tomorrow/tomorrow morning Update your report
/remind me today/tomorrow at 21:00 Update your report
/remind me on Tuesday at 22:00 Update your report
/remind me at 21:00 Update your report
/remind me in 5 days, 3 hours, 4 minutes Update your report
/remind me in 5 hours Update your report
/remind me in 3 hours, 4 minutes Update your report
/remind me in 4 minutes Update your report

_set a recurring reminder_
/remind me every 1st of december Update yearly report
/remind me every 1st of december at 8:23 Update yearly report
/remind me every 1st of the month Update monthly report
/remind me every 1st of the month at 8:23 Update monthly report
/remind me every Tuesday at 22:00 Update weekly report
/remind me every day at 8pm Update daily report
/remind me every 5 days, 3 hours, 4 minutes Update your report
/remind me every 3 hours, 4 minutes Update your report
/remind me every 2 minutes Update your report

_set timezone for chat reminders_
/gettimezone
/settimezone Asia/Ho_Chi_Minh
`
