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
	if m.Author.Bot {
		return
	}
	if m.Content[0] != byte(prefix[0]) {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

	// oh gosh
	help := []string{
		"ping",
		"help",
		"time",
		"timer",
		"comandos",
		"users",
		"fact",
		"rand",
		"me",
		"mine",
		"calc",
	}

	// TODO: separate this
	switch args[0] {
	case "ping":
		s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
	case "time":
		time := time.Now().Format("15h : 04m : 05s")
		s.ChannelMessageSendReply(m.ChannelID, time, m.Reference())
	case "timer":
		if len(args) != 2 {
			s.ChannelMessageSendReply(m.ChannelID, "error: use **timer <segundos>**", m.Reference())
			return
		}

		n, err := strconv.Atoi(args[1])
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, "error: numero inválido\notempo é medido em segundos", m.Reference())
			return
		}
		go s.ChannelTyping(m.ChannelID)
		timer := time.NewTimer(time.Duration(n) * time.Second)
		<-timer.C
		s.ChannelMessageSendReply(m.ChannelID, "O tempo acabou", m.Reference())
	case "comandos", "commands", "help":
		s.ChannelMessageSendReply(m.ChannelID, "## commands available: ", m.Reference())
		for _, v := range interactions.Commands {
			s.ChannelMessageSendComplex(m.ChannelID, ui.CommandsResponse(v))
		}
		s.ChannelMessageSend(m.ChannelID, "## comandos de texto disponíveis: ")
		var tc string
		for _, cmd := range help {
			tc += fmt.Sprintf(".%s\n", cmd)
		}
		s.ChannelMessageSend(m.ChannelID, tc)

	case "users", "user", "usrs", "membros":
		members, _ := s.GuildMembers(m.GuildID, "", 100)
		s.ChannelMessageSendReply(m.ChannelID, "## Membros:\n`", m.Reference())
		for _, v := range members {
			response := ui.MembersResponse(v)
			go s.ChannelMessageSendComplex(m.ChannelID, response)
		}
	case "fact", "cat", "catfact":
		catFacts := catFacts()
		s.ChannelMessageSendReply(m.ChannelID, catFacts.Fact, m.Reference())
	case "rand", "random", "r":
		if len(args) == 1 {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "numero aleatório: 10",
			})
			r := rand.Intn(10)
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(r), m.Reference())
			return
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, "error: numero inválido", m.Reference())
			return
		}
		r := rand.Intn(num)
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("numero aleatório: %d", num),
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
				s.ChannelMessageSendReply(m.ChannelID, "error: numero inválido", m.Reference())
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
		for _, v := range args[1:] {
			if strings.Contains(v, "x") || strings.Contains(v, "X") {
				s.ChannelMessageSendReply(m.ChannelID, "use * ao invés de x; seu boboca", m.Reference())
				return
			}
		}
		program, err := expr.Compile(strings.Join(args[1:], " "))
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSendReply(m.ChannelID, "error: calculo não suportado e/ou inválido (ainda)\n", m.Reference())
			return
		}

		result, err := expr.Run(program, nil)
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSendReply(m.ChannelID, "error: expressão matemática inválida", m.Reference())
			return
		}
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(result), m.Reference())
	}
}

func catFacts() struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
} {
	catFactAPIURL := "https://catfact.ninja"
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

func calc() {
}
