package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// THIS IS HELL!!!!!!!!!!!

type Reminder struct {
	end, start time.Time
	message string
	author string
	target string
	request *discordgo.Message
	rqid string
	timer *time.Timer
	repeats int
	period int
}

type Token struct {
	key string
	txt string
	i int
}

func setReminder(s *discordgo.Session, m *discordgo.MessageCreate, t *[]Reminder) {
	// first we need to parse time, then we set a timer and record it (in case the bot goes out)

	rawCmd, _ := strings.CutPrefix(m.Content, ".remind")
	// if !slices.Contains(strings.Split(rawCmd, " "), "to"){
	// 	s.ChannelMessageSendReply(m.ChannelID, "I know what you are. (please use \"to\" at the beginning of your reminder text, e.g. .remind in 2h to do laundry)", m.Reference())
	// 	return
	// }
	cmd := strings.Split(strings.TrimSpace(rawCmd), " ")

	tokens := []string{"to", "in", "at", "on", "every", "forever", "times"}
	tokenIndicesMap := make(map[string][]int)
	tokenIndicesOrdered := make([]Token, 0)

	for i, t := range cmd {
		if slices.Contains(tokens, t) {
			tokenIndicesMap[t] = append(tokenIndicesMap[t], i)
			tokenIndicesOrdered = append(tokenIndicesOrdered, Token{t, "", i})
		} else if len(tokenIndicesOrdered)>0{
			tokenIndicesOrdered[len(tokenIndicesOrdered)-1].txt += t + " "
		}
	}
	
	fmt.Println(tokenIndicesMap, "\n", tokenIndicesOrdered)
	if len(tokenIndicesOrdered)==0 {
		s.ChannelMessageSendReply(m.ChannelID, "Invalid reminder formatting: try again! (you probably forgot the time interval)", m.Reference())
		return
	}

	timeInUnix := 0
	timeSet := false
	dateSet := false
	isRelative := false
	foreverFlag := false
	repeats := 1
	repeatTime := 0
	msg := strings.Join(cmd[:tokenIndicesOrdered[0].i], " ")+" "
	msg, _ = strings.CutPrefix(msg, "me ")
	targetUser := m.Author.ID
	for _, t := range tokenIndicesOrdered {
		switch t.key {
			case "in":
				if timeSet || dateSet{
					msg+=t.key+" "+t.txt
					continue
				}
				text := strings.Split(strings.TrimSpace(t.txt), " ")
				eot := false
				totalTime := 0
				for len(text) > 0 && !eot && !timeSet{
					timeIncr, err := strconv.Atoi(text[0])
					if err != nil {
						if text[0] == "a" || text[0] == "an" {
							timeIncr = 1
						} else if strings.Contains("smhdwyc", string(text[0][len(text[0])-1])){
							time, thing := text[0][:len(text[0])-1], text[0][len(text[0])-1]
							timeIncr, err = strconv.Atoi(time)
							if err != nil {
								break
							}
							together := append([]string {time}, string(thing))
							text = append(together, text[1:]...)
						} else {
							break
						}
					}
					switch strings.ToLower(text[1]) {
						case "s", "seconds", "second", "sec":
							totalTime += timeIncr * 1
						case "m", "minutes", "minute", "min":
							totalTime += timeIncr * 60
						case "h", "hours", "hour":
							totalTime += timeIncr * 60 * 60
						case "d", "days", "day":
							totalTime += int(time.Now().AddDate(0,0,timeIncr).Unix())-int(time.Now().Unix())
						case "w", "weeks", "week":
							totalTime += int(time.Now().AddDate(0,0,timeIncr*7).Unix())-int(time.Now().Unix())
						case "months", "month":
							totalTime += int(time.Now().AddDate(0,timeIncr,0).Unix())-int(time.Now().Unix())
						case "y", "years", "year":
							totalTime += int(time.Now().AddDate(timeIncr,0,0).Unix())-int(time.Now().Unix())
						case "c", "centuries", "century":
							totalTime += int(time.Now().AddDate(timeIncr*100,0,0).Unix())-int(time.Now().Unix())
						default:
							eot = true
					}
					if !eot {
						text = text[2:]
					}
				}
				if totalTime == 0 {
					continue
				}
				timeInUnix = int(time.Now().Unix()) + totalTime
				msg += strings.Join(text, " ")
				timeSet = true
				dateSet = true
				isRelative = true
			case "at":
				if isRelative || timeSet {
					msg+=t.key+" "+t.txt
					continue
				}
				var sec, min, hour int
				rawt := strings.Split(strings.TrimSpace(t.txt), ":")
				var cl []int
				for _, r := range rawt {
					val, err := strconv.Atoi(r)
					if err != nil {
						iKnowWhatYouAre(s,m)
						return
					}
					cl = append(cl, val)
				}
				if len(cl) == 2 {
					sec = 0
				} else if len(cl) == 3 {sec = cl[2]} else {iKnowWhatYouAre(s,m); return}
				min = cl[1]
				hour = cl[0]
				timeInUnix += hour*60*60 + min*60 + sec
				timeSet = true
			case "on":
				if isRelative || dateSet {
					msg+=t.key+" "+t.txt
					continue
				}
				var day, month, year int
				rawt := strings.Split(strings.TrimSpace(t.txt), "/")
				if rawt[0] == t.txt {
					rawt = strings.Split(strings.TrimSpace(t.txt), ".")
					if rawt[0] == t.txt {
						rawt = strings.Split(strings.TrimSpace(t.txt), "-")
						if rawt[0] == t.txt {
							iKnowWhatYouAre(s,m)
							return
						}
					}
				}
				var cl []int
				for _, r := range rawt {
					val, err := strconv.Atoi(r)
					if err != nil {
						iKnowWhatYouAre(s,m)
						return
					}
					cl = append(cl, val)
				}
				if len(cl) == 2 {
					year = time.Now().Year()
				} else if len(cl) == 3 {year = cl[2]} else {iKnowWhatYouAre(s,m); return}
				if year%100 == year {
					year += (time.Now().Year()/100)*100
				}
				month = cl[1]
				day = cl[0]
				timeInUnix = int(time.Unix(int64(timeInUnix),0).AddDate(year,month,day).Unix())
				dateSet = true
			case "every":
				if repeatTime != 0 {
					msg+=t.key+" "+t.txt
					continue
				}
				text := strings.Split(strings.TrimSpace(t.txt), " ")
				eot := false
				for len(text) > 1 && !eot {
					timeIncr, err := strconv.Atoi(text[0])
					if err != nil {
						if strings.Contains("smhdwyc", string(text[0][len(text[0])-1])){
							time, thing := text[0][:len(text[0])-1], text[0][len(text[0])-1]
							together := append([]string {time}, string(thing))
							text = append(together, text[1:]...)
							timeIncr, err = strconv.Atoi(text[0])
							if err != nil {
								iKnowWhatYouAre(s,m)
								return
							}
						} else {
							timeIncr = 1
						}
					}
					switch strings.ToLower(text[1]) {
					case "s", "seconds", "second", "sec":
						repeatTime += timeIncr * 1
					case "m", "minutes", "minute", "min":
						repeatTime += timeIncr * 60
					case "h", "hours", "hour":
						repeatTime += timeIncr * 60 * 60
					case "d", "days", "day":
						repeatTime += int(time.Now().AddDate(0,0,timeIncr).Unix())-int(time.Now().Unix())
					case "w", "weeks", "week":
						repeatTime += int(time.Now().AddDate(0,0,timeIncr*7).Unix())-int(time.Now().Unix())
					case "months", "month":
						repeatTime += int(time.Now().AddDate(0,timeIncr,0).Unix())-int(time.Now().Unix())
					case "y", "years", "year":
						repeatTime += int(time.Now().AddDate(timeIncr,0,0).Unix())-int(time.Now().Unix())
					case "c", "centuries", "century":
						repeatTime += int(time.Now().AddDate(timeIncr*100,0,0).Unix())-int(time.Now().Unix())
					default:
						eot = true
					}
					if !eot {
						text = text[2:]
					}
				}
				if repeatTime == 0 {
					continue
				}
				msg += strings.Join(text, " ")
			case "times":
				if repeats != 1{
					msg+=t.key+" "+t.txt
					continue
				}
				timesHopefully := cmd[t.i-1]
				if times, err := strconv.Atoi(timesHopefully); err == nil {
					repeats = times
					msg = msg[0:len(msg)-len(timesHopefully)]
					msg += t.txt
				} else {
					msg+=t.key+" "+t.txt
					continue
				}
			case "forever":
				foreverFlag = true
				// this one's a doozy
			default:
				msg+=t.key+" "+t.txt
		}
	}

	if !timeSet && dateSet {
		hrs, min, sec := time.Now().Clock()
		timeInUnix += hrs*60*60 + min*60 + sec
	}
	if timeSet && !dateSet{
		year, month, day := time.Now().Date()
		timeInUnix = int(time.Unix(int64(timeInUnix),0).AddDate(year,int(month),day).Unix())
	}

	if foreverFlag && repeatTime != 0 && repeats == 1 {
		repeats = -1
	}

	if aoaoao, rest, _ := strings.Cut(msg, " "); aoaoao == "to"{
		msg = rest
	}

	// TODO: timezones!!
	if !isRelative {
		tzfile, err := os.ReadFile("timezones.txt")
		if err != nil {
			sadness(s,m,err)
			return
		}
		
		timezones := strings.Split(string(tzfile), "\n")
		tzstring := ""
		for _, filetz := range timezones {
			if strings.Split(filetz, " ")[0] == targetUser {
				tzstring = strings.Split(filetz, " ")[1]
			}
		}

		if tzstring == ""{
			s.ChannelMessageSendReply(m.ChannelID, "You tried to set an absolute reminder, but haven't specified a timezone! To do that, use .settz.", m.SoftReference())
			return
		} 

		tzloc, err := time.LoadLocation(tzstring)
		if err!=nil {
			s.ChannelMessageSendReply(m.ChannelID, "ping snail: "+err.Error(), m.SoftReference())
			return
		}

		timeParsed, err := time.ParseInLocation(time.DateTime, time.Unix(int64(timeInUnix),0).Format(time.DateTime), tzloc)
		if err!=nil {
			s.ChannelMessageSendReply(m.ChannelID, "ping snail: "+err.Error(), m.SoftReference())
			return
		}

		timeInUnix = int(timeParsed.Unix())
	}


	if tagIndex := slices.IndexFunc(cmd[:tokenIndicesOrdered[0].i], func(str string) bool {return strings.HasPrefix(str, "<@")}); tagIndex >= 0 {
		targetUser = cmd[tagIndex][2:len(cmd[tagIndex])-1]
		msg = strings.Replace(msg, "<@"+targetUser+">", "",-1)
	}

	msg, _ = strings.CutPrefix(strings.TrimSpace(msg), "to ")
	
	timerFile, err := os.OpenFile("timers.txt", os.O_APPEND, 0666)
	if err != nil {
		sadness(s,m,err)
		return
	}
	defer timerFile.Close()

	timerFile.WriteString(strings.Join([]string{m.Message.ID, strconv.Itoa(timeInUnix), strconv.Itoa(int(time.Now().Unix())), targetUser, m.Message.ChannelID, m.Author.ID, strconv.Itoa(repeats-1), strconv.Itoa(repeatTime), msg}, " ") + "\n")
	*t = append(*t, Reminder{
		end: time.Unix(int64(timeInUnix),0),
		start: time.Now(),
		message: msg,
		author: m.Author.ID,
		target: targetUser,
		request: m.Message,
		rqid: m.Message.ID,
		timer: time.NewTimer(time.Duration(timeInUnix-int(time.Now().Unix())) * time.Second),
		repeats: repeats-1,
		period: repeatTime,
	})
	s.ChannelMessageSendReply(m.ChannelID, "...As you wish. I shall send a reminder at <t:"+strconv.Itoa(timeInUnix)+">.", m.Reference())

	// msg = strings.Join(cmd[slices.Index(cmd,"to")+1:], " ")
	// cmd = cmd[:slices.Index(cmd,"to")]
	// targetUser := m.Author.ID
	// if len(cmd) == 0 || len(msg) == 0 {
	// 	iKnowWhatYouAre(s,m)
	// 	return
	// }

	// if cmd[0] == "me" {
	// 	cmd = cmd[1:]
	// }
	// if cmd[0][0] == '<' {
	// 	if rawCmd[:2] != "me" && cmd[0][1] == '@' {
	// 		targetUser = cmd[0][2 : len(cmd[0])-1]
	// 		cmd = cmd[1:]
	// 	} else if cmd[0][1:3] != "t:"{
	// 		iKnowWhatYouAre(s,m)
	// 		return
	// 	}
	// }

	// timeInUnix = 0
	// eot := false
	// repeats = 1
	// repeatTime = 0
	// if strings.Contains(strings.Join(cmd, " "), "<t:"){
	// 	beg := strings.Index(rawCmd, "<t:")+3
	// 	timeInUnix,err = strconv.Atoi(rawCmd[beg:beg+strings.IndexAny(rawCmd[beg:], ">:")])
	// 	if err!=nil{
	// 		iKnowWhatYouAre(s,m)
	// 		return
	// 	}
	// } else if slices.Contains(cmd, "at") || slices.Contains(cmd, "on"){ // FIXME: add timezone support
	// 	var sec, min, hour, day, month, year int
	// 	if slices.Contains(cmd, "at"){
	// 		rawt := strings.Split(cmd[slices.Index(cmd,"at")+1], ":")
	// 		var cl []int
	// 		for _, r := range rawt {
	// 			val, err := strconv.Atoi(r)
	// 			if err != nil {
	// 				iKnowWhatYouAre(s,m)
	// 				return
	// 			}
	// 			cl = append(cl, val)
	// 		}
	// 		if len(cl) == 2 {
	// 			sec = 0
	// 		} else if len(cl) == 3 {sec = cl[2]} else {iKnowWhatYouAre(s,m); return}
	// 		min = cl[1]
	// 		hour = cl[0]
	// 	} else {
	// 		sec = time.Now().Second()
	// 		min = time.Now().Minute()
	// 		hour = time.Now().Hour()
	// 	}
	// 	if slices.Contains(cmd, "on"){
	// 		rawt := strings.Split(cmd[slices.Index(cmd,"on")+1], ".")
	// 		var cl []int
	// 		for _, r := range rawt {
	// 			val, err := strconv.Atoi(r)
	// 			if err != nil {
	// 				iKnowWhatYouAre(s,m)
	// 				return
	// 			}
	// 			cl = append(cl, val)
	// 		}
	// 		if len(cl) == 2 {
	// 			year = time.Now().Year()
	// 		} else if len(cl) == 3 {year = cl[2]} else {iKnowWhatYouAre(s,m); return}
	// 		if year%100 == year {
	// 			year += (time.Now().Year()/100)*100
	// 		}
	// 		month = cl[1]
	// 		day = cl[0]
	// 	} else {
	// 		year = time.Now().Year()
	// 		month = int(time.Now().Month())
	// 		day = time.Now().Day()
	// 	}
	// 	destTime := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Now().Location())
	// 	if destTime.Unix() <= time.Now().Unix() {
	// 		iKnowWhatYouAre(s,m)
	// 		return
	// 	}
	// 	timeInUnix = int(destTime.Unix())

	// 	if slices.Contains(cmd, "every"){
	// 		if slices.Contains(cmd, "times") || slices.Contains(cmd, "forever"){
	// 			cmd = cmd[slices.Index(cmd, "every")+1:]
	// 			cmd = cmd[:1+slices.IndexFunc(cmd, func (a string) bool {
	// 				if a == "times" || a == "forever" {
	// 					return true
	// 				}
	// 				return false
	// 			})]
	// 		} else {
	// 			cmd = cmd[slices.Index(cmd,"every")+1:]
	// 		}
	// 		for len(cmd) > 1 && !eot {
	// 			timeIncr, err := strconv.Atoi(cmd[0])
	// 			if err != nil {
	// 				if strings.Contains("smhdwyc", string(cmd[0][len(cmd[0])-1])){
	// 					time, thing := cmd[0][:len(cmd[0])-1], cmd[0][len(cmd[0])-1]
	// 					together := append([]string {time}, string(thing))
	// 					cmd = append(together, cmd[1:]...)
	// 					timeIncr, err = strconv.Atoi(cmd[0])
	// 					if err != nil {
	// 						iKnowWhatYouAre(s,m)
	// 						return
	// 					}
	// 				} else {
	// 					timeIncr = 1
	// 				}
	// 			}
	// 			switch strings.ToLower(cmd[1]) {
	// 			case "s", "seconds", "second", "sec":
	// 				repeatTime += timeIncr * 1
	// 			case "m", "minutes", "minute", "min":
	// 				repeatTime += timeIncr * 60
	// 			case "h", "hours", "hour":
	// 				repeatTime += timeIncr * 60 * 60
	// 			case "d", "days", "day":
	// 				repeatTime += int(time.Now().AddDate(0,0,timeIncr).Unix())-int(time.Now().Unix())
	// 			case "w", "weeks", "week":
	// 				repeatTime += int(time.Now().AddDate(0,0,timeIncr*7).Unix())-int(time.Now().Unix())
	// 			case "months", "month":
	// 				repeatTime += int(time.Now().AddDate(0,timeIncr,0).Unix())-int(time.Now().Unix())
	// 			case "y", "years", "year":
	// 				repeatTime += int(time.Now().AddDate(timeIncr,0,0).Unix())-int(time.Now().Unix())
	// 			case "c", "centuries", "century":
	// 				repeatTime += int(time.Now().AddDate(timeIncr*100,0,0).Unix())-int(time.Now().Unix())
	// 			default:
	// 				eot = true
	// 			}
	// 			if !eot {
	// 				cmd = cmd[2:]
	// 			}
	// 		}
	// 		if repeatTime == 0 {
	// 			iKnowWhatYouAre(s,m)
	// 			return
	// 		}

	// 		repeats = 0
	// 		if len(cmd) != 0 {
	// 			if cmd[0] != "forever"{
	// 				repeats, err = strconv.Atoi(cmd[0])
	// 				if err!=nil {
	// 					iKnowWhatYouAre(s,m)
	// 					return
	// 				}
	// 			}
	// 		} 
	// 	}
	// } else {
	// 	totalTime := 0
	// 	if cmd[0] == "in" {
	// 		cmd = cmd[1:]
	// 	}
	// 	for len(cmd) > 0 && !eot {
	// 		timeIncr, err := strconv.Atoi(cmd[0])
	// 		if err != nil {
	// 			if cmd[0] == "a" || cmd[0] == "an" {
	// 				timeIncr = 1
	// 			} else if strings.Contains("smhdwyc", string(cmd[0][len(cmd[0])-1])){
	// 				time, thing := cmd[0][:len(cmd[0])-1], cmd[0][len(cmd[0])-1]
	// 				together := append([]string {time}, string(thing))
	// 				cmd = append(together, cmd[1:]...)
	// 				timeIncr, err = strconv.Atoi(cmd[0])
	// 				if err != nil {
	// 					iKnowWhatYouAre(s,m)
	// 					return
	// 				}
	// 			} else {
	// 				break
	// 			}
	// 		}
	// 		switch strings.ToLower(cmd[1]) {
	// 		case "s", "seconds", "second", "sec":
	// 			totalTime += timeIncr * 1
	// 		case "m", "minutes", "minute", "min":
	// 			totalTime += timeIncr * 60
	// 		case "h", "hours", "hour":
	// 			totalTime += timeIncr * 60 * 60
	// 		case "d", "days", "day":
	// 			totalTime += int(time.Now().AddDate(0,0,timeIncr).Unix())-int(time.Now().Unix())
	// 		case "w", "weeks", "week":
	// 			totalTime += int(time.Now().AddDate(0,0,timeIncr*7).Unix())-int(time.Now().Unix())
	// 		case "months", "month":
	// 			totalTime += int(time.Now().AddDate(0,timeIncr,0).Unix())-int(time.Now().Unix())
	// 		case "y", "years", "year":
	// 			totalTime += int(time.Now().AddDate(timeIncr,0,0).Unix())-int(time.Now().Unix())
	// 		case "c", "centuries", "century":
	// 			totalTime += int(time.Now().AddDate(timeIncr*100,0,0).Unix())-int(time.Now().Unix())
	// 		default:
	// 			eot = true
	// 		}
	// 		if !eot {
	// 			cmd = cmd[2:]
	// 		}
	// 	}
	// 	if totalTime == 0 {
	// 		iKnowWhatYouAre(s, m)
	// 		return
	// 	}
	// 	timeInUnix = int(time.Now().Unix()) + totalTime
	// }

	// timerFile, err := os.OpenFile("timers.txt", os.O_APPEND, 0666)
	// if err != nil {
	// 	sadness(s,m,err)
	// 	return
	// }
	// defer timerFile.Close()

	// timerFile.WriteString(strings.Join([]string{m.Message.ID, strconv.Itoa(timeInUnix), strconv.Itoa(int(time.Now().Unix())), targetUser, m.Message.ChannelID, m.Author.ID, strconv.Itoa(repeats-1), strconv.Itoa(repeatTime), msg}, " ") + "\n")
	// *t = append(*t, Reminder{
	// 	end: time.Unix(int64(timeInUnix),0),
	// 	start: time.Now(),
	// 	message: msg,
	// 	author: m.Author.ID,
	// 	target: targetUser,
	// 	request: m.Message,
	// 	rqid: m.Message.ID,
	// 	timer: time.NewTimer(time.Duration(timeInUnix-int(time.Now().Unix())) * time.Second),
	// 	repeats: repeats-1,
	// 	period: repeatTime,
	// })
	// s.ChannelMessageSendReply(m.ChannelID, "...As you wish. I shall send a reminder at <t:"+strconv.Itoa(timeInUnix)+">.", m.Reference())
}


