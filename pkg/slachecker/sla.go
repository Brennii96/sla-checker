package slachecker

import (
	"fmt"
	"time"
)

type SLA struct {
	StartTime     time.Time
	SLALength     int    // SLA duration, e.g., 4
	TimeUnit      string // SLA time unit, e.g., "hours", "minutes"
	BusinessHours struct {
		StartHour int
		EndHour   int
	}
	ValidDays []time.Weekday // e.g., []time.Weekday{time.Monday, time.Tuesday, ...}
	Holidays  []time.Time    // Specific holidays when SLA is not applicable
}

// SLAResult contains the details about SLA status
type SLAResult struct {
	IsWithinSLA          bool      `json:"isWithinSLA"`
	Deadline             time.Time `json:"deadline"`
	Remaining            string    `json:"remaining"`
	Overage              string    `json:"overage,omitempty"`
	WorkingTimeRemaining string    `json:"workingTimeRemaining"`
}

// IsWithinSLA checks if the current time is within the SLA duration from the start time.
func (s SLA) IsWithinSLA(currentTime time.Time) bool {
	// Calculate the SLA deadline based on business hours, weekends, and holidays
	slaDeadline := s.calculateSLADeadline()

	// Check if the current time is before the calculated SLA deadline
	return currentTime.Before(slaDeadline)
}

// CheckSLA checks if the current time is within the SLA duration and returns additional details
func (s SLA) CheckSLA(currentTime time.Time) SLAResult {
	// Calculate the SLA deadline based on business hours, weekends, and holidays
	slaDeadline := s.calculateSLADeadline()

	// Calculate the time difference
	timeRemaining := slaDeadline.Sub(currentTime)
	isWithinSLA := currentTime.Before(slaDeadline)

	var overage time.Duration
	if !isWithinSLA {
		overage = currentTime.Sub(slaDeadline)
	}

	// Calculate working time remaining
	workingTimeRemaining := s.calculateWorkingTimeRemaining(currentTime, slaDeadline)

	// Convert durations to readable strings
	remainingStr := formatDuration(timeRemaining)
	overageStr := formatDuration(overage)

	return SLAResult{
		IsWithinSLA:          isWithinSLA,
		Deadline:             slaDeadline,
		Remaining:            remainingStr,
		Overage:              overageStr,
		WorkingTimeRemaining: workingTimeRemaining,
	}
}

// calculateWorkingTimeRemaining calculates the remaining working time considering business hours and days
func (s SLA) calculateWorkingTimeRemaining(startTime, endTime time.Time) string {
	remainingDuration := time.Duration(0)
	location := startTime.Location()
	currentTime := startTime

	for currentTime.Before(endTime) {
		if s.isValidDay(currentTime) && s.isWithinBusinessHours(currentTime) && !s.isHoliday(currentTime) {
			// Calculate the end of the business day
			endOfBusinessDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), s.BusinessHours.EndHour, 0, 0, 0, location)
			if endOfBusinessDay.After(endTime) {
				endOfBusinessDay = endTime
			}

			// Add the remaining business hours for the current day
			if endOfBusinessDay.After(currentTime) {
				remainingDuration += endOfBusinessDay.Sub(currentTime)
			}

			// Move to the end of the business day
			currentTime = endOfBusinessDay
		}

		// Move to the start of the next business day if current time is after business hours
		if !s.isWithinBusinessHours(currentTime) {
			currentTime = s.moveToNextBusinessDay(currentTime)
		}
	}

	return formatDuration(remainingDuration)
}

// calculateSLADeadline calculates the SLA deadline based on business hours, weekends, and holidays
func (s SLA) calculateSLADeadline() time.Time {
	remainingDuration := s.getSLADuration()

	// Start from the initial SLA start time
	currentTime := s.StartTime

	for remainingDuration > 0 {
		// If it's a valid business day and hour, reduce the remaining SLA time
		if s.isValidDay(currentTime) && s.isWithinBusinessHours(currentTime) && !s.isHoliday(currentTime) {
			// Reduce remaining SLA time by one hour
			if remainingDuration >= time.Hour {
				remainingDuration -= time.Hour
			} else {
				break
			}
		}

		// Move forward one hour to the next time slot
		currentTime = currentTime.Add(time.Hour)

		// If we've moved past business hours, skip to the start of the next business day
		if !s.isWithinBusinessHours(currentTime) {
			currentTime = s.moveToNextBusinessDay(currentTime)
		}
	}

	return currentTime
}

// formatDuration converts time.Duration to a human-readable format
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// moveToNextBusinessDay moves the given time to the start of the next business day
func (s SLA) moveToNextBusinessDay(t time.Time) time.Time {
	// Move to the start of the next day
	t = time.Date(t.Year(), t.Month(), t.Day()+1, s.BusinessHours.StartHour, 0, 0, 0, t.Location())

	// Keep moving forward until we hit a valid business day
	for !s.isValidDay(t) || s.isHoliday(t) {
		t = t.Add(24 * time.Hour)
	}

	// Return the time set to the start of the next valid business day
	return t
}

// getSLADuration converts the SLA length and time unit into a time.Duration
func (s SLA) getSLADuration() time.Duration {
	switch s.TimeUnit {
	case "seconds":
		return time.Duration(s.SLALength) * time.Second
	case "minutes":
		return time.Duration(s.SLALength) * time.Minute
	case "hours":
		return time.Duration(s.SLALength) * time.Hour
	case "days":
		return time.Duration(s.SLALength) * 24 * time.Hour
	default:
		return 0
	}
}

// isValidDay checks if the given time falls on a valid day according to the SLA
func (s SLA) isValidDay(t time.Time) bool {
	for _, day := range s.ValidDays {
		if t.Weekday() == day {
			return true
		}
	}
	return false
}

// isWithinBusinessHours checks if the given time is within the defined business hours
func (s SLA) isWithinBusinessHours(t time.Time) bool {
	hour := t.Hour()
	return hour >= s.BusinessHours.StartHour && hour < s.BusinessHours.EndHour
}

// isHoliday checks if the given time falls on a holiday
func (s SLA) isHoliday(t time.Time) bool {
	for _, holiday := range s.Holidays {
		if t.Year() == holiday.Year() && t.YearDay() == holiday.YearDay() {
			return true
		}
	}
	return false
}
