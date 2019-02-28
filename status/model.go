package main

import "time"

type serviceStatus struct {
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	Up             bool      `json:"up"`
	LastStatusCode int       `json:"lastStatusCode"`
	LastCheck      time.Time `json:"lastCheck"`
	LastError      error     `json:"lastError"`
	Infra          bool      `json:"infra"`
	CircuitBreaker int       `json:"-"`
}

type redisMessage struct {
	Text string  `json:"text"`
	Link *string `json:"link"`
}
