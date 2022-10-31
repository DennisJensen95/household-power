package power_usage

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DennisJensen95/household-power/communicators/eloverblik"
	"github.com/DennisJensen95/household-power/communicators/energydatahub"
	"github.com/DennisJensen95/household-power/internal/compute"
	log "github.com/sirupsen/logrus"
)

// CalculateThePeriodUseOfKWH calculates the total kWh used in a given period
func PeriodUseKWH(token string, metering_point_id string, start_date string, end_date string) (float64, error) {
	timeSeries, err := periodUseKwhHourly(token, metering_point_id, start_date, end_date)

	if err != nil {
		log.WithField("error", err).Error("Failed to get time series data")
		return 0, err
	}

	total_kwh := compute.CalculateThePeriodUseOfKWH(timeSeries)
	return total_kwh, nil
}

// CalculateThePeriodUseOfKWH calculates the cost of the total kWh used in a given period
func PeriodUsePriceVariable(token string, metering_point_id string, start_date string, end_date string, expense float64) (float64, error) {
	timeSeries, err := periodUseKwhHourly(token, metering_point_id, start_date, end_date)

	if err != nil {
		log.WithField("error", err).Error("Failed to get time series data")
		return 0, err
	}

	energyDatahubCommunicator := energydatahub.EnergyDatahubCommunicator{
		Client:  &http.Client{},
		BaseUrl: "https://api.energidataservice.dk",
	}

	end_date_fetch_data, err := energydatahub.AddDaysToTimeString(end_date, 1)
	if err != nil {
		log.WithField("error", err).Error("Failed to add days to time string")
		return 0, err
	}

	start_date_fetch_data, err := energydatahub.AddDaysToTimeString(start_date, -1)
	if err != nil {
		log.WithField("error", err).Error("Failed to add days to time string")
		return 0, err
	}
	priceSeries, err := energyDatahubCommunicator.FetchDataPrice(start_date_fetch_data, end_date_fetch_data)

	if err != nil {
		log.WithField("error", err).Error("Failed to get price series data")
		return 0, err
	}

	spend_money, err := compute.CalculateVariablePrice(timeSeries, priceSeries, expense)

	if err != nil {
		log.WithField("error", err).Error("Failed to calculate variable price")
		return 0, err
	}

	return spend_money, nil
}

func periodUseKwhHourly(token string, metering_point_id string, start_date string, end_date string) (eloverblik.TimeSeriesData, error) {
	elOverblikCommunicator := eloverblik.ElOverblikCommunicator{
		Client: &http.Client{},
	}

	response, err := elOverblikCommunicator.GetTimeSeriesData(token, metering_point_id, start_date, end_date)
	if err != nil {
		log.Error("Failed to get time series data, err: ", err)
		return eloverblik.TimeSeriesData{}, err
	}

	decoder := json.NewDecoder(strings.NewReader(response))
	var timeSeries eloverblik.TimeSeriesData
	err = decoder.Decode(&timeSeries)

	if err != nil {
		log.Println("Error while decoding the response bytes:", err)
	}

	return timeSeries, nil
}
