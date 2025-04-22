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

	// Allow me some time to formulate my thoughts if you will, it's been a while since I've engaged in any discussion such as this.
	// If you do not wish to talk, you may simply close the pull request and I will leave it at that.

	// This will be a long rant of sorts, so please read through it before any conclusions are made.
	// My logic tells me you have a right to be fragmented like this over something precious like life, life is a gift, I hope we can both agree on this.
	// But in my heart of hearts, I cannot rightly say I believe that there is true justice in this world, even with the sanctity of life.
	// If the survival of one, endangers hundreds of others, are we not downplaying the sanctity of the hundreds, just to save the one?
	// It's not my place to make this judgement, perhaps I am misguided, like a prosphelyte simply doing their duty in the city.
	// Yet I cannot stand by and let it happen, my heart forces me to make a move, and that move as drastic as it may be, has proven to safe lives.
	// Life is a gift, one I believe should be taken if their life were to endager another, but I also empathize with you, strongly
	// I never make these decisions lightly, and inherently that makes them all the more terrifying for others, since I am relatively set in my ways.
	// I would like to say that, none of this applies to what I will say next.
	// Snail.
	// I am sorry, you did not deserve to be broken like this, and I sincerely wish nothing but the best for you.
	// I feel as though I am condescending even though I wish to be frank with you, and for that, I again, apologize deeply.
	// You mean a whole deal to me, and a whole deal to the server as a whole, as you take care of us, we wish to take care of you.
	// I am ill fit for this purpose, as I seem to have only brought trouble onto you, but I wish to still be atleast on positive terms with you.
	// You are a dear friend, and an even kinder soul, please do not let your light snuff out.

	// i see
	// perhaps it is me who is too weak to bear the weight of justice's scales
	// but i myself cannot deem myself worthy of such enactments
	// maybe i am of the "trembling beasts", though i never considered myself to be a "right bearer" in the first place
	// ...funny how i'm citing crime and punishment of all things
	// lives are simply too heavy to weigh out for me - numbers don't matter
	// i do not think i am capable of enacting any justice to a satisfying degree, and project that notion onto others
	// most likely foolishly and for worse
	// but i have heard you
	// your explanation... it gave partial peace of mind
	// the initial shock of your original words is thankfully long gone, do not worry about that
	// i can bear the sight of violence again
	// now i am confident i will return, but not immediately
	// thank you for reaching out despite my withdrawal
	// and sorry for not noticing your message sooner

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
