package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	bots []*tb.Bot
)

func botInit() error {
	tokens, err := getBotTokens(redisClient)
	if err != nil {
		log.Printf("Error in retriving bot tokens: %v. Cannot start telebot without tokens.", err)
		return err
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
	return nil
}

func botStart() {
	err := botInit()
	if err != nil {
		log.Fatalf("Error in initializing bots: %v", err)
	}

	for _, bot := range bots {
		defer bot.Stop()
	}

	/*b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()*/
}
