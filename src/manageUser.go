package main

import (
	"encoding/json"
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func addUser(user *tb.User) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.SAdd(userSet, user.ID).Err()
	if err != nil {
		log.Printf("Error in adding user ID: %v", err)
		return ErrRedisAddSet
	}
	jsonUser, err := json.Marshal(&user)
	if err != nil {
		log.Printf("Error in marshalling user to json: %v", err)
		return ErrJSONMarshall
	}
	err = redisClient.HSet(userHash, strconv.Itoa(user.ID), jsonUser).Err()
	if err != nil {
		log.Printf("Error adding user info in hash: %v", err)
		return ErrRedisAddHash
	}

	return nil
}

func isUser(userID int) (bool, error) {
	if redisClient == nil {
		return false, ErrNilPointer
	}
	user, err := redisClient.SIsMember(userSet, strconv.Itoa(userID)).Result()
	if err != nil {
		log.Printf("Error checking if ID is bot user: %v", err)
		return false, ErrRedisCheckSet
	}
	return user, nil
}

func getUserInfo(userID int) (*tb.User, error) {
	if redisClient == nil {
		return nil, ErrNilPointer
	}
	user, err := redisClient.HGet(userHash, strconv.Itoa(userID)).Result()
	if err != nil {
		log.Printf("Error retriving user info from hash: %v", err)
		return nil, ErrRedisRetrieveHash
	}
	jsonUser := &tb.User{}
	err = json.Unmarshal([]byte(user), jsonUser)
	if err != nil {
		log.Printf("Error unmarshalling user info: %v", err)
		return nil, ErrJSONUnmarshall
	}
	return jsonUser, nil
}

func isAuthrizedUser(userID int) (bool, error) {
	if redisClient == nil {
		return false, ErrNilPointer
	}
	auth, err := redisClient.SIsMember(authUserSet, strconv.Itoa(userID)).Result()
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
		return false, ErrRedisCheckSet
	}
	return auth, nil
}

func authorizeUser(userID int, authorized bool) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	if authorized {
		err := redisClient.SAdd(authUserSet, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error adding token to set: %v", err)
			return ErrRedisAddSet
		}
	} else {
		err := redisClient.SRem(authUserSet, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error removing token from set: %v", err)
			return ErrRedisRemSet
		}
	}
	return nil
}
