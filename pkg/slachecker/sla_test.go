package slachecker

import (
	"testing"
	"time"
)

func setupSLAWithHolidays(holidays []time.Time) SLA {
	return SLA{
		StartTime: time.Date(2024, time.September, 1, 9, 0, 0, 0, time.UTC),
		SLALength: 4,
		TimeUnit:  "hours",
		BusinessHours: struct {
			StartHour int
			EndHour   int
		}{
			StartHour: 9,
			EndHour:   17,
		},
		ValidDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		Holidays:  holidays,
	}
}

func TestCalculateWorkingTimeRemainingWithHolidays(t *testing.T) {
	// Define holidays
	holidays := []time.Time{
		time.Date(2024, time.August, 26, 0, 0, 0, 0, time.UTC), // Holiday on August 31, 2024
	}

	// Define the SLA with the holiday scenario
	sla := setupSLAWithHolidays(holidays)

	// Define the current time: 1 hour has passed on Friday (2024-08-30T17:00:00Z)
	currentTime := time.Date(2024, time.August, 30, 17, 0, 0, 0, time.UTC) // 5:00 PM on Friday

	// Define the expected result: 4 hours of working time should remain (2 hours on August 30, 2 hours on September 2)
	expectedWorkingTimeRemaining := "04:00:00"

	// Call the calculateSLADeadline function
	endTime, err := sla.calculateSLADeadline() // The deadline calculated from the SLA
	if err != nil {
		t.Fatalf("Error calculating SLA deadline: %v", err)
	}

	// Call the calculateWorkingTimeRemaining function
	workingTimeRemaining := sla.calculateWorkingTimeRemaining(currentTime, endTime)

	// Check if the result matches the expected value
	if workingTimeRemaining != expectedWorkingTimeRemaining {
		t.Errorf("Expected workingTimeRemaining to be %s, but got %s", expectedWorkingTimeRemaining, workingTimeRemaining)
	}
}

