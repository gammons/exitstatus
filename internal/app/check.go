package app

import "time"

// arrive at a status for a service
// status should be healthy, unhealthy, or unknown

// needs to know when exitstatus started
// needs to know when the last healthy check was

type CheckTypeType int
type CheckStatusType int

const (
	Interval CheckTypeType = 0
)

type Check struct {
	Name      string
	Active    bool
	CheckType CheckTypeType

	LastCheckin time.Time
	Interval    string
}

const (
	Healthy   CheckStatusType = 0
	Unhealthy CheckStatusType = 1
	Unknown   CheckStatusType = 2
)

type CheckStatus struct {
}

func (c *CheckStatus) GetStatus(ProgramStart time.Time) CheckStatusType {
	return Healthy
}
