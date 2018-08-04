package main

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	adminInlineMenu   [][]tb.InlineButton
	authInlineMenu    [][]tb.InlineButton
	genericInlineMenu [][]tb.InlineButton
	startMenu         [][]tb.InlineButton
)

var (
	startBtn = tb.InlineButton{
		Unique: "start_btn",
		Text:   "\xE2\x96\xB6 Avvia il barandaBot",
	}
	stopBtn = tb.InlineButton{
		Unique: "stop_btn",
		Text:   "\xF0\x9F\x9A\xAB Ferma il barandaBot",
	}
	infoBtn = tb.InlineButton{
		Unique: "info_btn",
		Text:   "\xE2\x84\xB9 Info",
	}
)

func setBotMenus() error {
	adminInlineMenu = append(adminInlineMenu, []tb.InlineButton{stopBtn, infoBtn})

	authInlineMenu = append(authInlineMenu, []tb.InlineButton{stopBtn, infoBtn})

	genericInlineMenu = append(genericInlineMenu, []tb.InlineButton{stopBtn, infoBtn})

	startMenu = append(startMenu, []tb.InlineButton{startBtn})

	return nil
}

func setBotCallbacks() error {
	if bot == nil {
		return ErrNilPointer
	}

	bot.Handle(&startBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		//TODO: save last message id per user so it's possible to hide inline keyboard
		//bot.Edit(lastMsgID, tb.ReplyMarkup{})
		err := startUser(c.Sender.ID, true)
		if err != nil {
			log.Printf("Error starting user: %v", err)
		}
		err = sendMsgWithMenu(c.Sender, restartMsg)
		if err != nil {
			log.Printf("Error sending message to started user: %v", err)
		}
	})

	bot.Handle(&stopBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		admin, err := isBotAdmin(c.Sender.ID)
		if err != nil {
			log.Printf("Error checking if user is admin: %v", err)
		}
		if admin {
			msg := "Non ci siamo... Io l'ho nominata AMMINISTRATORE, cosa crede?! Questo ruolo esige impegno! Non può certo bloccarmi!"
			err := sendMsg(c.Sender, msg)
			if err != nil {
				log.Printf("Error sending message to unstoppable user: %v", err)
			}
		} else {
			err = startUser(c.Sender.ID, false)
			if err != nil {
				log.Printf("Error starting user: %v", err)
			}
			err := sendMsgWithSpecificMenu(c.Sender, stopMsg, startMenu)
			if err != nil {
				log.Printf("Error sending message to stopped user: %v", err)
			}
		}
	})

	return nil
}
