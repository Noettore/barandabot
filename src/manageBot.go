package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func setBotToken(newToken string) error {
	var err error
	if redisClient == nil {
		return ErrNilPointer
	}
	if newToken == "" && cmdFlags.interactive {
		fmt.Println("Add the new token:")
		reader := bufio.NewReader(os.Stdin)
		newToken, err = reader.ReadString('\n')
		if err != nil {
			log.Printf("Error in reading new bot token: %v", err)
			return ErrStdRead
		}
	}
	token := strings.TrimSpace(newToken)
	matched, err := regexp.MatchString("^\\d+\\:([0-9]|[A-z]|\\_|\\-)+", token)
	if err != nil {
		log.Printf("Error in parsing bot token: %v", err)
		return ErrTokenParsing
	}
	if !matched {
		return ErrTokenInvalid
	}

	err = redisClient.Set(botToken, token, 0).Err()
	if err != nil {
		log.Printf("Error in adding new bot token: %v", err)
		return ErrRedisAddSet
	}

	return nil
}

func getBotToken() (string, error) {
	if redisClient == nil {
		return "", ErrNilPointer
	}
	tokenExists, err := redisClient.Exists(botToken).Result()
	if err != nil {
		log.Printf("Error checking if token exists in db: %v", err)
		return "", ErrRedisCheckSet
	}
	if tokenExists == 0 {
		fmt.Println("No bot token found.")
		err := setBotToken("")
		if err != nil {
			log.Printf("Couldn't add new bot tokens: %v", err)
			return "", ErrAddToken
		}
	}
	token, err := redisClient.Get(botToken).Result()
	if err != nil {
		log.Printf("Couldn't retrieve bot token: %v", err)
		return "", ErrRedisRetrieveSet
	}
	return token, nil
}

func addBotInfo(botToken string, botUser string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.HSet(botInfo, botToken, botUser).Err()
	if err != nil {
		log.Printf("Error in adding bot info: %v", err)
		return ErrRedisAddHash
	}

	return nil
}

func removeBotInfo(botToken string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.HDel(botInfo, botToken).Err()
	if err != nil {
		log.Printf("Error in removing bot info: %v", err)
		return ErrRedisDelHash
	}
	return nil
}
