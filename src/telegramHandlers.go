package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func setBotHandlers() error {
	if bot == nil {
		return ErrNilPointer
	}
	bot.Handle("/start", func(m *tb.Message) {
		startCmd(m.Sender, true)
	})
	bot.Handle("/stop", func(m *tb.Message) {
		stopCmd(m.Sender)
	})
	bot.Handle("/menu", func(m *tb.Message) {
		sendMsgWithMenu(m.Sender, menuMsg, true)
	})
	bot.Handle("/userInfo", func(m *tb.Message) {
		msg, _ := getUserDescription(m.Sender)
		sendMsgWithSpecificMenu(m.Sender, msg, myInfoMenu, true)
	})
	bot.Handle("/botInfo", func(m *tb.Message) {
		sendMsgWithSpecificMenu(m.Sender, contactMsg, botInfoMenu, true)
	})
	bot.Handle("/help", func(m *tb.Message) {
		sendMsgWithSpecificMenu(m.Sender, contactMsg, botInfoMenu, true)
	})
	bot.Handle("/config", func(m *tb.Message) {
		msg, _ := getUserDescription(m.Sender)
		sendMsgWithSpecificMenu(m.Sender, msg, myInfoMenu, true)
	})
	bot.Handle("/authUser", func(m *tb.Message) {
		authUserCmd(m.Sender, m.Payload, true)
	})
	bot.Handle("/deAuthUser", func(m *tb.Message) {
		deAuthUserCmd(m.Sender, m.Payload, true)
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		sendMsgWithMenu(m.Sender, wrongCmdMsg, true)
	})

	return nil
}
