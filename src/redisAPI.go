package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dixonwille/wmenu"
	"github.com/go-redis/redis"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	tkSet       = "botToken"
	botHash     = "botInfo"
	userSet     = "userID"
	userHash    = "userInfo"
	authUserSet = "authUser"
)

var redisClient *redis.Client

var (
	//ErrRedisConnection is thrown when a redis connection error occurs
	ErrRedisConnection = errors.New("redis: couldn't connect to remote instance")
	//ErrRedisAddSet is thrown when it's not possible to add a key in a set
	ErrRedisAddSet = errors.New("redis: couldn't add key in set")
	//ErrRedisRemSet is thrown when it's not possible to remove a key from a given set
	ErrRedisRemSet = errors.New("redis: couldn't remove key from set")
	//ErrRedisRetriveSet is thrown when it's not possible to retrive keys from a set
	ErrRedisRetriveSet = errors.New("redis: couldn't retrive keys from set")
	//ErrRedisCheckSet is thrown when it's not possible to check if a key is in a given set
	ErrRedisCheckSet = errors.New("redis: couldn't check if key is in set")
	//ErrRedisAddHash is thrown when it's not possible to add a key in a hash
	ErrRedisAddHash = errors.New("redis: couldn't add key in hash")
	//ErrRedisDelHash is thrown when it's not possible to remove a key from a hash
	ErrRedisDelHash = errors.New("redis: couldn't remove key from hash")
	//ErrTokenParsing is thrown when it's not possible to parse the bot token
	ErrTokenParsing = errors.New("botToken: cannot parse token")
	//ErrTokenInvalid is thrown when the string parsed isn't a valid telegram bot token
	ErrTokenInvalid = errors.New("botToken: string isn't a valid telegram bot token")
	//ErrAddToken is thrown when one or more bot token hasn't been added
	ErrAddToken = errors.New("couldn't add one or more tokens")
	//ErrRemoveToken is thrown when one or more bot tokens hasn't been removed
	ErrRemoveToken = errors.New("couldn't remove one or more tokens")
	//ErrJSONMarshall is thrown when it's impossible to marshall a given struct
	ErrJSONMarshall = errors.New("json: couldn't marshall struct")
)

func redisInit(addr string, pwd string, db int) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	err := redisClient.Ping().Err()
	if err != nil {
		log.Printf("Error in connecting to redis instance: %v", err)
		return ErrRedisConnection
	}
	return nil
}

func addBotToken(newToken string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	token := strings.TrimSpace(newToken)
	matched, err := regexp.MatchString("^\\d+\\:([0-9]|[A-z]|\\_|\\-)+", token)
	if err != nil {
		log.Printf("Error in parsing bot token: %v", err)
		return ErrTokenParsing
	}
	if !matched {
		return ErrTokenInvalid
	}

	err = redisClient.SAdd(tkSet, token).Err()
	if err != nil {
		log.Printf("Error in adding new bot token: %v", err)
		return ErrRedisAddSet
	}

	return nil
}

func addBotTokens(newTokens []string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	errNum := 0
	if newTokens == nil && cmdFlags.interactive {
		fmt.Println("Add the new tokens, comma-separated:")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error in reading new bot tokens: %v", err)
			return ErrStdRead
		}
		newTokens = strings.Split(line, ",")
	}
	for _, newToken := range newTokens {
		err := addBotToken(newToken)
		if err != nil {
			errNum++
			log.Printf("Error in adding new bot token %s: %v", newToken, err)
		}
	}
	if errNum == len(newTokens) {
		return ErrAddToken
	}
	return nil
}

func removeBotToken(token string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.SRem(tkSet, token).Err()
	if err != nil {
		log.Printf("Error in removing bot token %s: %v", token, err)
		return ErrRemoveToken
	}
	return nil
}

func removeBotTokens() error {
	if redisClient == nil {
		return ErrNilPointer
	}
	//tokens, err := redisClient.SMembers(tkSet).Result()
	botsInfo, err := redisClient.HGetAll(botHash).Result()
	if err != nil {
		log.Printf("Couldn't retrive bot info: %v", err)
		return ErrRedisRetriveSet
	}
	menu := wmenu.NewMenu("Select the token(s) you want to remove:")
	menu.AllowMultiple()
	menu.LoopOnInvalid()
	menu.Action(func(opts []wmenu.Opt) error {
		var returnErr error
		for _, opt := range opts {
			if opt.Value == nil {
				log.Println("Couldn't remove bot: nil token")
				returnErr = ErrNilPointer
			} else {
				err := removeBotToken(opt.Value.(string))
				if err != nil {
					log.Printf("Couldn't remove bot token: %v", err)
				}
				err = removeBotInfo(opt.Value.(string))
				if err != nil {
					log.Printf("Couldn't remove bot info: %v", err)
				}
			}
		}
		return returnErr
	})
	//for _, token := range tokens {
	for token, jsonBotInfo := range botsInfo {
		botInfo := &tb.Bot{}
		json.Unmarshal([]byte(jsonBotInfo), &botInfo)
		menu.Option(botInfo.Me.Username, token, false, nil)
	}
	err = menu.Run()
	if err != nil {
		log.Printf("Error in removeToken menu: %v", err)
		return ErrRemoveToken
	}
	return nil
}

func getBotTokens() ([]string, error) {
	if redisClient == nil {
		return nil, ErrNilPointer
	}
	tkNum, err := redisClient.SCard(tkSet).Result()
	if err != nil {
		log.Printf("Couldn't retrive number of bot tokens: %v", err)
		return nil, ErrRedisRetriveSet
	}
	if tkNum == 0 {
		fmt.Println("No bot token found.")
		err := addBotTokens(nil)
		if err != nil {
			log.Printf("Couldn't add new bot tokens: %v", err)
			return nil, ErrAddToken
		}
	}

	tokens, err := redisClient.SMembers(tkSet).Result()
	if err != nil {
		log.Printf("Couldn't retrive bot tokens: %v", err)
		return nil, ErrRedisRetriveSet
	}

	return tokens, nil
}

func addBotInfo(botToken string, bot *tb.Bot) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	jsonBot, err := json.Marshal(&bot)
	if err != nil {
		log.Printf("Error marshalling bot info: %v", err)
		return ErrJSONMarshall
	}
	err = redisClient.HSet(botHash, botToken, string(jsonBot)).Err()
	if err != nil {
		log.Printf("Error in adding bot info: %v", err)
		return ErrRedisAddHash
	}

	return nil
}

func removeBotInfo(botToken string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.HDel(botHash, botToken).Err()
	if err != nil {
		log.Printf("Error in removing bot info: %v", err)
		return ErrRedisDelHash
	}
	return nil
}

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
