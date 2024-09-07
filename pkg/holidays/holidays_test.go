package holidays_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brennii96/sla-checker/pkg/holidays"
)

// Mock data for holidays
var mockHolidays = []holidays.Holiday{
	{Date: "2023-01-01", LocalName: "New Year's Day", Name: "New Year's Day", CountryCode: "DE"},
	{Date: "2023-12-25", LocalName: "Christmas Day", Name: "Christmas Day", CountryCode: "DE"},
}

// Helper function to create a mock server
func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with the mock holidays data
		json.NewEncoder(w).Encode(mockHolidays)
	}))
}

func TestFetchHolidays(t *testing.T) {
	// Setup the mock server
	server := setupMockServer()
	defer server.Close()

	// Override the API URL with the mock server URL
	holidays.APIBaseURL = server.URL // You will need to define this in the holidays package

	// Fetch holidays for the first time (should hit the mock API)
	holidaysData, err := holidays.FetchHolidays(2023, "DE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(holidaysData) != len(mockHolidays) {
		t.Errorf("expected %d holidays, got %d", len(mockHolidays), len(holidaysData))
	}

	// Fetch holidays for the second time (should hit the cache)
	holidaysData, err = holidays.FetchHolidays(2023, "DE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(holidaysData) != len(mockHolidays) {
		t.Errorf("expected %d holidays from cache, got %d", len(mockHolidays), len(holidaysData))
	}
}
