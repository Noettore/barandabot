package main

import (
	"log"
	"strconv"

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

	})
	bot.Handle(&sendMsgBtn, func(c *tb.Callback) {
		bot.Respond(c, &tb.CallbackResponse{})

	})
	bot.Handle(&authUGSopranoBtn, func(c *tb.Callback) {
		userID, err := strconv.Atoi(c.Data)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		err = addUserGroupCmd(userID, ugSoprano)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Impossibile autorizzare l'utente",
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Autorizzato utente " + c.Data + "e aggiunto al gruppo Soprani",
				ShowAlert: true,
			})
		}
	})
	bot.Handle(&authUGContraltoBtn, func(c *tb.Callback) {
		userID, err := strconv.Atoi(c.Data)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		err = addUserGroupCmd(userID, ugContralto)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Impossibile autorizzare l'utente",
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Autorizzato utente " + c.Data + "e aggiunto al gruppo Contralti",
				ShowAlert: true,
			})
		}

	})
	bot.Handle(&authUGTenoreBtn, func(c *tb.Callback) {
		userID, err := strconv.Atoi(c.Data)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		err = addUserGroupCmd(userID, ugTenore)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Impossibile autorizzare l'utente",
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Autorizzato utente " + c.Data + "e aggiunto al gruppo Tenori",
				ShowAlert: true,
			})
		}
	})
	bot.Handle(&authUGBassoBtn, func(c *tb.Callback) {
		userID, err := strconv.Atoi(c.Data)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		err = addUserGroupCmd(userID, ugBasso)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Impossibile autorizzare l'utente",
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Autorizzato utente " + c.Data + "e aggiunto al gruppo Bassi",
				ShowAlert: true,
			})
		}
	})
	bot.Handle(&authUGCommissarioBtn, func(c *tb.Callback) {
		userID, err := strconv.Atoi(c.Data)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		err = addUserGroupCmd(userID, ugCommissario)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Impossibile autorizzare l'utente",
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Autorizzato utente " + c.Data + "e aggiunto al gruppo Commissari",
				ShowAlert: true,
			})
		}

	})
	bot.Handle(&authUGReferenteBtn, func(c *tb.Callback) {
		userID, err := strconv.Atoi(c.Data)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		err = addUserGroupCmd(userID, ugReferente)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Impossibile autorizzare l'utente",
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Autorizzato utente " + c.Data + "e aggiunto al gruppo Referenti",
				ShowAlert: true,
			})
		}
	})
	bot.Handle(&authUGPreparatoreBtn, func(c *tb.Callback) {
		userID, err := strconv.Atoi(c.Data)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		err = addUserGroupCmd(userID, ugPreparatore)
		if err != nil {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Impossibile autorizzare l'utente",
				ShowAlert: true,
			})
		} else {
			bot.Respond(c, &tb.CallbackResponse{
				Text:      "Autorizzato utente " + c.Data + "e aggiunto al gruppo Preparatori",
				ShowAlert: true,
			})
		}
	})

	return nil
}
