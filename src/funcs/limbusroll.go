package funcs

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func LimbusRoll(inst *Instance, m *discordgo.MessageCreate) {
	var coins, basePower, coinPower int

	cmd := strings.Split(m.Content, " ")[1:]

	if len(cmd) < 3 || len(cmd) > 4 {
		inst.Session.ChannelMessageSendReply(m.ChannelID, "Incorrect roll formatting (<# of coins> <base power> <coin power> [sp])", m.SoftReference())
		return
	}

	coins, err := strconv.Atoi(cmd[0])
	if err == nil {
		basePower, err = strconv.Atoi(cmd[1])
		if err == nil {
			coinPower, err = strconv.Atoi(cmd[2])
		}
	}
	if err != nil {
		inst.Session.ChannelMessageSendReply(m.ChannelID, "Incorrect roll formatting (<# of coins> <base power> <coin power> [sp])", m.SoftReference())
		return
	}

	sp := 0
	if len(cmd) == 4 {
		sp, err = strconv.Atoi(strings.Replace(cmd[3], "sp", "", 1))
		if err != nil {
			inst.Session.ChannelMessageSendReply(m.ChannelID, "Incorrect roll formatting (<# of coins> <base power> <coin power> [sp])", m.SoftReference())
			return
		}
	}

	rawStr := ""
	rawDamage := 0

	for i := 0; i < coins; i++ {
		rollB, _ := rand.Int(rand.Reader, big.NewInt(int64(100)))
		roll := int(rollB.Int64())
		if roll < 50+sp {
			basePower += coinPower
			rawStr += ":yellow_circle: "
		} else {
			rawStr += ":brown_circle: "
		}
		rawDamage += basePower
	}

	if basePower<0 {
		basePower = 0
	}

	out := fmt.Sprintf("You rolled: %d (%s)\nRaw damage: %d", basePower, rawStr[:len(rawStr)-1], rawDamage)
	inst.Session.ChannelMessageSendReply(m.ChannelID, out, m.SoftReference())
}
