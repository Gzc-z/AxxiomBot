package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"axiom/src/config"
	"axiom/src/interactions"
	ui "axiom/src/responses"

	"github.com/bwmarrin/discordgo"
	"github.com/expr-lang/expr"
)

var prefix = "."

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if m.Content == "" {
		return
	}
	if m.Content[0] != byte(prefix[0]) {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

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
			s.ChannelMessageSendReply(m.ChannelID, "error: numero inválido\no tempo é medido em segundos", m.Reference())
			return
		}
		s.ChannelTyping(m.ChannelID)
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
			s.ChannelMessageSendComplex(m.ChannelID, response)
		}
	case "fact", "cat", "catfact":
		catFacts := catFacts()
		s.ChannelMessageSendReply(m.ChannelID, catFacts.Fact, m.Reference())
	case "rand", "random", "r":
		if len(args) == 1 {
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
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(r), m.Reference())

	case "me":
		usr, _ := s.User(m.Author.ID)
		s.ChannelMessageSendReply(m.ChannelID, "## you:", m.Reference())
		response := ui.UserResponse(usr)
		s.ChannelMessageSendComplex(m.ChannelID, response)
	case "mine":
		if len(args) != 2 {
			pings := 5
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("%v pings", pings),
				Flags:   discordgo.MessageFlagsEphemeral,
			})
			stats := serverPing(s, m)
			if !stats.Online {
				s.ChannelMessageSendReply(m.ChannelID, "🔴 Server Offline", m.Reference())
				return
			}
			for i := range pings {
				s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Content: fmt.Sprintf("Online: 🟢 %v seq=%d\n", stats.Time, i+1),
					Flags:   discordgo.MessageFlagsEphemeral,
				})
			}
			return
		}
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("%v pings", args[1]),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		count, err := strconv.Atoi(args[1])
		if err != nil {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "error: numero inválido",
				Flags:   discordgo.MessageFlagsEphemeral,
			})
			return
		}
		stats := serverPing(s, m, count)
		if !stats.Online {
			s.ChannelMessageSendReply(m.ChannelID, "🔴 Server Offline", m.Reference())
			return
		}
		for i := range count {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("Online: 🟢 %v seq=%d\n", stats.Time, i+1),
				Flags:   discordgo.MessageFlagsEphemeral,
			})
		}
	case "garçom":
		members, _ := s.GuildMembers(m.GuildID, "", 100)
		user := members[rand.Intn(len(members))]

		message := fmt.Sprintf("<@%s>", m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v 🤵‍♂️ garçom\ntrouxe cerveja :)", message))

		response := ui.MembersResponse(user)
		s.ChannelMessageSendComplex(m.ChannelID, response)

	case "banir", "ban":
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("%v banido", args[1]), m.Reference())

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

const catFactAPIURL = "https://catfact.ninja"

var fact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func catFacts() struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
} {
	r, err := http.Get(catFactAPIURL + "/fact")
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&fact)
	if err != nil {
		log.Println(err)
	}
	return fact
}

type PingMine struct {
	Online bool
	Time   time.Duration
}

func serverPing(s *discordgo.Session, m *discordgo.MessageCreate, count ...int) PingMine {
	host := config.GetMinecraft()
	start := time.Now()

	conn, err := net.DialTimeout("tcp", host.IP+":"+host.PORT, time.Second*5)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSendReply(m.ChannelID, "error: couldn't connect to server\ncheck server's ip", m.Reference())
		return PingMine{}
	}
	defer conn.Close()

	p := PingMine{
		Online: true,
		Time:   time.Since(start),
	}
	return p
}
