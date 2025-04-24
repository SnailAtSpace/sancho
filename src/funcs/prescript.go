package funcs

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func Prescript(inst *Instance, m *discordgo.MessageCreate) {
	orb := imagick.NewMagickWand()
	defer orb.Destroy()

	err := orb.ReadImage("img/prescript.gif")
	if err != nil {
		inst.ErrorChan <- err
		return
	}

	poster, _ := inst.Session.GuildMember(m.GuildID, m.Author.ID)
	nick := poster.Nick
	if nick == "" {
		nick = poster.User.GlobalName
		if nick == "" {
			nick = poster.User.Username
		}
	}

	authorIDNum, _ := strconv.Atoi(m.Author.ID)
	today := time.Now().YearDay() + time.Now().Year()*365
	saltedRand := rand.New(rand.NewSource(int64(today*(authorIDNum%100000))))

	prescriptIndex := saltedRand.Intn(int(orb.GetNumberImages()))
	orb.SetIteratorIndex(int(prescriptIndex))
	out, _ := orb.GetImageBlob()

	outReader := bytes.NewReader(out)
	inst.Session.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Content:   fmt.Sprintf("Prescript #**%d** delivered to **%s**.\n*\\- YV*", int64(today*(authorIDNum%100000)), nick),
		Files: []*discordgo.File{{
			Name:   "img.gif",
			Reader: outReader,
		}},
	})
}
