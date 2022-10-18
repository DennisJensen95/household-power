package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

type MeteringIdList struct {
	MeteringPoint []string `json:"meteringPoint"`
}

func main() {
	const token = "xxx"
	const meteringPoint = "xxx"

	// Setup body for metering point
	metering_points := MeteringIdList{[]string{meteringPoint}}
	values := map[string]MeteringIdList{"meteringPoints": metering_points}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	bearer := "Bearer " + token
	url := "https://api.eloverblik.dk/CustomerApi//api/MeterData/GetTimeSeries/2022-10-10/2022-10-18/Hour"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal("Error creating request")
	}

	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	log.Println("Statuscode is: %d", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string(body))

}
