package energydatahub

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

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
		base_url: "https://api.energidataservice.dk",
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

	datahubCommunicator := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		assert.Equal(t, req.URL.String(), "https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start="+start_date+"&end="+end_date+"T00:00&filter=%7B%22PriceArea%22:[%22DK2%22]%7D&sort=HourUTC%20DESC&timezone=dk")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader(getTestResponse("../../tests/test_data/elspot_prices.json"))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	price_data, err := datahubCommunicator.FetchDataPricePerHour(start_date, end_date)

	assert.Nil(t, err)
	assert.Equal(t, price_data.Records[0].PriceArea, "DK2")

	// 2
	// Error case of getting a bad response with no data and bad exit code
	datahubCommunicator = NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		assert.Equal(t, req.URL.String(), "https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start="+start_date+"&end="+end_date+"T00:00&filter=%7B%22PriceArea%22:[%22DK2%22]%7D&sort=HourUTC%20DESC&timezone=dk")
		return &http.Response{
			StatusCode: 500,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	price_data, err = datahubCommunicator.FetchDataPricePerHour(start_date, end_date)

	assert.NotNil(t, err)
	assert.Equal(t, price_data.Total, 0)
}
