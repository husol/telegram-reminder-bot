package cron

import "time"

type JobType int

const (
	HealthCheck JobType = 1
	Backup      JobType = 2
	Reminder    JobType = 3
)

func (j JobType) String() string {
	return [...]string{"", "HealthCheck", "Backup", "Reminder"}[j]
}

type JobStatus int

const (
	Active    JobStatus = 1
	Inactive  JobStatus = 2
	Completed JobStatus = 3
)

func (j JobStatus) String() string {
	return [...]string{"", "Active", "Inactive", "Completed"}[j]
}

type Job struct {
	ID             int                `json:"id"`
	CronID         int                `json:"cron_id"`
	ChatID         int                `json:"owner_id"`
	Schedule       string             `json:"schedule"`
	Type           JobType            `json:"type"`
	Status         JobStatus          `json:"status"`
	RunOnlyOnce    bool               `json:"run_only_once"`
	RepeatSchedule *JobRepeatSchedule `json:"repeat_schedule"`
	CompletedAt    *time.Time         `json:"completed_at"`
	NextRunAt      *time.Time         `json:"next_run_at"` // NextRunAt should match the schedule
	LastRunAt      *time.Time         `json:"last_run_at"`
	CreatedAt      time.Time          `json:"created_at"`
}

type JobRepeatSchedule struct {
	Minutes int `json:"minutes"`
	Hours   int `json:"hours"`
	Days    int `json:"days"`
	Months  int `json:"months"`
}
