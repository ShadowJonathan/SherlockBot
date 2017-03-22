package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"fmt"

	"github.com/bwmarrin/discordgo"
)

var info *clientinfo

var dg *discordgo.Session

type clientinfo struct {
	Logcreds     bool
	ID           string
	Password     string
	Email        string
	Sherlock     string //ID of the sherlock in contact
	connected    bool
	initialready *discordgo.Ready
}

var ERR *log.Logger

func main() {
	data, err := ioutil.ReadFile("donotdeletemepls.inf")
	info = &clientinfo{}
	if err != nil {
		ioutil.WriteFile("donotdeletemepls.inf", []byte{}, 9001)
	} else {
		err = json.Unmarshal(data, info)
		if err != nil {
			ioutil.WriteFile("donotdeletemepls.inf", []byte{}, 9001)
			info = &clientinfo{}
		}
	}
	SHC := http.NewServeMux()
	SHC.HandleFunc("/", login)
	SetLogger()
	launchbrowser()
	ERR.Fatal(http.ListenAndServe(":8080", SHC))
}

func SetLogger() {
	l, err := os.OpenFile("error.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		ioutil.WriteFile("log.txt", []byte{}, 9001)
	}
	ERR = log.New(l, "SH: ", log.LstdFlags)
}

func launchbrowser() {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], "http://localhost:8080/")...)
	cmd.Start()
	return
}

var readydata chan *discordgo.Ready

func login(w http.ResponseWriter, r *http.Request) {
	if info.Logcreds {
		if !info.connected {
			var err error
			dg, err = discordgo.New(info.Email, info.Password)
			readydata = make(chan *discordgo.Ready)
			if err != nil {
				ERR.Fatal(err)
			}
			dg.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
				fmt.Println("ready!")
				readydata <- m
			})
			dg.Open()
			ready := <-readydata
			rd, _ := json.Marshal(ready)
			info.initialready = ready
			w.Write(rd)
			info.connected = true
		} else {
			rd, _ := json.Marshal(info.initialready)
			w.Write(rd)
		}
	}
}
