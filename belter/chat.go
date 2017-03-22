package Belt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"errors"

	"runtime/debug"

	"github.com/bwmarrin/discordgo"
)

type ChatLog struct {
	sync.Mutex
	Buffer   map[string]*Chatmessage // IDS->messages
	Mess     chan *discordgo.Message // fresh messages
	Edits    chan *discordgo.Message
	Deletes  chan string
	Backlog  chan *Chatmessage
	Settings *CLSettings
}

const CHATERR = ">ERR<"

func (cl *ChatLog) init() {
	cl.Buffer = make(map[string]*Chatmessage)
	cl.Mess = make(chan *discordgo.Message, 100)
	cl.Edits = make(chan *discordgo.Message, 100)
	cl.Deletes = make(chan string, 100)
	cl.Backlog = make(chan *Chatmessage, 1000)
	cl.loadsettings()
}

func (cl *ChatLog) loadsettings() {
	data, err := ioutil.ReadFile("chatlog/sett.json")
	cl.Settings = &CLSettings{}
	if err != nil {
		cl.Settings = &CLSettings{
			Autologafter: time.Now(),
			ScrubAmount:  1000,
		}
	} else {
		err = json.Unmarshal(data, cl.Settings)
		if err != nil {
			fmt.Println("Error unmarshalling CL settings:", err)
		}
	}
}

func (cl *ChatLog) backlog() {
	cl.Save()
	var count int
	var amount = cl.Settings.ScrubAmount
	fmt.Println("Backlogging with " + strconv.Itoa(cl.Settings.ScrubAmount) + " max")
	for _, g := range sh.dg.State.Guilds {
	CH:
		for _, c := range g.Channels {
			fmt.Println("Logging new channel:", c.Name, count)
			if c.Type != "text" {
				fmt.Println("Is voice channel, skipping...")
				continue
			}
			var latestid string
			for count < amount {
				if count >= amount {
					fmt.Println("Hit limit of backlogging")
					return
				}
				mess, err := sh.dg.ChannelMessages(c.ID, 100, latestid, "")
				if err != nil {
					if strings.Contains(err.Error(), "HTTP 403") {
						fmt.Println("UNAUTH CHANNEL:", c.Name)
						continue CH
					} else {
						fmt.Println("Error processing backlog:", err)
						return
					}
				}
				var anyprocess bool
				var lowest int
				for _, m := range mess {
					id, err := strconv.Atoi(m.ID)
					if (id < lowest || lowest == 0) && err == nil {
						lowest = id
					}
					if ots, _ := m.Timestamp.Parse(); ots.Unix() < cl.Settings.Autologafter.Unix() {
						fmt.Println("Hit back of chat-channel")
						continue CH
					}
					sm, ok := cl.search(m.ID)
					if ok {
						if sm.TopView != m.Content {
							sm.TopView = m.Content
							ts, _ := m.EditedTimestamp.Parse()
							sm.Edits = append(sm.Edits, &MessageEdit{
								At:   ts,
								Edit: m.Content,
								Live: false,
							})
							cl.Backlog <- sm
							anyprocess = true
							count++
						}
					} else {
						anyprocess = true
						ts, err := m.EditedTimestamp.Parse()
						ots, _ := m.Timestamp.Parse()
						if err != nil {
							cl.Backlog <- &Chatmessage{Orig: m.Content, TimeStamp: ots, TopView: m.Content, ID: m.ID, Author: m.Author.ID, Channel: m.ChannelID, Live: false}
						} else {
							var M = &Chatmessage{Orig: CHATERR, TimeStamp: ots, TopView: m.Content, ID: m.ID, Author: m.Author.ID, Channel: m.ChannelID, Live: false}
							M.Edits = append(M.Edits, &MessageEdit{
								At:   ts,
								Edit: m.Content,
								Live: false,
							})
							cl.Backlog <- M
						}
						count++
					}
				}
				latestid = strconv.Itoa(lowest)
				if !anyprocess {
					if latestid == "0" {
						continue CH
					}
					fmt.Println("No unknown messages found around", latestid)
				}
			}
			if count >= amount {
				fmt.Println("Hit limit of backlogging")
				return
			}
		}
	}
PC:
	for _, c := range sh.dg.State.PrivateChannels {
		fmt.Println("Logging new channel:", c.Recipient.Username)
		if c.Type != "text" {
			fmt.Println("Is voice channel, skipping...")
			continue
		}
		var latestid string
		for count < amount {
			mess, err := sh.dg.ChannelMessages(c.ID, 100, latestid, "")
			if err != nil {
				fmt.Println("Error processing backlog:", err)
				return

			}
			var anyprocess bool
			var lowest int
			for _, m := range mess {
				id, err := strconv.Atoi(m.ID)
				if (id < lowest || lowest == 0) && err == nil {
					lowest = id
				}
				if ots, _ := m.Timestamp.Parse(); ots.Unix() < cl.Settings.Autologafter.Unix() {
					fmt.Println("Hit back of chat-channel")
					continue PC
				}
				sm, ok := cl.search(m.ID)
				if ok {
					if sm.TopView != m.Content {
						sm.TopView = m.Content
						ts, _ := m.EditedTimestamp.Parse()
						sm.Edits = append(sm.Edits, &MessageEdit{
							At:   ts,
							Edit: m.Content,
							Live: false,
						})
						cl.Backlog <- sm
						anyprocess = true
						count++
					}
				} else {
					anyprocess = true
					ts, err := m.EditedTimestamp.Parse()
					ots, _ := m.Timestamp.Parse()
					if err != nil {
						cl.Backlog <- &Chatmessage{Orig: m.Content, TimeStamp: ots, TopView: m.Content, ID: m.ID, Author: m.Author.ID, Channel: m.ChannelID, Live: false}
					} else {
						var M = &Chatmessage{Orig: CHATERR, TimeStamp: ots, TopView: m.Content, ID: m.ID, Author: m.Author.ID, Channel: m.ChannelID, Live: false}
						M.Edits = append(M.Edits, &MessageEdit{
							At:   ts,
							Edit: m.Content,
							Live: false,
						})
						cl.Backlog <- M
					}
					count++
				}
			}
			latestid = strconv.Itoa(lowest)
			if !anyprocess {
				if latestid == "0" {
					continue PC
				}
				fmt.Println("No unknown messages found around", latestid)
			}
		}
		if count >= amount {
			fmt.Println("Hit limit of backlogging")
			return
		}
	}
}

