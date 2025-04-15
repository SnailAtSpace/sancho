package funcs

import (
	cryptorand "crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Roll(i *Instance, m *discordgo.MessageCreate) {
	c, err := i.Session.State.Channel(m.ChannelID)
	if err != nil || strings.Contains(m.Content, "bread") {
		return
	}
	if roll, err := composeRoll(m.Content); roll == "" || err != nil{
		i.Session.ChannelMessageSendReply(c.ID, roll, m.SoftReference())
		i.ErrorChan <- err
	} else {
		i.Session.ChannelMessageSendReply(c.ID, "Your roll is "+roll+".", m.SoftReference())
	}
}

func EditRoll(i *Instance, m *discordgo.MessageUpdate, mymsg *discordgo.Message) {
	c, err := i.Session.State.Channel(m.ChannelID)
	if err != nil || strings.Contains(m.Content, "bread") {
		return
	}
	if roll, err := composeRoll(m.Content); roll == "" || err != nil{
		i.Session.ChannelMessageEdit(c.ID, mymsg.ID, "I know what you are.")
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		i.Session.ChannelMessageEdit(c.ID, mymsg.ID, "Your roll is "+roll+".")
	}
}

func composeRoll(i string) (string, error) {
	var mod int
	r, _ := strings.CutPrefix(i, ".roll ")
	if idk, err := strconv.Atoi(r); err == nil && idk > 0 {
		num, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(idk)))
		return strconv.Itoa(int(num.Int64() + 1)), nil
	} else {
		rest := r
		if !strings.Contains(rest, "d") {
			return "", fmt.Errorf("no d found! dummy")
		}

		rest = "+" + rest // magic 2
		operators := "+-*^_"
		rawStr := ""
		sum := 0

		for len(rest) > 0 {
			// the following code is magic! don't touch

			modRune := strings.IndexAny(rest[1:], operators) + 1 // next rune = end of our bit

			if modRune == 0 {
				modRune = len(rest)
			}

			bit := rest[1:modRune] // this killed sancho
			sign := rest[0]
			fmt.Println(bit)
			if strings.Contains(bit, "d") {
				mod = 0
				countS, maxS, _ := strings.Cut(bit, "d")
				var count int
				if len(countS) > 0 {
					count, err = strconv.Atoi(countS)
					if err != nil {
						return "", fmt.Errorf("\"%s\" is not a valid number for dice count", countS)
					}
				} else {
					count = 1
				}
				max, err := strconv.Atoi(maxS)
				if err != nil {
					return "", fmt.Errorf("\"%s\" is not a valid number for dice maximum", maxS)
				}

				if rawStr != "" {
					rawStr += "| "
				}

				for i := 0; i < count; i++ {
					vbig, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64((max))))
					v := int(vbig.Int64()) + 1
					rawStr += strconv.Itoa(v) + " "
					mod += v
				}
				fmt.Printf("count: %d, max: %d, total: %d\n", count, max, mod)
			} else {
				mod, err = strconv.Atoi(bit)
				if err != nil {
					return "", fmt.Errorf("\"%s\" is not a valid number for modifier", bit)
				}
			}

			switch sign {
			case '+':
				sum += mod
			case '*':
				sum *= mod
			case '-':
				sum -= mod
			case '^':
				sum = int(math.Pow(float64(sum), float64(mod)))
			case '_':
				sum = int(math.Pow(float64(mod), float64(sum)))
			}
			rest = rest[modRune:]
		}

		out := strconv.Itoa(sum) + " (" + rawStr[:len(rawStr)-1] + ")"
		outFull := "Your roll is " + out + "."
		if len(outFull) > 2000 {
			return strconv.Itoa(sum), nil
		} else {
			return out, nil
		}
	}
}
