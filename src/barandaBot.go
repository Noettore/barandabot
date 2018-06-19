package main

import (
	"log"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
)

func main() {

	cmdFlags, err := getFlags()
	if err != nil {
		log.Fatalln("Error in parsing command line flags. Abort!")
	}

	redisClient, err = redisInit(cmdFlags.redisAddr, cmdFlags.redisPwd, cmdFlags.redisDB)
	defer redisClient.Close()
	if err != nil {
		log.Fatalf("Error in initializing redis instance: %v", err)
	}

	mainMenu()
}
