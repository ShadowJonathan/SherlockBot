package Belt

import "time"
import "github.com/bwmarrin/discordgo"
import "fmt"
import "io/ioutil"
import "os"
import "strings"
import "strconv"
import "encoding/json"
import "sync"

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
	var count int
	for _, g := range sh.dg.State.Guilds {
	CH:
		for _, c := range g.Channels {
			var latestid string
			for count < cl.Settings.ScrubAmount {
				mess, err := sh.dg.ChannelMessages(c.ID, 100, latestid, "")
				if err != nil {
					fmt.Println("Error processing backlog:", err)
				}
				for _, m := range mess {
					if ots, _ := m.Timestamp.Parse(); ots.Unix() < cl.Settings.Autologafter.Unix() {
						break CH
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
							count++
						}
					} else {
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
					}
				}
			}
			if count >= cl.Settings.ScrubAmount {
				return
			}
		}
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
	fmt.Println("Unknown error trying to find " + ID)
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
		var latid string
		var ids []string
		for _, m := range AppendtoNew {
			t, _ := strconv.Atoi(m.ID)
			if t < l || l == 0 {
				l = t
			}
			if m.TimeStamp.Unix() < e.Unix() || e.IsZero() {
				e = m.TimeStamp
			}
			if t > h {
				h = t
				latid = m.ID
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
			LatestID: latid,
			From:     e,
			To:       lat,
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
		CC.Messages = append(CC.Messages, AppendtoDump...)
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
	From     time.Time      `json:"F,omitempty"`
	To       time.Time      `json:"T"`
	SavedIDs []string       `json:"AID"`
	Messages []*Chatmessage `json:"Ms"`
	LatestID string         `json:",omitempty"`
}

type CLSettings struct {
	Autologafter time.Time `json:"ALA"`
	ScrubAmount  int
}
