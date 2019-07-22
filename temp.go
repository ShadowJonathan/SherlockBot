package main

import (
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
)

var Token = "Hello i am a token"

func main() {
	dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
	}
	
	dg.AddHandler(Ready)

	// Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

	sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
}

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Discord: Ready message received\nI am '" + r.User.Username + "'!\nMy User ID: " + r.User.ID)

	// Discord set up and ready
}