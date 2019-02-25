package main

import "time"

type serviceStatus struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	Up             bool      `json:"up"`
	LastStatusCode int       `json:"lastStatusCode"`
	LastCheck      time.Time `json:"lastCheck"`
	Infra          bool      `json:"infra"`
	CircuitBreaker int       `json:"-"`
}

type redisMessage struct {
	Text string  `json:"text"`
	Link *string `json:"link"`
}
