package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
)

const (
	tkSet = "botTokens"
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

func addBotTokens(client *redis.Client) {
	fmt.Println("Add the new tokens, comma-separated:")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Panicf("Error in reading new bot tokens: %v", err)
	}
	newTokens := strings.Split(line, ",")
	for i, newToken := range newTokens {
		newToken = strings.TrimSpace(newToken)
		matched, err := regexp.MatchString("^\\d+\\:([0-9]|[A-z]|\\_|\\-)+", newToken)
		if err != nil {
			log.Printf("Error in parsing new bot token: %v", err)
			newTokens = append(newTokens[:i], newTokens[i+1:]...)
		}
		if !matched {
			log.Printf("%s is not a valid bot token and has not been added.", newToken)
			newTokens = append(newTokens[:i], newTokens[i+1:]...)
		}
	}
	err = client.SAdd(tkSet, newTokens).Err()
	if err != nil {

	}
}

func getBotTokens(client *redis.Client) ([]string, error) {
	tkNum, err := client.SCard(tkSet).Result()
	if err != nil {
		log.Panicf("Couldn't retrive number of bot tokens: %v", err)
	}
	if tkNum == 0 {
		fmt.Println("No bot token found.")
		addBotTokens(client)
	}

	tokens, err := client.SMembers(tkSet).Result()
	if err != nil {

	}

	return tokens, nil
}
