package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dixonwille/wmenu"
)

type stringSlice []string

type flags struct {
	interactive bool
	redisAddr   string
	redisPwd    string
	redisDB     int
	token       string
	mediaPath   string
}

var cmdFlags flags

var (
	welcomeMessage = "Welcome! Here you can control the bot options and configurations."
	//ErrStdRead is thrown when it's not possible to read from the standard input
	ErrStdRead = errors.New("stdin: couldn't read string from stdin")
	//ErrMainMenu is thrown when a menu couldn't be started
	ErrMainMenu = errors.New("menu: couldn't start menu")
	//ErrAtoiConv is thrown when a string couldn't be converted to int
	ErrAtoiConv = errors.New("atoi: couldn't convert string to int")
	//ErrJSONMarshall is thrown when it's impossible to marshall a given struct
	ErrJSONMarshall = errors.New("json: couldn't marshall struct")
	//ErrJSONUnmarshall is thworn when it's impossible to unmarshall a given struct
	ErrJSONUnmarshall = errors.New("json: couldn't unmarshall struct")
	//ErrRemoveAdmin is thrown when it's impossible to remove an admin
	ErrRemoveAdmin = errors.New("menu: cannot remove admin")
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
		interactiveUsage   = "true/false if you want interactive/non-interactive usage"
		defaultAddr        = "127.0.0.1:6379"
		addrUsage          = "The address of the redis instance"
		defaultPwd         = ""
		pwdUsage           = "The password of the redis instance"
		defaultDB          = 0
		dbUsage            = "The database to be selected after connecting to redis instance"
		defaultToken       = ""
		tokenUsage         = "A bot token to be added to the set of tokens"
		defaultMediaPath   = ""
		mediaPathUsage     = "A path to be used as media directory"
	)
	/* var defaultMediaPath string
	usr, err := user.Current()
	if err != nil {
		log.Printf("Error retriving current user home dir: %v. Using current directory", err)
		os.MkdirAll("./.barandaBot/img", os.ModePerm)
		defaultMediaPath = "./.barandaBot/img"

	} else {
		os.MkdirAll(usr.HomeDir+"/.barandaBot/img", os.ModePerm)
		defaultMediaPath = usr.HomeDir + "/.barandaBot/img"
	} */

	flag.BoolVar(&(cmdFlags.interactive), "interactive", defaultInteractive, interactiveUsage)
	flag.BoolVar(&(cmdFlags.interactive), "i", defaultInteractive, interactiveUsage+"(shorthand)")
	flag.StringVar(&(cmdFlags.redisAddr), "redisAddr", defaultAddr, addrUsage)
	flag.StringVar(&(cmdFlags.redisAddr), "a", defaultAddr, addrUsage+"(shorthand)")
	flag.StringVar(&(cmdFlags.redisPwd), "redisPwd", defaultPwd, pwdUsage)
	flag.StringVar(&(cmdFlags.redisPwd), "p", defaultPwd, pwdUsage+"(shorthand)")
	flag.IntVar(&(cmdFlags.redisDB), "redisDB", defaultDB, dbUsage)
	flag.IntVar(&(cmdFlags.redisDB), "d", defaultDB, dbUsage+"(shorthand)")
	flag.StringVar(&(cmdFlags.token), "token", defaultToken, tokenUsage)
	flag.StringVar(&(cmdFlags.token), "t", defaultToken, tokenUsage+"(shorthand")
	flag.StringVar(&(cmdFlags.mediaPath), "mediaPath", defaultMediaPath, mediaPathUsage)
	flag.StringVar(&(cmdFlags.mediaPath), "m", defaultMediaPath, mediaPathUsage+"(shorthand")

	flag.Parse()

	return nil
}

func exit() error {
	if botStatus.isStarted {
		log.Printf("Stopping %s", bot.Me.Username)
		bot.Stop()
		log.Println("Bot stopped")
	}
	log.Println("Closing redis instance")
	redisClient.Close()
	log.Println("Redis instance closed")
	log.Println("Exiting")
	os.Exit(0)

	return nil
}

func isValidPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func mainMenu() *wmenu.Menu {
	menu := wmenu.NewMenu("What do you want to do?")
	menu.LoopOnInvalid()
	if !botStatus.isStarted {
		menu.Option("Start bot", nil, true, func(opt wmenu.Opt) error {
			return botStart()
		})
		menu.Option("Set bot token", nil, false, func(opt wmenu.Opt) error {
			return setBotToken("")
		})
		menu.Option("Set media path", nil, false, func(opt wmenu.Opt) error {
			return setMediaDir("")
		})
	}
	if botStatus.isStarted {
		menu.Option("Stop bot", nil, false, func(opt wmenu.Opt) error {
			return botStop()
		})
	}
	menu.Option("Add bot admin(s)", nil, false, func(opt wmenu.Opt) error {
		return addBotAdmins(nil)
	})
	if botStatus.hasAdmin {
		menu.Option("Remove bot admin(s)", nil, false, func(opt wmenu.Opt) error {
			return removeBotAdmins()
		})
	}
	menu.Option("Exit", nil, false, func(opt wmenu.Opt) error {
		return exit()
	})

	return menu
}

func mainMenuLoop() error {
	menu := mainMenu()
	botStarted := botStatus.isStarted
	hasAdmin := botStatus.hasAdmin
	fmt.Println(welcomeMessage)
	for {
		err := menu.Run()
		if err != nil {
			log.Printf("Error in main menu: %v", err)
		}
		if botStarted != botStatus.isStarted || hasAdmin != botStatus.hasAdmin {
			botStarted = botStatus.isStarted
			hasAdmin = botStatus.hasAdmin
			menu = mainMenu()
		}
	}
}
