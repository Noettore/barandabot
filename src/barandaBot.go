package main

import (
	"log"
)

func main() {

	err := getFlags()
	if err != nil {
		log.Fatalf("Error in parsing command line flags: %v", err)
	}

	err = redisInit(cmdFlags.redisAddr, cmdFlags.redisPwd, cmdFlags.redisDB)
	defer redisClient.Close()
	if err != nil {
		log.Fatalf("Error in initializing redis instance: %v", err)
	}

	if cmdFlags.interactive {
		mainMenu()
	} else if cmdFlags.tokens != nil {
		err = addBotTokens(cmdFlags.tokens)
		if err == ErrAddToken {
			log.Printf("Error in adding bot tokens: %v", err)
		}
	}
}
