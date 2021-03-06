package main

import (
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

var genericCommands = map[string]string{
	"/start":          "Far (ri)partire il bot",
	"/stop":           "Fermare il bot",
	"/menu":           "Visualizzare il menu principale",
	"/userInfo":       "Visualizzare le informazioni sul mio utente",
	"/botInfo":        "Visualizzare le informazioni sul bot",
	"/help":           "Visualizzare questo menu di aiuto",
	"/prossimoEvento": "Visualizzare le informazioni sul prossimo evento del coro",
}
var authCommands = map[string]string{
	"/prossimaProva":        "Visualizzare le informazioni sulla prossima prova",
	"/prossimaProvaSezione": "Visualizzare le informazioni sulla prossima prova di sezione",
	"/prossimaProvaInsieme": "Visualizzare le informazioni sulla prossima prova d'insieme",
	"/ultimaMail":           "Visualizzare l'ultima mail inviata alla sezione o a tutto il coro",
}
var adminCommands = map[string]string{
	"/authUser":   "Autorizzare un nuovo utente o aggiungerlo a nuovi gruppi",
	"/deAuthUser": "Deautorizzare un utente o rimuoverlo da alcuni gruppi",
	"/sendMsg":    "Inviare un messaggio ad una o piú delle sezioni di appartenenza",
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

func stopCmd(u *tb.User, newMsg bool) {
	admin, err := isBotAdmin(u.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
		return
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
		err := sendMsgWithSpecificMenu(u, stopMsg, startMenu, newMsg)
		if err != nil {
			log.Printf("Error sending message to stopped user: %v", err)
		}
	}
}

func authUserCmd(sender *tb.User, payload string, newMsg bool) {
	if payload == "" {
		err := sendMsgWithMenu(sender, authHowToMsg, newMsg)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}
	} else {
		userID, err := strconv.Atoi(payload)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		user, err := getUserInfo(userID)
		if err != nil {
			log.Printf("Error getting user info: %v", err)
			return
		}

		desc, err := getUserDescription(user)
		if err != nil {
			log.Printf("Error retriving user description: %v", err)
		}

		userGroups, err := getUserGroups(user.ID)
		if err != nil {
			log.Printf("Error retriving user groups: %v", err)
		}

		menu := getUserGroupMenu()
		menu[0][0].Data = strconv.Itoa(user.ID) + "+auth"
		menu[0][1].Data = strconv.Itoa(user.ID) + "+auth"
		menu[1][0].Data = strconv.Itoa(user.ID) + "+auth"
		menu[1][1].Data = strconv.Itoa(user.ID) + "+auth"
		menu[2][0].Data = strconv.Itoa(user.ID) + "+auth"
		menu[2][1].Data = strconv.Itoa(user.ID) + "+auth"
		menu[2][2].Data = strconv.Itoa(user.ID) + "+auth"

		for _, group := range userGroups {
			switch group {
			case ugSoprano:
				menu[0][0].Text = ""
			case ugContralto:
				menu[0][1].Text = ""
			case ugTenore:
				menu[1][0].Text = ""
			case ugBasso:
				menu[1][1].Text = ""
			case ugCommissario:
				menu[2][0].Text = ""
			case ugReferente:
				menu[2][1].Text = ""
			case ugPreparatore:
				menu[2][2].Text = ""
			}
		}

		err = sendMsgWithSpecificMenu(sender, "Stai per autorizzare il seguente utente:\n"+
			desc+
			"\nSe le informazioni sono corrette fai 'tap' sui gruppi di appartenenza dell'utente da autorizzare, altrimenti *torna al menù principale ed annulla l'autorizzazione*",
			menu, newMsg)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}
	}
}

