package main

import "time"

type CronSpec struct {
	Minute uint64 `json:"minute,omitempty"`
	Hour   uint64 `json:"hour,omitempty"`
	Dom    uint64 `json:"dom,omitempty"`
	Month  uint64 `json:"month,omitempty"`
	Dow    uint64 `json:"dow,omitempty"`
}

type Job struct {
	Command   string    `json:"command,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	CronSpec
}
