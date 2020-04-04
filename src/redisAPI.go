package main

import (
	"errors"
	"log"

	"github.com/go-redis/redis"
)

const (
	botToken       = "botToken"
	botInfo        = "botInfo"
	usersID        = "usersID"
	usersInfo      = "usersInfo"
	usersGroups    = "usersGroups"
	startedUsers   = "startedUsers"
	authUsers      = "authUsers"
	adminUsers     = "adminUsers"
	lastMsgPerUser = "lastMsgPerUser"
	mediaPath      = "mediaPath"
)

var redisClient *redis.Client

var (
	//ErrRedisConnection is thrown when a redis connection error occurs
	ErrRedisConnection = errors.New("redis: couldn't connect to remote instance")
	//ErrRedisAddSet is thrown when it's not possible to add a key in a set
	ErrRedisAddSet = errors.New("redis: couldn't add key in set")
	//ErrRedisRemSet is thrown when it's not possible to remove a key from a given set
	ErrRedisRemSet = errors.New("redis: couldn't remove key from set")
	//ErrRedisRetrieveSet is thrown when it's not possible to retrieve keys from a set
	ErrRedisRetrieveSet = errors.New("redis: couldn't retrieve keys from set")
	//ErrRedisCheckSet is thrown when it's not possible to check if a key is in a given set
	ErrRedisCheckSet = errors.New("redis: couldn't check if key is in set")
	//ErrRedisAddHash is thrown when it's not possible to add a key in a hash
	ErrRedisAddHash = errors.New("redis: couldn't add key in hash")
	//ErrRedisDelHash is thrown when it's not possible to remove a key from a hash
	ErrRedisDelHash = errors.New("redis: couldn't remove key from hash")
	//ErrRedisAddString is thrown when it's not possible to add a string
	ErrRedisAddString = errors.New("redis: couldn't add string")
	//ErrRedisDelString is thrown when it's not possible to remove a string
	ErrRedisDelString = errors.New("redis: couldn't remove string")
	//ErrRedisRetrieveHash is thrown when it's not possible to retrieve a key from a hash
	ErrRedisRetrieveHash = errors.New("redis: couldn't retrieve key from hash")
	//ErrTokenParsing is thrown when it's not possible to parse the bot token
	ErrTokenParsing = errors.New("botToken: cannot parse token")
	//ErrTokenInvalid is thrown when the string parsed isn't a valid telegram bot token
	ErrTokenInvalid = errors.New("botToken: string isn't a valid telegram bot token")
	//ErrIDParsing is thrown when it's not possible to parse the user ID
	ErrIDParsing = errors.New("userID: cannot parse ID")
	//ErrIDInvalid is thrown when the string parsed isn't a valid telegram user ID
	ErrIDInvalid = errors.New("userID: string isn't a valid telegram user ID")
	//ErrAddToken is thrown when one or more bot token hasn't been added
	ErrAddToken = errors.New("couldn't add one or more tokens")
	//ErrAddUser is thrown when one or more user hasn't been added
	ErrAddUser = errors.New("couldn't add one or more users")
	//ErrAddAdmin is thrown when one or more admin IDs hasn't been added
	ErrAddAdmin = errors.New("couldn't add one or more admins")
	//ErrAddAuthUser is thrown when one or more users cannot be authorized
	ErrAddAuthUser = errors.New("couldn't authorize one or more users")
	//ErrGetUser is thrown when user info couldn't be retrieven
	ErrGetUser = errors.New("couldn't retrieve user info")
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
