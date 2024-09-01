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
	fetchedHolidays, err := holidays.FetchHolidays(2024, "GB")
	if err != nil {
		log.Fatalf("Error fetching holidays: %v", err)
	}

	// Start time is Friday, August 30th, 2024 at 16:00 UTC
	startTime := time.Date(2024, time.August, 30, 16, 0, 0, 0, time.UTC)

	// Define the SLA length and unit
	slaLength := 4
	timeUnit := "hours" // SLA length of 4 hours

	var validDays []time.Weekday

	validDays = append(validDays, time.Monday)
	validDays = append(validDays, time.Tuesday)
	validDays = append(validDays, time.Wednesday)
	validDays = append(validDays, time.Thursday)
	validDays = append(validDays, time.Friday)

	sla := slachecker.SLA{
		StartTime: startTime,
		SLALength: slaLength,
		TimeUnit:  timeUnit,
		BusinessHours: struct {
			StartHour int
			EndHour   int
		}{
			StartHour: 9,
			EndHour:   17,
		},
		ValidDays: validDays,
		Holidays:  fetchedHolidays,
	}

	// Current time is Monday, September 2nd, 2024 at 12:00 UTC
	currentTime := time.Now().UTC()

	result := sla.CheckSLA(currentTime)
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
