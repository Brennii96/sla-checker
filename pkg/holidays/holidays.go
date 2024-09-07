package holidays

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/brennii96/sla-checker/pkg/cache"
)

var APIBaseURL = "https://date.nager.at/Api/v3/PublicHolidays"

// Holiday represents the structure of the response from Date.nager.at API.
type Holiday struct {
	Date        string `json:"date"`
	LocalName   string `json:"localName"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

// Cache instance for holidays.
var holidayCache = cache.NewCache[[]time.Time](24 * 7 * time.Hour) // 1 Week TTL

// FetchHolidays dynamically fetches holidays for a specific year and country code,
// and caches the result to avoid redundant API calls.
func FetchHolidays(year int, countryCode string) ([]time.Time, error) {
	cacheKey := fmt.Sprintf("%d_%s", year, countryCode)

	// Try to get holidays from cache.
	if holidays, found := holidayCache.Get(cacheKey); found {
		return holidays, nil
	}

	// If not found in cache, fetch holidays from the API.
	url := fmt.Sprintf("%s/%d/%s", APIBaseURL, year, countryCode)
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

	// Store the fetched holidays in the cache.
	holidayCache.Set(cacheKey, holidays)

	return holidays, nil
}
