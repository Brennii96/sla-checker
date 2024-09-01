package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/brennii96/sla-checker/pkg/holidays"
	"github.com/brennii96/sla-checker/pkg/slachecker"
)

func main() {
	// Fetch holidays concurrently
	holidaysCh := make(chan []time.Time)
	errorCh := make(chan error)

	go func() {
		fetchedHolidays, err := holidays.FetchHolidays(2024, "GB")
		if err != nil {
			errorCh <- err
			return
		}
		holidaysCh <- fetchedHolidays
	}()

	// Initialize start time, SLA settings, and valid days
	startTime := time.Date(2024, time.August, 30, 16, 0, 0, 0, time.UTC)
	validDays := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
	}

	// Create SLA configuration
	sla := slachecker.SLA{
		StartTime: startTime,
		SLALength: 4,
		TimeUnit:  "hours", // SLA length of 4 hours
		BusinessHours: struct {
			StartHour int
			EndHour   int
		}{
			StartHour: 9,
			EndHour:   17,
		},
		ValidDays: validDays,
	}

	// Wait for the holidays to be fetched
	select {
	case err := <-errorCh:
		log.Fatalf("Error fetching holidays: %v", err)
	case fetchedHolidays := <-holidaysCh:
		sla.Holidays = fetchedHolidays
	}

	// Check SLA with current time
	result := sla.CheckSLA(time.Now().UTC())

	// Marshal and print JSON result
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
