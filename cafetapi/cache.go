package main

import (
	"time"

	"github.com/labstack/echo"
)

type menuCache struct {
	LastFetch   *CafeMenu `json:"lastFetch"`
	TimeFetched time.Time `json:"timeFetched"`
	FetchURL    string    `json:"fetchURL"`
}

func isCacheValid(from, to int, date time.Time) bool {
	now := date.Day()

	if from > to { // Fin de mes
		return now >= from || now <= to
	}

	return now >= from && now <= to
}

func (s *server) getCachedCafeMenu() *CafeMenu {
	if s.cache != nil && time.Now().Sub(s.cache.TimeFetched) < CacheTimeout*time.Hour {
		return s.cache.LastFetch
	}

	return nil
}

func (s *server) getLatestCafeMenu() (*CafeMenu, error) {
	if s.cache != nil && time.Now().Sub(s.cache.TimeFetched) < CacheTimeout*time.Hour {
		return s.cache.LastFetch, nil
	}

	pdfURL, err := getLatestPdfURL()

	if err != nil {
		return nil, err
	}

	if s.cache != nil && s.cache.FetchURL == pdfURL {
		s.cache.TimeFetched = time.Now()
		return s.cache.LastFetch, nil
	}

	menu, err := s.getLatestCafeTable(pdfURL)
	if err != nil {
		return nil, err
	}

	s.cache = &menuCache{
		LastFetch:   menu,
		FetchURL:    pdfURL,
		TimeFetched: time.Now(),
	}

	return menu, nil
}

func (s *server) updateCache(logger echo.Logger) func() {
	return func() {
		logger.Print("Updating cache")
		pdfURL, err := getLatestPdfURL()

		if err != nil {
			logger.Error(err)
			return
		}

		if s.cache != nil && s.cache.FetchURL == pdfURL {
			s.cache.TimeFetched = time.Now()
			return
		}

		menu, err := s.getLatestCafeTable(pdfURL)
		if err != nil {
			logger.Error(err)
			return
		}

		s.cache = &menuCache{
			LastFetch:   menu,
			FetchURL:    pdfURL,
			TimeFetched: time.Now(),
		}
	}
}