func (cl *ChatLog) validate() {
	cl.Mutex.Lock()
	defer cl.Mutex.Unlock()
	chatdir, err := ioutil.ReadDir("chatlog")
	if err != nil {
		os.Mkdir("chatlog", 0777)
		return
	}
	sett, _ := ioutil.ReadFile("chatlog/sett.json")
	var All = make(map[string]*Chatmessage)
	var Dumptill int
	for _, f := range chatdir {
		if !f.IsDir() {
			n := strings.Split(f.Name(), "_")
			if len(n) == 4 && n[0] == "CC" {
				if n[1] == "DF" {
					Dumptill, _ = strconv.Atoi(n[2])
				}
				d, err := ioutil.ReadFile("chatlog/" + f.Name())
				if err != nil {
					fmt.Println("Error reading file "+f.Name()+",", err)
					return
				}
				var CC = &BufferChunk{}
				err = json.Unmarshal(d, CC)
				if err != nil {
					fmt.Println("Error unmarshalling file "+f.Name()+",", err)
					return
				}
				for _, m := range CC.Messages {
					All[m.ID] = m
				}
			}
		}
	}
	fmt.Println("Read all messages\n" + strconv.Itoa(len(All)) + " messages found")
	var Saves = make(map[string]*BufferChunk)
	var DumpSave = &BufferChunk{}
	var Workingsave = &BufferChunk{}

	// first get the dump out of there
	for id, m := range All {
		if i, _ := strconv.Atoi(id); i <= Dumptill {
			if id == "" {
				continue
			}
			DumpSave.Messages = append(DumpSave.Messages, m)
			DumpSave.SavedIDs = append(DumpSave.SavedIDs, id)
		}
	}
	for _, id := range DumpSave.SavedIDs {
		delete(All, id)
	}
	DFNAME := fmt.Sprintf("CC_DF_%d_buff.clg", Dumptill)

	// then do the heavy work
	var latestlowest int
	var last bool
	for {
		var count int
		for count < 200 {
			var latestfound int
			for _, id := range All {
				if id.ID == "" {
					continue
				}
				var i int
				var err error
				if i, err = strconv.Atoi(id.ID); (i < latestfound || latestfound == 0) && i > latestlowest {
					latestfound = i
				}
				if err != nil {
					fmt.Println("ID ERROR:", id.ID, err)
				}
			}

			m, ok := All[strconv.Itoa(latestfound)]
			if !ok && latestfound != 0 {
				panic(strconv.Itoa(latestfound) + " " + strconv.Itoa(latestlowest) + " " + strconv.Itoa(len(Workingsave.Messages)) + " " + strconv.Itoa(len(Saves)) + " " + strconv.Itoa(len(All)))
			} else if !ok && latestfound == 0 {
				last = true
				break
			}
			Workingsave.Messages = append(Workingsave.Messages, m)
			latestlowest = latestfound
			count++
		}
		var lowest int
		var highest int
		var ids []string
		for _, m := range Workingsave.Messages {
			ids = append(ids, m.ID)
			i, _ := strconv.Atoi(m.ID)
			if lowest > i || lowest == 0 {
				lowest = i
			}
			if highest < i {
				highest = i
			}
		}
		Workingsave.SavedIDs = ids
		name := "CC_" + strconv.Itoa(lowest) + "_" + strconv.Itoa(highest) + "_buff.clg"
		Saves[name] = Workingsave
		Workingsave = &BufferChunk{}
		if last {
			break
		}
	}
	fmt.Println("Processed all messages\n" + strconv.Itoa(len(Saves)) + " chunks made")
	var Savefiles []*SaveFile
	for name, buffer := range Saves {
		tf, err := json.Marshal(buffer)
		if err != nil {
			fmt.Println("Error parsing marshal on "+name+":", err)
			return
		}
		var TF = &SaveFile{
			Data: tf,
			Name: name,
		}
		Savefiles = append(Savefiles, TF)
	}

	if len(DumpSave.SavedIDs) > 1 {
		tf, err := json.Marshal(DumpSave)
		if err != nil {
			fmt.Println("Error parsing marshal on "+DFNAME+":", err)
			return
		}
		var TF = &SaveFile{
			Data: tf,
			Name: DFNAME,
		}
		Savefiles = append(Savefiles, TF)
	}

	fmt.Println("Prepared all files")

	os.Mkdir("backup", 0777)
	for _, f := range chatdir {
		if !f.IsDir() {
			data, err := ioutil.ReadFile("chatlog/" + f.Name())
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile("backup/"+f.Name(), data, 0777)
			if err != nil {
				panic(err)
			}
		}
	}
	var haserror bool
	err = os.RemoveAll("chatlog")
	if err != nil {
		fmt.Println("Error while deleting original folder:", err)
		haserror = true
	}
	if !haserror {
		err = os.MkdirAll("chatlog", 0777)
		if err != nil {
			fmt.Println("Error while creating new folder:", err)
			haserror = true
		}
	}
	if haserror {
		fmt.Println("Validate failed, restoring backup")
		backdir, err := ioutil.ReadDir("backup")
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll("chatlog", 0777)
		if err != nil {
			panic(err)
		}
		for _, bf := range backdir {
			if !bf.IsDir() {
				data, err := ioutil.ReadFile("backup/" + bf.Name())
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile("chatlog/"+bf.Name(), data, 0777)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		for _, file := range Savefiles {
			ioutil.WriteFile("chatlog/"+file.Name, file.Data, 0777)
		}

		os.RemoveAll("backup")

		ioutil.WriteFile("chatlog/sett.json", sett, 0777)

		fmt.Println("Validate complete")
	}
}

func (cl *ChatLog) work() {
	go cl.saveloop()
	for {
		select {
		case m := <-cl.Mess:
			if m.Author == nil {
				continue
			}
			t, err := m.Timestamp.Parse()
			if err != nil {
				fmt.Println("Error parsing timestamp:", m.Timestamp)
				continue
			}
			cl.Mutex.Lock()
			cl.Buffer[m.ID] = &Chatmessage{
				Orig:      m.Content,
				TimeStamp: t,
				TopView:   m.Content,
				ID:        m.ID,
				Author:    m.Author.ID,
				Channel:   m.ChannelID,
				Live:      true,
			}
			cl.Mutex.Unlock()
		case e := <-cl.Edits:
			origm, ok := cl.search(e.ID)
			if !ok {
				t, err := e.Timestamp.Parse()
				if err != nil {
					fmt.Println("Error parsing timestamp:", e.Timestamp)
					continue
				}
				cl.Buffer[e.ID] = &Chatmessage{
					Orig:      CHATERR,
					TimeStamp: t,
					TopView:   e.Content,
					ID:        e.ID,
					Author:    e.Author.ID,
					Channel:   e.ChannelID,
					Live:      false,
					Edits: append([]*MessageEdit{}, &MessageEdit{
						At:   time.Now(),
						Edit: e.Content,
						Live: true,
					}),
				}
			} else {
				origm.TopView = e.Content
				origm.Edits = append(origm.Edits, &MessageEdit{
					At:   time.Now(),
					Edit: e.Content,
					Live: true,
				})
				cl.Mutex.Lock()
				cl.Buffer[e.ID] = origm
				cl.Mutex.Unlock()
			}

		case d := <-cl.Deletes:
			origm, ok := cl.search(d)
			if !ok {
				fmt.Println("Message deleted, but never seen:", d)
				cl.Mutex.Lock()
				cl.Buffer[d] = &Chatmessage{
					Orig:       CHATERR,
					Live:       false,
					Deleted:    true,
					ID:         d,
					DeleteTime: time.Now(),
				}
				cl.Mutex.Unlock()
			} else {
				origm.Deleted = true
				origm.DeleteTime = time.Now()
				cl.Mutex.Lock()
				cl.Buffer[d] = origm
				cl.Mutex.Unlock()
			}

		case bl := <-cl.Backlog:
			cl.Mutex.Lock()
			cl.Buffer[bl.ID] = bl
			cl.Mutex.Unlock()
			fmt.Println("Processed backlog:", bl.ID)
		}
	}
}

func (cl *ChatLog) search(ID string) (*Chatmessage, bool) {
	cl.Mutex.Lock()
	m, ok := cl.Buffer[ID]
	cl.Mutex.Unlock()
	if ok {
		return m, true
	}
	chatdir, err := ioutil.ReadDir("chatlog")
	if err != nil {
		os.Mkdir("chatlog", 0777)
		return &Chatmessage{}, false
	}
	id, _ := strconv.Atoi(ID)
FINDLOOP:
	for _, f := range chatdir {
		if !f.IsDir() {
			n := strings.Split(f.Name(), "_")
			if len(n) == 4 && n[0] == "CC" {
				F, _ := strconv.Atoi(n[1])
				T, _ := strconv.Atoi(n[2])
				if id >= F && id <= T {
					d, err := ioutil.ReadFile("chatlog/" + f.Name())
					if err != nil {
						fmt.Println("Error reading file "+f.Name()+",", err)
						return &Chatmessage{}, false
					}
					var CC = &BufferChunk{}
					err = json.Unmarshal(d, CC)
					if err != nil {
						fmt.Println("Error unmarshalling file "+f.Name()+",", err)
						return &Chatmessage{}, false
					}
					for _, m := range CC.Messages {
						if m.ID == ID {
							return m, true
						}
					}
					break FINDLOOP
				}
			}
		}
	}
	return &Chatmessage{}, false
}

func (cl *ChatLog) saveloop() {
	for {
		time.Sleep(1 * time.Hour)
		cl.Save()
	}
}

func (cl *ChatLog) Save() {
	cl.Mutex.Lock()
	defer cl.Mutex.Unlock()
	var Saves = make(map[string]SaveFile)
	var NewSave SaveFile
	var AppendtoNew []*Chatmessage
	var AppendtoDump []*Chatmessage
	var DumpSave SaveFile
	var DS bool
	chatdir, err := ioutil.ReadDir("chatlog")
	if err != nil {
		os.Mkdir("chatlog", 0777)
	}
	var lowestid int
	for _, f := range chatdir {
		if !f.IsDir() {
			n := strings.Split(f.Name(), "_")
			if len(n) == 4 && n[0] == "CC" && n[1] != "DF" {
				F, _ := strconv.Atoi(n[1])
				T, _ := strconv.Atoi(n[2])
				if F < lowestid || lowestid == 0 {
					lowestid = F
				}
				d, _ := ioutil.ReadFile("chatlog/" + f.Name())
				Saves[f.Name()] = SaveFile{
					from: F,
					to:   T,
					Name: f.Name(),
					Data: d,
				}
			} else if len(n) == 4 && n[0] == "CC" && n[1] == "DF" {
				T, _ := strconv.Atoi(n[2])
				lowestid = T
				d, _ := ioutil.ReadFile("chatlog/" + f.Name())
				DumpSave = SaveFile{
					to:   T,
					Data: d,
					Name: f.Name(),
				}
				DS = true
			}
		}
	}
	for id, m := range cl.Buffer {
		ID, _ := strconv.Atoi(id)
		var Found = false
	FIND:
		for n, sf := range Saves {
			if ID >= sf.from && ID <= sf.to {
				sf.Save = true
				var CC = new(BufferChunk)
				err := json.Unmarshal(sf.Data, CC)
				if err != nil {
					panic(err)
				}
				var found = false
				for _, M := range CC.Messages {
					if M.ID == m.ID {
						found = true
						M = m
					}
				}
				if !found {
					CC.Messages = append(CC.Messages, m)
				}
				sf.Data, _ = json.Marshal(CC)
				Saves[n] = sf
				break FIND
			}
		}
		if !Found {
			if ID < lowestid {
				AppendtoDump = append(AppendtoDump, m)
			} else {
				AppendtoNew = append(AppendtoNew, m)
			}
		}
	}
	if len(AppendtoNew) != 0 {
		NewSave.Save = true
		var lowest string
		var highest string
		var e time.Time
		var lat time.Time
		var l int
		var h int
		var ids []string
		for _, m := range AppendtoNew {
			t, _ := strconv.Atoi(m.ID)
			if t < l || l == 0 {
				l = t
			}
			if t > h {
				h = t
			}
			if m.TimeStamp.Unix() < e.Unix() || e.IsZero() {
				e = m.TimeStamp
			}
			if m.TimeStamp.Unix() > lat.Unix() {
				lat = m.TimeStamp
			}
			ids = append(ids, m.ID)
		}
		lowest = strconv.Itoa(l)
		highest = strconv.Itoa(h)
		NewSave.Name = "CC_" + lowest + "_" + highest + "_buff"
		var CC = &BufferChunk{
			SavedIDs: ids,
			Messages: AppendtoNew,
		}
		data, _ := json.Marshal(CC)
		NewSave.Data = data
		Saves[NewSave.Name] = NewSave
	}
	if len(AppendtoDump) != 0 {
		var CC = new(BufferChunk)
		if !DS {
			CC = &BufferChunk{}
		} else {
			err := json.Unmarshal(DumpSave.Data, CC)
			if err != nil {
				fmt.Println("Error unmarchaling dumpfile,", err)
			}
		}
		CC.Messages = AppendtoDump
		var MIDs []string
		for _, m := range CC.Messages {
			MIDs = append(MIDs, m.ID)
		}
		CC.SavedIDs = MIDs
		data, _ := json.Marshal(CC)
		DumpSave.Data = data
		DumpSave.Save = true
		DumpSave.Name = fmt.Sprintf("CC_DF_%d_buff", lowestid)
		Saves[DumpSave.Name] = DumpSave
	}
	for _, sf := range Saves {
		if sf.Save {
			if !strings.Contains(sf.Name, ".clg") {
				ioutil.WriteFile("chatlog/"+sf.Name+".clg", sf.Data, 0777)
			} else {
				ioutil.WriteFile("chatlog/"+sf.Name, sf.Data, 0777)
			}
		}
	}

	cl.Buffer = make(map[string]*Chatmessage)

	d, err := json.Marshal(cl.Settings)
	if err == nil {
		ioutil.WriteFile("chatlog/sett.json", d, 0777)
	} else {
		fmt.Println("Error Marshal-ing settings,", err)
	}
}

type SaveFile struct {
	Data     []byte
	Name     string
	from, to int
	Save     bool
}

type Chatmessage struct {
	Orig       string         // if ">ERR<" then the original message has never been seen
	TimeStamp  time.Time      `json:"TS,omitempty"`
	Edits      []*MessageEdit `json:",omitempty"`
	TopView    string         // for fast processing
	Deleted    bool           `json:",omitempty"`
	DeleteTime time.Time      `json:"DT,omitempty"`
	ID         string
	Author     string // referencing ID
	Channel    string
	Live       bool // if the message has been edited when the bot was online
}

type MessageEdit struct {
	At   time.Time
	Edit string
	Live bool
}

type BufferChunk struct {
	SavedIDs []string       `json:"AID"`
	Messages []*Chatmessage `json:"Ms"`
}

type CLSettings struct {
	Autologafter time.Time `json:"ALA"`
	ScrubAmount  int
}

type ChatBufferResolve struct {
	ResolveChat
	TotalBuffer []*Chatmessage
	Chs         map[string]*ChannelBuffer
	//Gs          map[string]*GuildDetail
	AssetPrefix string
	Chatlogpath string
	Detailed    bool
}

type ResolveChat struct {
	TempBuffer []*Chatmessage
	TempSolved map[string][]*CompressedMessage
	Orig       *ChatBufferResolve
}

type ChannelBuffer struct {
	Messages []*CompressedMessage
	Count    int
	Name     string
	Private  bool
}

/*
type ChannelDetail struct {
	Topic       string
	Permissions map[string]*DetailOR
	Position    int
}


type GuildDetail struct {
	Channels map[int]string
}

type DetailOR struct {
	Type  string
	Deny  int
	Allow int
}
*/

type CompressedMessage struct {
	TopDown  string               `json:"TD"`
	Versions map[time.Time]string `json:"Vs,omitempty"`
	Del      bool                 `json:"D,omitempty"`
	Time     time.Time            `json:"T"`
	Author   string               `json:"Au"`
	Detail   *MessageDetail       `json:"Det,omitempty"`
	ID       string               `json:",omitempty"`
}

type MessageDetail struct {
	DeleteTime   time.Time          `json:"DT,omitempty"`
	Capturedlive bool               `json:"CL,omitempty"`
	LiveEdits    map[time.Time]bool `json:"LE,omitempty"`
}

func (c *ChatBufferResolve) Assetize(CLP string, Aprefix string, Detailed bool) error {
	if CLP == "" {
		return errors.New("No CLP defined")
	} else if Aprefix == "" {
		Aprefix = "assets"
	}
	sh.cl.Save()
	sh.cl.Mutex.Lock()
	defer sh.cl.Mutex.Unlock()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Error parsing assets:", err, "\n"+string(debug.Stack()))
		}
	}()
	c.AssetPrefix = Aprefix
	c.Chatlogpath = CLP
	c.Detailed = Detailed
	c.init()
	err := c.load()
	if err != nil {
		return err
	}
	if Err := c.ResolveChat.Work(); Err != nil {
		return Err
	}
	if Err := c.store(); Err != nil {
		return Err
	}
	if Err := c.inform(); Err != nil {
		return Err
	}
	return nil
}

