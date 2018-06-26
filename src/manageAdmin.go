package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dixonwille/wmenu"
	tb "gopkg.in/tucnak/telebot.v2"
)

func addBotAdmins(newAdminIDs []string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	errNum := 0
	if newAdminIDs == nil && cmdFlags.interactive {
		fmt.Println("Add the new admin IDs, comma-separated:")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error in reading new admin IDs: %v", err)
			return ErrStdRead
		}
		newAdminIDs = strings.Split(line, ",")
	}
	for _, newAdminID := range newAdminIDs {
		err := addBotAdmin(newAdminID)
		if err != nil {
			errNum++
			log.Printf("Error in adding new admin ID %s: %v", newAdminID, err)
		}
	}
	if errNum == len(newAdminIDs) {
		return ErrAddAdmin
	}
	return nil
}

func removeBotAdmins() error {
	if redisClient == nil {
		return ErrNilPointer
	}
	botAdmins, err := redisClient.SMembers(adminSet).Result()
	if err != nil {
		log.Printf("Couldn't retrieve admins: %v", err)
		return ErrRedisRetrieveSet
	}
	menu := wmenu.NewMenu("Select the admin(s) you want to remove:")
	menu.AllowMultiple()
	menu.LoopOnInvalid()
	menu.Action(func(opts []wmenu.Opt) error {
		var returnErr error
		for _, opt := range opts {
			if opt.Value == nil {
				log.Println("Couldn't remove admin: nil token")
				returnErr = ErrNilPointer
			} else {
				err := removeBotAdmin(opt.Value.(int))
				if err != nil {
					log.Printf("Couldn't remove bot admin: %v", err)
				}
			}
		}
		return returnErr
	})
	//for _, token := range tokens {
	for _, botAdmin := range botAdmins {
		adminID, err := strconv.Atoi(botAdmin)
		if err != nil {
			log.Printf("Error converting adminID from string to int: %v", err)
			return ErrAtoiConv
		}
		adminInfo, err := getUserInfo(adminID)
		menu.Option(adminInfo.Username+": "+adminInfo.FirstName+" "+adminInfo.LastName, adminID, false, nil)
	}
	err = menu.Run()
	if err != nil {
		log.Printf("Error in removeToken menu: %v", err)
		return ErrRemoveAdmin
	}
	return nil
}

func addBotAdmin(newAdminID string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	adminID := strings.TrimSpace(newAdminID)
	matched, err := regexp.MatchString("^\\d+$", adminID)
	if err != nil {
		log.Printf("Error in parsing admin ID: %v", err)
		return ErrIDParsing
	}
	if !matched {
		return ErrIDInvalid
	}

	ID, err := strconv.Atoi(adminID)
	if err != nil {
		log.Printf("Error converting user ID: %v", err)
		return ErrAtoiConv
	}
	chat, err := bot.ChatByID(adminID)
	if err != nil {
		log.Printf("Error retriving chat by id: %v", err)
		return ErrChatRetrieve
	}
	if chat.Type != tb.ChatPrivate {
		log.Printf("Admin must be a user!")
		return ErrAddAdmin
	}
	isUser, err := isUser(ID)
	if err != nil {
		log.Printf("Error checking if ID is bot user: %v", err)
		return ErrAddAdmin
	}
	if !isUser {
		err = addUser(&tb.User{int(chat.ID), chat.FirstName, chat.LastName, chat.Username})
		if err != nil {
			log.Printf("Error adding user: %v", err)
			return ErrAddUser
		}
	}
	err = redisClient.SAdd(adminSet, adminID).Err()
	if err != nil {
		log.Printf("Error in adding new admin ID: %v", err)
		return ErrRedisAddSet
	}

	err = authorizeUser(ID, true)
	if err != nil {
		log.Printf("Error in adding new admin ID in authorized users: %v", err)
		return ErrAddAuthUser
	}
	user, err := getUserInfo(ID)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return ErrGetUser
	}
	err = sendMessage(user, "Sei stato aggiunto come amministratore del BarandaBot")
	if err != nil {
		log.Printf("Error sending message to new admin: %v", err)
		return ErrSendMsg
	}

	return nil
}

func removeBotAdmin(adminID int) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.SRem(adminSet, strconv.Itoa(adminID)).Err()
	if err != nil {
		log.Printf("Error removing admin from set: %v", err)
		return ErrRedisRemSet
	}
	return nil
}
