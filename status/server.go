package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/globalsign/mgo"

	"github.com/labstack/echo"

	"github.com/go-redis/redis"
)

type server struct {
	redis             *redis.Client
	mongo             *mgo.Session
	status            []*serviceStatus
	circuitBreakLimit int
	log               echo.Logger
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

	redisb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	s.redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       redisb,
	})

	info := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("MONGO_HOST")},
		Database: os.Getenv("MONGO_DB"),
		Username: os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASS"),
		Timeout:  10 * time.Second,
	}

	database, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(err)
	}

	database.DB("etsisi-telegram-bot").C("status_history").EnsureIndexKey("url")

	s.mongo = database

	return s
}
