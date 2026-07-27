package main

import (
	"os"
	"os/signal"

	bot "axiom/src"
)

func main() {
	discordBot := bot.NewBot()

	defer discordBot.Close()
	go discordBot.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
