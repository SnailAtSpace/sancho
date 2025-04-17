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

var echoChan string = ""
var retch = make(chan int) // i hate this. i hate this with all my heart. but i cannot bring myself to do any better.

var inst *Instance
var status bool = false

var greedID, myID, femmoID, whoopsID, mattagerID, nachoBowl, enderID string
var badChannels []string

func main() {
	listenChannelID = ""
	fmt.Printf("[%s] I shall pronounce the bot started.\n", time.Now().Format(time.TimeOnly))

	rawSecrets, err := os.ReadFile("secrets.txt")
	if err != nil {
		log.Fatalf("no secrets file: %s", err.Error())
	}
	secrets := strings.Split(string(rawSecrets), "\r\n") // fucking stupid why 
	
	auth_token := secrets[0]
	nachoBowl = secrets[1]
	mattagerID = secrets[2]
	whoopsID = secrets[3]
	myID = secrets[4]
	femmoID = secrets[5]
	greedID = secrets[6]
	enderID = secrets[7]
	badChannels = strings.Split(secrets[8], " ")

	echoChan = nachoBowl
	listenChannelID = nachoBowl

	// firing up the discord session
	discord, err := discordgo.New("Bot " + auth_token)

	if err != nil {
		log.Fatalf("couldn't initialize discord session: %s", err.Error())
	}

	inst = &Instance{Session: discord, RManager: &ReminderManager{Reminders: make([]Reminder, 0)}, ErrorChan: make(chan error)}

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
	discord.AddHandler(presenceUpdate)

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
	femmoTicker := time.NewTicker(24 * time.Hour)
	femmoTimes := 0

	defer panicMsg(discord)

	for {
		<-ticker.C
		if status {
			femmoTicker.Reset(24 * time.Hour)
			femmoTimes = 0
		}
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
		case err = <-inst.ErrorChan:
			if err != nil {
				log.Println(err.Error())
			}
		case <-femmoTicker.C:
			femmoTimes++
			greed, err := inst.Session.UserChannelCreate(greedID)
			if err != nil {
				log.Println(err.Error())
				break
			}
			inst.Session.ChannelMessageSend(greed.ID, fmt.Sprintf("Breadfemmo has been offline for %d consecutive hours.", 24*femmoTimes))
		default:
		}
		iterateReminders(inst)
	}
}

func iterateReminders(inst *Instance) {
	for i := 0; i < len(inst.RManager.Reminders); i++ {
		select {
		case <-inst.RManager.Reminders[i].GetTimer().C:
			defer Remind(inst, &(inst.RManager.Reminders[i]))
		default:
		}
	}
}

func panicMsg(s *discordgo.Session) {
	a := recover()
	if a != nil {
		s.ChannelMessageSend(nachoBowl, fmt.Sprintf("<@%s> FATAL CRASH: "+a.(error).Error(), myID))
		panic(a.(error))
	}
}

func presenceUpdate(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	if m.User.ID == femmoID {
		if m.Status == discordgo.StatusOffline && status {
			status = false
			fmt.Println(status)
		}
		if m.Status != discordgo.StatusOffline && !status {
			status = true
			fmt.Println(status)
		}
	}
}
