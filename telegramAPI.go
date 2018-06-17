package main

import (
	"time"

	"github.com/go-redis/redis"
	tb "gopkg.in/tucnak/telebot.v2"
)

func botInit(redisClient *redis.Client) ([]*tb.Bot, []error) {
	var bots []*tb.Bot
	var errors []error

	tokens, err := getBotTokens(redisClient)
	if err != nil {

	}

	for _, token := range tokens {
		tmpBot, tmpErr := tb.NewBot(tb.Settings{
			Token:  token,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		})

		bots = append(bots, tmpBot)
		errors = append(errors, tmpErr)
	}

	return bots, errors
}
