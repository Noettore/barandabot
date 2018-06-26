package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *tb.Bot

var (
	//ErrNilPointer is thrown when a pointer is nil
	ErrNilPointer = errors.New("pointer is nil")
	//ErrIDFromMsg is thrown when the message doesn't contain user infos
	ErrIDFromMsg = errors.New("telegram: couldn't retrieve user ID from message")
	//ErrSendMsg is thrown when the message couldn't be send
	ErrSendMsg = errors.New("telegram: cannot send message")
	//ErrChatRetrive is thrown when the chat cannot be retrieved
	ErrChatRetrieve = errors.New("telegram: cannot retrieve chat")
)

func botInit() error {
	token, err := getBotToken()
	if err != nil {
		log.Printf("Error in retriving bot token: %v. Cannot start telebot without token.", err)
		return err
	}

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

	bot, err = tb.NewBot(tb.Settings{
		Token:  token,
		Poller: middlePoller,
	})

	if err != nil {
		log.Printf("Error in enstablishing connection for bot %s: %v", bot.Me.Username, err)
	} else {
		err = addBotInfo(token, bot)
		if err != nil {
			log.Printf("Error: bot %s info couldn't be added: %v", bot.Me.Username, err)
		}
	}
	return nil
}

func sendMessage(user *tb.User, msg string) error {
	_, err := bot.Send(user, msg)
	if err != nil {
		log.Printf("Error sending message to user: %v", err)
		return ErrSendMsg
	}
	return nil
}

func botStart() error {
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
