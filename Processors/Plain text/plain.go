package main

import (
	"time"
	"io/ioutil"
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
}

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