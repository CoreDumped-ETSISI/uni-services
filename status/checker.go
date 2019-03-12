package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func checkEndpoint(URL string) (bool, int) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(URL)

	if err != nil {
		return false, 0
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 == 5 || resp.StatusCode/100 == 4 {
		return false, resp.StatusCode
	}

	return true, resp.StatusCode
}

func (s *server) checkService(service *serviceStatus) {
	ok, status := checkEndpoint(service.URL)

	lastStatus := service.CircuitBreaker
	newStatus := lastStatus

	if ok {
		if service.Up {
			newStatus = 0
		} else {
			newStatus--
		}
	} else {
		if service.Up {
			newStatus++
		} else {
			newStatus = s.circuitBreakLimit
		}
	}

	if newStatus < 0 {
		newStatus = 0
	} else if newStatus > s.circuitBreakLimit {
		newStatus = s.circuitBreakLimit
	}

	var err error = nil

	if !service.Up && newStatus == 0 {
		// Service went up
		service.Up = true
		channel := "SERVICE_STATUS_CHANNEL"
		if service.Infra {
			channel = "SERVICE_STATUS_CHANNEL_INTERNAL"
		}

		err = s.publishToRedis(service, channel)
	} else if service.Up && newStatus == s.circuitBreakLimit {
		// Service went down
		service.Up = false
		channel := "SERVICE_STATUS_CHANNEL"
		if service.Infra {
			channel = "SERVICE_STATUS_CHANNEL_INTERNAL"
		}

		err = s.publishToRedis(service, channel)
	}

	if err != nil {
		s.log.Print("Error when publishing to redis", err)
	}

	service.LastCheck = time.Now()
	service.LastStatusCode = status

	// Save data
	service.CircuitBreaker = newStatus
}

func (s *server) checkAllServices() {
	for i := range s.status {
		s.checkService(s.status[i])
	}
}

func (s *server) launchPrecheck() {
	s.log.Print("Starting launch precheck...")
	for i := range s.status {
		ok, status := checkEndpoint(s.status[i].URL)

		s.status[i].LastCheck = time.Now()
		s.status[i].LastStatusCode = status

		// Save data
		if ok {
			s.status[i].CircuitBreaker = 0
			s.status[i].Up = true
		} else {
			s.status[i].CircuitBreaker = s.circuitBreakLimit
			s.status[i].Up = false
		}
	}

	s.log.Print("Precheck done.")
}

func (s *server) publishToRedis(service *serviceStatus, channel string) error {
	message := ""

	if service.Up {
		// Service went up
		message = "<b>" + service.Name + "</b>" + " ha vuelto! 🎉"
	} else {
		// Service went down
		message = "<b>" + service.Name + "</b>" + " acaba de morir 💀"
	}

	pubmsg := redisMessage{
		Text: message,
		Link: &service.URL,
	}

	data, err := json.Marshal(pubmsg)

	if err != nil {
		return err
	}

	err = s.redis.Publish(channel, string(data)).Err()

	if err == nil {
		s.log.Print("Published new message:", message)
	}

	return err
}