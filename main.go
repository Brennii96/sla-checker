package main

import (
	"fmt"
	"log"
	"time"

	"github.com/brennii96/sla-checker/pkg/holidays"
	"github.com/brennii96/sla-checker/pkg/slachecker"
)

func main() {
	holidays, err := holidays.FetchHolidays(2024, "GB")
	if err != nil {
		log.Fatalf("Error fetching holidays: %v", err)
	}

	// Start time is Friday, August 30th, 2024 at 16:00 UTC
	startTime := time.Date(2024, time.August, 30, 16, 0, 0, 0, time.UTC)

	// Define the SLA length and unit
	slaLength := 4
	timeUnit := "hours" // SLA length of 4 hours

	// Define the SLA configuration
	sla := slachecker.SLA{
		StartTime: startTime,
		SLALength: slaLength,
		TimeUnit:  timeUnit,
		BusinessHours: struct {
			StartHour int
			EndHour   int
		}{StartHour: 9, EndHour: 17},
		ValidDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		Holidays:  holidays,
	}

	// Current time is Monday, September 2nd, 2024 at 12:00 UTC
	currentTime := time.Now()

	// Check if the current time is within the SLA
	if sla.IsWithinSLA(currentTime) {
		fmt.Println("The current time is within the SLA")
	} else {
		fmt.Println("The current time is not within the SLA")
	}
}