func (c *ChatBufferResolve) init() {
	c.Chs = make(map[string]*ChannelBuffer)
	c.ResolveChat.Orig = c
	c.ResolveChat.TempSolved = make(map[string][]*CompressedMessage)
}

func (c *ChatBufferResolve) load() error {
	chatdir, err := ioutil.ReadDir(c.Chatlogpath)
	if err != nil {
		return errors.New("Error loading CLF")
	}
	for _, f := range chatdir {
		if !f.IsDir() {
			n := strings.Split(f.Name(), "_")
			if len(n) == 4 && n[0] == "CC" {
				var T = &BufferChunk{}
				d, _ := ioutil.ReadFile(c.Chatlogpath + "/" + f.Name())
				err := json.Unmarshal(d, T)
				if err != nil {
					return errors.New("Error loading from file " + f.Name() + ": " + err.Error())
				}
				c.TotalBuffer = append(c.TotalBuffer, T.Messages...)
			}
		}
	}
	return nil
}

func (c *ChatBufferResolve) store() error {
	for channel, mss := range c.Chs {
		fmt.Println("Storing", mss.Name)
		if channel == "" {
			continue
		}
		var Path string
		if mss.Private {
			Path = c.AssetPrefix + "/DMs/" + channel + " - " + mss.Name
		} else {
			Path = c.AssetPrefix + "/" + channel + " - " + mss.Name
		}
		os.MkdirAll(Path, 0777)
		var days = make(map[int][]*CompressedMessage)
		for _, mc := range mss.Messages {
			days[mc.Time.YearDay()] = append(days[mc.Time.YearDay()], mc)
		}
		for day, messages := range days {
			d, err := json.Marshal(messages)
			if err != nil {
				return errors.New("Error parsing json data: " + err.Error())
			}
			err = ioutil.WriteFile(Path+"/"+strconv.Itoa(day)+"-messages.cpm", d, 0777)
			if err != nil {
				return errors.New("Error saving file " + Path + "/" + strconv.Itoa(day) + "-messages.cpm: " + err.Error())
			}
		}
	}
	return nil
}

