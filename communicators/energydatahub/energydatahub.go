package energydatahub

import (
	"encoding/json"
	"errors"
	"net/http"
)

type EnergyDatahub interface {
	FetchDataPricePerHour(start_date string, end_date string) (string, error)
}

type EnergyDatahubCommunicator struct {
	base_url string
	Client   *http.Client
}

type PriceDataRecords struct {
	Total   int    `json:"total"`
	Filters string `json:"filters"`
	Sort    string `json:"sort"`
	Dataset string `json:"dataset"`
	Records []struct {
		HourUTC      string  `json:"HourUTC"`
		HourDK       string  `json:"HourDK"`
		PriceArea    string  `json:"PriceArea"`
		SpotPriceDKK float64 `json:"SpotPriceDKK"`
		SpotPriceEUR float64 `json:"SpotPriceEUR"`
	} `json:"records"`
}

func (communicator EnergyDatahubCommunicator) FetchDataPricePerHour(start_date string, end_date string) (PriceDataRecords, error) {
	// https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start=2022-10-01T00:00&end=2022-10-27T00:00&filter=%7B%22PriceArea%22:[%22DK2%22]%7D&sort=HourUTC%20DESC&timezone=dk
	url := communicator.base_url + "/dataset/Elspotprices?offset=0&start=" + start_date + "&end=" + end_date + "T00:00&filter=%7B%22PriceArea%22:[%22DK2%22]%7D&sort=HourUTC%20DESC&timezone=dk"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := communicator.Client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return PriceDataRecords{}, errors.New("error while sending GET request for access token")
	}
	defer resp.Body.Close()
	var data PriceDataRecords
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return PriceDataRecords{}, err
	}
	return data, nil
}
