package main

import (
	. "github.com/snailatspace/sancho/src/funcs"

	"github.com/bwmarrin/discordgo"
)

type CliCommand struct {
	Alias string
	Func  func(*discordgo.Session, string) error
}

type BotCommand struct {
	Aliases []string
	Func    func(*Instance, *discordgo.MessageCreate)
}

var cliCmds = []CliCommand{
	{"gn", gn},
	{"chan", changeChannel},
	{"say", sayCli},
	{"sayr", sayCliReply},
	{"sayi", sayCliAttachment},
	{"listen", listenChannel},
	{"channels", listChannels},
	{"dump", dump},
	{"edit", edit},
	{"dm", dm},
	{"del", delete},
	{"pff", postFromFile},
}

var botCmds = []BotCommand{
	{[]string{"help"}, Help},
	{[]string{"roll"}, Roll},
	{[]string{"bod"}, BoD},
	{[]string{"nacho", "badword", "rye", "ryeldhunt", "pet", "sanitize", "sanitise"}, SendImg},
	{[]string{"remind", "remindme"}, SetReminder},
	{[]string{"reminders"}, ListReminders},
	{[]string{"deremind", "forget"}, DeleteReminder},
	{[]string{"lmd"}, LamentMournAndDespair},
	{[]string{"said", "speechbubble"}, SpeechBubble},
	{[]string{"sanchoball", "8ball"}, Sanchoball},
	{[]string{"settz"}, SetTimezone},
	{[]string{"yesod", "jpeg", "corru"}, ApplyImageProcessing},
	{[]string{"limbusroll", "skill", "skillroll"}, LimbusRoll},
	{[]string{"prescript"}, Prescript},
}
