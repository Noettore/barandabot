package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func botInit() ([]*tb.Bot, []error) {
	var bots []*tb.Bot
	var errors []error
	for i, token := range tokens {
		var timeout int
		if i < len(timeouts) {
			timeout = timeouts[i]
		} else {
			timeout = 10
		}
		tmpBot, tmpErr := tb.NewBot(tb.Settings{
			Token:  token,
			Poller: &tb.LongPoller{Timeout: time.Duration(timeout) * time.Second},
		})

		bots = append(bots, tmpBot)
		errors = append(errors, tmpErr)
	}

	return bots, errors
}
