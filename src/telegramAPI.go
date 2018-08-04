package main

import (
	"errors"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type botBool struct {
	isStarted bool
	hasAdmin  bool
}

const (
	startMsg          string = "Salve, sono Stefano, il Magister! Come posso esservi d'aiuto?"
	alreadyStartedMsg string = "Si, mi dica, che c'è?! Sono qui!"
	restartMsg        string = "Eccomi, sono tornato! Ha bisogno? Mi dica pure!"
	stopMsg           string = "Mi assenterò per qualche istante, d'altra parte anch'io ho pur diritto alla mia vita privata. Masino mi attende!"
	newAdminMsg       string = "Beh allora, vediamo... Ah si, la nomino amministratore! Da grandi poteri derivano grandi responsabilità. Mi raccomando, non me ne faccia pentire!"
	delAdminMsg       string = "Ecco, che le avevo detto?! Mi sembrava di essere stato chiaro! Dovrò sollevarla dall'incarico... Mi spiace molto ma da ora in avanti non sarà più amministratore"
)

var genericCommands = map[string]bool{
	"/start":          true,
	"/stop":           true,
	"/menu":           true,
	"/prossimoEvento": true,
}
var authCommands = map[string]bool{
	"/prossimaProvaSezione": true,
	"/prossimaProvaInsieme": true,
}
var adminCommands = map[string]bool{
	"/autorizzaUtente": true,
	"/aggiungiAdmin":   true,
	"/rimuoviAdmin":    true,
}

var bot *tb.Bot
var botStatus botBool
var (
	//ErrNilPointer is thrown when a pointer is nil
	ErrNilPointer = errors.New("pointer is nil")
	//ErrIDFromMsg is thrown when the message doesn't contain user infos
	ErrIDFromMsg = errors.New("telegram: couldn't retrieve user ID from message")
	//ErrSendMsg is thrown when the message couldn't be send
	ErrSendMsg = errors.New("telegram: cannot send message")
	//ErrChatRetrieve is thrown when the chat cannot be retrieved
	ErrChatRetrieve = errors.New("telegram: cannot retrieve chat")
	//ErrTokenMissing is thrown when neither a token is in the db nor one is passed with -t on program start
	ErrTokenMissing = errors.New("telegram: cannot start bot without a token")
	//ErrBotInit is thrown when a bot couldn't be initialized
	ErrBotInit = errors.New("telegram: error in bot initialization")
	//ErrBotConn is thrown when there is a connection problem
	ErrBotConn = errors.New("telegram: cannot connect to bot")
)

func botInit() error {
	token, err := getBotToken()
	if err != nil {
		log.Printf("Error in retriving bot token: %v. Cannot start telebot without token.", err)
		return ErrTokenMissing
	}

	poller := &tb.LongPoller{Timeout: 15 * time.Second}
	middlePoller := tb.NewMiddlewarePoller(poller, setBotPoller)

	bot, err = tb.NewBot(tb.Settings{
		Token:  token,
		Poller: middlePoller,
	})
	if err != nil {
		log.Printf("Error in enstablishing connection for bot %s: %v", bot.Me.Username, err)
		return ErrBotConn
	}

	err = setBotHandlers()
	if err != nil {
		log.Printf("Error setting bot handlers: %v", err)
		return ErrBotInit
	}

	err = setBotMenus()
	if err != nil {
		log.Printf("Error setting bot menus: %v", err)
		return ErrBotInit
	}

	err = setBotCallbacks()
	if err != nil {
		log.Printf("Error setting bot callbacks: %v", err)
		return ErrBotInit
	}

	err = addBotInfo(token, bot.Me.Username)
	if err != nil {
		log.Printf("Error: bot %s info couldn't be added: %v", bot.Me.Username, err)
		return ErrBotInit
	}

	hasAdmin, err := hasBotAdmins()
	if err != nil {
		log.Printf("Error checking if bot has admins: %v", err)
		return ErrBotInit
	}
	botStatus.hasAdmin = hasAdmin

	return nil
}

func sendMsg(user *tb.User, msg string) error {
	_, err := bot.Send(user, msg)
	if err != nil {
		log.Printf("Error sending message to user: %v", err)
		return ErrSendMsg
	}
	return nil
}

func sendMsgWithMenu(user *tb.User, msg string) error {
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
	_, err = bot.Send(user, msg, &tb.ReplyMarkup{
		InlineKeyboard: menu,
	})
	if err != nil {
		log.Printf("Error sending message to user: %v", err)
		return ErrSendMsg
	}
	return nil
}

func sendMsgWithSpecificMenu(user *tb.User, msg string, menu [][]tb.InlineButton) error {
	_, err := bot.Send(user, msg, &tb.ReplyMarkup{
		InlineKeyboard: menu,
	})
	if err != nil {
		log.Printf("Error sending message to user: %v", err)
		return ErrSendMsg
	}
	return nil
}

func setBotPoller(upd *tb.Update) bool {
	if upd.Message == nil {
		return true
	}
	if upd.Message.Sender != nil {
		err := addUser(upd.Message.Sender)
		if err != nil {
			log.Printf("Error in adding user info: %v", err)
		}
	} else {
		log.Printf("%v", ErrIDFromMsg)
	}
	_, isGenericCmd := genericCommands[upd.Message.Text]
	_, isAuthCmd := authCommands[upd.Message.Text]
	_, isAdminCmd := adminCommands[upd.Message.Text]

	started, err := isStartedUser(upd.Message.Sender.ID)
	if err != nil {
		log.Printf("Error checking if user is started: %v", err)
	}
	if !started && upd.Message.Text != "/start" {
		return false
	}

	auth, err := isAuthrizedUser(upd.Message.Sender.ID)
	if err != nil {
		log.Printf("Error checking if user is authorized: %v", err)
	}
	admin, err := isBotAdmin(upd.Message.Sender.ID)
	if err != nil {
		log.Printf("Error checking if user is admin: %v", err)
	}
	if isAdminCmd && !admin {
		return false
	}
	if isAuthCmd && !auth {
		return false
	}
	if !isGenericCmd {
		return false
	}
	return true
}

func botStart() error {
	if bot == nil {
		return ErrNilPointer
	}

	go bot.Start()
	botStatus.isStarted = true
	log.Printf("Started %s", bot.Me.Username)

	return nil
}

func botStop() error {
	if bot == nil {
		return ErrNilPointer
	}
	log.Printf("Stopping %s", bot.Me.Username)
	bot.Stop()
	botStatus.isStarted = false
	log.Println("Bot stopped")

	return nil
}
