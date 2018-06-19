package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/dixonwille/wmenu"
)

type stringSlice []string

type flags struct {
	redisAddr string
	redisPwd  string
	redisDB   int
	tokens    stringSlice
}

var (
	//ErrStdRead it thrown when it's not possible to read from the standard input
	ErrStdRead = errors.New("couldn't read string from stdin")
)

func (i *stringSlice) Set(value string) error {

	*i = append(*i, value)
	return nil
}

func getFlags() (flags, error) {

	var cmdFlags flags

	const (
		defaultAddr = "127.0.0.1:6379"
		addrUsage   = "The address of the redis instance"
		defaultPwd  = ""
		pwdUsage    = "The password of the redis instance"
		defaultDB   = 0
		dbUsage     = "The database to be selected after connecting to redis instance"
		tokenUsage  = "A bot token to be added to the set of tokens"
	)

	flag.StringVar(&(cmdFlags.redisAddr), "redisAddr", defaultAddr, addrUsage)
	flag.StringVar(&(cmdFlags.redisAddr), "a", defaultAddr, addrUsage+"(shorthand)")
	flag.StringVar(&(cmdFlags.redisPwd), "redisPwd", defaultPwd, pwdUsage)
	flag.StringVar(&(cmdFlags.redisPwd), "p", defaultPwd, pwdUsage+"(shorthand)")
	flag.IntVar(&(cmdFlags.redisDB), "redisDB", defaultDB, dbUsage)
	flag.IntVar(&(cmdFlags.redisDB), "d", defaultDB, dbUsage+"(shorthand)")
	flag.Var(&(cmdFlags.tokens), "token", tokenUsage)
	flag.Var(&(cmdFlags.tokens), "t", tokenUsage+"(shorthand")

	flag.Parse()

	return cmdFlags, nil
}

func mainMenu() {
	fmt.Println("Welcome in barandaBot! Here you can control the bot(s) options and configurations.")
	menu := wmenu.NewMenu("What do you want to do?")
	menu.LoopOnInvalid()
	menu.Option("Start Bot(s)", nil, true, nil)
	menu.Option("Add bot token(s)", nil, false, func(opt wmenu.Opt) error {
		return addBotTokens(redisClient, nil)
	})
	menu.Option("Remove bot token(s)", nil, false, func(opt wmenu.Opt) error {
		return removeBotTokens(redisClient)
	})

	err := menu.Run()
	if err != nil {
		log.Printf("Error in main menu: %v", err)
	}
}
