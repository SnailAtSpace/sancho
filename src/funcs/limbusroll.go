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

	// idk anymore
	// i am nauseous and lightheaded
	// i'm not afraid of you - cautious would be the right term
	// i know how dangerous information can be just as much as you do
	// and the only thing that stopped others in my mind from using it against me was the inherent taboo of doing so
	// but you have already crossed that line long ago, alone on its opposite side
	// therefore, i am unable to trust you at all, and there is no way to mend this
	// i am, alertness aside, dejected
	// a human life is sacred to me - any, no matter the beholder
	// to know a person who killed is to also imagine yourself in such a situation, something i cannot fathom
	// does this make me a bad person? being unable to forgo my morals in order to simply accept?
	// if so... i am truly, genuinely sorry

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

	out := fmt.Sprintf("**You rolled: %d **(%s)\nRaw damage: %d", basePower, rawStr[:len(rawStr)-1], rawDamage)
	if len(out) > 2000 {
		out = fmt.Sprintf("**You rolled: %d**\nRaw damage: %d", basePower, rawDamage)
		if len(out) > 2000 {
			out = fmt.Sprintf("**You rolled: %d**", basePower)
			if len(out) > 2000 {
				out = "Sorry, your roll is literally too big to display."
			}
		}
	}
	inst.Session.ChannelMessageSendReply(m.ChannelID, out, m.SoftReference())
}
