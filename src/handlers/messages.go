package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"axiom/src/config"
	"axiom/src/interactions"
	"axiom/src/ui"

	"github.com/bwmarrin/discordgo"
	ping "github.com/prometheus-community/pro-bing"
)

var prefix = "."

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || len(m.Content) == 0 || m.Content[0] != '.' {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, "."))

	switch args[0] {
	case "time":
		time := time.Now().Format("15h : 04m : 05s")
		s.ChannelMessageSendReply(m.ChannelID, time, m.Reference())
	case "ping":
		s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
	case "comandos", "commands":
		cmds := interactions.Commands
		s.ChannelMessageSendReply(m.ChannelID, "## commands available: ", m.Reference())
		for _, v := range cmds {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("Nome: %s\nDescrição: %s\n%v",
					v.Name,
					v.Description,
					v.Version,
				),
			})
		}
	case "users", "user", "usrs":
		members, _ := s.GuildMembers(m.GuildID, "", 100)
		s.ChannelMessageSendReply(m.ChannelID, "## users:\n`3 first members`", m.Reference())
		for _, v := range members[0:3] {
			response := ui.MembersResponse(v)
			s.ChannelMessageSendComplex(m.ChannelID, response)
		}
	case "fact", "cat", "catfact":
		catFacts := catFacts(s, m)
		s.ChannelMessageSendReply(m.ChannelID, catFacts.Fact, m.Reference())
	case "rand", "random", "r":
		r := rand.Intn(100)
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(r), m.Reference())
	case "me":
		usr, _ := s.User(m.Author.ID)
		s.ChannelMessageSendReply(m.ChannelID, "## you:", m.Reference())
		response := ui.UserResponse(usr)
		s.ChannelMessageSendComplex(m.ChannelID, response)
	case "mine":
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "10 pings",
		})

		serverPing(s, m)
	}
}

const catFactAPIURL = "https://catfact.ninja"

func catFacts(s *discordgo.Session, m *discordgo.MessageCreate) struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
} {
	r, err := http.Get(catFactAPIURL + "/fact")
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	var fact struct {
		Fact   string `json:"fact"`
		Length int    `json:"length"`
	}
	err = json.NewDecoder(r.Body).Decode(&fact)
	if err != nil {
		log.Println(err)
	}
	return fact
}

func serverPing(s *discordgo.Session, m *discordgo.MessageCreate) string {
	host := config.GetMinecraft()
	pinger, err := ping.NewPinger(host.IP)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSendReply(m.ChannelID, "error: couldn't connect to server\ncheck server's ip", m.Reference())
		return ""
	}

	pinger.Count = 10 // -> (count of pings)
	pinger.Interval = 1 * time.Second
	// pinger.Timeout = 20 * time.Second

	// for now, SendReply and MessageEdit will stay here
	p, _ := s.ChannelMessageSendReply(m.ChannelID, "___", m.Reference())
	pinger.OnRecv = func(pkt *ping.Packet) {
		stats := fmt.Sprintf(
			"Resposta de %s: tempo=%v seq=%d\n",
			pkt.Addr,
			pkt.Rtt,
			pkt.Seq+1,
		)
		_, err := s.ChannelMessageEdit(m.ChannelID, p.ID, string(stats))
		if err != nil {
			log.Println(err)
		}
	}

	err = pinger.Run()
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSendReply(m.ChannelID, "error: couldn't run server", m.Reference())
		return ""
	}
	return "ok"
}
