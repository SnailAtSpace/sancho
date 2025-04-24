package funcs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var imgAltAliasLUT = map[string]string{
	"sanitise" : "sanitize",
}

func SendImg(i *Instance, m *discordgo.MessageCreate) {
	var ok bool

	c, err := i.Session.State.Channel(m.ChannelID)
	if err != nil {
		i.ErrorChan <- err
		return
	}

	name := strings.ToLower(strings.Split(m.Content[1:], " ")[0])

	files, err := filepath.Glob("./img/"+ name + ".*")
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
		return
	}

	if files == nil {
		name, ok = imgAltAliasLUT[name]
		files, err = filepath.Glob("./img/"+ name + ".*")
		if err != nil || !ok {
			Sadness(i, m)
			i.ErrorChan <- err
			return
		}
	}
	//fmt.Println("./img/"+ name + ".*")

	img, err := os.Open("./" + files[0])
	if err != nil {
		Sadness(i, m)
		i.ErrorChan <- err
		return
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
