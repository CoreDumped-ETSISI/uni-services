package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var linkre = regexp.MustCompile(`\/sites\/default\/files\/(.*?)menu_campus(.*?).pdf`)
var ErrLinkNotFound = errors.New("No se ha podido encontrar el link")

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

	url := linkre.FindString(string(b))

	if url == "" {
		return "", ErrLinkNotFound
	}

	return fmt.Sprintf("https://www.etsisi.upm.es%s", url), nil
}
