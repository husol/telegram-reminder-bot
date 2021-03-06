# Telegram Reminder Bot

Telegram bot that allows you to set reminders for a personal or group chat. Inspired by [Slack reminders](https://slack.com/intl/en-gb/help/articles/208423427-Set-a-reminder).  
Uses [robfig/cron](https://github.com/robfig/cron) for scheduling tasks, [tucnak/telebot](https://github.com/tucnak/telebot/tree/v2) and [enrico5b1b4/tbwrap](https://github.com/enrico5b1b4/tbwrap) for interacting with the telegram api, [etcd-io/bbolt](https://github.com/etcd-io/bbolt) for local storage of reminders.

![screen](screen.png "Screen")

## Usage

```
make build  

export TELEGRAM_REMINDER_DB_FILE=local.db  
export TELEGRAM_REMINDER_BOT_TOKEN=<TELEGRAM_BOT_TOKEN>
export TELEGRAM_ALLOWED_CHATS=<CHAT_IDS_SEPARATED_BY_COMMA>

./bin/build/telegram-reminder-bot
```

## Commands

### Remind help
List all commands  
`/remindhelp`

### Remind list
Retrieve the list of active and completed reminders  
`/remindlist`

### Remind detail
Display details of a reminder  
`/reminddetail 1`

### Remind on a date
Set a reminder in the format `[who] [when] [what]`

#### Fixed times and dates
- `/remind me on the 1 of december Update your report`  
  `/remind me on the 1st of december Update your report`
- `/remind me on the 1 of december at 8:23 Update your report`  
  `/remind me on the 1st of december at 8:23 Update your report`
- `/remind me tonight Update your report`  
  `/remind me tonight at 21:20 Update your report`  
  `/remind me tomorrow morning Update your report`  
  `/remind me tomorrow at 16:45 Update your report`
- `/remind me on Tuesday Update your report`
- `/remind me at 21:00 Update your report`
- `/remind me in 3 days Update your report`
- `/remind me in 5 days, 3 hours, 4 minutes Update your report`
- `/remind me in 5 hours Update your report`
- `/remind me in 3 hours, 4 minutes Update your report`
- `/remind me in 4 minutes Update your report`

#### Recurring
- `/remind me every 1st of december Update yearly report`  
  `/remind me every 1 of december Update yearly report`
- `/remind me every 1st of december at 8:23 Update yearly report`  
  `/remind me every 1 of december at 8:23 Update yearly report`
- `/remind me every 1st of the month Update monthly report`  
  `/remind me every 1 of the month Update monthly report`
- `/remind me every 1st of the month at 8:23 Update monthly report`  
  `/remind me every 1 of the month at 8:23 Update monthly report`
- `/remind me every Tuesday Update weekly report`
  `/remind me every Tuesday at 8:23 Update weekly report`  
- `/remind me every day at 8pm Update daily report`
- `/remind me every 5 days, 3 hours, 4 minutes Update your report`
- `/remind me every 3 hours, 4 minutes Update your report`
- `/remind me every 2 minutes Update your report`

#### Timezone management
- `/gettimezone`
- `/settimezone Asia/Ho_Chi_Minh`
