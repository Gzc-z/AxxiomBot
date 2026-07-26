package main

import (
	"log"
	"os"
	"os/signal"

	bot "axiom/src"
)

func main() {
	discordBot := bot.NewBot()
	defer discordBot.Close()

	go discordBot.Start()

	ds := discordBot.Session
	log.Printf("Logged in as: %v#%v", ds.State.User.Username, ds.State.User.Discriminator)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
