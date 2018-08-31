package main

import (
	"encoding/json"
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func modifyPrevMsg(userID int, storedMsg *tb.StoredMessage, newMsg string, newOptions *tb.SendOptions) error {
	msg, err := bot.Edit(storedMsg, newMsg, newOptions)
	if err != nil {
		log.Printf("Error modifying previous message: %v", err)
		return ErrSendMsg
	}
	err = setLastMsgPerUser(userID, msg)
	if err != nil {
		log.Printf("Error setting last msg per user: %v", err)
		return ErrSetLastMsg
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

func sendMsg(user *tb.User, msg string, new bool) error {
	sendMsgWithSpecificMenu(user, msg, nil, new)

	return nil
}

func sendMsgWithMenu(user *tb.User, msg string, new bool) error {
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
	sendMsgWithSpecificMenu(user, msg, menu, new)

	return nil
}

func sendMsgWithSpecificMenu(user *tb.User, msg string, menu [][]tb.InlineButton, new bool) error {
	if !new {
		storedMsg, err := getLastMsgPerUser(user.ID)
		if err != nil {
			log.Printf("Error retriving last message per user: %v", err)
			sentMsg, err := bot.Send(user, msg, &tb.SendOptions{
				ReplyMarkup:           &tb.ReplyMarkup{InlineKeyboard: menu},
				DisableWebPagePreview: true,
				ParseMode:             tb.ModeMarkdown,
			})
			if err != nil {
				log.Printf("Error sending message to user: %v", err)
				return ErrSendMsg
			}
			err = setLastMsgPerUser(user.ID, sentMsg)
			if err != nil {
				log.Printf("Error setting last msg per user: %v", err)
				return ErrSetLastMsg
			}
		}
		err = modifyPrevMsg(user.ID, storedMsg, msg, &tb.SendOptions{
			ReplyMarkup:           &tb.ReplyMarkup{InlineKeyboard: menu},
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
		if err != nil {
			log.Printf("Error sending message to user: %v", err)
			return ErrSendMsg
		}
	} else {
		sentMsg, err := bot.Send(user, msg, &tb.SendOptions{
			ReplyMarkup:           &tb.ReplyMarkup{InlineKeyboard: menu},
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
		if err != nil {
			log.Printf("Error sending message to user: %v", err)
			return ErrSendMsg
		}
		err = setLastMsgPerUser(user.ID, sentMsg)
		if err != nil {
			log.Printf("Error setting last msg per user: %v", err)
			return ErrSetLastMsg
		}
	}

	return nil
}
