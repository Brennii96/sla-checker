package holidays

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HolidayAPIResponse represents the structure of the response from Date.nager.at API
type Holiday struct {
	Date        string `json:"date"`
	LocalName   string `json:"localName"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

// FetchHolidays dynamically fetches holidays for Germany and a specific year
func FetchHolidays(year int, countryCode string) ([]time.Time, error) {
	url := fmt.Sprintf("https://date.nager.at/Api/v3/PublicHolidays/%d/%s", year, countryCode)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch holidays, status code: %d", resp.StatusCode)
	}

	var holidaysResp []Holiday
	err = json.NewDecoder(resp.Body).Decode(&holidaysResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	holidays := make([]time.Time, 0)
	for _, holiday := range holidaysResp {
		date, err := time.Parse("2006-01-02", holiday.Date)
		if err != nil {
			return nil, fmt.Errorf("error parsing date %s: %v", holiday.Date, err)
		}
		holidays = append(holidays, date)
	}

	return holidays, nil
}
