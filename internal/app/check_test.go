package app

import (
	"testing"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/poy/onpar"
	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

// all temporal checks are assumed to run *after* the check window.
func TestTemporalChecks(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	nextTime := cronexpr.MustParse("*/20 * * * *").Next(time.Now())
	nextTimeMinus10 := nextTime.Add(-10 * time.Minute)
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	o.Group("When a check has successfully checked in before window expiry", func() {
		o.Spec("Inactive remains Inactive", func(*testing.T) {
			check := &Check{
				Name:        "inactive check",
				LastCheckin: nextTimeMinus10,
				Schedule:    "*/20 * * * *",
				Status:      InactiveStatus,
			}

			check.UpdateTemporalStatus(time.Now())

			Expect(t, check.Status).To(Equal(InactiveStatus))
		})

		o.Spec("Healthy remains healthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: nextTimeMinus10,
				Schedule:    "*/20 * * * *",
				Status:      HealthyStatus,
			}

			check.UpdateTemporalStatus(time.Now())

			Expect(t, check.Status).To(Equal(HealthyStatus))
		})

		// this might not be necessary.  This mechanism should be covered by the checkin functionality.
		o.Spec("Unhealthy transitions to Healthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: nextTimeMinus10,
				Schedule:    "*/20 * * * *",
				Status:      UnhealthyStatus,
			}

			check.UpdateTemporalStatus(time.Now())

			Expect(t, check.Status).To(Equal(HealthyStatus))
		})
	})

	o.Group("When a check failed to check in before window expiry", func() {
		o.Spec("Inactive remains Inactive", func(*testing.T) {
			check := &Check{
				Name:        "inactive check",
				LastCheckin: oneDayAgo,
				Schedule:    "*/20 * * * *",
				Status:      InactiveStatus,
			}

			check.UpdateTemporalStatus(time.Now())

			Expect(t, check.Status).To(Equal(InactiveStatus))
		})

		o.Spec("Healthy transitions to Uhealthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: oneDayAgo,
				Schedule:    "*/20 * * * *",
				Status:      HealthyStatus,
			}

			check.UpdateTemporalStatus(time.Now())

			Expect(t, check.Status).To(Equal(UnhealthyStatus))
		})

		o.Spec("Unhealthy remains unhealthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: oneDayAgo,
				Schedule:    "*/20 * * * *",
				Status:      UnhealthyStatus,
			}

			check.UpdateTemporalStatus(time.Now())

			Expect(t, check.Status).To(Equal(UnhealthyStatus))
		})

		o.Spec("Unknown remains Unknown", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: oneDayAgo,
				Schedule:    "*/20 * * * *",
				Status:      UnknownStatus,
			}

			check.UpdateTemporalStatus(time.Now())

			Expect(t, check.Status).To(Equal(UnknownStatus))
		})
	})
}

func TestCheckIn(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	schedule := "*/20 * * * *"
	runTime, _ := time.Parse("2006-Jan-02", "2013-Feb-03")

	currentCheckInDueDate := cronexpr.MustParse(schedule).Next(runTime)
	nextCheckInDueDate := cronexpr.MustParse(schedule).Next(currentCheckInDueDate)
	check := &Check{
		Name:        "active check",
		LastCheckin: runTime,
		Schedule:    "*/20 * * * *",
		Status:      UnknownStatus,
	}

	o.Spec("It transitions from Unknown -> Healthy", func(*testing.T) {

		check.CheckIn(currentCheckInDueDate)

		Expect(t, check.Status).To(Equal(HealthyStatus))
		Expect(t, check.LastCheckin).To(Equal(currentCheckInDueDate))
		Expect(t, check.NextCheckinDueDate()).To(Equal(nextCheckInDueDate))
	})
}
