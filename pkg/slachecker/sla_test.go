package slachecker

import (
	"testing"
	"time"
)

func TestCalculateWorkingTimeRemaining(t *testing.T) {
	// Define the SLA with the test scenario
	sla := SLA{
		StartTime: time.Date(2024, time.August, 30, 16, 0, 0, 0, time.UTC), // 4:00 PM on Friday
		SLALength: 4,                                                       // 4 hours SLA
		TimeUnit:  "hours",
		BusinessHours: struct {
			StartHour int
			EndHour   int
		}{
			StartHour: 9,  // Business hours start at 9:00 AM
			EndHour:   17, // Business hours end at 5:00 PM
		},
		ValidDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		Holidays:  []time.Time{}, // No holidays for simplicity
	}

	// Define the current time: 1 hour has passed on Friday (2024-08-30T17:00:00Z)
	currentTime := time.Date(2024, time.August, 30, 17, 0, 0, 0, time.UTC) // 5:00 PM on Friday

	// Define the expected result: 3 hours of working time should remain
	expectedWorkingTimeRemaining := "03:00:00"

	// Call the calculateWorkingTimeRemaining function
	endTime := sla.calculateSLADeadline() // The deadline calculated from the SLA
	workingTimeRemaining := sla.calculateWorkingTimeRemaining(currentTime, endTime)

	// Check if the result matches the expected value
	if workingTimeRemaining != expectedWorkingTimeRemaining {
		t.Errorf("Expected workingTimeRemaining to be %s, but got %s", expectedWorkingTimeRemaining, workingTimeRemaining)
	}
}

func TestIsWithinSLA(t *testing.T) {
	// Define common SLA configuration for the tests
	sla := SLA{
		BusinessHours: struct {
			StartHour int
			EndHour   int
		}{StartHour: 9, EndHour: 17},
		ValidDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		Holidays:  []time.Time{time.Date(2024, time.July, 4, 0, 0, 0, 0, time.UTC)}, // Example holiday
	}

	tests := []struct {
		startTime   time.Time
		slaLength   int
		timeUnit    string
		currentTime time.Time
		expected    bool
	}{
		{
			// Test case: within business hours on a working day
			startTime:   time.Date(2024, time.August, 30, 9, 0, 0, 0, time.UTC), // Friday 9 AM
			slaLength:   4,
			timeUnit:    "hours",
			currentTime: time.Date(2024, time.August, 30, 12, 0, 0, 0, time.UTC), // Friday 12 PM
			expected:    true,
		},
		{
			// Test case: outside business hours on a working day
			startTime:   time.Date(2024, time.August, 30, 9, 0, 0, 0, time.UTC), // Friday 9 AM
			slaLength:   4,
			timeUnit:    "hours",
			currentTime: time.Date(2024, time.August, 30, 18, 0, 0, 0, time.UTC), // Friday 6 PM
			expected:    false,
		},
		{
			// Test case: during weekend (SLA should not count weekends)
			startTime:   time.Date(2024, time.August, 30, 16, 0, 0, 0, time.UTC), // Friday 4 PM
			slaLength:   4,
			timeUnit:    "hours",
			currentTime: time.Date(2024, time.September, 2, 12, 0, 0, 0, time.UTC), // Monday 12 PM
			expected:    false,
		},
		{
			// Test case: holiday (SLA should not count holidays)
			startTime:   time.Date(2024, time.July, 3, 16, 0, 0, 0, time.UTC), // Wednesday 4 PM
			slaLength:   4,
			timeUnit:    "hours",
			currentTime: time.Date(2024, time.July, 5, 10, 0, 0, 0, time.UTC), // Friday 10 AM
			expected:    true,
		},
	}

	for _, test := range tests {
		// Set up SLA with the test parameters
		sla.StartTime = test.startTime
		sla.SLALength = test.slaLength
		sla.TimeUnit = test.timeUnit

		// Run the IsWithinSLA check
		result := sla.IsWithinSLA(test.currentTime)

		// Compare the result with the expected outcome
		if result != test.expected {
			t.Errorf("Failed test: startTime=%v, slaLength=%d, timeUnit=%s, currentTime=%v, expected=%v, got=%v",
				test.startTime, test.slaLength, test.timeUnit, test.currentTime, test.expected, result)
		}
	}
}
