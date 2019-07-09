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

			Expect(t, check.PerformTemporalCheck(time.Now())).To(Equal(InactiveStatus))
		})

		o.Spec("Healthy remains healthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: nextTimeMinus10,
				Schedule:    "*/20 * * * *",
				GracePeriod: 20,
				Status:      HealthyStatus,
			}

			Expect(t, check.PerformTemporalCheck(time.Now())).To(Equal(HealthyStatus))
		})

		// this might not be necessary.  This mechanism should be covered by the checkin functionality.
		o.Spec("Unhealthy transitions to Healthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: nextTimeMinus10,
				Schedule:    "*/20 * * * *",
				Status:      UnhealthyStatus,
			}

			Expect(t, check.PerformTemporalCheck(time.Now())).To(Equal(HealthyStatus))
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

			Expect(t, check.PerformTemporalCheck(time.Now())).To(Equal(InactiveStatus))
		})

		o.Spec("Healthy transitions to Uhealthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: oneDayAgo,
				Schedule:    "*/20 * * * *",
				Status:      HealthyStatus,
			}

			Expect(t, check.PerformTemporalCheck(time.Now())).To(Equal(UnhealthyStatus))
		})

		o.Spec("Unhealthy remains unhealthy", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: oneDayAgo,
				Schedule:    "*/20 * * * *",
				Status:      UnhealthyStatus,
			}

			Expect(t, check.PerformTemporalCheck(time.Now())).To(Equal(UnhealthyStatus))
		})

		o.Spec("Unknown remains Unknown", func(*testing.T) {
			check := &Check{
				Name:        "active check",
				LastCheckin: oneDayAgo,
				Schedule:    "*/20 * * * *",
				Status:      UnknownStatus,
			}

			Expect(t, check.PerformTemporalCheck(time.Now())).To(Equal(UnknownStatus))
		})
	})
}

func TestNextCheckinDueDate(t *testing.T) {
}
