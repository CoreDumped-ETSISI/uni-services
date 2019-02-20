package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var linkre = regexp.MustCompile(`\/sites\/default\/files\/cafeteria\/menu_campus(.*?).pdf`)

func getLatestPdfURL() (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get("https://www.etsisi.upm.es/escuela/servicios/cafeteria")

	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://www.etsisi.upm.es%s", linkre.FindString(string(b))), nil
}
