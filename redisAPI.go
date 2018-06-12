package main

import (
	"log"

	"github.com/go-redis/redis"
)

func redisInit(addr string, pwd string, db int) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	err := redisClient.Ping().Err()
	if err != nil {
		log.Panicf("Error in connecting to redis instance: %v", err)
	}
	return redisClient, nil
}
