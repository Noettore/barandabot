package main

import "log"

func main() {

	cmdFlags, err := getFlags()
	if err != nil {
		log.Fatalln("Error in parsing command line flags. Abort!")
	}

	redisClient, err := redisInit(cmdFlags.redisAddr, cmdFlags.redisPwd, cmdFlags.redisDB)
	defer redisClient.Close()
	if err != nil {
		log.Fatalf("Error in initializing redis instance: %v", err)
	}

	startMenu()

	bots, err := botInit(redisClient)
	if err != nil {
		log.Fatalf("Error in initializing bots: %v", err)
	}

	for _, bot := range bots {
		defer bot.Stop()
	}

	/*b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()*/
}
