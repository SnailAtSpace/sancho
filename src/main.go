package main

import (
	. "github.com/snailatspace/sancho/src/funcs"

	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var listenChannelID string

var echoChan string = "1331332284372222074"
var retch = make(chan int) // i hate this. i hate this with all my heart. but i cannot bring myself to do any better.

var inst *Instance

func main() {
	listenChannelID = ""
	fmt.Printf("[%s] I shall pronounce the bot started.\n", time.Now().Format(time.TimeOnly))

	tokens, err := os.Open("secrets.txt")
	if err != nil {
		log.Fatalf("no secrets file: %s", err.Error())
	}
	defer tokens.Close()

	scanner := bufio.NewScanner(tokens)
	var auth_token string
	if scanner.Scan() {
		auth_token = scanner.Text()
	} else {
		log.Fatalf("empty secrets file")
	}

	// firing up the discord session
	discord, err := discordgo.New("Bot " + auth_token)

	if err != nil {
		log.Fatalf("couldn't initialize discord session: %s", err.Error())
	}

	// setting intents and adding event handlers
	discord.Identify.Intents = 335666240
	discord.Identify.Intents |= discordgo.IntentGuildMembers
	discord.Identify.Intents |= discordgo.IntentGuilds
	discord.Identify.Intents |= discordgo.IntentGuildMessages
	discord.Identify.Intents |= discordgo.IntentGuildMessageTyping
	discord.Identify.Intents |= discordgo.IntentGuildMessageReactions
	discord.Identify.Intents |= discordgo.IntentGuildPresences
	discord.AddHandler(guildCreate)
	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)
	discord.AddHandler(messageUpdate)

	// knocking on discord's window
	err = discord.Open()
	if err != nil {
		log.Fatalf("error opening connection: %s", err.Error())
	}
	defer discord.Close()
	fmt.Printf("[%s] The onus has fallen onto me. Started on API version %s\n", time.Now().Format(time.TimeOnly), discordgo.APIVersion)
	// finally, enabling imagick for image processing
	imagick.Initialize()
	defer imagick.Terminate()

	inst = &Instance{Session: discord, Reminder: &ReminderManager{Reminders: make([]Reminder, 0)}, ErrorChan: make(chan error)}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)

	defer panicMsg(discord)

	for {
		<-ticker.C
		select {
		case text := <-ch:
			// process the input asynchronously
			go func() {
				cmd := strings.Split(text, " ")[0]
				for _, c := range cliCmds {
					if cmd == c.Alias {
						err = c.Func(discord, text)
						if err != nil {
							log.Println(err.Error())
						}
					}
				}
				retch <- 0
			}()
			if <-retch == 1 {
				return
			}
		case <-sc:
			return
		default:
		}
		iterateReminders(inst)
	}
}

func iterateReminders (inst *Instance) {
	for i := 0; i < len(inst.Reminder.Reminders); i++ {
		select {
		case <-inst.Reminder.Reminders[i].GetTimer().C:
			defer Remind(inst, &(inst.Reminder.Reminders[i]))
		default:
		}
	}
}

func panicMsg(s *discordgo.Session) {
	a := recover()
	if a != nil {
		s.ChannelMessageSend("1331332284372222074", "<@479126092330827777> FATAL CRASH: "+a.(error).Error())
		panic(a.(error))
	}
}
