package app

import (
	"time"

	"github.com/gorhill/cronexpr"
)

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

func (c *Check) CheckIn(evaluationTime time.Time) {
	c.LastCheckin = evaluationTime
	c.UpdateTemporalStatus(evaluationTime)
}

func (c *Check) UpdateTemporalStatus(evaluationTime time.Time) {
	if c.Status == InactiveStatus {
		return
	}

	if evaluationTime.After(c.NextCheckinDueDate()) {
		if c.Status == UnknownStatus {
			return
		}
		c.Status = UnhealthyStatus
		return
	}

	c.Status = HealthyStatus
}
