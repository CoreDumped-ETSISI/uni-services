package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type emtSession struct {
	tim   time.Time
	exp   time.Time
	token string
}

type EMT struct {
	email    string
	password string
	session  *emtSession
}

func New(email, password string) *EMT {
	e := &EMT{
		email:    email,
		password: password,
	}

	go e.maintainSessionOpen()

	return e
}

func (e *EMT) maintainSessionOpen() {
	for {
		err := e.RefreshSession()

		if err != nil {
			log.Println("Error when refreshing session: ", err)
			time.Sleep(15 * time.Second)
		} else {
			time.Sleep(10 * time.Minute)
		}
	}
}

func (e *EMT) RefreshSession() error {
	c := &http.Client{}
	r, err := http.NewRequest("GET", "https://openapi.emtmadrid.es/v1/mobilitylabs/user/login/", http.NoBody)

	if err != nil {
		return err
	}

	r.Header.Set("email", e.email)
	r.Header.Set("password", e.password)

	resp, err := c.Do(r)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Unsuccessful status code")
	}

	var d struct {
		Data []struct {
			Exp   int    `json: "tokenSecExpiration"`
			Token string `json:"accessToken"`
		} `json: "data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&d)

	if err != nil {
		return err
	}

	s := &emtSession{
		tim:   time.Now(),
		exp:   time.Now().Add(time.Duration(d.Data[0].Exp) * time.Second),
		token: d.Data[0].Token,
	}

	e.session = s

	return nil
}

func (e *EMT) GetStopEstimates(stop int) ([]Bus, error) {
	if e.session == nil {
		return nil, errors.New("No session open")
	}

	api := fmt.Sprintf("https://openapi.emtmadrid.es/v1/transport/busemtmad/stops/%v/arrives/", stop)
	var data struct {
		CultureInfo         string `json: "cultureInfo`
		StopRequired        string `json:"Text_StopRequired_YN"`
		EstimationsRequired string `json:"Text_EstimationsRequired_YN"`
		IncidencesRequired  string `json:"Text_IncidencesRequired_YN"`
	}

	data.CultureInfo = "ES"
	data.StopRequired = "Y"
	data.EstimationsRequired = "Y"
	data.IncidencesRequired = "N"

	b := &bytes.Buffer{}

	_ = json.NewEncoder(b).Encode(data)

	c := &http.Client{}
	r, err := http.NewRequest("POST", api, b)

	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("accessToken", e.session.token)

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
	errcode, err := strconv.Atoi(m["code"].(string))

	if errcode != 0 || err != nil {
		return nil, fmt.Errorf("server responded unsuccessfuly. errorCode: %v, err: %v", errcode, err)
	}

	var arrives []Bus

	dataObj := m["data"].([]interface{})
	arrData := dataObj[0].(map[string]interface{})
	arrMap := arrData["Arrive"].([]interface{})

	for _, a := range arrMap {
		arrive := a.(map[string]interface{})
		arrives = append(arrives, Bus{
			LineID:      arrive["line"].(string),
			Destination: arrive["destination"].(string),
			BusID:       arrive["bus"].(string),
			TimeLeft:    time.Second * time.Duration((math.Min(25*60, arrive["estimateArrive"].(float64)))),
			Distance:    int(arrive["DistanceBus"].(float64)),
			Latitude:    arrive["geometry"].(map[string]interface{})["coordinates"].([]interface{})[0].(float64),
			Longitude:   arrive["geometry"].(map[string]interface{})["coordinates"].([]interface{})[1].(float64),
		})
	}

	return arrives, nil
}
