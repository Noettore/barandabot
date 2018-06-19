package main

import (
	"log"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	bots []*tb.Bot
)

func botsInit() error {
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

func botsStart() error {
	err := botsInit()
	if err != nil {
		log.Fatalf("Error in initializing bots: %v", err)
	}

	for _, bot := range bots {
		defer bot.Stop()
	}

	var wg sync.WaitGroup
	for i := range bots {
		defer wg.Done()
		botStart(bots[i])
	}

	return nil
}

func botStart(bot *tb.Bot) error {
	log.Printf("Started bot %s", bot.Me.Username)
	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "hello world")
	})

	bot.Start()

	return nil
}
