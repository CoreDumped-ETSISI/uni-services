package main

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

type server struct {
	redis *redis.Client
}

func New() *server {
	s := &server{}

	redisb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	s.redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       redisb,
	})

	return s
}
