package main

import (
	"fmt"
	"math"
	"math/big"
	cryptorand "crypto/rand"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var err error

// call if something goes wrong
func sadness(s *discordgo.Session, m *discordgo.MessageCreate, errl error) {
	if m != nil{
		s.ChannelMessageSendReply(m.ChannelID, "Sorry, my creator must have fucked something up.\nPlease pierce him with a sanguine lance and drink his blood.", m.Reference())
	}
	fmt.Println(errl)
}

func help(s *discordgo.Session, m *discordgo.MessageCreate){
	ht, err := os.ReadFile("help.md")
	if err!=nil{
		sadness(s,m,err)
	}
	s.ChannelMessageSendReply(m.ChannelID, 
		string(ht),
		m.Reference())
}

func roll(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil || strings.Contains(m.Content, "bread"){
		return
	}
	if roll := composeRoll(m.Content); roll == ""{
		s.ChannelMessageSendReply(c.ID, roll, m.SoftReference())
	} else {
		s.ChannelMessageSendReply(c.ID, "Your roll is "+ roll+".", m.SoftReference())
	}
}


func editRoll(s *discordgo.Session, m *discordgo.MessageUpdate, mymsg *discordgo.Message){
	c, err := s.State.Channel(m.ChannelID)
	if err != nil || strings.Contains(m.Content, "bread"){
		return
	}
	if roll := composeRoll(m.Content); roll == ""{
		s.ChannelMessageEdit(c.ID, mymsg.ID, "I know what you are.")
	} else {
		s.ChannelMessageEdit(c.ID, mymsg.ID, "Your roll is "+ roll +".")
	}
}

func composeRoll(i string) string{
	var mod int
	r, _ := strings.CutPrefix(i, ".roll ")
	if idk, err := strconv.Atoi(r); err == nil && idk > 0 {
		num, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(idk)))
		return strconv.Itoa(int(num.Int64()+1))
	} else {
		rest := r
		if !strings.Contains(rest, "d") {
			return ""
		}

		rest = "+"+rest // magic 2
		operators := "+-*^_"
		rawStr := ""
		sum := 0

		for len(rest)>0 {
			// the following code is magic! don't touch
			
			modRune := strings.IndexAny(rest[1:], operators)+1 // next rune = end of our bit

			if modRune == 0 {
				modRune = len(rest)
			}

			bit := rest[1:modRune] // this killed sancho
			sign := rest[0]
			fmt.Println(bit)
			if strings.Contains(bit, "d") {
				mod = 0
				countS, maxS, found := strings.Cut(bit,"d")
				if !found{
					return ""
				}
				var count int
				if len(countS)>0{
					count, err = strconv.Atoi(countS)
					if err!=nil{
						return ""
					}
				} else {
					count = 1
				}
				max, err := strconv.Atoi(maxS)
				if err!=nil{
					return ""
				}
				for i := 0; i < count; i++ {
					vbig, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64((max))))
					v := int(vbig.Int64()) + 1
					rawStr += strconv.Itoa(v) + " "
					mod += v
				}
				fmt.Printf("count: %d, max: %d, total: %d\n", count, max, mod)
			} else {
				mod, err = strconv.Atoi(bit)
				if err!=nil{
					return ""
				}
			}

			switch sign {
			case '+':
				sum += mod
			case '*':
				sum *= mod
			case '-':
				sum -= mod
			case '^':
				sum = int(math.Pow(float64(sum), float64(mod)))
			case '_':
				sum = int(math.Pow(float64(mod), float64(sum)))
			}
			rest = rest[modRune:]
		}

		out := strconv.Itoa(sum) + " (" + rawStr[:len(rawStr)-1] + ")"
		outFull := "Your roll is " + out + "."
		if len(outFull) > 2000 {
			return strconv.Itoa(sum)
		} else {return out}
	}
}

// call if user fucked up command use - shift the guilt, don't tell a single soul that you're the one to blame
func iKnowWhatYouAre(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, "I know what you are.", m.Reference())
}

func bod(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	resp, err := http.Get("https://tiphereth.zasz.su/static/assets/cards/Shi1.png")
	if err != nil {
		sadness(s,m,err)
	}
	defer resp.Body.Close()
	yujin := resp.Body

	resp, err = http.Get("https://tiphereth.zasz.su/static/assets/cards/Roland4Phase_Yujin.png")
	if err != nil {
		sadness(s,m,err)
	}
	defer resp.Body.Close()
	yujinDead := resp.Body

	pingId, _ := strings.CutPrefix(m.Content, ".bod ")
	//pingId = pingId[1:len(pingId)-1]
	ping := ""
	if _, err := strconv.Atoi(pingId[2 : len(pingId)-1]); pingId[1] == '@' && err == nil {
		ping = pingId
	}
	rollbig,_ := cryptorand.Int(cryptorand.Reader, big.NewInt(4))
	roll := int(rollbig.Int64())+1
	if roll == 4 {
		s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
			Content: "**4**" + " " + ping,
			Files: []*discordgo.File{
				{
					Name:   "yujin.png",
					Reader: yujin,
				},
			},
		})
	} else {
		s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
			Content: strconv.Itoa(roll),
			Files: []*discordgo.File{
				{
					Name:   "yujinDead.jpg",
					Reader: yujinDead,
				},
			},
		})
	}
}

func conceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, "What... is it this time?", m.Reference())
}

func fut(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages := []string{"...I won't even have to call Father for this.",
		"...Hold your tongue; I will no longer tolerate any more \"ingenious ideas\".",
		"You're nothing before a Second Kindred, let alone Father.",
		"...Did you learn that from that Knight? How humorous.",
	}
	s.ChannelMessageSendReply(m.ChannelID, messages[rand.Intn(len(messages))], m.Reference())
}

func sanchoball(s *discordgo.Session, m *discordgo.MessageCreate){
	outcomes := []string{"No.","Perhaps.","If Father wills it.","Most definitely.","Absolutely.","Maybe.","Clearly not.", "You'd be stupid to try.", "Is it not obvious?", "...", "Depends on you, and you alone.", "...You remind me of that arrogantly hopeful Fixer.", "For the Family, of course.", "V-verily, tis t-true-\n...You get the point. Yes.", "That is... simply impossible."}
	today := time.Now().YearDay()+time.Now().Year()*365

	_,msg,_ := strings.Cut(m.Content, " ")
	msg = strings.ToLower(msg)
	msgSalt := 0
	for i, l := range msg {
		msgSalt += int(l)*(i+1) * (((i+1)%2)*2-1)
	}
	var pseudo int
	id, _ := strconv.Atoi(m.Author.ID)
	if msgSalt != 0{
		saltedRand := rand.New(rand.NewSource(int64(msgSalt)))
		pseudo = int(math.Abs(float64((today + id + saltedRand.Int())%len(outcomes))))
	} else {
		pseudo = int(math.Abs(float64((today + id + rand.Int())%len(outcomes))))
	}

	s.ChannelMessageSendReply(m.ChannelID, outcomes[pseudo], m.SoftReference())
}