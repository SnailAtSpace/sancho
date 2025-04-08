package main

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func gn(discord *discordgo.Session, s string) error {
	img, err := os.Open("img/goodnight.png")
	if err != nil {
		return err
	}
	defer img.Close()
	discord.ChannelMessageSendComplex("1331332284372222074", &discordgo.MessageSend{
		Content: "Good night, Family. Tomorrow we shall take part in the banquet... again. For now, however, I will rest.",
		Files: []*discordgo.File{
			{
				Name:   "goodnight.png",
				Reader: img,
			},
		},
	})
	retch <- 1 // idk how go works, holy shit!
	return nil
}

func changeChannel(discord *discordgo.Session, s string) error {
	echoChan, _ = strings.CutPrefix(s, "chan ")
	return nil
}

func sayCli(discord *discordgo.Session, s string) error {
	raw, _ := strings.CutPrefix(s, "say ")
	discord.ChannelMessageSend(echoChan, raw)
	return nil
}

func sayCliReply(discord *discordgo.Session, s string) error {
	raw, _ := strings.CutPrefix(s, "sayr ")
	repId, msg, found := strings.Cut(raw, " ")
	if !found {
		return fmt.Errorf("Wow. You somehow fucked that up.")
	}
	discord.ChannelMessageSendReply(echoChan, msg, &discordgo.MessageReference{MessageID: repId})
	return nil
}

func sayCliAttachment(discord *discordgo.Session, s string) error {
	raw, _ := strings.CutPrefix(s, "sayi ")
	name, msg, found := strings.Cut(raw, " ")
	if !found {
		msg = ""
	}
	var msgId string
	if strings.Contains(msg, " ") {
		msgId, msg, _ = strings.Cut(msg, " ")
	}
	img, err := os.Open("img/" + name)
	if err != nil {
		return err
	}
	defer img.Close()
	if msgId == "" {
		discord.ChannelMessageSendComplex(echoChan, &discordgo.MessageSend{
			Content: msg,
			Files: []*discordgo.File{
				{
					Name:   name,
					Reader: img,
				},
			},
		})
	} else {
		discord.ChannelMessageSendComplex(echoChan, &discordgo.MessageSend{
			Content:   msg,
			Reference: &discordgo.MessageReference{MessageID: msgId},
			Files: []*discordgo.File{
				{
					Name:   name,
					Reader: img,
				},
			},
		})
	}
	return nil
}

func getPfp(discord *discordgo.Session, s string) error {
	id, _ := strings.CutPrefix(s, "pfp ")
	user, err := discord.User(id)
	if err != nil {
		return fmt.Errorf("user not found "+err.Error())
	}
	pfp, err := discord.UserAvatarDecode(user)
	if err != nil {
		return fmt.Errorf("could not pull pfp "+err.Error())
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, pfp)
	if err != nil {
		return fmt.Errorf("error encoding: weird "+err.Error())
	}
	err = os.WriteFile(fmt.Sprintf("img/%s.png", id), buf.Bytes(), os.FileMode(os.O_RDWR))
	if err != nil {
		return fmt.Errorf("error when saving the file "+err.Error())
	}
	return nil
}

func listChannels(discord *discordgo.Session, s string) error {
	id, _ := strings.CutPrefix(s, "channels ")
	chans, err := discord.GuildChannels(id)
	if err != nil {
		return fmt.Errorf("server not found: "+err.Error())
	}
	serv, _ := discord.Guild(id)
	fmt.Printf("Channels of server %s:\n", serv.Name)
	for i, channel := range chans {
		fmt.Printf("[%d] %s (%s)\n", i, channel.Name, channel.ID)
	}
	return nil
}

func listenChannel(discord *discordgo.Session, s string) error {
	id, found := strings.CutPrefix(s, "listen ")
	if !found {
		listenChannelID = ""
	} else {
		listenChannelID = id
		fmt.Println(listenChannelID)
	}
	return nil
}