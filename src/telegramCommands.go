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
	"/config":         true,
	"/botInfo":        true,
	"/help":           true,
	"/prossimoEvento": true,
}
var authCommands = map[string]bool{
	"/prossimaProvaSezione": true,
	"/prossimaProvaInsieme": true,
}
var adminCommands = map[string]bool{
	"/authUser":   true,
	"/deAuthUser": true,
	"/addAdmin":   true,
	"/delAdmin":   true,
}

func startCmd(u *tb.User, newMsg bool) {
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
			err := addUser(u)
			if err != nil {
				log.Printf("Error adding user: %v", err)
			}
			msg = startMsg
		}
	} else {
		msg = alreadyStartedMsg
	}

	err = sendMsgWithMenu(u, msg, newMsg)
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
		//img := &tb.Photo{File: tb.FromDisk()}
		err := sendMsgWithMenu(u, unstoppableMsg, false)
		if err != nil {
			log.Printf("Error sending message to unstoppable user: %v", err)
		}
	} else {
		err = startUser(u.ID, false)
		if err != nil {
			log.Printf("Error starting user: %v", err)
		}
		err := sendMsgWithSpecificMenu(u, stopMsg, startMenu, false)
		if err != nil {
			log.Printf("Error sending message to stopped user: %v", err)
		}
	}
}

func authUserCmd(u *tb.User, payload string) {
	if payload == "" {
		err := sendMsg(u, authHowToMsg, true)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}
	} else {
		//TODO check if payload is valid ID
		desc, err := getUserDescription(u)
		if err != nil {
			log.Printf("Error retriving user description: %v", err)
		}
		menu := authUserMenu
		menu[0][0].Data = payload
		menu[0][1].Data = payload
		menu[1][0].Data = payload
		menu[1][1].Data = payload
		menu[2][0].Data = payload
		menu[2][1].Data = payload
		menu[2][2].Data = payload
		err = sendMsgWithSpecificMenu(u, "Stai per autorizzare il seguente utente:\n"+
			desc+
			"\nSe le informazioni sono corrette fai 'tap' sui gruppi di appartenenza dell'utente da autorizzare, altrimenti *torna al menù principale ed annulla l'autorizzazione*",
			menu, true)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}
	}
}

func deAuthUserCmd(u *tb.User, payload string) {
	if payload == "" {
		err := sendMsg(u, deAuthHowToMsg, true)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}
	} else {
		userID, err := strconv.Atoi(payload)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		authorizeUser(userID, false)
		//TODO
	}
}

func addUserGroupCmd(userID int, group userGroup) error {
	userGroups, err := getUserGroups(userID)
	if err != nil {
		log.Printf("Error retriving user groups: %v", err)
	}
	is, err := isUserInGroup(userID, group)
	if err != nil {
		log.Printf("Error checking if user is in group: %v", err)
	}
	if is {
		return ErrAddUser
	}
	userGroups = append(userGroups, group)
	err = setUserGroups(userID, userGroups...)
	if err != nil {
		log.Printf("Error adding user in group: %v", err)
		return ErrAddAuthUser
	}

	return nil
}
