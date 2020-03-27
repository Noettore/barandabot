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
	stopMsg           string = "Mi assenterò per qualche istante, d'altra parte anch'io ho pur diritto alla mia vita privata. Masino mi attende \xF0\x9F\x90\xB1"
	unstoppableMsg    string = "Non ci siamo... Io l'ho nominata AMMINISTRATORE, cosa crede?! Questo ruolo esige impegno! Non può certo bloccarmi!"
	wrongCmdMsg       string = "Non capisco, si spieghi meglio! Per cortesia, basta basta! La prego! Non so di cosa sta parlando!"
	authHowToMsg      string = "Per autorizzare un utente invia un messaggio con scritto \n`/authUser ID_UTENTE`\n sostituendo `ID_UTENTE` con l'ID che ti é stato comunicato dall'utente da autorizzare"
	deAuthHowToMsg    string = "Per deautorizzare un utente invia un messaggio con scritto \n`/authUser USERNAME`\n sostituendo `USERNAME` con il nome utente da deautorizzare"
	newAuthMsg        string = "Benvenuto! Da ora in poi lei fa ufficialmente parte del magnifico *Coro dell'Università di Pisa*! Deve sentirsi onorato."
	delAuthMsg        string = "Capisco, quindi se ne sta andando... Beh un po' mi dispiace, devo ammetterlo. Se ripassa da queste parti sarà sempre il benvenuto! Arrivederci."
	newAdminMsg       string = "Beh allora, vediamo... Ah si, la nomino amministratore! Da grandi poteri derivano grandi responsabilità. Mi raccomando, non me ne faccia pentire!"
	delAdminMsg       string = "Ecco, che le avevo detto?! Mi sembrava di essere stato chiaro! Dovrò sollevarla dall'incarico... Mi spiace molto ma da ora in avanti non sarà più amministratore"
	menuMsg           string = "Ecco a lei, questo è l'elenco di tutto ciò che può chiedermi. Non mi disturbi con altre richieste!"
	contactMsg        string = "*BarandaBot*\xE2\x84\xA2" +
		"\nSe hai domande, suggerimenti o se vuoi segnalare bug e altri malfunzionamenti puoi contattare l'Altissimo con i seguenti mezzi di comunicazione:" +
		"\n- \xF0\x9F\x90\xA6 _Piccione viaggiatore_: [Palazzo Ricci, Pisa](https://goo.gl/maps/gMUbV2eqJiL2)" +
		"\n- \xF0\x9F\x93\xA7 _Mail_: telebot.corounipi@gmail.com" +
		"\n- \xF0\x9F\x93\x82 _GitHub_: https://github.com/Noettore/barandaBot"
)

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
	//ErrSetLastMsg is thrown when it's not possible to set last message per user in hash
	ErrSetLastMsg = errors.New("cannot set last message per user")
	//ErrInvalidPath is thrown when a path isn't valid
	ErrInvalidPath = errors.New("path is not valid")
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

func setBotPoller(upd *tb.Update) bool {
	if upd.Message == nil {
		return true
	}
	/* if upd.Message.Sender != nil {
		err := addUser(upd.Message.Sender)
		if err != nil {
			log.Printf("Error in adding user info: %v", err)
		}
	} else {
		log.Printf("%v", ErrIDFromMsg)
	} */
	isUser, err := isUser(upd.Message.Sender.ID)
	if err != nil {
		log.Printf("Error checking if message come from a bot user: %v", err)
	}
	if !isUser && upd.Message.Text != "/start" {
		return false
	}

	//_, isGenericCmd := genericCommands[upd.Message.Text]
	_, isAuthCmd := authCommands[upd.Message.Text]
	_, isAdminCmd := adminCommands[upd.Message.Text]

	started, err := isStartedUser(upd.Message.Sender.ID)
	if err != nil {
		log.Printf("Error checking if user is started: %v", err)
	}
	if !started && upd.Message.Text != "/start" {
		sendMsgWithSpecificMenu(upd.Message.Sender, "ZzZzZzZzZzZ", startMenu, true)
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
	if isAdminCmd && admin == false {
		return false
	}
	if isAuthCmd && auth == false {
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
