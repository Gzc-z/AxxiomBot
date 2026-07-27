package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"axiom/src/config"
	"axiom/src/interactions"
	"axiom/src/ui"

	"github.com/bwmarrin/discordgo"
	"github.com/expr-lang/expr"
	ping "github.com/prometheus-community/pro-bing"
)

var prefix = "."

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || len(m.Content) == 0 || m.Content[0] != '.' {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, "."))
	var help []string = []string{
		"time",
		"ping",
		"comandos",
		"users",
		"fact",
		"rand",
		"me",
		"mine",
		"calc",
		"math",
		// "expr",
	}
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, "### error: args == 0 ??", m.Reference())
		return

	}
	switch args[0] {
	case "help":
		s.ChannelMessageSendReply(m.ChannelID, "## comandos de texto disponíveis: ", m.Reference())
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", help))
	case "time":
		time := time.Now().Format("15h : 04m : 05s")
		s.ChannelMessageSendReply(m.ChannelID, time, m.Reference())
	case "ping":
		s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
	case "comandos", "commands":
		s.ChannelMessageSendReply(m.ChannelID, "## commands available: ", m.Reference())
		for _, v := range interactions.Commands {
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
		if len(args) == 2 {
			num, err := strconv.Atoi(args[1])
			if err != nil {
				s.ChannelMessageSendReply(m.ChannelID, "error: invalid number", m.Reference())
				return
			}
			r := rand.Intn(num)
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("numero aleatório: %d", num),
			})
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(r), m.Reference())
			return
		}
		r := rand.Intn(10)
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "numero aleatório: 10",
		})
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(r), m.Reference())
	case "me":
		usr, _ := s.User(m.Author.ID)
		s.ChannelMessageSendReply(m.ChannelID, "## you:", m.Reference())
		response := ui.UserResponse(usr)
		s.ChannelMessageSendComplex(m.ChannelID, response)
	case "mine":
		if len(args) == 2 {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("%v pings", args[1]),
			})
			count, err := strconv.Atoi(args[1])
			if err != nil {
				s.ChannelMessageSendReply(m.ChannelID, "error: invalid number", m.Reference())
				return
			}
			serverPing(s, m, count)
			return
		}
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "infinite ping",
			// Flags:   discordgo.MessageFlagsEphemeral,
		})

		serverPing(s, m)
	case "calc", "math", "c":
		if len(args) == 1 {
			s.ChannelMessageSendReply(m.ChannelID, "use **calc <expressão>**", m.Reference())
			return
		}
		if len(args) >= 2 {
			for _, v := range args[1:] {
				if strings.Contains(v, "x") || strings.Contains(v, "X") {
					s.ChannelMessageSendReply(m.ChannelID, "use * instead of x", m.Reference())
					return
				}
			}
			program, err := expr.Compile(strings.Join(args[1:], " "))
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSendReply(m.ChannelID, "error: invalid compile", m.Reference())
				return
			}

			result, err := expr.Run(program, nil)
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSendReply(m.ChannelID, "error: invalid math expression", m.Reference())
				return
			}
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(result), m.Reference())
		}
		// case "expr":
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

func serverPing(s *discordgo.Session, m *discordgo.MessageCreate, count ...int) string {
	host := config.GetMinecraft()
	pinger, err := ping.NewPinger(host.IP)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSendReply(m.ChannelID, "error: couldn't connect to server\ncheck server's ip", m.Reference())
		return ""
	}
	if len(count) != 0 {
		pinger.Count = count[0] // -> (count of pings)
	} else {
		pinger.Count = 0 // -> (count of pings)
	}
	pinger.Interval = 1 * time.Second // -> (interval of pings)
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
