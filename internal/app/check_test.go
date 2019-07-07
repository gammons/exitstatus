package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/stretchr/testify/suite"
)

type CheckStatusSuite struct {
	suite.Suite
}

// active/inactive statuses are easy, since they only check the Active flag.
func (suite *CheckStatusSuite) TestInactiveStatus() {

	inactiveCheck := &Check{
		Name:     "inactive check",
		Schedule: "*/20 * * * *",
		Status:   Inactive,
	}

	suite.Equal(Inactive, inactiveCheck.EvaluateStatus(time.Now()))
}

type UnknownStatusSuite struct {
	suite.Suite
}

// Unknown will remain unknown if:
// * the evaluation time is after the next due checkin
// * the current status is unknown
func (suite *UnknownStatusSuite) TestUnknownToUnknown() {
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	check := &Check{
		Name:        "active check",
		LastCheckin: oneDayAgo,
		Schedule:    "*/20 * * * *",
		Status:      Unknown,
	}

	suite.Equal(Unknown, check.EvaluateStatus(time.Now()))
}

// If the last status was healthy, this will switch to unhealthy.
func (suite *HealthyStatusSuite) TestHealthyToUnhealthy() {
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	check := &Check{
		Name:        "active check",
		LastCheckin: oneDayAgo,
		Schedule:    "*/20 * * * *",
		Status:      Healthy,
	}

	suite.Equal(Unhealthy, check.EvaluateStatus(time.Now()))
}

// if the evaluation window is before the next required checkin time,
// then the status remains healthy
func (suite *HealthyStatusSuite) TestRemainsHealthy() {
	nextTime := cronexpr.MustParse("*/20 * * * *").Next(time.Now())
	nextTimeMinus10 := nextTime.Add(-10 * time.Minute)

	check := &Check{
		Name:        "active check",
		LastCheckin: nextTimeMinus10,
		Schedule:    "*/20 * * * *",
		Status:      Healthy,
	}

	suite.Equal(Healthy, check.EvaluateStatus(time.Now()))
}

// If a status is currently unhealthy, and we evaluate
func (suite *UnknownStatusSuite) TestStatusRemainsUnhealthy() {
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	fmt.Println("oneDayAgo  is ", oneDayAgo)

	check := &Check{
		Name:        "active check",
		LastCheckin: oneDayAgo,
		Schedule:    "*/20 * * * *",
		Status:      Unknown,
	}

	suite.Equal(Unknown, check.EvaluateStatus(time.Now()))
}

type HealthyStatusSuite struct {
	suite.Suite
}

// required for a healthy status:
// * status is currently healthy
// * evaluation time is less than the next checkin time
func TestChecks(t *testing.T) {
	fmt.Println("here")
	suite.Run(t, new(CheckStatusSuite))
	suite.Run(t, new(UnknownStatusSuite))
	suite.Run(t, new(HealthyStatusSuite))
}
