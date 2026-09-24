package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	r2 "math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"axxiom/src/config"
	"axxiom/src/interactions"
	ui "axxiom/src/interactions/texts"

	"github.com/bwmarrin/discordgo"
	"github.com/expr-lang/expr"
)

var prefix = "."

type PingMine struct {
	Online bool
	Time   time.Duration
}

type Order struct {
	s          *discordgo.Session
	m          *discordgo.MessageCreate
	args       []string
	argsLength int
}

var BuiltInScripts = map[string]func(*Order){
	"timer":   Timer,
	"members": Members,
	// "color":   Func: RandomColor, idk what i can do here
	"calc":    Calc,
	"catfact": CatFacts,
	"random":  Random,
	"ping":    Mine,
	"help":    Help,
	"apostar": Apostar,
}

// resolve
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	} else if m.Content == "" {
		return
	} else if isImage(s, m) {
		return
	} else if m.Content[0] != byte(prefix[0]) {
		return
	}

	// roles := strings.Split(config.GetGuildID(), ",")
	// if !slices.Contains(roles, m.GuildID) {
	// 	return
	// }

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

	cmd, exist := BuiltInScripts[args[0]]
	if exist {
		s.ChannelTyping(m.ChannelID)
		ctx := &Order{
			s:          s,
			m:          m,
			args:       args,
			argsLength: len(args),
		}
		cmd(ctx)
		return
	}
	fmt.Printf("%s from: %s on: %v\n", args, m.Author, m.GuildID)

	// other default additional commands
	s.ChannelTyping(m.ChannelID)
	switch args[0] {
	case "ping":
		s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
	case "time":
		time := time.Now().Format("15h : 04m : 05s")
		s.ChannelMessageSendReply(m.ChannelID, time, m.Reference())
	case "me":
		usr, _ := s.User(m.Author.ID)
		if len(args) == 2 {
			var err error
			usr, err = s.User(args[1])
			if err != nil {
				s.ChannelMessageSendReply(m.ChannelID, "error: usuário não encontrado", m.Reference())
				return
			}
		}

		s.ChannelMessageSendReply(m.ChannelID, "## you:", m.Reference())
		response := ui.UserResponse(usr)
		s.ChannelMessageSendComplex(m.ChannelID, response)

	case "garçom":
		members, _ := s.GuildMembers(m.GuildID, "", 100)
		user := members[rand.Intn(len(members))]
		response := ui.MembersProfile(user)

		author := fmt.Sprintf("<@%s>", m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v 🤵‍♂️ \ntrouxe cerveja :)", author))
		s.ChannelMessageSendComplex(m.ChannelID, response)

	case "welcome":
		if len(args) == 2 {
			s.ChannelMessageSendReply(m.ChannelID, "wel come: "+args[1], m.Reference())
			return
		}
		usr, _ := s.User(m.Author.ID)
		s.ChannelMessageSendReply(m.ChannelID, "wel come: "+usr.GlobalName, m.Reference())

	case "banir", "ban":
		s.ChannelMessageSendReply(m.ChannelID, args[1]+" foi banido(a)", m.Reference())
	case "axx", "axxiom", "Axxiom":
		s.ChannelMessageSendReply(m.ChannelID, "Use /axx", m.Reference())
	case "jogar":
		if len(args) == 2 {
			s.ChannelMessageSendReply(m.ChannelID, args[1]+" joga com o menino pô", m.Reference())
			return
		}
	}
}

func Help(ctx *Order) {
	s := ctx.s
	m := ctx.m
	s.ChannelMessageSendReply(m.ChannelID, "## commands available: ", m.Reference())
	for _, v := range interactions.Commands {
		s.ChannelMessageSendComplex(m.ChannelID, ui.CommandsResponse(v))
	}
	// txtc := fmt.Sprintf("## comandos de texto: \n")
	// for _, cmd := range BuiltInScripts {
	// 	txtc += fmt.Sprintf(".%s\n", cmd)
	// }
	// s.ChannelMessageSend(m.ChannelID, txtc)
}

