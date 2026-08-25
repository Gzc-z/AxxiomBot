package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
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

type Context struct {
	s    *discordgo.Session
	m    *discordgo.MessageCreate
	args []string
}

var BuiltInCmds = map[string]func(*Context){
	"ping":  Ping,
	"time":  Time,
	"timer": Timer,
	"me":    Me,
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	} else if m.Content == "" {
		return
	} else if isImage(s, m) {
		return
	}

	if m.GuildID != config.GetGuildID() {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

	if m.Content[0] != byte(prefix[0]) {
		return
	}
	cmd, ok := BuiltInCmds[args[0]]
	if ok {
		cmd(&Context{s, m, args})
		// Help(&Context{s, m, args})
		// return
	}

	switch args[0] {
	case "users", "user", "usrs", "membros":
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

	case "forge":
		params := url.Values{}
		params.Add("gameId", "432")
		params.Add("index", "1")
		params.Add("pageSize", "10")

		endpoint := "https://api.curseforge.com/v1/mods/search?" + params.Encode()
		req, err := http.NewRequest(
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("x-api-key", config.GetMinecraft().Forge)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))

	case "calc", "math", "c":
		// ×÷π√∆£^✓%
		if len(args) == 1 {
			s.ChannelMessageSendReply(m.ChannelID, "use **calc <expressão>**", m.Reference())
			return
		}
		expression := strings.Join(args[1:], " ")
		expression = strings.Replace(expression, "×", "*", -1)
		expression = strings.Replace(expression, "÷", "/", -1)
		//
		// strings.Replace(v, "π", "3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067982148086513282306647093844609550582231725359408128481117450284102701938521105559644622948954930381964428810975665933446128475648233786783165271201909145648566923460348610454326648213393607260249141273724587006606315588174881520920962829254091715364367892590360011330530548820466521384146951941511609433057270365759591953092186117381932611793105118548074462379962749567351885752724891227938183011", -1)
		comp, err := expr.Compile(expression)
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSendReply(m.ChannelID, "error: calculo não suportado e/ou inválido (ainda)\n", m.Reference())
			return
		}

		result, err := expr.Run(comp, nil)
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSendReply(m.ChannelID, "error: expressão matemática inválida", m.Reference())
			return
		}
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(result), m.Reference())
	}
}

func Ping(ctx *Context) {
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, "pong", ctx.m.Reference())
}

func Help(ctx *Context) {
	s := ctx.s
	m := ctx.m
	s.ChannelMessageSendReply(m.ChannelID, "## commands available: ", m.Reference())
	for _, v := range interactions.Commands {
		s.ChannelMessageSendComplex(m.ChannelID, ui.CommandsResponse(v))
	}
	s.ChannelMessageSend(m.ChannelID, "## comandos de texto disponíveis: ")
	var txtc string
	for cmd := range BuiltInCmds {
		txtc += fmt.Sprintf(".%s\n", cmd)
	}
	s.ChannelMessageSend(m.ChannelID, txtc)
}

func Time(ctx *Context) {
	time := time.Now().Format("15h : 04m : 05s")
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, time, ctx.m.Reference())
}

func Timer(ctx *Context) {
	s := ctx.s
	m := ctx.m
	args := ctx.args

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
}

func Me(ctx *Context) {
	usr, _ := ctx.s.User(ctx.m.Author.ID)
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, "## you:", ctx.m.Reference())
	response := ui.UserResponse(usr)
	ctx.s.ChannelMessageSendComplex(ctx.m.ChannelID, response)
}

func Members(ctx *Context) {
	members, _ := ctx.s.GuildMembers(ctx.m.GuildID, "", 100)
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, "## Membros:\n`", ctx.m.Reference())
	for _, v := range members {
		response := ui.MembersResponse(v)
		ctx.s.ChannelMessageSendComplex(ctx.m.ChannelID, response)
	}
}

var fact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func catFacts() struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
} {
	r, err := http.Get("https://catfact.ninja" + "/fact")
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

func isImage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, a := range m.Attachments {
		if a.ContentType != "" && strings.HasPrefix(a.ContentType, "image/") {
			return true
		}
	}
	return false
}
