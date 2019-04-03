package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EMT struct {
	idClient string
	passKey  string
}

func New(idClient, passKey string) EMT {
	return EMT{
		idClient: idClient,
		passKey:  passKey,
	}
}

func (e EMT) GetStopEstimates(stop int) ([]Bus, error) {
	api := "https://openbus.emtmadrid.es:9443/emt-proxy-server/last/media/GetEstimatesIncident.php"
	data := url.Values{}

	data.Set("idClient", e.idClient)
	data.Set("passKey", e.passKey)
	data.Set("idStop", strconv.Itoa(stop))
	data.Set("Text_StopRequired_YN", "N")
	data.Set("Audio_StopRequired_YN", "N")
	data.Set("Text_EstimationsRequired_YN", "Y")
	data.Set("Audio_EstimationsRequired_YN", "N")
	data.Set("Audio_IncidencesRequired_YN", "N")
	data.Set("Audio_StopRequired_YN", "N")
	data.Set("cultureInfo", "EN")

	c := &http.Client{}
	r, err := http.NewRequest("POST", api, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server responded unsuccessfuly")
	}

	var m map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&m)

	if err != nil {
		return nil, err
	}

	// First, check error code
	errcode, err := strconv.Atoi(m["errorCode"].(string))

	if errcode != 0 || err != nil {
		return nil, fmt.Errorf("server responded unsuccessfuly. errorCode: %v, err: %v", errcode, err)
	}

	var arrives []Bus

	arrMap := m["arrives"].(map[string]interface{})["arriveEstimationList"].(map[string]interface{})

	if arrive, ok := arrMap["arrive"].(map[string]interface{}); ok {
		// Only one arrive.
		arrives = append(arrives, Bus{
			LineID:      arrive["lineId"].(string),
			Destination: arrive["destination"].(string),
			BusID:       arrive["busId"].(string),
			TimeLeft:    time.Second * time.Duration((math.Min(25*60, arrive["busTimeLeft"].(float64)))),
			Distance:    int(arrive["busDistance"].(float64)),
			Latitude:    arrive["latitude"].(float64),
			Longitude:   arrive["longitude"].(float64),
		})
	} else if arriveList, ok := arrMap["arrive"].([]interface{}); ok {
		for _, a := range arriveList {
			arrive := a.(map[string]interface{})
			arrives = append(arrives, Bus{
				LineID:      arrive["lineId"].(string),
				Destination: arrive["destination"].(string),
				BusID:       arrive["busId"].(string),
				TimeLeft:    time.Second * time.Duration((math.Min(25*60, arrive["busTimeLeft"].(float64)))),
				Distance:    int(arrive["busDistance"].(float64)),
				Latitude:    arrive["latitude"].(float64),
				Longitude:   arrive["longitude"].(float64),
			})
		}
	}

	return arrives, nil
}
