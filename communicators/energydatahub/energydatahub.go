package energydatahub

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type EnergyDatahub interface {
	FetchDataPricePerHour(start_date string, end_date string) (string, error)
}

type EnergyDatahubCommunicator struct {
	BaseUrl string
	Client  *http.Client
}

type PriceDataRecords struct {
	Total   int      `json:"total"`
	Filters string   `json:"filters"`
	Sort    string   `json:"sort"`
	Dataset string   `json:"dataset"`
	Records []Record `json:"records"`
}
type Record struct {
	HourUTC      string  `json:"HourUTC"`
	HourDK       string  `json:"HourDK"`
	PriceArea    string  `json:"PriceArea"`
	SpotPriceDKK float64 `json:"SpotPriceDKK"`
	SpotPriceEUR float64 `json:"SpotPriceEUR"`
}

func AddDaysToTimeString(timeString string, days int) (string, error) {
	timestamp, err := time.Parse("2006-01-02", timeString)
	if err != nil {
		log.WithField("error", err).Error("Failed to parse time string")
		return "", err
	}
	timestamp = timestamp.AddDate(0, 0, days)
	return timestamp.Format("2006-01-02"), nil
}

func (communicator EnergyDatahubCommunicator) FetchDataPrice(start_date string, end_date string) (PriceDataRecords, error) {
	url := communicator.BaseUrl + "/dataset/Elspotprices?offset=0&start=" + start_date + "&end=" + end_date + "&filter=%7B%22PriceArea%22:[%22DK2%22]%7D&sort=HourUTC%20DESC&timezone=dk"
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
