package main

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	tb "gopkg.in/tucnak/telebot.v2"
)

type userGroup int

const (
	ugSoprano userGroup = iota
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

func authorizeUser(userID int, authorize bool) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	isAuthUser, err := isAuthrizedUser(userID)
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
	}
	if isAuthUser && authorize {
		return nil
	}

	user, err := getUserInfo(userID)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return ErrGetUser
	}
	if authorize {
		err := redisClient.SAdd(authUsers, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error adding token to set: %v", err)
			return ErrRedisAddSet
		}
		err = sendMsgWithMenu(user, newAuthMsg, true)
		if err != nil {
			log.Printf("Error sending message to new authorized user: %v", err)
			return ErrSendMsg
		}
	} else {
		err := redisClient.SRem(authUsers, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error removing token from set: %v", err)
			return ErrRedisRemSet
		}
		err = sendMsg(user, delAuthMsg, true)
		if err != nil {
			log.Printf("Error sending message to removed authorized user: %v", err)
			return ErrSendMsg
		}
	}
	return nil
}

func addUserGroups(userID int, groups ...userGroup) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i] < groups[j] })
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

func remUserGroups(userID int, newGroups []userGroup, remGroups ...userGroup) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	sort.Slice(newGroups, func(i, j int) bool { return newGroups[i] < newGroups[j] })

	for _, remGroup := range remGroups {
		err := redisClient.SRem("ug"+strconv.Itoa(int(remGroup)), strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error removing user from usergroup set: %v", err)
			return ErrRedisAddSet
		}
	}

	if len(newGroups) > 0 {
		var csvGroups string
		for _, group := range newGroups {
			csvGroups += strconv.Itoa(int(group)) + ","
		}
		err := redisClient.HSet(usersGroups, strconv.Itoa(userID), csvGroups).Err()
		if err != nil {
			log.Printf("Error adding user groups to hash: %v", err)
			return ErrRedisAddHash
		}
	} else {
		err := redisClient.HDel(usersGroups, strconv.Itoa(userID)).Err()
		if err != nil {
			log.Printf("Error removing user from usersGroups hash: %v", err)
			return ErrRedisAddHash
		}
	}

	return nil
}

func getUserGroups(userID int) ([]userGroup, error) {
	if redisClient == nil {
		return nil, ErrNilPointer
	}

	csvGroups, err := redisClient.HGet(usersGroups, strconv.Itoa(userID)).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error retriving user groups: %v", err)
		return nil, ErrRedisRetrieveHash
	}
	if err == redis.Nil {
		return nil, nil
	}
	var retGroups []userGroup
	groups := strings.Split(csvGroups, ",")
	for _, group := range groups {
		if group != "" {
			intGroup, err := strconv.Atoi(group)
			if err != nil {
				log.Printf("Error converting user group: %v", err)
				return nil, ErrAtoiConv
			}
			retGroups = append(retGroups, userGroup(intGroup))
		}
	}
	return retGroups, nil
}

func getUsersInGroup(group userGroup) ([]int, error) {
	if redisClient == nil {
		return nil, ErrNilPointer
	}
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

func isUserInGroup(userID int, group userGroup) (bool, error) {
	if redisClient == nil {
		return false, ErrNilPointer
	}
	is, err := redisClient.SIsMember("ug"+strconv.Itoa(int(group)), strconv.Itoa(userID)).Result()
	if err != nil {
		log.Printf("Error checking if user is in group: %v", err)
		return false, ErrRedisCheckSet
	}
	return is, nil
}

func convertUserGroups(groups []userGroup) []string {
	var stringGroups []string
	for _, group := range groups {
		switch group {
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

func getGroupName(group userGroup) (string, error) {
	switch group {
	case ugSoprano:
		return "Soprano", nil
	case ugContralto:
		return "Contralto", nil
	case ugTenore:
		return "Tenore", nil
	case ugBasso:
		return "Basso", nil
	case ugCommissario:
		return "Commissario", nil
	case ugReferente:
		return "Referente", nil
	case ugPreparatore:
		return "Preparatore", nil
	default:
		return "", ErrGroupInvalid
	}
}

func getUserDescription(u *tb.User) (string, error) {
	userGroups, err := getUserGroups(u.ID)
	if err != nil {
		log.Printf("Error retriving user groups: %v", err)
		return "", ErrRedisRetrieveHash
	}
	stringGroups := convertUserGroups(userGroups)

	isAdmin, err := isBotAdmin(u.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
		return "", ErrRedisCheckSet
	}
	isAuth, err := isAuthrizedUser(u.ID)
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
		return "", ErrRedisCheckSet
	}

	msg := "\xF0\x9F\x91\xA4 *INFORMAZIONI UTENTE*" +
		"\n- *Nome*: " + u.FirstName +
		"\n- *Username*: " + u.Username +
		"\n- *ID*: " + strconv.Itoa(u.ID)

	if len(stringGroups) > 0 {
		msg += "\n- *Gruppi*: "
		for i, group := range stringGroups {
			msg += group
			if i <= len(stringGroups)-2 {
				msg += ", "
			}
		}
	}

	msg += "\n- *Tipo utente*: "

	if isAdmin {
		msg += "Admin"
	} else if isAuth {
		msg += "Autorizzato"
	} else {
		msg += "Utente semplice"
	}

	return msg, nil
}
