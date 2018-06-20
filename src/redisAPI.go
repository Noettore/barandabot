package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/dixonwille/wmenu"
	"github.com/go-redis/redis"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	tkSet   = "botTokens"
	botHash = "botInfos"
)

var redisClient *redis.Client

var (
	//ErrRedisConnection is thrown when a redis connection error occurs
	ErrRedisConnection = errors.New("redis: couldn't connect to remote instance")
	//ErrRedisAddSet is thrown when it's not possible to add a key in a set
	ErrRedisAddSet = errors.New("redis: couldn't add key in set")
	//ErrRedisRetriveSet is thrown when it's not possible to retrive keys from a set
	ErrRedisRetriveSet = errors.New("redis: couldn't retrive keys from set")
	//ErrRedisAddHash is thrown when it's not possible to add a key in a hash
	ErrRedisAddHash = errors.New("redis: couldn't add key in hash")
	//ErrTokenParsing is thrown when it's not possible to parse the bot token
	ErrTokenParsing = errors.New("botToken: cannot parse token")
	//ErrTokenInvalid is thrown when the string parsed isn't a valid telegram bot token
	ErrTokenInvalid = errors.New("botToken: string isn't a valid telegram bot token")
	//ErrAddToken is thrown when one or more bot token hasn't been added
	ErrAddToken = errors.New("couldn't add one or more tokens")
	//ErrRemoveToken is thrown when one or more bot tokens hasn't been removed
	ErrRemoveToken = errors.New("couldn't remove one or more tokens")
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
		for _, opt := range opts {
			if opt.Value == nil {
				log.Println("Couldn't remove bot: nil token")
				return ErrNilPointer
			}
			err := removeBotToken(opt.Value.(string))
			if err != nil {
				log.Printf("Couldn't remove bot: %v", err)
			}
		}
		return nil
	})
	//for _, token := range tokens {
	for token, botInfo := range botsInfo {
		menu.Option(botInfo, token, false, nil)
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

func addBotInfo(bot *tb.Bot, botToken string) error {
	if redisClient == nil {
		return ErrNilPointer
	}
	err := redisClient.HSet(botHash, botToken, bot.Me.Username).Err()
	if err != nil {
		log.Printf("Error in adding bot info: %v", err)
		return ErrRedisAddHash
	}

	return nil
}