func TestIsWithinSLAWithHolidays(t *testing.T) {
	// Define holidays
	holidays := []time.Time{
		time.Date(2024, time.August, 26, 0, 0, 0, 0, time.UTC), // Holiday on August 26, 2024
	}

	// Define common SLA configuration for the tests with holidays
	sla := setupSLAWithHolidays(holidays)

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
			currentTime: time.Date(2024, time.September, 2, 13, 0, 0, 0, time.UTC), // Monday 13
			expected:    false,
		},
		{
			// Test case: holiday (SLA should not count holidays)
			startTime:   time.Date(2024, time.August, 23, 16, 0, 0, 0, time.UTC), // Thursday 4 PM
			slaLength:   4,
			timeUnit:    "hours",
			currentTime: time.Date(2024, time.August, 27, 10, 0, 0, 0, time.UTC), // Monday 10 AM
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

func TestCheckSLAWithHolidays(t *testing.T) {
	// Define holidays
	holidays := []time.Time{
		time.Date(2024, time.August, 26, 0, 0, 0, 0, time.UTC), // Holiday on August 26, 2024
	}

	// Define common SLA configuration for the tests with holidays
	sla := setupSLAWithHolidays(holidays)

	tests := []struct {
		startTime           time.Time
		slaLength           int
		timeUnit            string
		currentTime         time.Time
		expectedIsWithinSLA bool
		expectedDeadline    time.Time
		expectedRemaining   string
		expectedOverage     string
	}{
		{
			// Test case: within business hours on a working day
			startTime:           time.Date(2024, time.August, 30, 9, 0, 0, 0, time.UTC), // Friday 9 AM
			slaLength:           4,
			timeUnit:            "hours",
			currentTime:         time.Date(2024, time.August, 30, 12, 0, 0, 0, time.UTC), // Friday 12 PM
			expectedIsWithinSLA: true,
			expectedDeadline:    time.Date(2024, time.August, 30, 13, 0, 0, 0, time.UTC), // Deadline 1 PM
			expectedRemaining:   "01:00:00",                                              // 1 hour remaining
			expectedOverage:     "00:00:00",                                              // No overage
		},
		{
			// Test case: outside business hours on a working day
			startTime:           time.Date(2024, time.August, 30, 9, 0, 0, 0, time.UTC), // Friday 9 AM
			slaLength:           4,
			timeUnit:            "hours",
			currentTime:         time.Date(2024, time.August, 30, 18, 0, 0, 0, time.UTC), // Friday 6 PM
			expectedIsWithinSLA: false,
			expectedDeadline:    time.Date(2024, time.August, 30, 13, 0, 0, 0, time.UTC), // Next valid deadline
			expectedRemaining:   "00:00:00",                                              // Adjust based on your `formatDuration` output
			expectedOverage:     "05:00:00",                                              // No overage
		},
		{
			// Test case: during weekend (SLA should not count weekends)
			startTime:           time.Date(2024, time.August, 30, 16, 0, 0, 0, time.UTC), // Friday 4 PM
			slaLength:           4,
			timeUnit:            "hours",
			currentTime:         time.Date(2024, time.September, 2, 13, 0, 0, 0, time.UTC), // Monday 1 PM
			expectedIsWithinSLA: false,
			expectedDeadline:    time.Date(2024, time.September, 02, 12, 0, 0, 0, time.UTC), // Next valid deadline
			expectedRemaining:   "00:00:00",                                                 // Adjust based on your `formatDuration` output
			expectedOverage:     "01:00:00",                                                 // No overage
		},
		{
			// Test case: holiday (SLA should not count holidays)
			startTime:           time.Date(2024, time.August, 23, 16, 0, 0, 0, time.UTC), // Friday 4 PM
			slaLength:           4,
			timeUnit:            "hours",
			currentTime:         time.Date(2024, time.August, 27, 9, 0, 0, 0, time.UTC), // Tuesday 10 AM
			expectedIsWithinSLA: true,
			expectedDeadline:    time.Date(2024, time.August, 27, 12, 0, 0, 0, time.UTC), // Deadline on Tuesday 4 PM
			expectedRemaining:   "03:00:00",                                              // Adjust based on your `formatDuration` output
			expectedOverage:     "00:00:00",                                              // No overage
		},
	}

	for _, test := range tests {
		// Set up SLA with the test parameters
		sla.StartTime = test.startTime
		sla.SLALength = test.slaLength
		sla.TimeUnit = test.timeUnit

		// Run the CheckSLA method
		result := sla.CheckSLA(test.currentTime)

		// Compare the result with the expected outcome
		if result.IsWithinSLA != test.expectedIsWithinSLA {
			t.Errorf("Failed test: startTime=%v, slaLength=%d, timeUnit=%s, currentTime=%v, expectedIsWithinSLA=%v, got=%v",
				test.startTime, test.slaLength, test.timeUnit, test.currentTime, test.expectedIsWithinSLA, result.IsWithinSLA)
		}
		if !result.Deadline.Equal(test.expectedDeadline) {
			t.Errorf("Failed test: startTime=%v, slaLength=%d, timeUnit=%s, currentTime=%v, expectedDeadline=%v, got=%v",
				test.startTime, test.slaLength, test.timeUnit, test.currentTime, test.expectedDeadline, result.Deadline)
		}
		if result.Remaining != test.expectedRemaining {
			t.Errorf("Failed test: startTime=%v, slaLength=%d, timeUnit=%s, currentTime=%v, expectedRemaining=%v, got=%v",
				test.startTime, test.slaLength, test.timeUnit, test.currentTime, test.expectedRemaining, result.Remaining)
		}
		if result.Overage != test.expectedOverage {
			t.Errorf("Failed test: startTime=%v, slaLength=%d, timeUnit=%s, currentTime=%v, expectedOverage=%v, got=%v",
				test.startTime, test.slaLength, test.timeUnit, test.currentTime, test.expectedOverage, result.Overage)
		}
	}
}

func TestCalculateWorkingTimeIgnoringHolidays(t *testing.T) {
	// Define holidays (which will be ignored)
	holidays := []time.Time{
		time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC), // Holiday on September 2, 2024
	}

	// Define the SLA, ignoring holidays
	sla := setupSLAWithHolidays(holidays)
	sla.IgnoreHolidays = true

	// Define the start time: Monday, September 2, 2024, at 9:00 AM
	sla.StartTime = time.Date(2024, time.September, 2, 9, 0, 0, 0, time.UTC)
	sla.SLALength = 4 // SLA length is 4 hours
	sla.TimeUnit = "hours"

	// Define the current time: Monday, September 2, 2024, at 11:00 AM
	currentTime := time.Date(2024, time.September, 2, 11, 0, 0, 0, time.UTC)

	// Expected results
	expectedDeadline := time.Date(2024, time.September, 2, 13, 0, 0, 0, time.UTC) // Deadline should be 1:00 PM
	expectedIsWithinSLA := true
	expectedRemaining := "02:00:00" // 2 hours remaining
	expectedOverage := "00:00:00"   // No overage

	// Run the CheckSLA method
	result := sla.CheckSLA(currentTime)

	// Validate the results
	if result.IsWithinSLA != expectedIsWithinSLA {
		t.Errorf("Expected IsWithinSLA to be %v, but got %v", expectedIsWithinSLA, result.IsWithinSLA)
	}
	if !result.Deadline.Equal(expectedDeadline) {
		t.Errorf("Expected deadline to be %v, but got %v", expectedDeadline, result.Deadline)
	}
	if result.Remaining != expectedRemaining {
		t.Errorf("Expected remaining time to be %v, but got %v", expectedRemaining, result.Remaining)
	}
	if result.Overage != expectedOverage {
		t.Errorf("Expected overage time to be %v, but got %v", expectedOverage, result.Overage)
	}
}
