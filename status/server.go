package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/go-pg/pg/orm"

	"github.com/labstack/echo"

	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
)

type server struct {
	redis             *redis.Client
	postgres          *pg.DB
	status            []*serviceStatus
	circuitBreakLimit int
	cron              string
	log               echo.Logger
	historyCache      []*serviceHistory
}

func New() *server {
	path := "/endpoints.json"

	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	var endpoints []*serviceStatus

	err = json.NewDecoder(f).Decode(&endpoints)

	if err != nil {
		panic(err)
	}

	s := &server{}

	s.status = endpoints
	s.circuitBreakLimit, _ = strconv.Atoi(os.Getenv("CIRCUIT_BREAK_LIMIT"))
	s.cron = os.Getenv("CHECK_INTERVAL")

	redisb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	s.redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       redisb,
	})

	if _, ok := os.LookupEnv("DB_HOST"); ok {
		s.postgres = pg.Connect(&pg.Options{
			Addr:     os.Getenv("DB_HOST"),
			Database: os.Getenv("DB_DB"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
		})

		s.postgres.CreateTable(&serviceHistory{}, &orm.CreateTableOptions{
			IfNotExists: true,
		})
	}

	return s
}
