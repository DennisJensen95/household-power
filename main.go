package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// access_token = self._get_access_token()
// headers = self._create_headers(access_token)
// body = '{"meteringPoints": {"meteringPoint": ["' + metering_point + '"]}}'
// url = self._base_url + '/api/meteringpoints/meteringpoint/getcharges'
//
// response = requests.post(url,
//  data=body,
//  headers=headers,
//  timeout=5
//  )

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

var (
	token           = flag.String("token", "", "access-token")
	meteringPointId = flag.String("meteringPointId", "", "meter-point-id")
)

func _get_access_token(authorization_token string) string {
	url := eloverblik_base_api + "/api/Token"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	bearer := "Bearer " + authorization_token
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	decoder := json.NewDecoder(resp.Body)
	var token ShortAccessToken
	err = decoder.Decode(&token)

	if err != nil {
		log.Println("Error while decoding the response bytes:", err)
	}

	return token.AccessToken
}

func get_time_series_data(token string, metering_point_id string) {
	// Setup body for metering point
	metering_points := MeteringIdList{[]string{*meteringPointId}}
	values := map[string]MeteringIdList{"meteringPoints": metering_points}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	// Setup request
	access_token := _get_access_token(token)
	bearer := "Bearer " + access_token
	url := eloverblik_base_api + "/api/MeterData/GetTimeSeries/2022-10-10/2022-10-18/Hour"

	fmt.Println("Bearer: ", bearer)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal("Error creating request")
	}

	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	decoder := json.NewDecoder(resp.Body)
	var timeSeries TimeSeriesData
	err = decoder.Decode(&timeSeries)

	if err != nil {
		log.Println("Error while decoding the response bytes:", err)
	}

	fmt.Println(timeSeries)
}

func init_parser() {
	flag.Parse()

	if *token == "" {
		flag.Usage()
		log.Fatal("Please provide an access-token")
	}

	if *meteringPointId == "" {
		flag.Usage()
		log.Fatal("Please provide a metering-point-id")
	}
}

func main() {
	init_parser()

	get_time_series_data(*token, *meteringPointId)
}