func remind(s *discordgo.Session, r *Reminder){
	r = &Reminder{
		end: r.end,
		start: r.start,
		message: r.message,
		author: r.author,
		target: r.target,
		request: r.request,
		rqid: r.rqid,
		timer: r.timer,
		repeats: r.repeats,
		period: r.period,
	}
	timerFile, err := os.OpenFile("timers.txt", os.O_RDONLY, 0666)
	if err != nil {
		sadness(s,nil,err)
		return
	}
	defer timerFile.Close()

	scanner := bufio.NewScanner(timerFile)

	request := r.request
	var reqChan string
	if request == nil { // if the original message was deleted and we have restored at least once...
		for scanner.Scan() {
			reminderText := strings.SplitN(scanner.Text(), " ", 6)
			if reminderText[0] == r.rqid {
				reqChan = reminderText[4]
			}
		}
		if reqChan != "" {
			if len("<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)") >2000 {
				_, err = s.ChannelMessageSend(reqChan, "<@"+r.target+">: YOUR MESSAGE DIDN'T FIT")
			} else {
				_, err = s.ChannelMessageSend(reqChan, "<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)")
			}
			if err != nil {
				sadness(s,nil,err)
				return
			}
		} else{
			sadness(s,nil,err)
			print(r.message)
			//SHIT WE DIDN'T FIND IT ABORT ABORT
			return
		}
	} else {
		if len("<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)") >2000 {
			_, err = s.ChannelMessageSendReply(r.request.ChannelID, "<@"+r.target+">: YOUR MESSAGE DIDN'T FIT", r.request.SoftReference())
		} else {
			_, err = s.ChannelMessageSendReply(r.request.ChannelID, "<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)", r.request.SoftReference())
		}
		if err != nil {
			sadness(s,nil,err)
			return
		}
	}
	
	if err!=nil {
		sadness(s,nil,err)
		return
	}

	// timerFile.Close()
	// timerFile, err = os.OpenFile("timers.txt", os.O_RDWR, 0666)
	// if err != nil {
	// 	sadness(s,nil,err)
	// 	return
	// }
	defer timerFile.Close()
	head, err := timerFile.Seek(0,0)
	if err!=nil {
		log.Fatal(head, err)
		return
	}

	newFileData := ""
	timerFile.Seek(0,0)
	scanner = bufio.NewScanner(timerFile)

	for scanner.Scan() {
		if strings.Split(scanner.Text(), " ")[0] != r.rqid {
			newFileData += scanner.Text() +"\n"
		} else if strings.Split(scanner.Text(), " ")[6] != "0" {
			newRepeats := r.repeats-1
			if newRepeats < 0 {newRepeats = -1}
			newFileData += strings.Join([]string{
				r.rqid,
				strconv.Itoa(int(r.end.Unix())+r.period),
				strconv.Itoa(int(r.end.Unix())),
				r.target,
				strings.Split(scanner.Text(), " ")[4],
				r.author,
				strconv.Itoa(newRepeats),
				strconv.Itoa(r.period),
				r.message}," ")

			reminders = append(reminders, Reminder{
				end: time.Unix(int64(int(r.end.Unix())+r.period),0),
				start: time.Unix(int64(r.end.Unix()),0),
				message: r.message,
				author: r.author,
				target: r.target,
				request: r.request,
				rqid: r.rqid,
				timer: time.NewTimer(time.Duration(r.period) * time.Second),
				repeats: newRepeats,
				period: r.period,
			})
		}
	}

	i := slices.Index(reminders, *r)
	reminders = slices.Delete(reminders, i, i+1)
	// help!!!!

	err = os.WriteFile("timers.txt", []byte(newFileData), 0666)
	if err != nil {
		sadness(s,nil,err)
		return
	}
}

