package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	var thedir string
	for _, f := range dir {
		if f.IsDir() {
			subdir, err := ioutil.ReadDir(f.Name())
			if err != nil {
				panic(err)
			}
			var yes bool
			for _, f2 := range subdir {
				if f2.Name() == "asset.info" {
					yes = true
					break
				}
			}
			if yes {
				thedir = f.Name()
				break
			}
		}
	}
	if thedir == "" {
		fmt.Println("No usable asset directory found")
		return
	}
	var DMfolder string
	var channels = make(map[string]string)
	var DMchannels = make(map[string]string)
	var AI = new(assetinfo)
	assetdir, err := ioutil.ReadDir(thedir)
	if err != nil {
		panic(err)
	}
	for _, f := range assetdir {
		if f.Name() == "DMs" && f.IsDir() {
			DMfolder = f.Name()
		} else if !f.IsDir() && f.Name() == "asset.info" {
			data, err := ioutil.ReadFile(thedir + "/" + f.Name())
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(data, AI)
			if err != nil {
				panic(err)
			}
		} else {
			if f.IsDir() {
				s := strings.Split(f.Name(), " - ")
				channels[s[0]] = f.Name()
			}
		}
	}
	if DMfolder != "" {
		dm, err := ioutil.ReadDir(thedir + "/" + DMfolder)
		if err != nil {
			panic(err)
		}
		for _, f := range dm {
			s := strings.Split(f.Name(), " - ")
			DMchannels[s[0]] = f.Name()
		}
	}
	channelswithdays := make(map[string]map[int]string)

	privatewithdays := make(map[string]map[int]string)

	for ch, folder := range channels {
		path := thedir + "/" + folder
		channelswithdays[ch] = make(map[int]string)
		localdir, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range localdir {
			s := strings.Split(file.Name(), "-")
			n, err := strconv.Atoi(s[0])
			if err != nil {
				panic(err)
			}
			channelswithdays[ch][n] = file.Name()
		}
	}

	for ch, folder := range DMchannels {
		path := thedir + "/" + DMfolder + "/" + folder
		privatewithdays[ch] = make(map[int]string)
		localdir, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range localdir {
			s := strings.Split(file.Name(), "-")
			n, err := strconv.Atoi(s[0])
			if err != nil {
				panic(err)
			}
			channelswithdays[ch][n] = file.Name()
		}
	}

	res := getresolve(thedir, AI)

	res.pwd = privatewithdays
	res.cwd = channelswithdays
	res.work.make()
}

type resolve struct {
	rchannels map[string][]*file // channel name + files
	rdms      map[string][]*file
	assetdir  string
	*assetinfo
	cwd map[string]map[int]string // channel name -> number + filename
	pwd map[string]map[int]string
	work
}

type file map[string][]byte

type work struct {
	orig *resolve
	cf   *stringify
}

func (w *work) make() {
	for ch, files := range w.orig.cwd {
		for _, messagefile := range files {
			f, err := ioutil.ReadFile(w.orig.assetdir + "/" + ch + " - " + w.orig.assetinfo.Channels[ch] + "/" + messagefile)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(f, w.cf.messages)
			if err != nil {
				panic(err)
			}
			w.cf.parse()
			err = os.MkdirAll("plain/"+w.orig.assetinfo.Channels[ch], 0777)
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile("plain/"+w.orig.assetinfo.Channels[ch]+"/"+messagefile, w.cf.data, 0777)
		}
	}
}

type stringify struct {
	AI           *assetinfo
	messages     []*CompressedMessage
	lines        []string
	completeline string
	data         []byte
}

func (s *stringify) parse() {
	s.str()
}

func (s *stringify) str() {
	var lastauthor string
	var lastid int
	for {
		var theid int
		var idint int
		for I, m := range s.messages {
			i, _ := strconv.Atoi(m.ID)
			if (i < theid || theid == 0) && theid > lastid {
				theid = i
				idint = I
			}
		}

		m := s.messages[idint]
		if lastauthor == m.Author {
			s.lines = append(s.lines, m.clang())
		} else {
			hourminutesecond := m.Time.Format("04:05:45")
			y, mon, d := m.Time.Date()
			daymonthyear := strconv.Itoa(d) + "-" + strconv.Itoa(int(mon)) + "-" + strconv.Itoa(y)
			s.lines = append(s.lines, m.Author+" - "+hourminutesecond+" "+daymonthyear)
		}
	}

}

func (cm *CompressedMessage) clang() string {
	var temp []string
	temp = append(temp, "    "+cm.TopDown)
	for t, ver := range cm.Versions {
		if ver == cm.TopDown {
			continue
		} else {
			hourminutesecond := t.Format("04:05:45")
			y, m, d := t.Date()
			daymonthyear := strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
			temp = append(temp, fmt.Sprintf("(Edit %s %s: %s)", hourminutesecond, daymonthyear, ver))
		}
	}
	if cm.Del {
		temp = append(temp, "(DELETED)")
	}
	return strings.Join(temp, "\n")
}

func (s *stringify) fullstr() {
	s.completeline = strings.Join(s.lines, "\n")
}

func (s *stringify) d() {
	s.data = []byte(s.completeline)
}

func getresolve(ad string, ai *assetinfo) *resolve {
	r := &resolve{
		rchannels: make(map[string][]*file),
		rdms:      make(map[string][]*file),
		assetdir:  ad,
		assetinfo: ai,
		cwd:       make(map[string]map[int]string), pwd: make(map[string]map[int]string),
	}
	r.work = work{
		orig: r,
		cf:   &stringify{AI: r.assetinfo},
	}
	return r
}

/* format:
Shadowjonathan - 08:40:10 22-3-2017
    I didnt do that!
(Edit 08:40:30 22-3-2017: wut)
    hm? :3
*/

type assetinfo struct {
	Channels map[string]string
	Users    map[string]string
}

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
