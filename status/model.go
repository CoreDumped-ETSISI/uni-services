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
	Timestamp  time.Time `json:"timestamp" bson:"_id" sql:"timestamp,pk"`
	URL        string    `json:"url" bson:"url" sql:"url,pk"`
	Up         bool      `json:"up" bson:"up" sql:"up"`
	StatusCode int       `json:"statusCode" bson:"statusCode" sql:"status"`
	Error      string    `json:"error,omitempty" bson:"error" sql:"error"`
	Latency    float64   `json:"latency" sql:"latency"`
}