func Apostar(ctx *Order) {
	s := ctx.s
	m := ctx.m
	args := ctx.args

	if len(args) != 2 {
		s.ChannelMessageSendReply(m.ChannelID, "use **apostar <aposta>**", m.Reference())
		return
	}
	aposta := args[1]
	r := rand.Float32()
	if r > 0.58 {
		r += 1
	}
	apost, _ := strconv.ParseFloat(aposta, 8)
	calc := float32(apost) * r
	if float32(apost) > calc {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Você apostou %v, e perdeu kkkk", apost), m.Reference())
		return
	}
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Você apostou %v, e ganhou %.2f", apost, calc), m.Reference())
}

func Timer(ctx *Order) {
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

func Members(ctx *Order) {
	members, _ := ctx.s.GuildMembers(ctx.m.GuildID, "", 100)
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, "## Membros:\n`", ctx.m.Reference())
	if len(members) > 10 {
		return
	}
	for _, v := range members {
		response := ui.MembersProfile(v)
		ctx.s.ChannelMessageSendComplex(ctx.m.ChannelID, response)
	}
}

var fact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func CatFacts(ctx *Order) {
	r, err := http.Get("https://catfact.ninja" + "/fact")
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&fact)
	if err != nil {
		log.Println(err)
	}

	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, fact.Fact, ctx.m.Reference())
}

func RandomColor(ctx *Order) {
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, strconv.Itoa(rand.Intn(0xffffff)), ctx.m.Reference())
}

func Calc(ctx *Order) {
	// ×÷π√∆£^✓%
	m := ctx.m
	if len(ctx.args) == 1 {
		ctx.s.ChannelMessageSendReply(m.ChannelID, "use **calc <expressão>**", m.Reference())
		return
	}
	expression := strings.Join(ctx.args[1:], " ")
	expression = strings.Replace(expression, "×", "*", -1)
	expression = strings.Replace(expression, "÷", "/", -1)
	if len(expression) > 100 {
		fmt.Println("expression too long")
		ctx.s.ChannelMessageSendReply(m.ChannelID, "error: calculo muuuito grande\n", m.Reference())
		return
	}

	comp, err := expr.Compile(expression)
	if err != nil {
		fmt.Println(err)
		ctx.s.ChannelMessageSendReply(m.ChannelID, "error: calculo não suportado e/ou inválido (ainda)\n", m.Reference())
		return
	}

	result, err := expr.Run(comp, nil)
	if err != nil {
		fmt.Println(err)
		ctx.s.ChannelMessageSendReply(m.ChannelID, "error: expressão matemática inválida", m.Reference())
		return
	}
	ctx.s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(result), m.Reference())
}

func Random(ctx *Order) {
	args := ctx.args
	m := ctx.m
	s := ctx.s
	if len(args) == 1 {
		r := rand.Intn(10)
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(r), m.Reference())
		return
	}
	if len(args) != 2 {
		s.ChannelMessageSendReply(m.ChannelID, "numero inválido", m.Reference())
		return
	}
	num, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil || num < 1 || num > ^uint64(0) {
		s.ChannelMessageSendReply(m.ChannelID, "numero fora de alcance", m.Reference())
		return
	}

	r := r2.Uint64N(num)
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(r), m.Reference())
}

func Mine(ctx *Order) {
	m := ctx.m
	s := ctx.s

	pings := 5
	for i := range pings {
		stats := serverPing(s, m)
		if !stats.Online {
			// s.ChannelMessageSendReply(m.ChannelID, s, m.Reference)
			s.ChannelMessageSendReply(m.ChannelID, "Online: 🔴", m.Reference())
			return
		}
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("Online: 🟢 %v seq=%d\n", stats.Time, i+1),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
	}
}

func serverPing(s *discordgo.Session, m *discordgo.MessageCreate) PingMine {
	host := config.GetMinecraft()
	start := time.Now()

	conn, err := net.DialTimeout("tcp", host.IP+":"+host.PORT, time.Second*5)
	if err != nil {
		log.Println(err)
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
