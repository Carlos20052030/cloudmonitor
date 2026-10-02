package domain

import "time"

type Status string

const (
	StatusUP   Status = "UP"
	StatusDOWN Status = "DOWN"
)

type Target struct {
	Name string
	URL  string
}

type Result struct {
	TargetName string
	URL        string
	Status     Status
	StatusCode int
	LatencyMS  int64
	Error      string
	CheckedAt  time.Time
}
