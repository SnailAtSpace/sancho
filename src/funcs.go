package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/k4s/phantomgo"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var err error

// call if something goes wrong
func sadness(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m != nil{
		s.ChannelMessageSendReply(m.ChannelID, "Sorry, my creator must have fucked something up.\nPlease pierce him with a sanguine lance and drink his blood.", m.Reference())
	}
	fmt.Println(err)
}

func help(s *discordgo.Session, m *discordgo.MessageCreate){
	ht, err := os.ReadFile("help.md")
	if err!=nil{
		sadness(s,m)
	}
	s.ChannelMessageSendReply(m.ChannelID, 
		string(ht),
		m.Reference())
}

func roll(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil || strings.Contains(m.Content, "bread"){
		return
	}
	if roll := composeRoll(m.Content); roll == "I know what you are."{
		s.ChannelMessageSendReply(c.ID, roll, m.SoftReference())
	} else {
		s.ChannelMessageSendReply(c.ID, "Your roll is "+ roll+".", m.SoftReference())
	}
}


func editRoll(s *discordgo.Session, m *discordgo.MessageUpdate, mymsg *discordgo.Message){
	c, err := s.State.Channel(m.ChannelID)
	if err != nil || strings.Contains(m.Content, "bread"){
		return
	}
	if roll := composeRoll(m.Content); roll == ""{
		s.ChannelMessageEdit(c.ID, mymsg.ID, "I know what you are.")
	} else {
		s.ChannelMessageEdit(c.ID, mymsg.ID, "Your roll is "+ roll +".")
	}
}

func composeRoll(i string) string{
	var count, max, mod int
	r, _ := strings.CutPrefix(i, ".roll ")
	if idk, err := strconv.Atoi(r); err == nil && idk > 0 {
		count = 1
		max, _ = strconv.Atoi(r)
		mod = 0
		return strconv.Itoa(int(rand.Int63n(int64(max))+1))
	} else {
		countS, rest, found := strings.Cut(r, "d")
		if !found {

			return ""
		}

		if countS == "" {
			count = 1
		} else {
			count, err = strconv.Atoi(countS)
			if err != nil {
				return ""
			}
		}

		modRune := strings.IndexAny(rest, "+-*^_")
		if modRune == -1 {
			modRune = len(rest)
		}

		max, err = strconv.Atoi(rest[:modRune])
		if err != nil {
			return ""
		}

		if max < 1 || count < 1 {
			fmt.Println(max, count)
			return ""
		}

		rawStr := ""
		sum := 0
		for i := 0; i < count; i++ {
			v := rand.Intn(max) + 1
			rawStr += strconv.Itoa(v) + " "
			sum += v
		}

		for modRune < len(rest) {
			// the following code is magic! don't touch
			rest = rest[modRune:]
			sign := rest[0]
			modRune = strings.IndexAny(rest[1:], "+-*^_") + 1
			if modRune == 0 {
				modRune = len(rest)
			}
			mod, err = strconv.Atoi(rest[1:modRune])
			if err != nil {
				return ""
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

		}

		out := strconv.Itoa(sum) + " (" + rawStr[:len(rawStr)-1] + ")"
		outFull := "Your roll is " + out + "."
		if len(outFull) > 2000 {
			return strconv.Itoa(sum)
		} else {return out}
	}
}

// call if user fucked up command use
func iKnowWhatYouAre(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, "I know what you are.", m.Reference())
}

func bod(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	resp, err := http.Get("https://tiphereth.zasz.su/static/assets/cards/Shi1.png")
	if err != nil {
		sadness(s, m)
	}
	defer resp.Body.Close()
	yujin := resp.Body

	resp, err = http.Get("https://tiphereth.zasz.su/static/assets/cards/Roland4Phase_Yujin.png")
	if err != nil {
		sadness(s, m)
	}
	defer resp.Body.Close()
	yujinDead := resp.Body

	pingId, _ := strings.CutPrefix(m.Content, ".bod ")
	//pingId = pingId[1:len(pingId)-1]
	ping := ""
	if _, err := strconv.Atoi(pingId[2 : len(pingId)-1]); pingId[1] == '@' && err == nil {
		ping = pingId
	}
	roll := rand.Intn(4) + 1
	if roll == 4 {
		s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
			Content: "**4**" + " " + ping,
			Files: []*discordgo.File{
				{
					Name:   "yujin.png",
					Reader: yujin,
				},
			},
		})
	} else {
		s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
			Content: strconv.Itoa(roll),
			Files: []*discordgo.File{
				{
					Name:   "yujinDead.jpg",
					Reader: yujinDead,
				},
			},
		})
	}
}

func conceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, "What... is it this time?", m.Reference())
}

func fut(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages := []string{"...I won't even have to call Father for this.",
		"...Hold your tongue; I will no longer tolerate any more \"ingenious ideas\".",
		"You're nothing before a Second Kindred, let alone Father.",
		"...Did you learn that from that Knight? How humorous.",
	}
	s.ChannelMessageSendReply(m.ChannelID, messages[rand.Intn(len(messages))], m.Reference())
}

func sendimg(s *discordgo.Session, m *discordgo.MessageCreate, name string){
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	img, err := os.Open("img/"+name)
	if err!=nil{
		sadness(s,m)
	}
	defer img.Close()
	s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   name,
				Reader: img,
			},
		},
	})
}

func jpegify(s *discordgo.Session, m *discordgo.MessageCreate, orb *imagick.MagickWand, quality int) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	var resp *http.Response
	if len(m.Attachments) == 0 || !strings.Contains(m.Attachments[0].ContentType, "image") {
		if m.ReferencedMessage == nil {
			s.ChannelMessageSendReply(m.ChannelID, "Please send an actual image.", m.Reference())
			return
		}
		if len(m.ReferencedMessage.Attachments) == 0 || !strings.Contains(m.ReferencedMessage.Attachments[0].ContentType, "image") {
			s.ChannelMessageSendReply(m.ChannelID, "Please send an actual image.", m.Reference())
			return
		}
		resp, err = http.Get(m.ReferencedMessage.Attachments[0].URL)
	} else {
		resp, err = http.Get(m.Attachments[0].URL)
	}

	//now we pipe the first (for now) image into imagemagick and wait for the result - how?

	//fn := m.Attachments[0].Filename
	//ext := fn[strings.LastIndex(fn, ".")+1:]
	if err != nil {
		fmt.Println("couldn't get image from internet")
		sadness(s, m)
		return
	}
	defer resp.Body.Close()
	orig, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("couldn't extract from html")
		sadness(s, m)
		return
	}

	err = orb.ReadImageBlob(orig)
	defer orb.Clear()
	if err != nil {
		fmt.Println("couldn't read img into orb", err)
		sadness(s, m)
		return
	}
	orb.SetImageFormat("JPEG")
	orb.SetImageCompressionQuality(uint(quality))
	orb.SetCompressionQuality(uint(quality))
	if quality < 2 {
		x, y := orb.GetImageWidth(), orb.GetImageHeight()
		scalingFactor := math.Max(float64(x/160), float64(y/100))
		orb.ModulateImage(100, 135, 100)
		orb.ResizeImage(uint(float64(x)/scalingFactor), uint(float64(y)/scalingFactor), imagick.FILTER_BOX)
		out, _ := orb.GetImageBlob()
		orb.Clear()
		orb.ReadImageBlob(out)
		orb.ResizeImage(x, y, imagick.FILTER_BOX)
		orb.SetImageFormat("JPEG")
		orb.PosterizeImage(16, imagick.DITHER_METHOD_FLOYD_STEINBERG)
		orb.SetImageCompressionQuality(uint(quality))
		orb.SetCompressionQuality(uint(quality))
	}
	out, err := orb.GetImageBlob()
	if err != nil {
		fmt.Println("couldn't shove it back in")
		sadness(s, m)
		return
	}
	outReader := bytes.NewReader(out)
	s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Files: []*discordgo.File{
			{
				Name:   "img.jpg",
				Reader: outReader,
			},
		},
	})
}

// IN DEVELOPMENT
func getPrescript(s *discordgo.Session, m *discordgo.MessageCreate){
	prescript, err := browser.Download(&phantomgo.Param{
		Method: "GET",
		Url: "https://nyos.dev/prescript",
		Header:       http.Header{"Cookie": []string{"your cookies"}},
		UsePhantomJS: true,
	})
	if err!=nil{
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(prescript.Body)
	for scanner.Scan(){
		fmt.Println(scanner.Text())
	}
}