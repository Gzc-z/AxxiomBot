package main

import (
	"os"
	"os/signal"

	bot "axiom/src"
)

func main() {
	discordBot := bot.NewBot()

	discordBot.Start()
	defer discordBot.Close()

	make := make(chan os.Signal, 1)
	signal.Notify(make, os.Interrupt)
	<-make
}
