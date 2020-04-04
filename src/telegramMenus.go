package main

import (
	"log"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	superAdminInlineMenu [][]tb.InlineButton
	adminInlineMenu      [][]tb.InlineButton
	authInlineMenu       [][]tb.InlineButton
	genericInlineMenu    [][]tb.InlineButton
	startMenu            [][]tb.InlineButton
	myInfoMenu           [][]tb.InlineButton
	botInfoMenu          [][]tb.InlineButton
	authUserMenu         [][]tb.InlineButton
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
		Text:   "\xE2\x84\xB9 Bot info",
	}
	userBtn = tb.InlineButton{
		Unique: "user_btn",
		Text:   "\xF0\x9F\x91\xA4 My info",
	}
	authBtn = tb.InlineButton{
		Unique: "auth_btn",
		Text:   "\xE2\x9C\x85 Autorizza utente",
	}
	deAuthBtn = tb.InlineButton{
		Unique: "de_auth_btn",
		Text:   "\xE2\x9D\x8C Deautorizza utente",
	}
	adminBtn = tb.InlineButton{
		Unique: "admin_btn",
		Text:   "\xF0\x9F\x91\x91 Nomina amministratore",
	}
	deAdminBtn = tb.InlineButton{
		Unique: "de_admin_btn",
		Text:   "\xF0\x9F\x92\x80 Rimuovi amministratore",
	}
	sendMsgBtn = tb.InlineButton{
		Unique: "send_msg_btn",
		Text:   "\xF0\x9F\x93\xA3 Invia messaggio alla sezione",
	}
	authUGSopranoBtn = tb.InlineButton{
		Unique: "auth_ugSoprano_btn",
		Text:   "\xF0\x9F\x91\xA7 Soprani",
	}
	authUGContraltoBtn = tb.InlineButton{
		Unique: "auth_ugContralto_btn",
		Text:   "\xF0\x9F\x91\xA9 Contralti",
	}
	authUGTenoreBtn = tb.InlineButton{
		Unique: "auth_ugTenore_btn",
		Text:   "\xF0\x9F\x91\xA6 Tenori",
	}
	authUGBassoBtn = tb.InlineButton{
		Unique: "auth_ugBasso_btn",
		Text:   "\xF0\x9F\x91\xA8 Bassi",
	}
	authUGCommissarioBtn = tb.InlineButton{
		Unique: "auth_ugCommissario_btn",
		Text:   "\xF0\x9F\x93\x9D Commissari",
	}
	authUGReferenteBtn = tb.InlineButton{
		Unique: "auth_ugReferente_btn",
		Text:   "\xF0\x9F\x93\x8B Referenti",
	}
	authUGPreparatoreBtn = tb.InlineButton{
		Unique: "auth_ugPreparatori_btn",
		Text:   "\xF0\x9F\x8E\xB9 Preparatori",
	}
)

func setBotMenus() error {

	genericInlineMenu = append(genericInlineMenu, []tb.InlineButton{userBtn, infoBtn})

	authInlineMenu = genericInlineMenu
	//authInlineMenu = append(authInlineMenu, []tb.InlineButton{, })

	adminInlineMenu = authInlineMenu
	adminInlineMenu = append(adminInlineMenu,
		[]tb.InlineButton{authBtn, deAuthBtn},
		[]tb.InlineButton{sendMsgBtn},
	)

	superAdminInlineMenu = adminInlineMenu
	superAdminInlineMenu = append(superAdminInlineMenu, []tb.InlineButton{adminBtn, deAdminBtn})

	startMenu = append(startMenu, []tb.InlineButton{startBtn})
	myInfoMenu = append(myInfoMenu, []tb.InlineButton{backBtn})
	botInfoMenu = append(botInfoMenu, []tb.InlineButton{stopBtn}, []tb.InlineButton{backBtn})
	authUserMenu = append(authUserMenu,
		[]tb.InlineButton{authUGSopranoBtn, authUGContraltoBtn},
		[]tb.InlineButton{authUGTenoreBtn, authUGBassoBtn},
		[]tb.InlineButton{authUGCommissarioBtn, authUGReferenteBtn, authUGPreparatoreBtn},
		[]tb.InlineButton{backBtn},
	)

	return nil
}

func groupCallback(c *tb.Callback, groupName string) {
	dataContent := strings.Split(c.Data, "+")
	userID, err := strconv.Atoi(dataContent[0])
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
		return
	}
	var errAlert, authAlert string
	var add bool
	if len(dataContent) > 1 && dataContent[1] == "remove" {
		add = false
		errAlert = "Impossibile deautorizzare l'utente per il gruppo " + groupName
		authAlert = "Utente " + dataContent[0] + " rimosso dal gruppo " + groupName
	} else {
		add = true
		errAlert = "Impossibile aggiungere l'utente al gruppo " + groupName
		authAlert = "Utente " + dataContent[0] + " aggiunto al gruppo " + groupName
	}
	err = addUserGroupCmd(userID, ugContralto, add)
	if err != nil {
		bot.Respond(c, &tb.CallbackResponse{
			Text:      errAlert,
			ShowAlert: true,
		})
	} else {
		bot.Respond(c, &tb.CallbackResponse{
			Text:      authAlert,
			ShowAlert: true,
		})
	}
}

func setBotCallbacks() error {
	if bot == nil {
		return ErrNilPointer
	}

	bot.Handle(&startBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		startCmd(c.Sender, false)
	})

	bot.Handle(&stopBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		stopCmd(c.Sender)
	})

	bot.Handle(&userBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		msg, _ := getUserDescription(c.Sender)
		sendMsgWithSpecificMenu(c.Sender, msg, myInfoMenu, false)
	})
	bot.Handle(&infoBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithSpecificMenu(c.Sender, contactMsg, botInfoMenu, false)
	})
	bot.Handle(&backBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithMenu(c.Sender, menuMsg, false)
	})

	bot.Handle(&authBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithMenu(c.Sender, authHowToMsg, false)

	})
	bot.Handle(&deAuthBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithMenu(c.Sender, deAuthHowToMsg, false)

	})
	bot.Handle(&sendMsgBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})

	})
	bot.Handle(&authUGSopranoBtn, func(c *tb.Callback) {
		groupCallback(c, "Soprani")
	})
	bot.Handle(&authUGContraltoBtn, func(c *tb.Callback) {
		groupCallback(c, "Contralti")
	})
	bot.Handle(&authUGTenoreBtn, func(c *tb.Callback) {
		groupCallback(c, "Tenori")
	})
	bot.Handle(&authUGBassoBtn, func(c *tb.Callback) {
		groupCallback(c, "Bassi")
	})
	bot.Handle(&authUGCommissarioBtn, func(c *tb.Callback) {
		groupCallback(c, "Commissari")
	})
	bot.Handle(&authUGReferenteBtn, func(c *tb.Callback) {
		groupCallback(c, "Referenti")
	})
	bot.Handle(&authUGPreparatoreBtn, func(c *tb.Callback) {
		groupCallback(c, "Preparatori")
	})

	return nil
}
