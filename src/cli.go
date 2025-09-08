package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/go-tts/tts/pkg/speech"
	"github.com/hajimehoshi/go-mp3"
)

var connection *discordgo.VoiceConnection

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
	raw = strings.ReplaceAll(raw, "\\n", "\n")
	discord.ChannelMessageSend(echoChan, raw)
	return nil
}

func sayCliReply(discord *discordgo.Session, s string) error {
	raw, _ := strings.CutPrefix(s, "sayr ")
	repId, msg, found := strings.Cut(raw, " ")
	if !found {
		return fmt.Errorf("Wow. You somehow fucked that up.")
	}
	msg = strings.ReplaceAll(msg, "\\n", "\n")
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
	msg = strings.ReplaceAll(msg, "\\n", "\n")
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
		return fmt.Errorf("user not found " + err.Error())
	}
	pfp, err := discord.UserAvatarDecode(user)
	if err != nil {
		return fmt.Errorf("could not pull pfp " + err.Error())
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, pfp)
	if err != nil {
		return fmt.Errorf("error encoding: weird " + err.Error())
	}
	err = os.WriteFile(fmt.Sprintf("img/%s.png", id), buf.Bytes(), os.FileMode(os.O_RDWR))
	if err != nil {
		return fmt.Errorf("error when saving the file " + err.Error())
	}
	return nil
}

func listChannels(discord *discordgo.Session, s string) error {
	id, _ := strings.CutPrefix(s, "channels ")
	chans, err := discord.GuildChannels(id)
	if err != nil {
		return fmt.Errorf("server not found: " + err.Error())
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

func dump(discord *discordgo.Session, s string) error {
	id, _ := strings.CutPrefix(s, "dump ")
	msg, err := discord.ChannelMessages(id, 100, "", "", "")

	if err != nil {
		return err
	}

	channel, err := discord.Channel(id)
	if err != nil {
		return err
	}

	chanName := channel.Name
	if channel.Type == discordgo.ChannelTypeDM {
		chanName = channel.Recipients[0].Username
	}

	out := fmt.Sprintf("# Dump of %s (%s):\n\n", chanName, channel.ID)
	slices.Reverse(msg)
	for _, r := range msg {
		out += fmt.Sprintf("[%s] %s:\n%s\n\n", r.Timestamp.Local().Format(time.TimeOnly), r.Author.Username, r.ContentWithMentionsReplaced())
	}

	err = os.WriteFile(fmt.Sprintf("dump-%s-%s.txt", chanName, time.Now().Local().Format(time.DateOnly)), []byte(out), 0777)
	if err != nil {
		return err
	}
	fmt.Println("done")
	return nil
}

func edit(discord *discordgo.Session, s string) error {
	raw, _ := strings.CutPrefix(s, "edit ")
	id, raw, _ := strings.Cut(raw, " ")
	channel, newMsg, _ := strings.Cut(raw, " ")
	_, err := discord.ChannelMessageEdit(channel, id, newMsg)
	if err != nil {
		return err
	}
	return nil
}

func dm(discord *discordgo.Session, s string) error {
	id, _ := strings.CutPrefix(s, "dm ")
	channel, err := discord.UserChannelCreate(id)
	if err != nil {
		return err
	}
	echoChan = channel.ID
	fmt.Println(echoChan)
	return nil
}

func delete(discord *discordgo.Session, s string) error {
	id, _ := strings.CutPrefix(s, "del ")
	err := discord.ChannelMessageDelete(echoChan, id)
	if err != nil {
		return err
	}
	return nil
}

func postFromFile(discord *discordgo.Session, s string) error {
	fileName, _ := strings.CutPrefix(s, "pff ")
	text, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	_, err = discord.ChannelMessageSend(echoChan, string(text))
	if err != nil {
		return err
	}
	return nil
}

func joinVC(discord *discordgo.Session, s string) error {
	var err error
	channel := strings.Split(s, " ")[1]
	connection, err = discord.ChannelVoiceJoin("1250579779837493278", channel, false, false)
	return err
}

func leaveVC(discord *discordgo.Session, s string) error {
	if connection != nil {
		err := connection.Disconnect()
		return err
	} else {
		return fmt.Errorf("no connection")
	}
}

func sayInVc(discord *discordgo.Session, s string) error {
	var audioBufferPCM []int16
	audioBufferChan := make(chan []int16)

	text := strings.Join(strings.Split(s, " ")[1:], " ")
	fmt.Println(text)
	audioIn, err := speech.FromText(text, speech.LangUs)
	if err != nil {
		return err
	}

	decoder, err := mp3.NewDecoder(audioIn)
	if err != nil {
		return err
	}
	fmt.Println("mp3 decoded!")
	fmt.Println(decoder.SampleRate())

	audioBufferBytes, _ := io.ReadAll(decoder)
	for i:=0;i<len(audioBufferBytes)-1;i+=2 {
		sample := int16(audioBufferBytes[i+1])<<8 | int16(audioBufferBytes[i])
		audioBufferPCM = append(audioBufferPCM, int16(sample))
	}

	os.WriteFile("sound.wav", audioBufferBytes, 0777)

	fmt.Println(len(audioBufferPCM))
	if len(audioBufferPCM) == 0 {
		return nil
	}
	go dgvoice.SendPCM(connection, audioBufferChan)
	audioBufferChan <- audioBufferPCM
	return nil
}
