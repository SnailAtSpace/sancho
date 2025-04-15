package funcs

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func EnderApologyReaction(inst *Instance, m *discordgo.MessageCreate){
	var envname string = "SanchoEnderApologiesCounter"

	rawLines, err := os.ReadFile("apologylines.txt")
	if err != nil {
		inst.ErrorChan <- err
		return
	}
	lines := strings.Split(string(rawLines), "\r\n")

	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(lines))))
	if err != nil {
		inst.ErrorChan <- err
	}

	env := os.Getenv(envname)
	if env == "" {
		env = "0"
	}
	times, _ := strconv.Atoi(env)
	times++

	os.Setenv(envname, strconv.Itoa(times))

	inst.Session.ChannelMessageSendReply(m.ChannelID, lines[int(index.Int64())]+"\nTimes Ender has apologized: "+strconv.Itoa(times), m.SoftReference())
}