package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

type userGroup int

const (
	ugEsterno userGroup = iota
	ugSoprano
	ugContralto
	ugTenore
	ugBasso
	ugCommissario
	ugReferente
	ugPreparatore
)

func addUser(user *tb.User) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.SAdd(usersID, user.ID).Err()
	if err != nil {
		log.Printf("Error in adding user ID: %v", err)
		return ErrRedisAddSet
	}
	jsonUser, err := json.Marshal(&user)
	if err != nil {
		log.Printf("Error in marshalling user to json: %v", err)
		return ErrJSONMarshall
	}
	err = redisClient.HSet(usersInfo, strconv.Itoa(user.ID), jsonUser).Err()
	if err != nil {
		log.Printf("Error adding user info in hash: %v", err)
		return ErrRedisAddHash
	}

	err = setUserGroups(user.ID, ugEsterno)
	if err != nil {
		log.Printf("Error setting user default group: %v", err)
		return ErrRedisAddSet
	}

	return nil
}

func isUser(userID int) (bool, error) {
	if redisClient == nil {
		return false, ErrNilPointer
	}
	user, err := redisClient.SIsMember(usersID, strconv.Itoa(userID)).Result()
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
	user, err := redisClient.HGet(usersInfo, strconv.Itoa(userID)).Result()
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

func isStartedUser(userID int) (bool, error) {
	if redisClient == nil {
		return false, ErrNilPointer
	}
	started, err := redisClient.SIsMember(startedUsers, strconv.Itoa(userID)).Result()
	if err != nil {
		log.Printf("Error checking if user is started: %v", err)
		return false, ErrRedisCheckSet
	}
	return started, nil
}

func startUser(userID int, start bool) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	if start {
		err := redisClient.SAdd(startedUsers, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error adding token to set: %v", err)
			return ErrRedisAddSet
		}
	} else {
		err := redisClient.SRem(startedUsers, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error removing token from set: %v", err)
			return ErrRedisRemSet
		}
	}
	return nil
}

func isAuthrizedUser(userID int) (bool, error) {
	if redisClient == nil {
		return false, ErrNilPointer
	}
	auth, err := redisClient.SIsMember(authUsers, strconv.Itoa(userID)).Result()
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
		err := redisClient.SAdd(authUsers, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error adding token to set: %v", err)
			return ErrRedisAddSet
		}
	} else {
		err := redisClient.SRem(authUsers, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error removing token from set: %v", err)
			return ErrRedisRemSet
		}
	}
	return nil
}

func setUserGroups(userID int, groups ...userGroup) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	var csvGroups string
	for _, group := range groups {
		csvGroups += strconv.Itoa(int(group)) + ","
		err := redisClient.SAdd("ug"+strconv.Itoa(int(group)), strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error adding user to usergroup set: %v", err)
			return ErrRedisAddSet
		}
	}
	err := redisClient.HSet(usersGroups, strconv.Itoa(userID), csvGroups).Err()
	if err != nil {
		log.Printf("Error adding user groups to hash: %v", err)
		return ErrRedisAddHash
	}

	return nil
}

func getUserGroups(userID int) ([]userGroup, error) {
	if redisClient == nil {
		return nil, ErrNilPointer
	}

	csvGroups, err := redisClient.HGet(usersGroups, strconv.Itoa(userID)).Result()
	if err != nil {
		log.Printf("Error retriving user groups: %v", err)
		return nil, ErrRedisRetrieveHash
	}
	var retGroups []userGroup
	groups := strings.Split(csvGroups, ",")
	for _, group := range groups {
		intGroup, err := strconv.Atoi(group)
		if err != nil {
			log.Printf("Error converting user group: %v", err)
			return nil, ErrAtoiConv
		}
		retGroups = append(retGroups, userGroup(intGroup))
	}
	return retGroups, nil
}

func getUsersInGroup(group userGroup) ([]int, error) {
	users, err := redisClient.SMembers("ug" + strconv.Itoa(int(group))).Result()
	if err != nil {
		log.Printf("Error retriving users in group: %v", err)
		return nil, ErrRedisRetrieveSet
	}
	var retUsers []int
	for _, user := range users {
		intUser, err := strconv.Atoi(user)
		if err != nil {
			log.Printf("Error converting user ID: %v", err)
			return nil, ErrAtoiConv
		}
		retUsers = append(retUsers, intUser)
	}
	return retUsers, nil
}

func convertUserGroups(groups []userGroup) []string {
	var stringGroups []string
	for _, group := range groups {
		switch group {
		case ugEsterno:
			stringGroups = append(stringGroups, "Esterno al coro")
		case ugSoprano:
			stringGroups = append(stringGroups, "Soprano")
		case ugContralto:
			stringGroups = append(stringGroups, "Contralto")
		case ugTenore:
			stringGroups = append(stringGroups, "Tenore")
		case ugBasso:
			stringGroups = append(stringGroups, "Basso")
		case ugCommissario:
			stringGroups = append(stringGroups, "Commissario")
		case ugReferente:
			stringGroups = append(stringGroups, "Referente")
		case ugPreparatore:
			stringGroups = append(stringGroups, "Preparatore")
		}
	}

	return stringGroups
}
