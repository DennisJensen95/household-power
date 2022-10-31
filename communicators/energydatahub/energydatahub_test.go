package energydatahub

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) EnergyDatahubCommunicator {
	return EnergyDatahubCommunicator{
		Client: &http.Client{
			Transport: RoundTripFunc(fn),
		},
		BaseUrl: "https://api.energidataservice.dk",
	}
}

func getTestResponse(file string) string {
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return string(byteValue)
}

func TestFetchDataPricePerHour(t *testing.T) {
	// 1
	// Happy case of getting a nice response with data and good exit code
	start_date := "2022-10-24"
	end_date := "2022-10-25"

	start_date, _ = AddDaysToTimeString(start_date, -1)
	end_date, _ = AddDaysToTimeString(end_date, 1)

	datahubCommunicator := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		assert.Equal(t, req.URL.String(), "https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start="+start_date+"&end="+end_date+"&filter=%7B%22PriceArea%22:[%22DK2%22]%7D&sort=HourUTC%20DESC&timezone=dk")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader(getTestResponse("../../tests/test_data/elspot_prices.json"))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	price_data, err := datahubCommunicator.FetchDataPrice(start_date, end_date)

	assert.Nil(t, err)
	assert.Equal(t, price_data.Records[0].PriceArea, "DK2")

	// 2
	// Error case of getting a bad response with no data and bad exit code
	datahubCommunicator = NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		assert.Equal(t, req.URL.String(), "https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start="+start_date+"&end="+end_date+"&filter=%7B%22PriceArea%22:[%22DK2%22]%7D&sort=HourUTC%20DESC&timezone=dk")
		return &http.Response{
			StatusCode: 500,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	price_data, err = datahubCommunicator.FetchDataPrice(start_date, end_date)

	assert.NotNil(t, err)
	assert.Equal(t, len(price_data.Records), 0)
}

func parseDateFromEnergyDataHub(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

// Not executed unless integration tests are called
func IntegrationTestFetchingPriceData(t *testing.T) {
	communicator := EnergyDatahubCommunicator{
		BaseUrl: "https://api.energidataservice.dk",
		Client:  &http.Client{},
	}

	start_date := "2022-09-01"
	end_date := "2022-09-30"
	data, err := communicator.FetchDataPrice(start_date, end_date)
	assert.Nil(t, err)

	log.Info("Total number of records: ", len(data.Records))

	start_date_time, _ := parseDateFromEnergyDataHub(start_date)
	end_date_time, _ := parseDateFromEnergyDataHub(end_date)

	total_hours := end_date_time.Sub(start_date_time).Hours()

	assert.Nil(t, err)
	assert.Equalf(t, 24*29, int(total_hours), "Total hours should be %f, but was %d", total_hours, 24*29)

	// Add two days as the fetch adds buffer in price data fetch on a day on each side.
	assert.Equalf(t, int(total_hours), len(data.Records), "There should be the same number of hours %f as price data points %d", total_hours, len(data.Records))

	log.Info(len(data.Records))
}
