package bot

import (
	"fmt"
	"os"
	"os/signal"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestPing(t *testing.T) {
	ping := NewBot()
	ping.Start()
	defer ping.Close()

	fmt.Println("created")
	s := ping.Session
	s.Identify.Intents |= discordgo.IntentGuilds
	s.Identify.Intents |= discordgo.IntentGuildMembers
	s.Identify.Intents |= discordgo.IntentGuildMessages
	s.Identify.Intents |= discordgo.IntentDirectMessages

	s.ChannelMessageSend("1531108984751128736", "ping")

	make := make(chan os.Signal, 1)
	signal.Notify(make, os.Interrupt)
	<-make
}
