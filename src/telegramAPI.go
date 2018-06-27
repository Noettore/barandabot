package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	newAdminMsg string = "Sei stato aggiunto come amministratore. Adesso hai a disposizione una serie aggiuntiva di comandi e controlli per il bot."
	delAdminMsg string = "Sei stato rimosso da amministratore."
)

var bot *tb.Bot
var isStartedBot bool

var (
	//ErrNilPointer is thrown when a pointer is nil
	ErrNilPointer = errors.New("pointer is nil")
	//ErrIDFromMsg is thrown when the message doesn't contain user infos
	ErrIDFromMsg = errors.New("telegram: couldn't retrieve user ID from message")
	//ErrSendMsg is thrown when the message couldn't be send
	ErrSendMsg = errors.New("telegram: cannot send message")
	//ErrChatRetrieve is thrown when the chat cannot be retrieved
	ErrChatRetrieve = errors.New("telegram: cannot retrieve chat")
	//ErrTokenMissing is thrown when neither a token is in the db nor one is passed with -t on program start
	ErrTokenMissing = errors.New("telegram: cannot start bot without a token")
	//ErrBotInit is thrown when a bot couldn't be initialized
	ErrBotInit = errors.New("telegram: error in bot initialization")
	//ErrBotConn is thrown when there is a connection problem
	ErrBotConn = errors.New("telegram: cannot connect to bot")
)

func botInit() error {
	token, err := getBotToken()
	if err != nil {
		log.Printf("Error in retriving bot token: %v. Cannot start telebot without token.", err)
		return ErrTokenMissing
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
		return ErrBotConn
	}

	err = setBotHandlers()
	if err != nil {
		log.Printf("Error setting bot handlers: %v", err)
		return ErrBotInit
	}

	err = addBotInfo(token, bot.Me.Username)
	if err != nil {
		log.Printf("Error: bot %s info couldn't be added: %v", bot.Me.Username, err)
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

func setBotHandlers() error {
	if bot == nil {
		return ErrNilPointer
	}

	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "hello world")
	})
	bot.Handle("/userID", func(m *tb.Message) {
		bot.Send(m.Sender, strconv.Itoa(m.Sender.ID))
	})

	return nil
}

func botStart() error {
	if bot == nil {
		return ErrNilPointer
	}

	go bot.Start()
	isStartedBot = true
	log.Printf("Started %s", bot.Me.Username)

	return nil
}

func botStop() error {
	if bot == nil {
		return ErrNilPointer
	}
	log.Printf("Stopping %s", bot.Me.Username)
	bot.Stop()
	isStartedBot = false
	log.Println("Bot stopped")

	return nil
}
