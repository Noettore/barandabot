package main

import "log"

func main() {

	cmdFlags, err := getFlags()
	if err != nil {
		log.Fatal("Error in parsing command line flags. Abort!")
	}

	redisClient, err := redisInit(cmdFlags.redisAddr, cmdFlags.redisPwd, cmdFlags.redisDB)
	defer redisClient.Close()

	if err != nil {
		log.Panicf("Error in initializing redis instance: %v", err)
	}

	bots, errors := botInit(redisClient)
	for i, err := range errors {
		if err != nil {
			log.Printf("Error in initializing bot: %v", err)
			bots = append(bots[:i], bots[i+1:]...)
		}
	}

	for _, bot := range bots {
		defer bot.Stop()
	}

	/*b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()*/
}
