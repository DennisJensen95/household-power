package compute

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/DennisJensen95/household-power/communicators/eloverblik"
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