func (c *ChatBufferResolve) inform() error {
	var users = make(map[string]string)
	var channels = make(map[string]string)

	for _, c := range c.TotalBuffer {
		if c.Author == "" || c.Channel == "" {
			continue
		}
		users[c.Author] = ""
		channels[c.Channel] = ""
	}

	for u := range users {
		var U *discordgo.User
		var err error
		for _, g := range sh.dg.State.Guilds {
			var GM *discordgo.Member
			GM, err = sh.dg.State.Member(g.ID, u)
			if err == nil {
				U = GM.User
				break
			}
		}
		var name string
		if err != nil {
			fmt.Println("Tried to find user, got error", err)
			name = "ERR"
		} else {
			name = U.Username
		}
		users[u] = name
	}

	for c := range channels {
		CH, err := sh.dg.State.Channel(c)
		var name string
		if err != nil {
			fmt.Println("Tried to find channel, got error", err)
			name = "ERR"
		} else {
			if CH.IsPrivate {
				name = CH.Recipient.Username
			} else {
				name = CH.Name
			}
		}
		channels[c] = name
	}
	ai := new(assetinfo)
	ai.Channels = channels
	ai.Users = users

	data, err := json.Marshal(ai)
	if err != nil {
		panic(err)
	} else {
		ioutil.WriteFile(c.AssetPrefix+"/asset.info", data, 0777)
	}
	return nil
}

