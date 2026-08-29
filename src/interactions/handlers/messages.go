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

type PingMine struct {
	Online bool
	Time   time.Duration
}

type Context struct {
	s          *discordgo.Session
	m          *discordgo.MessageCreate
	args       []string
	argsLength int
}

type Script struct {
	Name  string
	Alias []string
	Func  func(*Context)
}

var BuiltInScripts = map[string]Script{
	"timer":   {Name: "timer", Func: Timer},
	"members": {Name: "members", Func: Members},
	"color":   {Name: "color", Func: RandomColor},
	"calc":    {Name: "calc", Func: Calc},
	"catfact": {Name: "catfact", Func: CatFacts},
	"random":  {Name: "random", Func: Random},
	"mine":    {Name: "mine", Func: Mine},
	"help":    {Name: "help", Func: Help},
}

// this would call API -> external scripts
// var ExtScript = map[string]API{}

func (src *Script) mkalias(alias ...string) {
	for _, v := range alias {
		src.Alias = append(src.Alias, v)
	}
}

// resolve
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

	cmd, ok := BuiltInScripts[args[0]]
	if ok {
		cmd.Func(&Context{
			s:          s,
			m:          m,
			args:       args,
			argsLength: len(args),
		})
	} else {
		// command not found
	}

	// other default additional commands
	switch args[0] {
	case "ping":
		s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
	case "time":
		time := time.Now().Format("15h : 04m : 05s")
		s.ChannelMessageSendReply(m.ChannelID, time, m.Reference())
	case "me":
		usr, _ := s.User(m.Author.ID)
		s.ChannelMessageSendReply(m.ChannelID, "## you:", m.Reference())
		response := ui.UserResponse(usr)
		s.ChannelMessageSendComplex(m.ChannelID, response)

	case "garçom":
		members, _ := s.GuildMembers(m.GuildID, "", 100)
		user := members[rand.Intn(len(members))]
		response := ui.MembersResponse(user)

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
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("%v foi banido", args[1]), m.Reference())

	}
}

func Help(ctx *Context) {
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

func Members(ctx *Context) {
	members, _ := ctx.s.GuildMembers(ctx.m.GuildID, "", 100)
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, "## Membros:\n`", ctx.m.Reference())
	if len(members) > 10 {
		return
	}
	for _, v := range members {
		response := ui.MembersResponse(v)
		ctx.s.ChannelMessageSendComplex(ctx.m.ChannelID, response)
	}
}

var fact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func CatFacts(ctx *Context) {
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

func RandomColor(ctx *Context) {
	ctx.s.ChannelMessageSendReply(ctx.m.ChannelID, strconv.Itoa(rand.Intn(0xffffff)), ctx.m.Reference())
}

func Calc(ctx *Context) {
	// ×÷π√∆£^✓%
	m := ctx.m
	if len(ctx.args) == 1 {
		ctx.s.ChannelMessageSendReply(m.ChannelID, "use **calc <expressão>**", m.Reference())
		return
	}
	expression := strings.Join(ctx.args[1:], " ")
	expression = strings.Replace(expression, "×", "*", -1)
	expression = strings.Replace(expression, "÷", "/", -1)
	//
	// strings.Replace(v, "π", "3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067982148086513282306647093844609550582231725359408128481117450284102701938521105559644622948954930381964428810975665933446128475648233786783165271201909145648566923460348610454326648213393607260249141273724587006606315588174881520920962829254091715364367892590360011330530548820466521384146951941511609433057270365759591953092186117381932611793105118548074462379962749567351885752724891227938183011", -1)
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

func Random(ctx *Context) {
	args := ctx.args
	m := ctx.m
	s := ctx.s
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
}

func Mine(ctx *Context) {
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
