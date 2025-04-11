package funcs

import (
	cryptorand "crypto/rand"
	"math"
	"math/big"
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
func Sadness(i *Instance, m *discordgo.MessageCreate) {
	if m != nil {
		i.Session.ChannelMessageSendReply(m.ChannelID, "My creator must have fucked something up.\nPlease ping him immediately and repeatedly until he sees this problem.", m.Reference())
	}
}

func Help(i *Instance, m *discordgo.MessageCreate) {
	ht, err := os.ReadFile("help.md")
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
		return
	}
	_, err = i.Session.ChannelMessageSendReply(m.ChannelID,
		string(ht),
		m.Reference())
	i.ErrorChan <- err
}

// call if user fucked up command use - shift the guilt, don't tell a single soul that you're the one to blame
func iKnowWhatYouAre(i *Instance, m *discordgo.MessageCreate) {
	i.Session.ChannelMessageSendReply(m.ChannelID, "I know what you are.", m.Reference())
}

func BoD(i *Instance, m *discordgo.MessageCreate) {
	c, err := i.Session.State.Channel(m.ChannelID)
	if err != nil {
		i.ErrorChan <- err
		return
	}
	resp, err := http.Get("https://tiphereth.zasz.su/static/assets/cards/Shi1.png")
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
		return
	}
	defer resp.Body.Close()
	yujin := resp.Body

	resp, err = http.Get("https://tiphereth.zasz.su/static/assets/cards/Roland4Phase_Yujin.png")
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
		return
	}
	defer resp.Body.Close()
	yujinDead := resp.Body

	pingId, _ := strings.CutPrefix(m.Content, ".bod ")
	//pingId = pingId[1:len(pingId)-1]
	ping := ""
	if _, err := strconv.Atoi(pingId[2 : len(pingId)-1]); pingId[1] == '@' && err == nil {
		ping = pingId
	}
	rollbig, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(4))
	roll := int(rollbig.Int64()) + 1
	if roll == 4 {
		i.Session.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
			Content: "**4**" + " " + ping,
			Files: []*discordgo.File{
				{
					Name:   "yujin.png",
					Reader: yujin,
				},
			},
		})
	} else {
		i.Session.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
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

func fut(i *Instance, m *discordgo.MessageCreate) {
	messages := []string{"...I won't even have to call Father for this.",
		"...Hold your tongue; I will no longer tolerate any more \"ingenious ideas\".",
		"You're nothing before a Second Kindred, let alone Father.",
		"...Did you learn that from that Knight? How humorous.",
	}
	i.Session.ChannelMessageSendReply(m.ChannelID, messages[rand.Intn(len(messages))], m.Reference())
}

func Sanchoball(i *Instance, m *discordgo.MessageCreate) {
	outcomes := []string{"No.", "Perhaps.", "If Father wills it.", "Most definitely.", "Absolutely.", "Maybe.", "Clearly not.", "You'd be stupid to try.", "Is it not obvious?", "...", "Depends on you, and you alone.", "...You remind me of that arrogantly hopeful Fixer.", "For the Family, of course.", "V-verily, tis t-true-\n...You get the point. Yes.", "That is... simply impossible."}
	today := time.Now().YearDay() + time.Now().Year()*365

	_, msg, _ := strings.Cut(m.Content, " ")
	msg = strings.ToLower(msg)
	msgSalt := 0
	for i, l := range msg {
		msgSalt += int(l) * (i + 1) * (((i+1)%2)*2 - 1)
	}
	var pseudo int
	id, _ := strconv.Atoi(m.Author.ID)
	if msgSalt != 0 {
		saltedRand := rand.New(rand.NewSource(int64(msgSalt)))
		pseudo = int(math.Abs(float64((today + id + saltedRand.Int()) % len(outcomes))))
	} else {
		pseudo = int(math.Abs(float64((today + id + rand.Int()) % len(outcomes))))
	}

	i.Session.ChannelMessageSendReply(m.ChannelID, outcomes[pseudo], m.SoftReference())
}
