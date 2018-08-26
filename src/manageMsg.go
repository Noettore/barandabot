package main

import (
	"encoding/json"
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func closeMsgMenu(storedMsg *tb.StoredMessage) error {
	_, err := bot.EditReplyMarkup(storedMsg, &tb.ReplyMarkup{
		InlineKeyboard: nil,
	})
	if err != nil {
		log.Printf("Error modifying the message: %v", err)
	}

	return nil
}

func setLastMsgPerUser(userID int, msg *tb.Message) error {
	storedMsg := tb.StoredMessage{
		MessageID: strconv.Itoa(msg.ID),
		ChatID:    msg.Chat.ID}

	jsonMsg, err := json.Marshal(storedMsg)
	if err != nil {
		log.Printf("Error in marshalling msg to json: %v", err)
		return ErrJSONMarshall
	}
	err = redisClient.HSet(lastMsgPerUser, strconv.Itoa(userID), jsonMsg).Err()
	if err != nil {
		log.Printf("Error adding last message per user info in hash: %v", err)
		return ErrRedisAddHash
	}

	return nil
}

func getLastMsgPerUser(userID int) (*tb.StoredMessage, error) {
	msg, err := redisClient.HGet(lastMsgPerUser, strconv.Itoa(userID)).Result()
	if err != nil {
		log.Printf("Error retriving last msg per user info from hash: %v", err)
		return nil, ErrRedisRetrieveHash
	}
	jsonMsg := &tb.StoredMessage{}
	err = json.Unmarshal([]byte(msg), jsonMsg)
	if err != nil {
		log.Printf("Error unmarshalling last msg per user info: %v", err)
		return nil, ErrJSONUnmarshall
	}
	return jsonMsg, nil
}

func sendMsg(user *tb.User, msg string) error {
	sentMsg, err := bot.Send(user, msg, &tb.SendOptions{
		ParseMode: "Markdown",
	})
	if err != nil {
		log.Printf("Error sending message to user: %v", err)
		return ErrSendMsg
	}
	storedMsg, err := getLastMsgPerUser(user.ID)
	if err != nil {
		log.Printf("Error retriving last message per user: %v", err)
	} else {
		err = closeMsgMenu(storedMsg)
		if err != nil {
			log.Printf("Error modifying the message: %v", err)
		}
	}
	err = setLastMsgPerUser(user.ID, sentMsg)
	if err != nil {
		log.Printf("Error setting last msg per user: %v", err)
		return ErrSetLastMsg
	}
	return nil
}

func sendMsgWithMenu(user *tb.User, msg string) error {
	var menu [][]tb.InlineButton

	auth, err := isAuthrizedUser(user.ID)
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
	}
	admin, err := isBotAdmin(user.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
	}

	if admin {
		menu = adminInlineMenu
	} else if auth {
		menu = authInlineMenu
	} else {
		menu = genericInlineMenu
	}
	sentMsg, err := bot.Send(user, msg, &tb.SendOptions{
		ReplyMarkup: &tb.ReplyMarkup{
			InlineKeyboard: menu,
		},
		ParseMode: "Markdown",
	})
	if err != nil {
		log.Printf("Error sending message to user: %v", err)
		return ErrSendMsg
	}
	storedMsg, err := getLastMsgPerUser(user.ID)
	if err != nil {
		log.Printf("Error retriving last message per user: %v", err)
	} else {
		err = closeMsgMenu(storedMsg)
		if err != nil {
			log.Printf("Error modifying the message: %v", err)
		}
	}
	err = setLastMsgPerUser(user.ID, sentMsg)
	if err != nil {
		log.Printf("Error setting last msg per user: %v", err)
		return ErrSetLastMsg
	}
	return nil
}

func sendMsgWithSpecificMenu(user *tb.User, msg string, menu [][]tb.InlineButton) error {
	sentMsg, err := bot.Send(user, msg, &tb.SendOptions{
		ReplyMarkup: &tb.ReplyMarkup{
			InlineKeyboard: menu,
		},
		ParseMode: "Markdown",
	})
	if err != nil {
		log.Printf("Error sending message to user: %v", err)
		return ErrSendMsg
	}
	storedMsg, err := getLastMsgPerUser(user.ID)
	if err != nil {
		log.Printf("Error retriving last message per user: %v", err)
	} else {
		err = closeMsgMenu(storedMsg)
		if err != nil {
			log.Printf("Error modifying the message: %v", err)
		}
	}
	err = setLastMsgPerUser(user.ID, sentMsg)
	if err != nil {
		log.Printf("Error setting last msg per user: %v", err)
		return ErrSetLastMsg
	}

	return nil
}
