package power_usage

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DennisJensen95/household-power/communicators/eloverblik"
	"github.com/DennisJensen95/household-power/internal/compute"
	log "github.com/sirupsen/logrus"
)

// CalculateThePeriodUseOfKWH calculates the total kWh used in a given period
func PeriodUseKWH(token string, metering_point_id string, start_date string, end_date string) (float64, error) {

	elOverblikCommunicator := eloverblik.ElOverblikCommunicator{
		Client: &http.Client{},
	}

	response, err := elOverblikCommunicator.GetTimeSeriesData(token, metering_point_id, start_date, end_date)
	if err != nil {
		log.Error("Failed to get time series data, err: ", err)
		return 0, err
	}

	decoder := json.NewDecoder(strings.NewReader(response))
	var timeSeries eloverblik.TimeSeriesData
	err = decoder.Decode(&timeSeries)

	if err != nil {
		log.Println("Error while decoding the response bytes:", err)
	}

	total_kwh := compute.CalculateThePeriodUseOfKWH(timeSeries)
	return total_kwh, nil
}
