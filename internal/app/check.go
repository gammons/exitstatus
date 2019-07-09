package app

import (
	"time"

	"github.com/gorhill/cronexpr"
)

// arrive at a status for a service
// status should be healthy, unhealthy, or unknown

// needs to know when exitstatus started
// needs to know when the last healthy check was

type CheckStatusType string

type Check struct {
	GracePeriod int
	LastCheckin time.Time
	Name        string
	Shortcode   string
	Schedule    string
	Status      CheckStatusType
}

const (
	HealthyStatus   CheckStatusType = "healthy"
	UnhealthyStatus CheckStatusType = "unhealthy"
	UnknownStatus   CheckStatusType = "unknown"
	InactiveStatus  CheckStatusType = "inactive"
)

func (c *Check) NextCheckinDueDate() time.Time {
	nextTime := cronexpr.MustParse(c.Schedule).Next(c.LastCheckin)
	return nextTime.Add(time.Duration(c.GracePeriod) * time.Second)
}

func (c *Check) PerformTemporalCheck(evaluationTime time.Time) CheckStatusType {
	if c.Status == InactiveStatus {
		return InactiveStatus
	}

	if evaluationTime.After(c.NextCheckinDueDate()) {
		if c.Status == UnknownStatus {
			return UnknownStatus
		}
		return UnhealthyStatus
	}

	return HealthyStatus
}
