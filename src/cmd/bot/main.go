package main

import (
	"os"
	"os/signal"

	bot "axiom/src"
)

func main() {
	discordBot := bot.NewBot()

	go discordBot.Start()
	defer discordBot.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
