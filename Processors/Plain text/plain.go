package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	for ch, folder := range channels {
		path := thedir + "/" + folder

	}
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
