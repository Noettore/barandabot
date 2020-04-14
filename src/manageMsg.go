package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type groupMsg struct {
	SenderID int       `sql:"sender_id" json:"sender_id"`
	Group    userGroup `sql:"group" json:"group"`
	Msg      string    `sql:"msg" json:"msg"`
	Date     time.Time `sql:"date" json:"date"`
	Sent     bool      `sql:"sent" json:"sent"`
}

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
		ChatID:    msg.Chat.ID,
	}

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

func addNewGroupMsg(sender *tb.User, group userGroup, msg string) (int64, error) {
	newGroupMsg := groupMsg{sender.ID, group, msg, time.Now(), false}
	jsonMsg, err := json.Marshal(newGroupMsg)
	if err != nil {
		log.Printf("Error in marshalling groupMsg to json: %v", err)
		return -1, ErrJSONMarshall
	}
	//err = redisClient.HSet(lastMsgPerUser, strconv.Itoa(userID), jsonMsg).Err()
	msgID, err := redisClient.RPush(groupMsgs, jsonMsg).Result()
	if err != nil {
		log.Printf("Error adding new group message in hash: %v", err)
		return -1, ErrRedisAddList
	}
	return msgID - 1, nil
}

func setGroupMsg(msg *groupMsg, index int64) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error in marshalling groupMsg to json: %v", err)
		return ErrJSONMarshall
	}
	err = redisClient.LSet(groupMsgs, index, jsonMsg).Err()
	if err != nil {
		log.Printf("Error modifying a groupMsg: %v", err)
		return ErrRedisSetList
	}
	return nil
}

func sendMsgToGroup(msgID string) error {
	ID, err := strconv.ParseInt(msgID, 10, 64)
	if err != nil {
		log.Printf("Error converting msgID to int64: %v", err)
		return ErrAtoiConv
	}
	msg, err := redisClient.LIndex(groupMsgs, ID).Result()
	if err != nil {
		log.Printf("Error retriving group message from hash: %v", err)
		return ErrRedisRetrieveHash
	}
	jsonMsg := &groupMsg{}
	err = json.Unmarshal([]byte(msg), jsonMsg)
	if err != nil {
		log.Printf("Error unmarshalling groupMsg: %v", err)
		return ErrJSONUnmarshall
	}
	if jsonMsg.Sent {
		return ErrSendMsg
	}
	sender, err := getUserInfo(jsonMsg.SenderID)
	if err != nil {
		log.Printf("Error retrieving sender info: %v", err)
		return ErrGetUser
	}
	users, err := getUsersInGroup(jsonMsg.Group)
	if err != nil {
		log.Printf("Error retrieving users in sendTo group: %v", err)
		return ErrGroupInvalid
	}
	for _, userID := range users {
		user, err := getUserInfo(userID)
		if err != nil {
			log.Printf("Error retrieving user info from id: %v", err)
			continue
		}
		groupName, _ := getGroupName(jsonMsg.Group)
		msg = "*Messaggio inviato da " + sender.FirstName + " a tutta la sezione " + groupName + "*\n" + jsonMsg.Msg
		err = sendMsg(user, msg, true)
		if err != nil {
			log.Printf("Error sending msg to user: %v", err)
		}
		err = sendMsgWithMenu(user, msgReceivedMsg, true)
		if err != nil {
			log.Printf("Error sending msg to user: %v", err)
		}
	}

	jsonMsg.Sent = true
	jsonMsg.Date = time.Now()
	err = setGroupMsg(jsonMsg, ID)
	if err != nil {
		log.Printf("Error updating groupMsg after send: %v", err)
	}

	err = sendMsgWithMenu(sender, "Messaggio inviato a tutti i componenti della sezione", false)
	if err != nil {
		log.Printf("Error sending msg to sender: %v", err)
	}

	return nil
}
