package main

import (
	"log"
)

func main() {

	err := getFlags()
	if err != nil {
		log.Fatalln("Error in parsing command line flags. Abort!")
	}

	err = redisInit(cmdFlags.redisAddr, cmdFlags.redisPwd, cmdFlags.redisDB)
	defer redisClient.Close()
	if err != nil {
		log.Fatalf("Error in initializing redis instance: %v", err)
	}

	if cmdFlags.interactive {
		mainMenu()
	} else if cmdFlags.tokens != nil {
		err = addBotTokens(redisClient, cmdFlags.tokens)
		if err == ErrAddToken {
			log.Printf("Error in adding bot tokens: %v", err)
		}
	}
}
