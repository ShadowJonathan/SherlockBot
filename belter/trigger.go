package Belt

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

// incomplete

/*
0 member;
    0 username {ID, old, new}
    1 nickname {ID, old, new}
    2 avatar {ID, old, new}
    3 join {ID}
    4 leave {ID}
    5 new role {ID, role ID}
    6 remove role {ID, role ID}
    7 chat message {ID, Message ID}
    8 got muzzle {ID}
    9 muzzle removed {ID}
1 channel;
    0 name {ID, old, new}
    1 topic {ID, old, new}
    2 chat message {ID, *discordgo.Message}
2 OR;
    0 add {channel ID, ID, type}
    1 delete {channel ID, ID, type}
    2 allow {channel ID, ID, type, old, new}
    3 deny {channel ID, ID, type, old, new}
3 server;
    0 name {old, new}
    1 owner {old (ID), new (ID)}
    2 icon {old, new}
    3 region {old, new}
4 special;
    0 got kicked from server {}
    1 mentioned {User mentioned by (ID), channelID}
    2 mass-leave (triggered by 10 or more people leaving within 10 minutes)
    3 mass-join (raid, 10 or more people joining (and getting muzzled))
*/

var (
	gID string
	uID string
)

var (
	logroot    = "../log/"
	guildroot  = func(gID string) string { return logroot + gID }
	memberroot = logroot + "m/"
	memberlog  = func(mID, gID string) string {
		return memberroot + gID + "/" + mID + "/"
	}
	memberinfo = func(mID, gID string) string {
		return memberroot + gID + "/" + mID + "/" + "latest.inf"
	}
	channelroot = logroot + "c/"
	channellog  = func(cID, gID string) string {
		return channelroot + gID + "/" + cID + "/"
	}
	channelinfo = func(cID, gID string) string {
		return channelroot + gID + "/" + cID + "/" + "latest.inf"
	}
)

//note to self, this must be "go"-able

func trigger(statepointer int, subpointer int, data interface{}) error {
	//statepointer is the type of data that has to be saved, and data is the data that has to be saved
	if data != nil {
		switch statepointer {
		case 0:
			switch subpointer {
			case 0:
			case 1:
			case 2:
			case 3:
			case 4:
			case 5:
			case 6:
			case 7:
				//chat, member-side
			case 8:
			case 9:
			}
		case 1:
			switch subpointer {
			case 0:
			case 1:
			case 2:
				//chat channel-side
				var data *triggerchannelchat
				store("cmc", data.chat)
			}
		case 2:
			switch subpointer {
			case 0:
			case 1:
			case 2:
			case 3:
			}
		case 3:
			switch subpointer {
			case 0:
			case 1:
			case 2:
			case 3:
			}
		case 4:
			switch subpointer {
			case 0:
			case 1:
			case 2:
			case 3:
			}
		}
	} else {
		return errors.New("Nil data")
	}
	return errors.New("Unknown error")
}

var cmc *log.Logger
var cmm *log.Logger

func installfeminists() {
	// dont judge me by the name, okay?
}

func store(T string, data interface{}, instr ...string) {
	switch T {
	case "cmc":
	}
}

type triggerchannelchat struct {
	ID   string
	chat *discordgo.Message
}
