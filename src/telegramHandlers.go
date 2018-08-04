package main

import (
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func startHandler(m *tb.Message) {
	var msg string

	isUser, err := isUser(m.Sender.ID)
	if err != nil {
		log.Printf("Error checking if ID is bot user: %v", err)
	}

	started, err := isStartedUser(m.Sender.ID)
	if err != nil {
		log.Printf("Error checking if user is started: %v", err)
	}
	if !started {
		err = startUser(m.Sender.ID, true)
		if err != nil {
			log.Printf("Error starting user: %v", err)
		}
		if isUser {
			msg = restartMsg
		} else {
			msg = startMsg
		}
	} else {
		msg = alreadyStartedMsg
	}

	err = sendMsgWithMenu(m.Sender, msg)
	if err != nil {
		log.Printf("Error sending message to started user: %v", err)
	}
}

func stopHandler(m *tb.Message) {
	admin, err := isBotAdmin(m.Sender.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
	}
	if admin {
		msg := "Non ci siamo... Io l'ho nominata AMMINISTRATORE, cosa crede?! Questo ruolo esige impegno! Non può certo bloccarmi!"
		err := sendMsg(m.Sender, msg)
		if err != nil {
			log.Printf("Error sending message to unstoppable user: %v", err)
		}
	} else {
		err = startUser(m.Sender.ID, false)
		if err != nil {
			log.Printf("Error starting user: %v", err)
		}
		err := sendMsgWithSpecificMenu(m.Sender, stopMsg, startMenu)
		if err != nil {
			log.Printf("Error sending message to stopped user: %v", err)
		}
	}
}

func setBotHandlers() error {
	if bot == nil {
		return ErrNilPointer
	}
	bot.Handle("/start", startHandler)
	bot.Handle("/stop", stopHandler)
	bot.Handle("/menu", func(m *tb.Message) {
		bot.Send(m.Sender, "hello world")
	})
	bot.Handle("/userID", func(m *tb.Message) {
		bot.Send(m.Sender, strconv.Itoa(m.Sender.ID))
	})

	return nil
}
