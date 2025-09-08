package main

import (
	. "github.com/snailatspace/sancho/src/funcs"

	"fmt"
	"log"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ready(s *discordgo.Session, m *discordgo.Ready) {
	server, err := s.State.Guild("1250579779837493278")
	var num int
	if err != nil {
		num = 13
	} else {
		num = server.MemberCount
	}
	s.UpdateCustomStatus("Allow me to regale thee... that, in this... adventure of mine... Verily, I was blessed with a family of " + strconv.Itoa(num-1) + ".")

	err = inst.RManager.ReviseRemindersAfterStartup(inst)
	if err != nil {
		log.Panicln(err)
	}
}

func guildCreate(s *discordgo.Session, m *discordgo.GuildCreate) {
	if m.Guild.Unavailable {
		return
	}

	channels := m.Guild.Channels

	fmt.Println("Joined server", m.Guild.Name, m.Guild.ID)
	for i := 0; i < len(channels); i++ {
		perms, _ := s.State.UserChannelPermissions("1330935741018276022", channels[i].ID)
		if channels[i].Type == 0 && (perms&2048 == 2048) && time.Now().Unix()-m.JoinedAt.Unix() < 30 {
			s.ChannelMessageSend(channels[i].ID, "The Server will be well-cared for.\n...After all, the onus always fell on me to give roles that you abandoned.")
			return
		}
	}

	if m.Guild.ID == "1250579779837493278" {
		s.UpdateCustomStatus("Allow me to regale thee... that, in this... adventure of mine... Verily, I was blessed with a family of " + strconv.Itoa(m.Guild.MemberCount-2) + ".")
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer panicMsg(s)

	// i shouldn't be doing this
	if m == nil {
		log.Println("fucked up MESSAGE_CREATE event, moving on...")
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)

	if err!=nil {
		log.Println(err)
		return
	}

	if channel.Type == discordgo.ChannelTypeDM {
		fmt.Printf("[%s]{%s} DM from %s: %s\n", m.Timestamp.Local().Format(time.TimeOnly), m.ID, m.Author.Username, m.ContentWithMentionsReplaced())
	} 
	var refid string
	if m.ReferencedMessage != nil {
		refid = m.ReferencedMessage.Author.ID
	}
	normMsg := strings.TrimSpace(strings.ToLower(m.ContentWithMentionsReplaced()))
	if len(normMsg) == 0 {
		return
	}
	if strings.HasPrefix(normMsg, "((") && strings.HasSuffix(normMsg, "))") {
		lastI := strings.LastIndex(normMsg, "))")
		normMsg = strings.TrimSpace(normMsg[2:lastI])
	}
	if len(normMsg) == 0 {
		return
	}
	if normMsg[0] == '.' {
		cmd := strings.Split(normMsg[1:], " ")[0]
		for _, c := range botCmds {
			if slices.Contains(c.Aliases, cmd) {
				go c.Func(inst, m)
				return
			}
		}
	}
	if (m.ChannelID == listenChannelID) || strings.Contains(strings.ToLower(m.ContentWithMentionsReplaced()), "sancho") {
		fmt.Printf("[%s]{%s} @%s %s: %s\n", m.Timestamp.Local().Format(time.TimeOnly), m.ID, channel.Name, m.Author.Username, m.ContentWithMentionsReplaced())
	}
	if strings.Contains(normMsg, "mwah") && (refid == s.State.User.ID || strings.Contains(normMsg, "sancho")) && (m.Author.ID == whoopsID) {
		whoopsMessages := []string{
			"...Not now.\n-# We can hold hands though.",
			"You know what?\n-# mwah",
			"T-thy company is m-most ap-appreciated...",
		}
		pick := rand.IntN(3)
		for i := 0; i < 2; i++ {
			if pick == 1 {
				pick = rand.IntN(3)
			}
		}
		s.ChannelMessageSendReply(m.ChannelID, whoopsMessages[pick], m.Reference())
	} else if strings.Contains(normMsg, "conceived") && m.Author.ID == mattagerID {
		inst.Session.ChannelMessageSendReply(m.ChannelID, "What... is it this time?", m.Reference())
	} else if strings.Contains(normMsg, "sorry") && m.Author.ID == enderID && !slices.Contains(badChannels, m.ChannelID){
		go EnderApologyReaction(inst, m.Message)
	}
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	normMsg := strings.TrimSpace(strings.ToLower(m.ContentWithMentionsReplaced()))
	if len(normMsg) == 0 {
		return
	}

	if m.BeforeUpdate!=nil{
		if !strings.Contains(strings.ToLower(m.BeforeUpdate.Content), "sorry") && strings.Contains(strings.ToLower(m.Content), "sorry") && m.Author.ID == enderID {
			EnderApologyReaction(inst, m.Message)
		}
	}
	if normMsg[0] == '.' {
		cmd := strings.Split(normMsg[1:], " ")[0]
		if cmd == "roll" {
			msgs, err := s.ChannelMessages(m.ChannelID, 100, "", m.ID, "")
			if err != nil {
				log.Println(err)
				return
			}
			var mymsg *discordgo.Message
			for _, r := range msgs {
				if r.ReferencedMessage != nil {
					if r.Author.ID == s.State.User.ID && r.ReferencedMessage.ID == m.ID {
						mymsg = r
					}
				}
			}
			if mymsg == nil {
				log.Println("couldn't find it :(")
				return
			} // ok if we HAVE the message, it must be right
			EditRoll(inst, m, mymsg)
		}
	} 
}
