package main

import "github.com/bwmarrin/discordgo"

type CliCommand struct {
	alias string
	meth  func(*discordgo.Session, string) error
}

type BotCommand struct {
	aliases []string
	meth    func() error
}

var cliCmds = []CliCommand{
	CliCommand{"gn", gn},
	CliCommand{"chan", changeChannel},
	CliCommand{"say", sayCli},
	CliCommand{"sayr", sayCliReply},
	CliCommand{"sayi", sayCliAttachment},
	CliCommand{"listen", listenChannel},
	CliCommand{"channels", listChannels},
}

// var botCmds = []BotCommand{
// 	BotCommand{[]string{"roll"}, roll},
// }
