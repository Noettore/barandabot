package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	tb "gopkg.in/tucnak/telebot.v2"
)

func botInit(redisClient *redis.Client) ([]*tb.Bot, error) {
	var bots []*tb.Bot

	tokens, err := getBotTokens(redisClient)
	if err != nil {
		log.Printf("Error in retriving bot tokens: %v. Cannot start telebot without tokens.", err)
		return nil, err
	}

	for _, token := range tokens {
		bot, err := tb.NewBot(tb.Settings{
			Token:  token,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Printf("Error in enstablishing connection for bot %s: %v", bot.Me.Username, err)
		} else {
			bots = append(bots, bot)
		}
	}
	return bots, nil
}
