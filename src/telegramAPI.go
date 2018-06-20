package main

import (
	"errors"
	"log"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var bots []*tb.Bot

var (
	//ErrNilPointer is thrown when a pointer is nil
	ErrNilPointer = errors.New("pointer is nil")
)

func botsInit() error {
	tokens, err := getBotTokens()
	if err != nil {
		log.Printf("Error in retriving bot tokens: %v. Cannot start telebot without tokens.", err)
		return err
	}

	if tokens == nil {
		log.Println("Error: pointer is nil")
		return ErrNilPointer
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
			err = addBotInfo(bot, token)
			if err != nil {
				log.Printf("Error: bot %s info couldn't be added: %v", bot.Me.Username, err)
			}
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
		if bots[i] != nil {
			go botStart(bots[i])
		}
	}

	return nil
}

func botStart(bot *tb.Bot) error {
	if bot == nil {
		return ErrNilPointer
	}
	log.Printf("Started bot %s", bot.Me.Username)
	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "hello world")
	})

	bot.Start()

	return nil
}
