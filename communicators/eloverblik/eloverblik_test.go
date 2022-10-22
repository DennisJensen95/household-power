package eloverblik

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type MockHeaderElOverblikCommunicator struct {
	ElOverblikCommunicator
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) ElOverblikCommunicator {
	return ElOverblikCommunicator{
		Client: &http.Client{
			Transport: RoundTripFunc(fn),
		},
	}
}

func TestGetToken(t *testing.T) {
	// 1
	// Happy case of getting a token with expected response format
	elOverblikCommunicator := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		assert.Equal(t, req.URL.String(), "https://api.eloverblik.dk/CustomerApi/api/Token")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("{\"result\":\"testAccessToken\"}")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	token, err := elOverblikCommunicator.getAccessToken("testToken")
	assert.Nil(t, err)
	assert.Equal(t, token, "testAccessToken")

	// 2
	// Error case of getting a token with unexpected response format
	elOverblikCommunicator = NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		assert.Equal(t, req.URL.String(), "https://api.eloverblik.dk/CustomerApi/api/Token")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("{\"SomethingWeird\":\"testAccessToken\"")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	token, err = elOverblikCommunicator.getAccessToken("testToken")
	assert.NotNil(t, err)
	assert.Equal(t, token, "")

	// 3
	// HTTP request failed
	elOverblikCommunicator = NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		assert.Equal(t, req.URL.String(), "https://api.eloverblik.dk/CustomerApi/api/Token")
		return &http.Response{
			StatusCode: 500,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("{\"SomethingWeird\":\"testAccessToken\"")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	token, err = elOverblikCommunicator.getAccessToken("testToken")
	assert.NotNil(t, err)
	assert.Equal(t, token, "")
}

func getTestResponse(file string) string {
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return string(byteValue)
}

func TestGetMeteringPoints(t *testing.T) {
	// 1
	// Happy case of getting Data with expected response format
	start_date := "2020-01-01"
	end_date := "2020-01-02"
	metering_point := "testMeteringPointId"
	calledTimes := 0
	accessToken := "testAccessToken"
	testDataFile := "../../tests/test_data/time_series_response.json"

	elOverblikCommunicator := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters

		if calledTimes == 0 {
			assert.Equal(t, req.URL.String(), "https://api.eloverblik.dk/CustomerApi/api/Token")
			calledTimes++
			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(strings.NewReader("{\"result\":\"testAccessToken\"}")),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		assert.Equal(t, req.URL.String(), "https://api.eloverblik.dk/CustomerApi/api/MeterData/GetTimeSeries/"+start_date+"/"+end_date+"/Hour")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader(getTestResponse(testDataFile))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	response_body, err := elOverblikCommunicator.GetTimeSeriesData(accessToken, metering_point, start_date, end_date)
	assert.Nil(t, err)
	assert.Equal(t, response_body, getTestResponse(testDataFile))

	// 2
	// Error case of getting accessToken
	calledTimes = 0
	elOverblikCommunicator2 := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		calledTimes++
		log.Error(calledTimes)
		assert.Equal(t, req.URL.String(), "https://api.eloverblik.dk/CustomerApi/api/Token")
		return &http.Response{
			StatusCode: 500,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("{\"result\":\"testAccessToken\"")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	response_body, err = elOverblikCommunicator2.GetTimeSeriesData(accessToken, metering_point, start_date, end_date)
	assert.NotNil(t, err)
	assert.Equal(t, response_body, "")
}
