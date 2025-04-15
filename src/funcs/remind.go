package funcs

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

// PERHAPS, THIS IS HELL.

// AND NOW, 	   REMIN-
// 				    DERS

type ReminderManager struct {
	Reminders []Reminder
}

type Reminder struct {
	end, start time.Time
	message    string
	author     string
	target     string
	request    *discordgo.Message
	rqid       string
	timer      *time.Timer
	repeats    int
	period     int
}

type Token struct {
	key  string
	txt  string
	inst int
}

func SetReminder(inst *Instance, m *discordgo.MessageCreate) {
	rawCmd, _ := strings.CutPrefix(m.Content, ".remind")
	cmd := strings.Split(strings.TrimSpace(rawCmd), " ")

	tokens := []string{"to", "in", "at", "on", "every", "forever", "times"}
	tokenIndicesMap := make(map[string][]int)
	tokenIndicesOrdered := make([]Token, 0)

	for inst, t := range cmd {
		if slices.Contains(tokens, t) {
			tokenIndicesMap[t] = append(tokenIndicesMap[t], inst)
			tokenIndicesOrdered = append(tokenIndicesOrdered, Token{t, "", inst})
		} else if len(tokenIndicesOrdered) > 0 {
			tokenIndicesOrdered[len(tokenIndicesOrdered)-1].txt += t + " "
		}
	}

	fmt.Println(tokenIndicesMap, "\n", tokenIndicesOrdered)
	if len(tokenIndicesOrdered) == 0 {
		inst.Session.ChannelMessageSendReply(m.ChannelID, "Invalid reminder formatting: try again! (you probably forgot the time interval)", m.Reference())
		return
	}

	timeInUnix := 0
	timeSet := false
	dateSet := false
	isRelative := false
	foreverFlag := false
	repeats := 1
	repeatTime := 0
	msg := strings.Join(cmd[:tokenIndicesOrdered[0].inst], " ") + " "
	msg, _ = strings.CutPrefix(msg, "me ")
	targetUser := m.Author.ID
	for _, t := range tokenIndicesOrdered {
		switch t.key {
		case "in":
			if timeSet || dateSet {
				msg += t.key + " " + t.txt
				continue
			}
			text := strings.Split(strings.TrimSpace(t.txt), " ")
			eot := false
			totalTime := 0
			for len(text) > 0 && !eot && !timeSet {
				timeIncr, err := strconv.Atoi(text[0])
				if err != nil {
					if text[0] == "a" || text[0] == "an" {
						timeIncr = 1
					} else if strings.Contains("smhdwyc", string(text[0][len(text[0])-1])) {
						time, thing := text[0][:len(text[0])-1], text[0][len(text[0])-1]
						timeIncr, err = strconv.Atoi(time)
						if err != nil {
							break
						}
						together := append([]string{time}, string(thing))
						text = append(together, text[1:]...)
					} else {
						break
					}
				}
				switch strings.ToLower(text[1]) {
				case "inst.Session", "seconds", "second", "sec":
					totalTime += timeIncr * 1
				case "m", "minutes", "minute", "min":
					totalTime += timeIncr * 60
				case "h", "hours", "hour":
					totalTime += timeIncr * 60 * 60
				case "d", "days", "day":
					totalTime += int(time.Now().AddDate(0, 0, timeIncr).Unix()) - int(time.Now().Unix())
				case "w", "weeks", "week":
					totalTime += int(time.Now().AddDate(0, 0, timeIncr*7).Unix()) - int(time.Now().Unix())
				case "months", "month":
					totalTime += int(time.Now().AddDate(0, timeIncr, 0).Unix()) - int(time.Now().Unix())
				case "y", "years", "year":
					totalTime += int(time.Now().AddDate(timeIncr, 0, 0).Unix()) - int(time.Now().Unix())
				case "c", "centuries", "century":
					totalTime += int(time.Now().AddDate(timeIncr*100, 0, 0).Unix()) - int(time.Now().Unix())
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
				msg += t.key + " " + t.txt
				continue
			}
			var sec, min, hour int
			rawt := strings.Split(strings.TrimSpace(t.txt), ":")
			var cl []int
			for _, r := range rawt {
				val, err := strconv.Atoi(r)
				if err != nil {
					iKnowWhatYouAre(inst, m)
					inst.ErrorChan <- err
					return
				}
				cl = append(cl, val)
			}
			if len(cl) == 2 {
				sec = 0
			} else if len(cl) == 3 {
				sec = cl[2]
			} else {
				iKnowWhatYouAre(inst, m)
				inst.ErrorChan <- fmt.Errorf("skill issue: bad time formatting (found %d>3 fragments)", len(cl))
				return
			}
			min = cl[1]
			hour = cl[0]
			timeInUnix += hour*60*60 + min*60 + sec
			timeSet = true
		case "on":
			if isRelative || dateSet {
				msg += t.key + " " + t.txt
				continue
			}
			var day, month, year int
			rawt := strings.Split(strings.TrimSpace(t.txt), "/")
			if rawt[0] == t.txt {
				rawt = strings.Split(strings.TrimSpace(t.txt), ".")
				if rawt[0] == t.txt {
					rawt = strings.Split(strings.TrimSpace(t.txt), "-")
					if rawt[0] == t.txt {
						iKnowWhatYouAre(inst, m)
						inst.ErrorChan <- fmt.Errorf("skill issue: bad date separator \"%inst.Session\"", t.txt)
						return
					}
				}
			}
			var cl []int
			for _, r := range rawt {
				val, err := strconv.Atoi(r)
				if err != nil {
					iKnowWhatYouAre(inst, m)
					inst.ErrorChan <- fmt.Errorf("skill issue: bad date formatting (%inst.Session is not number)", r)
					return
				}
				cl = append(cl, val)
			}
			if len(cl) == 2 {
				year = time.Now().Year()
			} else if len(cl) == 3 {
				year = cl[2]
			} else {
				iKnowWhatYouAre(inst, m)
				inst.ErrorChan <- fmt.Errorf("skill issue: bad date formatting (found %d>3 fragments)", len(cl))
				return
			}
			if year%100 == year {
				year += (time.Now().Year() / 100) * 100
			}
			month = cl[1]
			day = cl[0]
			timeInUnix = int(time.Unix(int64(timeInUnix), 0).AddDate(year, month, day).Unix())
			dateSet = true
		case "every":
			if repeatTime != 0 {
				msg += t.key + " " + t.txt
				continue
			}
			text := strings.Split(strings.TrimSpace(t.txt), " ")
			eot := false
			for len(text) > 1 && !eot {
				timeIncr, err := strconv.Atoi(text[0])
				if err != nil {
					if strings.Contains("smhdwyc", string(text[0][len(text[0])-1])) {
						time, thing := text[0][:len(text[0])-1], text[0][len(text[0])-1]
						together := append([]string{time}, string(thing))
						text = append(together, text[1:]...)
						timeIncr, err = strconv.Atoi(text[0])
						if err != nil {
							iKnowWhatYouAre(inst, m)
							inst.ErrorChan <- err
							return
						}
					} else {
						timeIncr = 1
					}
				}
				switch strings.ToLower(text[1]) {
				case "inst.Session", "seconds", "second", "sec":
					repeatTime += timeIncr * 1
				case "m", "minutes", "minute", "min":
					repeatTime += timeIncr * 60
				case "h", "hours", "hour":
					repeatTime += timeIncr * 60 * 60
				case "d", "days", "day":
					repeatTime += int(time.Now().AddDate(0, 0, timeIncr).Unix()) - int(time.Now().Unix())
				case "w", "weeks", "week":
					repeatTime += int(time.Now().AddDate(0, 0, timeIncr*7).Unix()) - int(time.Now().Unix())
				case "months", "month":
					repeatTime += int(time.Now().AddDate(0, timeIncr, 0).Unix()) - int(time.Now().Unix())
				case "y", "years", "year":
					repeatTime += int(time.Now().AddDate(timeIncr, 0, 0).Unix()) - int(time.Now().Unix())
				case "c", "centuries", "century":
					repeatTime += int(time.Now().AddDate(timeIncr*100, 0, 0).Unix()) - int(time.Now().Unix())
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
			if repeats != 1 {
				msg += t.key + " " + t.txt
				continue
			}
			timesHopefully := cmd[t.inst-1]
			if times, err := strconv.Atoi(timesHopefully); err == nil {
				repeats = times
				msg = msg[0 : len(msg)-len(timesHopefully)]
				msg += t.txt
			} else {
				msg += t.key + " " + t.txt
				continue
			}
		case "forever":
			foreverFlag = true
			// this one's a doozy
		default:
			msg += t.key + " " + t.txt
		}
	}

	if !timeSet && dateSet {
		hrs, min, sec := time.Now().Clock()
		timeInUnix += hrs*60*60 + min*60 + sec
	}
	if timeSet && !dateSet {
		year, month, day := time.Now().Date()
		timeInUnix = int(time.Unix(int64(timeInUnix), 0).AddDate(year, int(month), day).Unix())
	}

	if foreverFlag && repeatTime != 0 && repeats == 1 {
		repeats = -1
	}

	if aoaoao, rest, _ := strings.Cut(msg, " "); aoaoao == "to" {
		msg = rest
	}

	// TODO: timezones!!
	if !isRelative {
		tzfile, err := os.ReadFile("timezones.txt")
		if err != nil {
			Sadness(inst, m)
			inst.ErrorChan <- err
			return
		}

		timezones := strings.Split(string(tzfile), "\n")
		tzstring := ""
		for _, filetz := range timezones {
			if strings.Split(filetz, " ")[0] == targetUser {
				tzstring = strings.Split(filetz, " ")[1]
			}
		}

		if tzstring == "" {
			inst.Session.ChannelMessageSendReply(m.ChannelID, "You tried to set an absolute reminder, but haven't specified a timezone! To do that, use .settz.", m.SoftReference())
			return
		}

		tzloc, err := time.LoadLocation(tzstring)
		if err != nil {
			Sadness(inst, m)
			inst.ErrorChan <- err
			return
		}

		timeParsed, err := time.ParseInLocation(time.DateTime, time.Unix(int64(timeInUnix), 0).Format(time.DateTime), tzloc)
		if err != nil {
			Sadness(inst, m)
			inst.ErrorChan <- err
			return
		}

		timeInUnix = int(timeParsed.Unix())
	}

	if tagIndex := slices.IndexFunc(cmd[:tokenIndicesOrdered[0].inst], func(str string) bool { return strings.HasPrefix(str, "<@") }); tagIndex >= 0 {
		targetUser = cmd[tagIndex][2 : len(cmd[tagIndex])-1]
		msg = strings.Replace(msg, "<@"+targetUser+">", "", -1)
	}

	msg, _ = strings.CutPrefix(strings.TrimSpace(msg), "to ")

	inst.RManager.AppendReminder(
		time.Unix(int64(timeInUnix), 0),
		time.Now(),
		msg,
		m.Author.ID,
		targetUser,
		m.Message,
		m.Message.ID,
		time.NewTimer(time.Duration(timeInUnix-int(time.Now().Unix()))*time.Second),
		repeats-1,
		repeatTime,
		m.ChannelID,
	)
	inst.Session.ChannelMessageSendReply(m.ChannelID, "...As you wish. I shall send a reminder at <t:"+strconv.Itoa(timeInUnix)+">.", m.Reference())
}

func Remind(inst *Instance, r *Reminder) {
	r = r.Clone()
	timerFile, err := os.OpenFile("timers.txt", os.O_RDONLY, 0666)
	if err != nil {
		Sadness(inst, nil)
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
			if len("<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)") > 2000 {
				_, err = inst.Session.ChannelMessageSend(reqChan, "<@"+r.target+">: YOUR MESSAGE DIDN'T FIT")
			} else {
				_, err = inst.Session.ChannelMessageSend(reqChan, "<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)")
			}
			if err != nil {
				Sadness(inst, nil)
				return
			}
		} else {
			Sadness(inst, nil)
			print(r.message)
			//SHIT WE DIDN'T FIND IT ABORT ABORT
			return
		}
	} else {
		if len("<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)") > 2000 {
			_, err = inst.Session.ChannelMessageSendReply(r.request.ChannelID, "<@"+r.target+">: YOUR MESSAGE DIDN'T FIT", r.request.SoftReference())
		} else {
			_, err = inst.Session.ChannelMessageSendReply(r.request.ChannelID, "<@"+r.target+">: "+r.message+" (set at <t:"+strconv.Itoa(int(r.start.Unix()))+">)", r.request.SoftReference())
		}
		if err != nil {
			Sadness(inst, nil)
			return
		}
	}

	if err != nil {
		Sadness(inst, nil)
		return
	}

	// timerFile.Close()
	// timerFile, err = os.OpenFile("timers.txt", os.O_RDWR, 0666)
	// if err != nil {
	// 	Sadness(inst,nil,err)
	// 	return
	// }
	defer timerFile.Close()
	head, err := timerFile.Seek(0, 0)
	if err != nil {
		log.Fatal(head, err)
		return
	}

	newFileData := ""
	timerFile.Seek(0, 0)
	scanner = bufio.NewScanner(timerFile)

	for scanner.Scan() {
		if strings.Split(scanner.Text(), " ")[0] != r.rqid {
			newFileData += scanner.Text() + "\n"
		} else if strings.Split(scanner.Text(), " ")[6] != "0" {
			newRepeats := r.repeats - 1
			if newRepeats < 0 {
				newRepeats = -1
			}
			newFileData += strings.Join([]string{
				r.rqid,
				strconv.Itoa(int(r.end.Unix()) + r.period),
				strconv.Itoa(int(r.end.Unix())),
				r.target,
				strings.Split(scanner.Text(), " ")[4],
				r.author,
				strconv.Itoa(newRepeats),
				strconv.Itoa(r.period),
				r.message}, " ")

			inst.RManager.Reminders = append(inst.RManager.Reminders, Reminder{
				end:     time.Unix(int64(int(r.end.Unix())+r.period), 0),
				start:   time.Unix(int64(r.end.Unix()), 0),
				message: r.message,
				author:  r.author,
				target:  r.target,
				request: r.request,
				rqid:    r.rqid,
				timer:   time.NewTimer(time.Duration(r.period) * time.Second),
				repeats: newRepeats,
				period:  r.period,
			})
		}
	}

	i := slices.Index(inst.RManager.Reminders, *r)
	inst.RManager.Reminders = slices.Delete(inst.RManager.Reminders, i, i+1)
	// help!!!!

	err = os.WriteFile("timers.txt", []byte(newFileData), 0666)
	if err != nil {
		Sadness(inst, nil)
		return
	}
}

func ListReminders(inst *Instance, m *discordgo.MessageCreate) {
	fullResponse := ""
	count := 0
	for _, rem := range inst.RManager.Reminders {
		if rem.target == m.Author.ID {
			count += 1
			fullResponse += strconv.Itoa(count) + ": " + rem.message + " @ <t:" + strconv.Itoa(int(rem.end.Unix())) + ">\n"
		}
	}
	if fullResponse == "" {
		inst.Session.ChannelMessageSendReply(m.ChannelID, "No reminders set.", m.SoftReference())
		return
	}
	_, err := inst.Session.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:   fullResponse,
		Reference: m.Reference(),
		Flags:     discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		Sadness(inst, m)
		return
	}
}

func DeleteReminder(inst *Instance, m *discordgo.MessageCreate) {
	_, rawind, found := strings.Cut(m.Content, " ")
	if !found {
		iKnowWhatYouAre(inst, m)
		return
	}
	ind, err := strconv.Atoi(rawind)
	if err != nil {
		iKnowWhatYouAre(inst, m)
		return
	}

	timerFile, err := os.OpenFile("timers.txt", os.O_RDWR, 0666)
	if err != nil {
		Sadness(inst, nil)
		return
	}
	defer timerFile.Close()

	counter := 0
	scanner := bufio.NewScanner(timerFile)
	newFileData := ""
	for scanner.Scan() {
		if strings.SplitN(scanner.Text(), " ", 9)[5] != m.Author.ID {
			newFileData += scanner.Text() + "\n"
		} else {
			counter++
			if counter != ind {
				newFileData += scanner.Text() + "\n"
			}
		}
	}
	err = os.WriteFile("timers.txt", []byte(newFileData), 0666)
	if err != nil {
		Sadness(inst, nil)
		return
	}

	counter = 0
	for i, rem := range inst.RManager.Reminders {
		if rem.author == m.Author.ID {
			counter++
		}
		if counter == ind {
			inst.Session.ChannelMessageSendReply(m.ChannelID, "...Reminder to "+rem.message+" successfully deleted.", m.Reference())
			inst.RManager.Reminders = slices.Delete(inst.RManager.Reminders, i, i+1)
			return
		}
	}
}

func SetTimezone(inst *Instance, m *discordgo.MessageCreate) {
	rawtz, _ := strings.CutPrefix(m.Content, ".settz ")

	if strings.Contains(rawtz, "dst") {

	}

	if !strings.Contains(rawtz, "/") {
		rawtz = strings.ReplaceAll(strings.ToUpper(rawtz), "UTC", "GMT")
		rawtz = "Etc/" + rawtz
	}

	tzloc, err := time.LoadLocation(rawtz)
	if err != nil {
		inst.Session.ChannelMessageSendReply(m.ChannelID, "\""+rawtz+"\" is apparently not a valid timezone", m.SoftReference())
		return
	}

	tzFile, err := os.OpenFile("timezones.txt", os.O_APPEND, 0666)
	if err != nil {
		Sadness(inst, m)
		return
	}
	defer tzFile.Close()

	if _, err := tzFile.WriteString(strings.Join([]string{m.ID, tzloc.String(), "\n"}, " ")); err != nil {
		log.Fatal(err)
	}

	inst.Session.ChannelMessageSendReply(m.ChannelID, "Timezone set to "+tzloc.String(), m.SoftReference())
}

func (r *ReminderManager) AppendReminder(end time.Time, start time.Time, message string, author string, target string, request *discordgo.Message, rqid string, timer *time.Timer, repeats int, period int, channelID string) error {
	timerFile, err := os.OpenFile("timers.txt", os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer timerFile.Close()

	timerFile.WriteString(strings.Join([]string{rqid, strconv.Itoa(int(end.Unix())), strconv.Itoa(int(start.Unix())), target, channelID, author, strconv.Itoa(repeats), strconv.Itoa(period), message}, " ") + "\n")
	r.Reminders = append(r.Reminders, Reminder{
		end:     end,
		start:   start,
		message: message,
		author:  author,
		target:  target,
		request: request,
		rqid:    rqid,
		timer:   timer,
		repeats: repeats,
		period:  period,
	})
	return nil
}

func (r *Reminder) GetTimer() *time.Timer {
	return r.timer
}

func (r *Reminder) Clone() *Reminder {
	return &Reminder{
		end:     r.end,
		start:   r.start,
		message: r.message,
		author:  r.author,
		target:  r.target,
		request: r.request,
		rqid:    r.rqid,
		timer:   r.timer,
		repeats: r.repeats,
		period:  r.period,
	}
}

func (r *ReminderManager) ReviseRemindersAfterStartup(inst *Instance) error {
	reminderFile, err := os.OpenFile("timers.txt", os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer reminderFile.Close()

	newFileData := ""
	scanner := bufio.NewScanner(reminderFile)

	for scanner.Scan() {
		reminderText := strings.SplitN(scanner.Text(), " ", 9)
		remTime, _ := strconv.Atoi(reminderText[1])
		if int64(remTime) <= time.Now().Unix() {
			// the order is: request message ID (0), end time (1), start time (2), target user ID (3), channel ID (4), author ID (5), repeats (6), repeat period (7), message (8)
			_, err := inst.Session.ChannelMessageSend(reminderText[4], "<@"+reminderText[3]+">: "+reminderText[8]+" (set at <t:"+reminderText[2]+">) (SORRY I'M LATE I WAS BEING LOBOTOMIZED)")
			if err != nil {
				return err
			}
			if reminderText[6] != "0" {
				endInt, _ := strconv.Atoi(reminderText[1])
				repeatsInt, _ := strconv.Atoi(reminderText[6])
				periodInt, _ := strconv.Atoi(reminderText[7])
				skippedRepeats := (int(time.Now().Unix()) - endInt) / periodInt
				if repeatsInt > skippedRepeats || repeatsInt < 0 {
					newFileData += strings.Join([]string{
						reminderText[0],
						strconv.Itoa(endInt + periodInt*(skippedRepeats+1)),
						strconv.Itoa(endInt + periodInt*skippedRepeats),
						reminderText[3],
						reminderText[4],
						reminderText[5],
						strconv.Itoa(repeatsInt - skippedRepeats),
						strconv.Itoa(periodInt),
						reminderText[8]}, " ")
				}
				r.Reminders = append(r.Reminders, Reminder{
					end:     time.Unix(int64(endInt+periodInt*(skippedRepeats+1)), 0),
					start:   time.Unix(int64(endInt+periodInt*skippedRepeats), 0),
					message: reminderText[8],
					author:  reminderText[5],
					target:  reminderText[3],
					request: nil,
					rqid:    reminderText[0],
					timer:   time.NewTimer(time.Duration(endInt+periodInt*(skippedRepeats+1)-int(time.Now().Unix())) * time.Second),
					repeats: repeatsInt - skippedRepeats,
					period:  periodInt,
				})
			}
		} else {
			newFileData += scanner.Text() + "\n"
			endInt, _ := strconv.Atoi(reminderText[1])
			totalTime := endInt - int(time.Now().Unix())
			startInt, _ := strconv.Atoi(reminderText[2])
			repeatsInt, _ := strconv.Atoi(reminderText[6])
			periodInt, _ := strconv.Atoi(reminderText[7])
			r.Reminders = append(r.Reminders, Reminder{
				end:     time.Unix(int64(endInt), 0),
				start:   time.Unix(int64(startInt), 0),
				message: reminderText[8],
				author:  reminderText[5],
				target:  reminderText[3],
				request: nil,
				rqid:    reminderText[0],
				timer:   time.NewTimer(time.Duration(totalTime) * time.Second),
				repeats: repeatsInt,
				period:  periodInt,
			})
		}
	}
	err = os.WriteFile("timers.txt", []byte(newFileData), 0666)
	if err != nil {
		return err
	}
	return nil
}
