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

	name := strings.Split(m.Content[1:], " ")[0]

	files, err := filepath.Glob("./img/"+ name + ".*")
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
	}
	//fmt.Println("./img/"+ name + ".*")

	img, err := os.Open("./" + files[0])
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
	}
	defer img.Close()
	i.Session.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   strings.Split(files[0], "\\")[1],
				Reader: img,
			},
		},
	})
}
