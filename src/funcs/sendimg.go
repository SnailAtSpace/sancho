package funcs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SendImg(i *Instance, m *discordgo.MessageCreate) {
	c, err := i.Session.State.Channel(m.ChannelID)
	if err != nil {
		i.ErrorChan <- err
	}

	name := strings.Split(m.Content, " ")[0]

	files, err := filepath.Glob(name + ".*")
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
	}

	img, err := os.Open("img/" + files[0])
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
	}
	defer img.Close()
	i.Session.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   name,
				Reader: img,
			},
		},
	})
}
