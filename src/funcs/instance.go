package funcs

import "github.com/bwmarrin/discordgo"

type Instance struct {
	Session   *discordgo.Session
	RManager  *ReminderManager
	ErrorChan chan error
}