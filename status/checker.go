package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func checkEndpoint(URL string) (bool, int, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(URL)

	if err != nil {
		return false, 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 == 5 || resp.StatusCode/100 == 4 {
		return false, resp.StatusCode, nil
	}

	return true, resp.StatusCode, nil
}

func (s *server) checkService(service *serviceStatus) {
	start := time.Now()
	ok, status, lastErr := checkEndpoint(service.URL)
	delta := time.Now().Sub(start)

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
	service.LastError = lastErr
	service.LastLatency = delta.Seconds()

	// Save data
	service.CircuitBreaker = newStatus
	if err = s.saveServiceStatus(service); err != nil {
		s.log.Print("Error when saving: ", err)
	}
}

func (s *server) saveServiceStatus(service *serviceStatus) error {
	var errText string

	if service.LastError != nil {
		errText = service.LastError.Error()
	}

	err := s.postgres.Insert(&serviceHistory{
		Timestamp:  time.Now(),
		Up:         service.Up,
		URL:        service.URL,
		StatusCode: service.LastStatusCode,
		Error:      errText,
		Latency:    service.LastLatency,
	})

	return err
}

func (s *server) checkAllServices() {
	ok, _, _ := checkEndpoint("https://google.com")

	if !ok {
		s.log.Error("Could not check services because there is no internet access")
		return
	}

	for i := range s.status {
		s.checkService(s.status[i])
	}

	err := s.invalidateCache()

	if err != nil {
		s.log.Error(err)
	}
}

func (s *server) launchPrecheck() {
	s.log.Print("Starting launch precheck...")
	for i := range s.status {
		ok, status, lastErr := checkEndpoint(s.status[i].URL)

		s.status[i].LastCheck = time.Now()
		s.status[i].LastStatusCode = status
		s.status[i].LastError = lastErr

		// Save data
		if ok {
			s.status[i].CircuitBreaker = 0
			s.status[i].Up = true
		} else {
			s.status[i].CircuitBreaker = s.circuitBreakLimit
			s.status[i].Up = false
		}
	}

	err := s.invalidateCache()

	if err != nil {
		panic(err)
	}

	s.log.Print("Precheck done.")
}

func (s *server) invalidateCache() error {
	var his []serviceHistory

	t := time.Now().Add(-2160 * time.Hour)

	err := s.postgres.Model(&his).
		Where("timestamp > ?", t).
		Order("timestamp ASC").
		Select(&his)

	if err != nil {
		return err
	}

	var thinhis []*serviceHistory
	prevs := map[string]*serviceHistory{}

	for i := range his {
		p := prevs[his[i].URL]
		if p == nil {
			thinhis = append(thinhis, &his[i])
		} else if !his[i].Up && p.Up || his[i].Up && !p.Up {
			thinhis = append(thinhis, p, &his[i])
		} else if p.Timestamp.Day() != his[i].Timestamp.Day() {
			thinhis = append(thinhis, p, &his[i])
		}

		prevs[his[i].URL] = &his[i]
	}

	// Add last points
	for key := range prevs {
		thinhis = append(thinhis, prevs[key])
	}

	s.historyCache = thinhis

	return nil
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
