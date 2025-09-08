package funcs

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
)

func EnderApologyReaction(inst *Instance, m *discordgo.Message) {
	// rawLines, err := os.ReadFile("apologylines.txt")
	// if err != nil {
	// 	inst.ErrorChan <- err
	// 	return
	// }
	// lines := strings.Split(string(rawLines), "\r\n")

	// //index, err := rand.Int(rand.Reader, big.NewInt(int64(len(lines))))
	// if err != nil {
	// 	inst.ErrorChan <- err
	// }

	inifile, err := ini.Load("randombullshit.ini")
	if err != nil {
		inst.ErrorChan <- err
		return
	}

	times := inifile.Section("").Key("endercount").MustInt()
	times++

	inifile.Section("").Key("endercount").SetValue(strconv.Itoa(times))
	inifile.SaveTo("randombullshit.ini")

	// inst.Session.ChannelMessageSendReply(m.ChannelID, lines[int(index.Int64())]+"\n*Times Ender has apologized: "+strconv.Itoa(times)+"*", m.SoftReference())
}

func DisplaySorryCount(inst *Instance, m *discordgo.MessageCreate) {
	
	inifile, err := ini.Load("randombullshit.ini")
	if err != nil {
		inst.ErrorChan <- err
		return
	}

	times := inifile.Section("").Key("endercount").MustInt()

	inst.Session.ChannelMessageSendReply(m.ChannelID, "Times Ender has apologized: "+strconv.Itoa(times), m.SoftReference())
}