
# SLA Checker



## Authors

- [@brennii96](https://www.github.com/brennii96)


# SLA Checker

The `slachecker` package allows you to manage and evaluate Service Level Agreements (SLAs) by taking into account business hours, valid business days, and holidays. It provides functionality to determine if a given time is within the SLA duration and offers detailed compliance information.

## Features

- **`SLA` Struct**: Defines the SLA configuration.
- **`SLAResult` Struct**: Contains the result of the SLA evaluation.
- **Key Methods**:
  - `IsWithinSLA(currentTime time.Time) bool`
  - `CheckSLA(currentTime time.Time) SLAResult`


## Installation

cd into your project and init the module.

```bash
  cd project
  go mod init github.com/brennii96/sla-checker
```
## Usage/Examples

The holidays modules is using the [Date.nager.at]((https://github.com/nager/Nager.Date)) API to get public holidays for a given year and country code.

Import the required modules
```go
import (
    "github.com/brennii96/sla-checker/pkg/holidays"
	"github.com/brennii96/sla-checker/pkg/slachecker"
)
```

```go
// Parse SLA start and end times
startTime, err := time.Parse(time.RFC3339, reqBody.SLAStartTime)
if err != nil {
    http.Error(w, "Invalid SLA start time format", http.StatusBadRequest)
    return
}

// Fetch holidays (if required)
holidays, err := holidays.FetchHolidays(time.Now().Year(), reqBody.CountryCode)
if err != nil {
    http.Error(w, "Error fetching holidays", http.StatusInternalServerError)
    return
}

// Convert validDays strings to time.Weekday
var validDays []time.Weekday
for _, day := range reqBody.ValidDays {
    switch day {
    case "Monday":
        validDays = append(validDays, time.Monday)
    case "Tuesday":
        validDays = append(validDays, time.Tuesday)
    case "Wednesday":
        validDays = append(validDays, time.Wednesday)
    case "Thursday":
        validDays = append(validDays, time.Thursday)
    case "Friday":
        validDays = append(validDays, time.Friday)
    case "Saturday":
        validDays = append(validDays, time.Saturday)
    case "Sunday":
        validDays = append(validDays, time.Sunday)
    default:
        http.Error(w, "Invalid valid day", http.StatusBadRequest)
        return
    }
}

// Create SLA object
sla := slachecker.SLA{
    StartTime: startTime,
    SLALength: reqBody.SLALength,
    TimeUnit:  reqBody.TimeUnit,
    BusinessHours: struct {
        StartHour int
        EndHour   int
    }{
        StartHour: reqBody.BusinessHours.StartHour,
        EndHour:   reqBody.BusinessHours.EndHour,
    },
    ValidDays: validDays,
    Holidays:  holidays,
}
// Check if current time is within SLA
currentTime := time.Now()
// isWithinSLA := sla.IsWithinSLA(currentTime) // returns simple true/false
result := sla.CheckSLA(currentTime)

```
SLA
```go
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
```


CheckSLA result will be:
```go
// SLAResult contains the details about SLA status
type SLAResult struct {
	IsWithinSLA          bool      `json:"isWithinSLA"`
	Deadline             time.Time `json:"deadline"`
	Remaining            string    `json:"remaining"`
	Overage              string    `json:"overage,omitempty"`
	WorkingTimeRemaining string    `json:"workingTimeRemaining"`
}
```


## License

[MIT](./LICENSE.txt)
