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
	Healthy   CheckStatusType = "healthy"
	Unhealthy CheckStatusType = "unhealthy"
	Unknown   CheckStatusType = "unknown"
	Inactive  CheckStatusType = "inactive"
)

func (c *Check) EvaluateStatus(evaluationTime time.Time) CheckStatusType {
	if c.Status == Inactive {
		return Inactive
	}

	nextTime := cronexpr.MustParse(c.Schedule).Next(c.LastCheckin)

	if evaluationTime.After(nextTime) {
		if c.Status == Unknown {
			return Unknown
		}
		return Unhealthy
	}

	return Healthy
}
