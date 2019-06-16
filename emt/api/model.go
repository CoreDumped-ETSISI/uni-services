package api

import "time"

type Bus struct {
	LineID      string        `json:"lineId"`
	Destination string        `json:"destination"`
	BusID       string        `json:"busId"`
	TimeLeft    time.Duration `json:"timeLeft"`
	Distance    int           `json:"distance"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
}

type apiResponse struct {
	Code string
	Data []struct {
		Arrive []struct {
			DistanceBus int
			BusID       string `json:"bus"`
			Destination string
			Line        string
			TimeLeft    float64 `json:"estimateArrive"`
			Geometry    struct {
				Coordinates []float64
			}
		}
	}
}