type assetinfo struct {
	Channels map[string]string
	Users    map[string]string
}

func (r *ResolveChat) Work() error {
	Mess := r.Orig.TotalBuffer
	for _, m := range Mess {
		M := &CompressedMessage{
			TopDown: m.TopView,
			Del:     m.Deleted,
			Time:    m.TimeStamp,
			Author:  m.Author,
			ID:      m.ID,
		}
		if len(m.Edits) > 0 {
			M.Versions = make(map[time.Time]string)
			M.Versions[m.TimeStamp] = m.Orig
			for _, e := range m.Edits {
				M.Versions[e.At] = e.Edit
			}
		}
		r.TempSolved[m.Channel] = append(r.TempSolved[m.Channel], M)
	}
	for ch, ms := range r.TempSolved {
		C, err := sh.dg.Channel(ch)
		if err != nil {
			r.Orig.Chs[ch] = &ChannelBuffer{
				Messages: ms,
				Name:     "ERR",
				Private:  false,
				Count:    len(ms),
			}
		} else {
			if C.Name == "" {
				if C.IsPrivate {
					C.Name = C.Recipient.Username
				} else {
					C.Name = "ERRNONAME"
				}
			}
			r.Orig.Chs[ch] = &ChannelBuffer{
				Messages: ms,
				Name:     C.Name,
				Private:  C.IsPrivate,
				Count:    len(ms),
			}
		}
	}
	return nil
}