func listReminders(s *discordgo.Session, m *discordgo.MessageCreate, r *[]Reminder) {
	fullResponse := ""
	count := 0
	for _, rem := range *r{
		if rem.target == m.Author.ID {
			count+=1
			fullResponse += strconv.Itoa(count)+": "+rem.message+" @ <t:"+strconv.Itoa(int(rem.end.Unix()))+">\n"
		}
	}
	if fullResponse==""{
		s.ChannelMessageSendReply(m.ChannelID,"No reminders set.", m.SoftReference())
		return
	}
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: fullResponse,
		Reference: m.Reference(),
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		sadness(s,m,err)
		return
	}
}

func deleteReminder(s *discordgo.Session, m *discordgo.MessageCreate, r *[]Reminder) {
	_, rawind, found := strings.Cut(m.Content, " ")
	if !found {
		iKnowWhatYouAre(s,m)
		return
	}
	ind, err := strconv.Atoi(rawind)
	if err!=nil {
		iKnowWhatYouAre(s,m)
		return
	}

	timerFile, err := os.OpenFile("timers.txt", os.O_RDWR, 0666)
	if err != nil {
		sadness(s,nil,err)
		return
	}
	defer timerFile.Close()

	counter := 0
	scanner := bufio.NewScanner(timerFile)
	newFileData := ""
	for scanner.Scan() {
		if strings.SplitN(scanner.Text()," ",9)[5] != m.Author.ID {
			newFileData += scanner.Text()+"\n"
		} else {
			counter++
			if counter != ind {
				newFileData += scanner.Text()+"\n"
			}
		}
	}
	err = os.WriteFile("timers.txt", []byte(newFileData), 0666)
	if err != nil {
		sadness(s,nil,err)
		return
	}

	counter = 0
	for i,rem := range *r {
		if rem.author == m.Author.ID {
			counter++
		}
		if counter == ind {
			s.ChannelMessageSendReply(m.ChannelID, "...Reminder to "+rem.message+" successfully deleted.", m.Reference())
			*r = slices.Delete(*r, i, i+1)
			return
		}
	}
}

func setTimezone(s *discordgo.Session, m *discordgo.MessageCreate) {
	rawtz, _ := strings.CutPrefix(m.Content, ".settz ")
	
	if strings.Contains(rawtz, "dst") {
		
	}

	if !strings.Contains(rawtz, "/") {
		rawtz = strings.ReplaceAll(strings.ToUpper(rawtz), "UTC", "GMT")
		rawtz = "Etc/"+rawtz
	}

	tzloc, err := time.LoadLocation(rawtz)
	if err!=nil {
		s.ChannelMessageSendReply(m.ChannelID, "\""+rawtz+"\" is apparently not a valid timezone", m.SoftReference())
		return
	}

	tzFile, err := os.OpenFile("timezones.txt", os.O_APPEND, 0666)
	if err != nil {
		sadness(s,m,err)
		return
	}
	defer tzFile.Close()

	if _, err := tzFile.WriteString(strings.Join([]string{m.ID, tzloc.String(),"\n"}, " ")); err != nil {
        log.Fatal(err)
    }

	s.ChannelMessageSendReply(m.ChannelID, "Timezone set to "+tzloc.String(), m.SoftReference())
}