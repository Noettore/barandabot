package main

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var bots []*tb.Bot

var (
	//ErrNilPointer is thrown when a pointer is nil
	ErrNilPointer = errors.New("pointer is nil")
	//ErrIDFromMsg is thrown when the message doesn't contain user infos
	ErrIDFromMsg = errors.New("telegram: couldn't retrive user ID from message")
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
		poller := &tb.LongPoller{Timeout: 15 * time.Second}
		middlePoller := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
			if upd.Message == nil {
				return true
			}
			if upd.Message.Sender != nil {
				err := addUser(upd.Message.Sender)
				if err != nil {
					log.Printf("Error in adding user info: %v", err)
				}
				err = authorizeUser(upd.Message.Sender.ID, true)
				if err != nil {
					log.Printf("Error in authorizing user: %v", err)
				}
			} else {
				log.Printf("%v", ErrIDFromMsg)
			}
			auth, err := isAuthrizedUser(upd.Message.Sender.ID)
			if err != nil {
				log.Printf("Error checking if user is authorized: %v", err)
			}
			if !auth {
				return false
			}

			return true
		})

		bot, err := tb.NewBot(tb.Settings{
			Token:  token,
			Poller: middlePoller,
		})

		if err != nil {
			log.Printf("Error in enstablishing connection for bot %s: %v", bot.Me.Username, err)
		} else {
			bots = append(bots, bot)
			err = addBotInfo(token, bot)
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
		if bots[i] != nil {
			wg.Add(1)
			go botStart(bots[i], &wg)
		}
	}
	wg.Wait()

	return nil
}

func botStart(bot *tb.Bot, wg *sync.WaitGroup) error {
	defer wg.Done()
	if bot == nil {
		return ErrNilPointer
	}
	log.Printf("Started %s", bot.Me.Username)
	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "hello world")
	})
	bot.Handle("/userID", func(m *tb.Message) {
		bot.Send(m.Sender, strconv.Itoa(m.Sender.ID))
	})

	bot.Start()

	return nil
}
