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
	backMenu             [][]tb.InlineButton
	botInfoMenu          [][]tb.InlineButton
	authUserMenu         [][]tb.InlineButton
	sendMsgMenu          [][]tb.InlineButton
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
		Text:   "\xF0\x9F\x94\x99 Torna indietro",
	}
	infoBtn = tb.InlineButton{
		Unique: "info_btn",
		Text:   "\xE2\x84\xB9 Bot info",
	}
	userBtn = tb.InlineButton{
		Unique: "user_btn",
		Text:   "\xF0\x9F\x91\xA4 My info",
	}
	helpBtn = tb.InlineButton{
		Unique: "help_btn",
		Text:   "\xF0\x9F\x86\x98 Aiuto",
	}
	authBtn = tb.InlineButton{
		Unique: "auth_btn",
		Text:   "\xE2\x9E\x95 Autorizza utente",
	}
	deAuthBtn = tb.InlineButton{
		Unique: "de_auth_btn",
		Text:   "\xE2\x9E\x96 Deautorizza utente",
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
	confirmSendBtn = tb.InlineButton{
		Unique: "confirm_send_btn",
		Text:   "\xE2\x9C\x85 Conferma",
	}
	ugSopranoBtn = tb.InlineButton{
		Unique: "ugSoprano_btn",
		Text:   "\xF0\x9F\x91\xA7 Soprani",
	}
	ugContraltoBtn = tb.InlineButton{
		Unique: "ugContralto_btn",
		Text:   "\xF0\x9F\x91\xA9 Contralti",
	}
	ugTenoreBtn = tb.InlineButton{
		Unique: "ugTenore_btn",
		Text:   "\xF0\x9F\x91\xA6 Tenori",
	}
	ugBassoBtn = tb.InlineButton{
		Unique: "ugBasso_btn",
		Text:   "\xF0\x9F\x91\xA8 Bassi",
	}
	ugCommissarioBtn = tb.InlineButton{
		Unique: "ugCommissario_btn",
		Text:   "\xF0\x9F\x93\x9D Commissari",
	}
	ugReferenteBtn = tb.InlineButton{
		Unique: "ugReferente_btn",
		Text:   "\xF0\x9F\x93\x8B Referenti",
	}
	ugPreparatoreBtn = tb.InlineButton{
		Unique: "ugPreparatori_btn",
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
	backMenu = append(backMenu, []tb.InlineButton{backBtn})
	botInfoMenu = append(botInfoMenu, []tb.InlineButton{helpBtn, stopBtn}, []tb.InlineButton{backBtn})
	/* authUserMenu = append(authUserMenu,
		[]tb.InlineButton{ugSopranoBtn, ugContraltoBtn},
		[]tb.InlineButton{ugTenoreBtn, ugBassoBtn},
		[]tb.InlineButton{ugCommissarioBtn, ugReferenteBtn, ugPreparatoreBtn},
		[]tb.InlineButton{backBtn},
	) */
	sendMsgMenu = append(sendMsgMenu, []tb.InlineButton{confirmSendBtn, backBtn})

	return nil
}

func getUserGroupMenu() [][]tb.InlineButton {
	var ugMenu [][]tb.InlineButton
	ugMenu = append(ugMenu,
		[]tb.InlineButton{ugSopranoBtn, ugContraltoBtn},
		[]tb.InlineButton{ugTenoreBtn, ugBassoBtn},
		[]tb.InlineButton{ugCommissarioBtn, ugReferenteBtn, ugPreparatoreBtn},
		[]tb.InlineButton{backBtn},
	)
	return ugMenu
}

func ugBtnCallback(c *tb.Callback, group userGroup) {
	dataContent := strings.Split(c.Data, "+")
	if len(dataContent) <= 1 {
		log.Printf("Error: too few arguments")
		return
	}
	var userID int
	var msgID int64
	var errAlert, successAlert string
	auth, deAuth, sendUg := false, false, false

	groupName, err := getGroupName(group)

	if dataContent[1] == "deAuth" {
		deAuth = true
		errAlert = "Impossibile deautorizzare l'utente per il gruppo " + groupName
		successAlert = "Utente " + dataContent[0] + " rimosso dal gruppo " + groupName
		userID, err = strconv.Atoi(dataContent[0])
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
			return
		}
	} else if dataContent[1] == "auth" {
		auth = true
		errAlert = "Impossibile aggiungere l'utente al gruppo " + groupName
		successAlert = "Utente " + dataContent[0] + " aggiunto al gruppo " + groupName
		userID, err = strconv.Atoi(dataContent[0])
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
			return
		}
	} else if dataContent[1] == "sendUg" {
		sendUg = true
		errAlert = "Impossibile aggiungere il gruppo " + groupName + " alla lista dei destinatari"
		successAlert = "Gruppo " + groupName + " aggiunto alla lista dei destinatari"
		msgID, err = strconv.ParseInt(dataContent[0], 10, 64)
		if err != nil {
			log.Printf("Error converting msgID to int64: %v", err)
			return
		}
	}
	if auth || deAuth {
		err = addUserGroupCmd(userID, group, auth)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      errAlert,
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      successAlert,
				ShowAlert: true,
			})
			if auth {
				authUserCmd(c.Sender, dataContent[0], false)
			} else if deAuth {
				deAuthUserCmd(c.Sender, dataContent[0], false)
			}
		}
	} else if sendUg {
		err = addUGToGroupMsg(msgID, group)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      errAlert,
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      successAlert,
				ShowAlert: true,
			})
		}
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
		stopCmd(c.Sender, false)
	})

	bot.Handle(&userBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		msg, _ := getUserDescription(c.Sender)
		sendMsgWithSpecificMenu(c.Sender, msg, backMenu, false)
	})

	bot.Handle(&infoBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithSpecificMenu(c.Sender, contactMsg, botInfoMenu, false)
	})

	bot.Handle(&helpBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		helpCmd(c.Sender, false)
	})

	bot.Handle(&backBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})
		sendMsgWithMenu(c.Sender, menuMsg, false)
	})

	bot.Handle(&confirmSendBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{
			Text:      sentStartedMsg,
			ShowAlert: true,
		})
		sendMsgToGroup(c.Data)
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
		sendMsgWithMenu(c.Sender, sendMsgHowToMsg, false)
	})

	bot.Handle(&ugSopranoBtn, func(c *tb.Callback) {
		ugBtnCallback(c, ugSoprano)
	})

	bot.Handle(&ugContraltoBtn, func(c *tb.Callback) {
		ugBtnCallback(c, ugContralto)
	})

	bot.Handle(&ugTenoreBtn, func(c *tb.Callback) {
		ugBtnCallback(c, ugTenore)
	})

	bot.Handle(&ugBassoBtn, func(c *tb.Callback) {
		ugBtnCallback(c, ugBasso)
	})

	bot.Handle(&ugCommissarioBtn, func(c *tb.Callback) {
		ugBtnCallback(c, ugCommissario)
	})

	bot.Handle(&ugReferenteBtn, func(c *tb.Callback) {
		ugBtnCallback(c, ugReferente)
	})

	bot.Handle(&ugPreparatoreBtn, func(c *tb.Callback) {
		ugBtnCallback(c, ugPreparatore)
	})

	return nil
}
