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
			Exp   int    `json:"tokenSecExpiration"`
			Token string `json:"accessToken"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&d)

	if err != nil {
		return err
	}

	if d.Data == nil || len(d.Data) == 0 {
		return errors.New("Unsuccessful status code")
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
		CultureInfo         string `json:"cultureInfo"`
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
	r.Header.Set("Connection", "Keep-Alive")
	r.Header.Set("Accept-Language", "en-US")
	r.Header.Set("User-Agent", "Mozilla/5.0")

	t := time.Now()

	resp, err := c.Do(r)

	el := time.Now().Sub(t)

	log.Println("Time spent waiting for response:", el)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server responded unsuccessfuly")
	}

	var m apiResponse

	err = json.NewDecoder(resp.Body).Decode(&m)

	if err != nil {
		return nil, err
	}

	// First, check error code
	errcode, err := strconv.Atoi(m.Code)

	if errcode != 0 || err != nil {
		return nil, fmt.Errorf("server responded unsuccessfuly. errorCode: %v, err: %v", errcode, err)
	}

	var arrives []Bus

	if m.Data == nil || len(m.Data) == 0 {
		return arrives, nil
	}

	for _, a := range m.Data[0].Arrive {
		arrives = append(arrives, Bus{
			LineID:      a.Line,
			Destination: a.Destination,
			BusID:       strconv.Itoa(a.BusID),
			TimeLeft:    time.Second * time.Duration((math.Min(25*60, a.TimeLeft))),
			Distance:    int(a.DistanceBus),
			Latitude:    a.Geometry.Coordinates[0],
			Longitude:   a.Geometry.Coordinates[1],
		})
	}

	return arrives, nil
}
