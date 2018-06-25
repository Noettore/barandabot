package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/dixonwille/wmenu"
)

type stringSlice []string

type flags struct {
	interactive bool
	redisAddr   string
	redisPwd    string
	redisDB     int
	tokens      stringSlice
}

var cmdFlags flags

var (
	welcomeMessage = "Welcome in barandaBot! Here you can control the bot(s) options and configurations."
	//ErrStdRead is thrown when it's not possible to read from the standard input
	ErrStdRead = errors.New("stdin: couldn't read string from stdin")
	//ErrMainMenu is thrown when a menu couldn't be started
	ErrMainMenu = errors.New("menu: couldn't start menu")
)

func (i *stringSlice) String() string {
	return fmt.Sprint(*i)
}

func (i *stringSlice) Set(values string) error {
	splittedValues := strings.Split(values, ",")
	for _, value := range splittedValues {
		*i = append(*i, value)
	}
	return nil
}

func getFlags() error {
	const (
		defaultInteractive = true
		interactiveUsage   = "False if the bot isn't executed on a tty"
		defaultAddr        = "127.0.0.1:6379"
		addrUsage          = "The address of the redis instance"
		defaultPwd         = ""
		pwdUsage           = "The password of the redis instance"
		defaultDB          = 0
		dbUsage            = "The database to be selected after connecting to redis instance"
		tokenUsage         = "A bot token to be added to the set of tokens"
	)

	flag.BoolVar(&(cmdFlags.interactive), "interactive", defaultInteractive, interactiveUsage)
	flag.BoolVar(&(cmdFlags.interactive), "i", defaultInteractive, interactiveUsage+"(shorthand)")
	flag.StringVar(&(cmdFlags.redisAddr), "redisAddr", defaultAddr, addrUsage)
	flag.StringVar(&(cmdFlags.redisAddr), "a", defaultAddr, addrUsage+"(shorthand)")
	flag.StringVar(&(cmdFlags.redisPwd), "redisPwd", defaultPwd, pwdUsage)
	flag.StringVar(&(cmdFlags.redisPwd), "p", defaultPwd, pwdUsage+"(shorthand)")
	flag.IntVar(&(cmdFlags.redisDB), "redisDB", defaultDB, dbUsage)
	flag.IntVar(&(cmdFlags.redisDB), "d", defaultDB, dbUsage+"(shorthand)")
	flag.Var(&(cmdFlags.tokens), "token", tokenUsage)
	flag.Var(&(cmdFlags.tokens), "t", tokenUsage+"(shorthand")

	flag.Parse()

	return nil
}

func mainMenu() error {
	fmt.Println(welcomeMessage)
	menu := wmenu.NewMenu("What do you want to do?")
	menu.LoopOnInvalid()
	menu.Option("Start Bot(s)", nil, true, func(opt wmenu.Opt) error {
		return botsStart()
	})
	menu.Option("Add bot token(s)", nil, false, func(opt wmenu.Opt) error {
		return addBotTokens(nil)
	})
	menu.Option("Remove bot token(s)", nil, false, func(opt wmenu.Opt) error {
		return removeBotTokens()
	})

	var returnErr error

	for {
		err := menu.Run()
		if err != nil {
			log.Printf("Error in main menu: %v", err)
			returnErr = ErrMainMenu
		}
	}
	return returnErr
}
