package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	adminInlineMenu   [][]tb.InlineButton
	authInlineMenu    [][]tb.InlineButton
	genericInlineMenu [][]tb.InlineButton
	startMenu         [][]tb.InlineButton
	goBackMenu        [][]tb.InlineButton
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
	backBtn = tb.InlineButton{
		Unique: "back_btn",
		Text:   "\xF0\x9F\x94\x99 Torna al menù principale",
	}
	infoBtn = tb.InlineButton{
		Unique: "info_btn",
		Text:   "\xE2\x84\xB9 Info",
	}
	userBtn = tb.InlineButton{
		Unique: "user_btn",
		Text:   "\xF0\x9F\x91\xA4 My info",
	}
)

func setBotMenus() error {

	genericInlineMenu = append(genericInlineMenu, []tb.InlineButton{userBtn, infoBtn}, []tb.InlineButton{stopBtn})
	authInlineMenu = genericInlineMenu
	adminInlineMenu = genericInlineMenu
	//adminInlineMenu = append(adminInlineMenu, []tb.InlineButton{stopBtn, infoBtn})
	//authInlineMenu = append(authInlineMenu, []tb.InlineButton{stopBtn, infoBtn})

	startMenu = append(startMenu, []tb.InlineButton{startBtn})
	goBackMenu = append(goBackMenu, []tb.InlineButton{backBtn})

	return nil
}

func setBotCallbacks() error {
	if bot == nil {
		return ErrNilPointer
	}

	bot.Handle(&startBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		startCmd(c.Sender)
	})

	bot.Handle(&stopBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		stopCmd(c.Sender)
	})

	bot.Handle(&userBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		userInfoCmd(c.Sender)
	})
	bot.Handle(&infoBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithSpecificMenu(c.Sender, contactMsg, goBackMenu)
	})
	bot.Handle(&backBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithMenu(c.Sender, menuMsg)
	})

	return nil
}