func deAuthUserCmd(sender *tb.User, payload string, newMsg bool) {
	if payload == "" {
		err := sendMsgWithMenu(sender, deAuthHowToMsg, newMsg)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}
	} else {
		userID, err := strconv.Atoi(payload)
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
		}
		user, err := getUserInfo(userID)
		if err != nil {
			log.Printf("Error getting user info: %v", err)
			return
		}

		desc, err := getUserDescription(user)
		if err != nil {
			log.Printf("Error retriving user description: %v", err)
		}

		menu := getUserGroupMenu()
		menu[0][0].Data = strconv.Itoa(user.ID) + "+deAuth"
		menu[0][1].Data = strconv.Itoa(user.ID) + "+deAuth"
		menu[1][0].Data = strconv.Itoa(user.ID) + "+deAuth"
		menu[1][1].Data = strconv.Itoa(user.ID) + "+deAuth"
		menu[2][0].Data = strconv.Itoa(user.ID) + "+deAuth"
		menu[2][1].Data = strconv.Itoa(user.ID) + "+deAuth"
		menu[2][2].Data = strconv.Itoa(user.ID) + "+deAuth"

		if is, _ := isUserInGroup(user.ID, ugSoprano); !is {
			menu[0][0].Text = ""
		}
		if is, _ := isUserInGroup(user.ID, ugContralto); !is {
			menu[0][1].Text = ""
		}
		if is, _ := isUserInGroup(user.ID, ugTenore); !is {
			menu[1][0].Text = ""
		}
		if is, _ := isUserInGroup(user.ID, ugBasso); !is {
			menu[1][1].Text = ""
		}
		if is, _ := isUserInGroup(user.ID, ugCommissario); !is {
			menu[2][0].Text = ""
		}
		if is, _ := isUserInGroup(user.ID, ugReferente); !is {
			menu[2][1].Text = ""
		}
		if is, _ := isUserInGroup(user.ID, ugPreparatore); !is {
			menu[2][2].Text = ""
		}

		err = sendMsgWithSpecificMenu(sender, "Stai per deautorizzare il seguente utente:\n"+
			desc+
			"\nSe le informazioni sono corrette fai 'tap' sui gruppi da cui deautorizzare l'utente, altrimenti *torna al menù principale ed annulla l'autorizzazione*",
			menu, newMsg)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}

	}
}

func addUserGroupCmd(userID int, group userGroup, add bool) error {
	userGroups, err := getUserGroups(userID)
	if err != nil {
		log.Printf("Error retriving user groups: %v", err)
		return ErrAddAuthUser
	}
	is, err := isUserInGroup(userID, group)
	if err != nil {
		log.Printf("Error checking if user is in group: %v", err)
		return ErrAddAuthUser
	}
	if is && !add {
		if len(userGroups) <= 1 {
			err = authorizeUser(userID, false)
			if err != nil {
				log.Printf("Error deauthorizing user: %v", err)
				return ErrAddAuthUser
			}
		}
		for i, ug := range userGroups {
			if ug == group {
				userGroups = append(userGroups[:i], userGroups[i+1:]...)
				break
			}
		}
		err = remUserGroups(userID, userGroups, group)
		if err != nil {
			log.Printf("Error adding user in group: %v", err)
			return ErrAddAuthUser
		}
	} else if !is && add {
		if len(userGroups) == 0 {
			err = authorizeUser(userID, true)
			if err != nil {
				log.Printf("Error authorizing user: %v", err)
				return ErrAddAuthUser
			}
		}
		userGroups = append(userGroups, group)
		err = addUserGroups(userID, userGroups...)
		if err != nil {
			log.Printf("Error adding user in group: %v", err)
			return ErrAddAuthUser
		}
	}
	return nil
}

func helpCmd(user *tb.User, newMsg bool) {
	var msg string
	msg += helpGenericMsg

	isAuth, err := isAuthrizedUser(user.ID)
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
	} else if isAuth {
		msg += helpAuthMsg
	}

	isAdmin, err := isBotAdmin(user.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
	} else if isAdmin {
		msg += helpAdminMsg
	}

	err = sendMsgWithSpecificMenu(user, msg, backMenu, newMsg)
	if err != nil {
		log.Printf("Error sending message to started user: %v", err)
	}
}

func sendMsgCmd(sender *tb.User, payload string, newMsg bool) {
	if payload == "" {
		err := sendMsgWithMenu(sender, sendMsgHowToMsg, newMsg)
		if err != nil {
			log.Printf("Error in sending message: %v", err)
		}
	} else {
		msgID, err := addNewGroupMsg(sender, map[userGroup]bool{}, payload)
		if err != nil {
			log.Printf("Error adding new groupMsg in db: %v", err)
			return
		}
		menu := getGroupMsgMenu(msgID, sender)

		msg := "*Il messaggio che stai per inviare é *\n" + payload + "\n*Seleziona i gruppi a cui vuoi inviarlo*"
		sendMsgWithSpecificMenu(sender, msg, menu, true)
	}
}
