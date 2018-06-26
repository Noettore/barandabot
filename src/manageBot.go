package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
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
	token, err := redisClient.Get(botToken).Result()
	if err != nil {
		log.Printf("Couldn't retrieve bot token: %v", err)
		return "", ErrRedisRetrieveSet
	}
	if token == "" {
		fmt.Println("No bot token found.")
		err := setBotToken("")
		if err != nil {
			log.Printf("Couldn't add new bot tokens: %v", err)
			return "", ErrAddToken
		}
	}

	return token, nil
}

func addBotInfo(botToken string, bot *tb.Bot) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	jsonBot, err := json.Marshal(&bot)
	if err != nil {
		log.Printf("Error marshalling bot info: %v", err)
		return ErrJSONMarshall
	}
	err = redisClient.HSet(botHash, botToken, string(jsonBot)).Err()
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
	err := redisClient.HDel(botHash, botToken).Err()
	if err != nil {
		log.Printf("Error in removing bot info: %v", err)
		return ErrRedisDelHash
	}
	return nil
}
