package Belt

import (
	fmt "fmt"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Version struct {
	Major               byte
	Minor               byte
	Build               byte
	Experimental        bool
	ExperimentalVersion byte
}

type Sherlock struct {
	dg      *discordgo.Session
	Debug   bool
	version Version
	OwnID   string
	OwnAV   string
	OwnName string
	Stop    bool
}

// Vars after this

var sh *Sherlock

// Functions after this

func BBReady(s *discordgo.Session, r *discordgo.Ready) {
	sh.OwnID = r.User.ID
	sh.OwnAV = r.User.Avatar
	sh.OwnName = r.User.Username
	fmt.Println("Discord: Ready message received\nSH: I am '" + sh.OwnName + "'!\nSH: My User ID: " + sh.OwnID)
}

func Initialize(Token string) {
	isdebug, err := ioutil.ReadFile("debugtoggle")
	sh = &Sherlock{
		version: Version{0, 1, 0, false, 0},
		Debug:   (err == nil && len(isdebug) > 0),
		Stop:    false,
	}
	sh.dg, err = discordgo.New(Token)
	if err != nil {
		fmt.Println("Discord Session error, check token, error message: " + err.Error())
		return
	}
	// handlers
	sh.dg.AddHandler(BBReady)

	fmt.Println("SH: Handlers installed")

	err = sh.dg.Open()
	if err == nil {
		fmt.Println("Discord: Connection established")
		for !sh.Stop {
			time.Sleep(400 * time.Millisecond)
		}
	} else {
		fmt.Println("Error opening websocket connection: ", err.Error())
	}
	fmt.Println("SH: Sherlock stopping...")
	sh.dg.Close()
}
