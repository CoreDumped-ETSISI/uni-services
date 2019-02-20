package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type pdfTableSettings struct {
	Flavor     string   `json:"flavor"`
	TableAreas []string `json:"table_areas"`
	Columns    []string `json:"columns"`
}

type pdfTableRequest struct {
	PDF      string           `json:"pdf"`
	Settings pdfTableSettings `json:"settings"`
}

type pdfTableReport struct {
	Accuracy   float64 `json:"accuracy"`
	Order      int     `json:"order"`
	Page       int     `json:"page"`
	Whitespace float64 `json:"whitespace"`
}

type pdfTable [][]string

type pdfTableResponse struct {
	Report pdfTableReport `json:"report"`
	Data   pdfTable       `json:"data"`
}

func (s server) getLatestCafeTable(pdfaddr string) (*CafeMenu, error) {
	client := &http.Client{
		Timeout: time.Second * 120,
	}

	settings := pdfTableSettings{
		Flavor: "stream",
		TableAreas: []string{
			"144.9118179245283,509.17870605892165,511.4923984276729,490.6510858254586",
			"115.7971147798742,461.5362540300167,766.9077487421383,206.11977509727626",
		},
		Columns: []string{
			"144",
			"243.5047899371069,373.197558490566,508.1839094339622,639.2000735849056",
		},
	}

	obj := pdfTableRequest{
		PDF:      pdfaddr,
		Settings: settings,
	}

	buf, err := json.Marshal(obj)

	if err != nil {
		return nil, err
	}

	resp, err := client.Post(fmt.Sprintf("http://%s:8080/table", s.tableServer), "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var pdf []pdfTableResponse
	err = json.NewDecoder(resp.Body).Decode(&pdf)

	if err != nil {
		return nil, err
	}

	return parsePdfTable(pdf), nil
}
