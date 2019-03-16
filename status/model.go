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
	LastLatency    float64   `json:"latency"`
	CircuitBreaker int       `json:"-"`
}

type redisMessage struct {
	Text string  `json:"text"`
	Link *string `json:"link"`
}

type serviceHistory struct {
	Timestamp  time.Time `json:"timestamp" bson:"_id"`
	URL        string    `json:"url" bson:"url"`
	Up         bool      `json:"up" bson:"up"`
	StatusCode int       `json:"statusCode" bson:"statusCode"`
	Error      string    `json:"error,omitempty" bson:"error"`
	Latency    float64   `json:"latency"`
}
