package main

import "log"

func main() {

	cmdFlags, err := getFlags()
	if err != nil {
		log.Fatal("Error in parsing command line flags. Abort!")
	}

	redisClient, err := redisInit(cmdFlags.redisAddr, cmdFlags.redisPwd, cmdFlags.redisDB)
	defer redisClient.Close()

	bots, errors := botInit()
	for _, bot := range bots {
		bot.Stop()
	}

	/*b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()*/
}
