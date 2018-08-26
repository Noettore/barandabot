package main

import (
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

var genericCommands = map[string]bool{
	"/start":          true,
	"/stop":           true,
	"/menu":           true,
	"/userInfo":       true,
	"/botInfo":        true,
	"/prossimoEvento": true,
}
var authCommands = map[string]bool{
	"/prossimaProvaSezione": true,
	"/prossimaProvaInsieme": true,
}
var adminCommands = map[string]bool{
	"/authUser": true,
	"/addAdmin": true,
	"/delAdmin": true,
}

func startCmd(u *tb.User) {
	var msg string

	isUser, err := isUser(u.ID)
	if err != nil {
		log.Printf("Error checking if ID is bot user: %v", err)
	}

	started, err := isStartedUser(u.ID)
	if err != nil {
		log.Printf("Error checking if user is started: %v", err)
	}
	if !started {
		err = startUser(u.ID, true)
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

	err = sendMsgWithMenu(u, msg)
	if err != nil {
		log.Printf("Error sending message to started user: %v", err)
	}
}

func stopCmd(u *tb.User) {
	admin, err := isBotAdmin(u.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
	}
	if admin {
		err := sendMsgWithMenu(u, unstoppableMsg)
		if err != nil {
			log.Printf("Error sending message to unstoppable user: %v", err)
		}
	} else {
		err = startUser(u.ID, false)
		if err != nil {
			log.Printf("Error starting user: %v", err)
		}
		err := sendMsgWithSpecificMenu(u, stopMsg, startMenu)
		if err != nil {
			log.Printf("Error sending message to stopped user: %v", err)
		}
	}
}

func userInfoCmd(u *tb.User) {
	userGroups, err := getUserGroups(u.ID)
	if err != nil {
		log.Printf("Error retriving user groups: %v", err)
	}
	stringGroups := convertUserGroups(userGroups)

	isAdmin, err := isBotAdmin(u.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
	}
	isAuth, err := isAuthrizedUser(u.ID)
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
	}

	msg := "\xF0\x9F\x91\xA4 *INFORMAZIONI UTENTE*" +
		"\n- *Nome*: " + u.FirstName +
		"\n- *Username*: " + u.Username +
		"\n- *ID*: " + strconv.Itoa(u.ID) +
		"\n- *Gruppi*: "

	for _, group := range stringGroups {
		msg += group + ", "
	}

	msg += "\n- *Tipo utente*: "

	if isAdmin {
		msg += "Admin"
	} else if isAuth {
		msg += "Autorizzato"
	} else {
		msg += "Utente semplice"
	}
	err = sendMsgWithSpecificMenu(u, msg, goBackMenu)
}
