package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func setMediaDir(newPath string) error {
	var err error
	if redisClient == nil {
		return ErrNilPointer
	}
	if newPath == "" && cmdFlags.interactive {
		fmt.Println("Add the new media path:")
		reader := bufio.NewReader(os.Stdin)
		newPath, err = reader.ReadString('\n')
		if err != nil {
			log.Printf("Error in reading new media path: %v", err)
			return ErrStdRead
		}
	}
	path := strings.TrimSpace(newPath)
	valid, err := isValidPath(path)
	if err != nil {
		log.Printf("Error in validating path: %v", err)
	}
	if !valid {
		return ErrInvalidPath
	}

	err = redisClient.Set(mediaPath, path, 0).Err()
	if err != nil {
		log.Printf("Error in adding new media path: %v", err)
		return ErrRedisAddSet
	}

	return nil
}

func getMediaDir() (string, error) {
	if redisClient == nil {
		return "", ErrNilPointer
	}
	mediaDirExists, err := redisClient.Exists(mediaPath).Result()
	if err != nil {
		log.Printf("Error checking if media path exists in db: %v", err)
		return "", ErrRedisCheckSet
	}
	if mediaDirExists == 0 {
		fmt.Println("No media path found.")
		err := setMediaDir("")
		if err != nil {
			log.Printf("Couldn't set new media path: %v", err)
			return "", ErrRedisAddSet
		}
	}
	path, err := redisClient.Get(mediaPath).Result()
	if err != nil {
		log.Printf("Couldn't retrieve mediaPath: %v", err)
		return "", ErrRedisRetrieveSet
	}
	return path, nil
}

func sendImg(user *tb.User, img *tb.Photo) error {
	_, err := bot.Send(user, img)
	if err != nil {
		log.Printf("Error sending img to user: %v", err)
		return ErrSendMsg
	}
	return nil
}
