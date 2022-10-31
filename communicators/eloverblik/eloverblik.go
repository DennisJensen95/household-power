package eloverblik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// API Access interface for eloverblik.dk
type ShortAccessToken struct {
	AccessToken string `json:"result"`
}

type MeteringIdList struct {
	MeteringPoint []string `json:"meteringPoint"`
}

// Time series types below
type TimeSeriesData struct {
	Result []Result `json:"result"`
}
type SenderMarketParticipantMRID struct {
	CodingScheme interface{} `json:"codingScheme"`
	Name         interface{} `json:"name"`
}
type PeriodTimeInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type MRID struct {
	CodingScheme string `json:"codingScheme"`
	Name         string `json:"name"`
}
type MarketEvaluationPoint struct {
	MRID MRID `json:"mRID"`
}
type TimeInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type Point struct {
	Position            string `json:"position"`
	OutQuantityQuantity string `json:"out_Quantity.quantity"`
	OutQuantityQuality  string `json:"out_Quantity.quality"`
}
type Period struct {
	Resolution   string       `json:"resolution"`
	TimeInterval TimeInterval `json:"timeInterval"`
	Point        []Point      `json:"Point"`
}
type TimeSeries struct {
	MRID                  string                `json:"mRID"`
	BusinessType          string                `json:"businessType"`
	CurveType             string                `json:"curveType"`
	MeasurementUnitName   string                `json:"measurement_Unit.name"`
	MarketEvaluationPoint MarketEvaluationPoint `json:"MarketEvaluationPoint"`
	Period                []Period              `json:"Period"`
}
type MyEnergyDataMarketDocument struct {
	MRID                        string                      `json:"mRID"`
	CreatedDateTime             string                      `json:"createdDateTime"`
	SenderMarketParticipantName string                      `json:"sender_MarketParticipant.name"`
	SenderMarketParticipantMRID SenderMarketParticipantMRID `json:"sender_MarketParticipant.mRID"`
	PeriodTimeInterval          PeriodTimeInterval          `json:"period.timeInterval"`
	TimeSeries                  []TimeSeries                `json:"TimeSeries"`
}
type Result struct {
	MyEnergyDataMarketDocument MyEnergyDataMarketDocument `json:"MyEnergyData_MarketDocument"`
	Success                    bool                       `json:"success"`
	ErrorCode                  int                        `json:"errorCode"`
	ErrorText                  string                     `json:"errorText"`
	ID                         string                     `json:"id"`
	StackTrace                 interface{}                `json:"stackTrace"`
}

var eloverblik_base_api = "https://api.eloverblik.dk/CustomerApi"

type ElOverblikCommunicator struct {
	Client *http.Client
}

// GetToken .
func (c *ElOverblikCommunicator) GetAccessToken(authorization_token string) (string, error) {
	url := eloverblik_base_api + "/api/Token"
	req, _ := http.NewRequest("GET", url, nil)

	bearer := "Bearer " + authorization_token
	req.Header.Set("Authorization", bearer)

	resp, err := c.Client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"error":      err,
			"respStatus": resp.StatusCode,
		}).Error("Error while getting access token")
		return "", fmt.Errorf("error while sending GET request for access token")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var token ShortAccessToken
	err = decoder.Decode(&token)

	if err != nil {
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error while decoding access token")
		return "", err
	}
	return token.AccessToken, nil
}

func (c *ElOverblikCommunicator) GetTimeSeriesData(token string, metering_point_id string, start_date string, end_date string) (string, error) {
	// Setup body for metering point
	metering_points := MeteringIdList{[]string{metering_point_id}}
	values := map[string]MeteringIdList{"meteringPoints": metering_points}
	json_data, _ := json.Marshal(values)

	// Setup request
	url := eloverblik_base_api + fmt.Sprintf("/api/MeterData/GetTimeSeries/%s/%s/Hour", start_date, end_date)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json_data))

	err := c.setupRequestObject(req, token)
	if err != nil {
		log.Error("Error setting up request object")
		return "", err
	}

	// Send request
	resp, err := c.Client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Error("Error on response.\n[ERROR] -", err)
		return "", fmt.Errorf("error while sending GET request for time series data")
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error while reading response body")
	}

	return string(bytes), nil
}

// Set the required headers and access token to access the eloverblik API.
func (c *ElOverblikCommunicator) setupRequestObject(req *http.Request, access_token string) error {
	bearer := "Bearer " + access_token

	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return nil
}
