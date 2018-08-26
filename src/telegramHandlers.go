package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func setBotHandlers() error {
	if bot == nil {
		return ErrNilPointer
	}
	bot.Handle("/start", func(m *tb.Message) {
		startCmd(m.Sender)
	})
	bot.Handle("/stop", func(m *tb.Message) {
		stopCmd(m.Sender)
	})
	bot.Handle("/menu", func(m *tb.Message) {
		sendMsgWithMenu(m.Sender, menuMsg)
	})
	bot.Handle("/userInfo", func(m *tb.Message) {
		userInfoCmd(m.Sender)
	})
	bot.Handle("/botInfo", func(m *tb.Message) {
		sendMsgWithSpecificMenu(m.Sender, contactMsg, goBackMenu)
	})

	return nil
}
