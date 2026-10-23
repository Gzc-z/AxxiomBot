// TODO: improve ts logic

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	r2 "math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"axxiom/internal/config"
	"axxiom/internal/interactions"
	ui "axxiom/internal/interactions/texts"

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

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	} else if m.Content == "" {
		return
	} else if m.Content[0] != byte(prefix[0]) {
		return
	} else if isImage(m) {
		return
	} else if strings.Contains(m.Content, "@everyone") {
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
		ctx := &Order{
			s:          s,
			m:          m,
			args:       args,
			argsLength: len(args),
		}
		cmd(ctx)
		return
	}

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

		response := ui.UserResponse(usr)
		_, err := s.ChannelMessageSendComplex(m.ChannelID, response)
		if err != nil {
			log.Println("error sending message: ", err)
		}

	case "garçom":
		members, _ := s.GuildMembers(m.GuildID, "", 100)
		if len(members) == 0 {
			s.ChannelMessageSendReply(m.ChannelID, "não há membros no servidor", m.Reference())
			return
		}
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
		if len(args) != 2 {
			s.ChannelMessageSendReply(m.ChannelID, "use **.ban <usuário>**", m.Reference())
			return
		}
		s.ChannelMessageSendReply(m.ChannelID, args[1]+" foi banido(a)", m.Reference())
	case "jogar":
		if len(args) != 2 {
			s.ChannelMessageSendReply(m.ChannelID, "use **.jogar <usuário>**", m.Reference())
			return
		}
		s.ChannelMessageSendReply(m.ChannelID, args[1]+" joga com o menino pô", m.Reference())
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
	apost, err := strconv.ParseFloat(aposta, 32)
	if apost <= 0 {
		s.ChannelMessageSendReply(m.ChannelID, "error: aposta inválida", m.Reference())
		return
	}
	r := rand.Float32()
	if r > 0.58 {
		r += 1
	}
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "error: numero inválido", m.Reference())
		return
	}
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

var catClient = &http.Client{Timeout: 5 * time.Second}

type catFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func CatFacts(ctx *Order) {
	r, err := http.Get("https://catfact.ninja/fact")
	if err != nil {
		log.Println("catfact request: ", err)
		return
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		log.Printf("catfact: unexpected status %d", r.StatusCode)
		return
	}

	var fact catFact
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&fact); err != nil {
		log.Println("catfact decode: ", err)
		return
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
	expression = strings.ReplaceAll(expression, "×", "*")
	expression = strings.ReplaceAll(expression, "÷", "/")
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
	server := config.MinecraftAddress()
	start := time.Now()

	conn, err := net.DialTimeout("tcp", server, time.Second*5)
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

func isImage(m *discordgo.MessageCreate) bool {
	for _, a := range m.Attachments {
		if a.ContentType != "" && strings.HasPrefix(a.ContentType, "image/") {
			return true
		}
	}
	return false
}
