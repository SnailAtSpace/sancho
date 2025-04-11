package funcs

import "github.com/bwmarrin/discordgo"

type Instance struct {
	Session   *discordgo.Session
	Reminder  *ReminderManager
	ErrorChan chan error
}
