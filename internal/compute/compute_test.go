package compute

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/DennisJensen95/household-power/communicators/eloverblik"
	"github.com/DennisJensen95/household-power/communicators/energydatahub"
	"github.com/stretchr/testify/assert"
)

func TestCalculationOfKWH(t *testing.T) {
	// Read the file
	jsonFile, err := os.Open("../../tests/test_data/time_series_response.json")
	if err != nil {
		t.Fatal("Failed to open test file, err: ", err)
	}

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read the mocked data")
	}

	var data eloverblik.TimeSeriesData
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal the mocked data")
	}

	total_kwh := CalculateThePeriodUseOfKWH(data)
	if total_kwh != 150 {
		t.Fatalf("Expected total kWh to be 0.0, got %f", total_kwh)
	}
}

func TestParsingDateFromEnergyDataHub(t *testing.T) {
	date, err := parseDateFromEnergyDataHub("2022-10-23T00:00:00Z")
	assert.Nil(t, err)
	assert.Equal(t, 0, date.Hour())
	assert.Equal(t, 0, date.Minute())
	assert.Equal(t, 0, date.Second())
	assert.Equal(t, 23, date.Day())
	assert.Equal(t, 10, int(date.Month()))
	assert.Equal(t, 2022, date.Year())
}

func TestParsingDateFromElOverblikl(t *testing.T) {
	date, err := parseDateFromElOverblik("2022-10-23T00:00:00")
	assert.Nil(t, err)
	assert.Equal(t, 0, date.Hour())
	assert.Equal(t, 0, date.Minute())
	assert.Equal(t, 0, date.Second())
	assert.Equal(t, 23, date.Day())
	assert.Equal(t, 10, int(date.Month()))
	assert.Equal(t, 2022, date.Year())
}

func TestGettingPricesBetweenDates(t *testing.T) {
	// Setup

	// Read the file
	jsonFile, err := os.Open("../../tests/test_data/elspot_prices.json")
	if err != nil {
		t.Fatal("Failed to open test file, err: ", err)
	}

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read the mocked data")
	}

	// Encode the data
	var data energydatahub.PriceDataRecords
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal the mocked data")
	}

	// 1
	// Happy case of a single day extraction

	prices, err := getPricesBetweenDates("2022-10-24T00:00:00Z", "2022-10-25T00:00:00Z", data)

	assert.Equal(t, 24, len(prices.Records))
	assert.Nil(t, err)

	// 2
	// Multiple days extraction
	prices, err = getPricesBetweenDates("2022-10-22T22:00:00Z", "2022-10-24T22:00:00Z", data)
	assert.Nil(t, err)
	assert.Equal(t, 48, len(prices.Records))
}

func TestCostCalculation(t *testing.T) {
	// Read the file
	jsonFile, err := os.Open("../../tests/test_data/power_consumption_24_25.json")
	if err != nil {
		t.Fatal("Failed to open test file, err: ", err)
	}

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read the mocked data")
	}

	var data eloverblik.TimeSeriesData
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal the mocked data")
	}

	jsonFilePrices, err := os.Open("../../tests/test_data/elspot_prices.json")
	if err != nil {
		t.Fatal("Failed to open test file, err: ", err)
	}

	c, err := ioutil.ReadAll(jsonFilePrices)
	if err != nil {
		t.Fatalf("Failed to read the mocked data")
	}

	var data_prices energydatahub.PriceDataRecords
	err = json.Unmarshal(c, &data_prices)
	if err != nil {
		t.Fatalf("Failed to unmarshal the mocked data")
	}

	total_cost, err := CalculateVariablePrice(data, data_prices, 1.35)
	assert.Equalf(t, 30.9, math.Round(total_cost*10)/10, "Expected total cost to be 30.7, got %f", math.Round(total_cost*10)/10)
	assert.Nil(t, err)
}
