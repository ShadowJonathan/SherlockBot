package Belt

import (
	fmt "fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Version struct {
	Major byte
	Minor byte
}

type Sherlock struct {
	dg       *discordgo.Session
	Debug    bool
	version  Version
	OwnID    string
	OwnAV    string
	OwnName  string
	Stop     bool
	StopLoop bool
}

// Vars after this

var sh *Sherlock

// Functions after this

var PrimeGuild string
var LastCheck = &LastChangeStatus{}

var notifiers []string

func CheckChange(G *GuildInfo, GID string) (bool, *discordgo.Guild, *FullChangeStruct) {
	ChGl, err := sh.dg.State.Guild(GID)
	if err != nil {
		fmt.Println("Error getting guild: " + err.Error())
	}
	Equal, FCS := DeepEqual(G.g, ChGl)
	return Equal, ChGl, FCS
}

// compare vals here

// rest funcs

func AppendChange(Gold *discordgo.Guild, Gnew *discordgo.Guild, TotC *FullChangeStruct) []string {
	var ChangeString []string
	var NewString string
	if TotC.Guild.Name {
		ChangeString = append(ChangeString, "Server '"+Gold.Name+" changed it's name to '"+Gnew.Name+"!")
	}
	if TotC.Guild.OwnerID {
		ChangeString = append(ChangeString, "`ALERT!`\n`ALERT!`\n`ALERT!`\nServer "+Gold.Name+"'s owner just changed from "+GetUserName(Gold.OwnerID, Gold)+" to "+GetUserName(Gnew.OwnerID, Gnew)+"!!!")
	}
	if TotC.Guild.Icon {
		ChangeString = append(ChangeString, "Server "+Gold.Name+" changed it's icon!")
	}
	if TotC.Guild.Region {
		ChangeString = append(ChangeString, "Server "+Gold.Name+" changed it's region!\nPreviously: "+Gold.Region+"\nNow: "+Gnew.Region)
	}
	if TotC.Guild.channels {
		for _, CCh := range TotC.Channels {
			if !CCh.ExistCrisis {
				Cold := GetChannel(CCh.ID, Gold)
				Cnew := GetChannel(CCh.ID, Gnew)
				if CCh.Name {
					ChangeString = append(ChangeString, "Channel "+Cold.Name+" changed it's name to "+Cnew.Name+"!")
				}
				if CCh.Topic {
					ChangeString = append(ChangeString, "Channel "+Cold.Name+" changed topics from "+Cold.Topic+" to "+Cnew.Topic+"!")
				}
				if CCh.perms {
					for _, P := range CCh.Perms {
						if !P.ExistCrisis {
							Oor := GetOR(P.ID, Cold)
							Nor := GetOR(P.ID, Cnew)
							if Oor.Type == "role" {
								Ror := GetRole(P.ID, Gold)

								if P.Allow {
									ChangeString = append(ChangeString, "Channel "+Cold.Name+" has changed it's allowed overwrite-permissions for the "+Ror.Name+" role changed from `"+strconv.FormatInt(int64(Oor.Allow), 10)+"` to `"+strconv.FormatInt(int64(Nor.Allow), 10)+"`")
								}
								if P.Deny {
									ChangeString = append(ChangeString, "Channel "+Cold.Name+" has changed it's denied overwrite-permissions for the "+Ror.Name+" role changed from `"+strconv.FormatInt(int64(Oor.Deny), 10)+"` to `"+strconv.FormatInt(int64(Nor.Deny), 10)+"`")
								}
							} else {
								Mor := GetUser(P.ID, Gold)

								if P.Allow {
									ChangeString = append(ChangeString, "Channel "+Cold.Name+" has changed it's allowed overwrite-permissions for "+Mor.Username+" changed from `"+strconv.FormatInt(int64(Oor.Allow), 10)+"` to `"+strconv.FormatInt(int64(Nor.Allow), 10)+"`")
								}
								if P.Deny {
									ChangeString = append(ChangeString, "Channel "+Cold.Name+" has changed it's denied overwrite-permissions for "+Mor.Username+" changed from `"+strconv.FormatInt(int64(Oor.Deny), 10)+"` to `"+strconv.FormatInt(int64(Nor.Deny), 10)+"`")
								}
							}
						}
					}
				}
			}
		}
	}
	if TotC.Guild.roles {

	}

	return ChangeString
}

func StartCheckLoop() {
	GLDB, err := ioutil.ReadFile("PrimeGuild")
	if err != nil {
		fmt.Println("Error reading PG file: " + err.Error())
		if GLDB == nil {
			ioutil.WriteFile("PrimeGuild", nil, 9000)
		}
		return
	}
	PrimeGuild = string(GLDB)
	ResumeCheck(PrimeGuild)
	go CheckLoop(LastCheck, PrimeGuild)
}

func CheckLoop(LC *LastChangeStatus, Gid string) {
	LastCheck = LC
	var G *GuildInfo
	for sh.StopLoop == false {
		time.Sleep(30 * time.Second)
		fmt.Println("Loop debug")
		G = LastCheck.GI
		IsEqual, NewGuild, AllChange := CheckChange(G, Gid)
		if !IsEqual {
			_, _, _, H, Mi, S := GetTime()
			h := strconv.FormatInt(int64(H), 10)
			mi := strconv.FormatInt(int64(Mi), 10)
			s := strconv.FormatInt(int64(S), 10)

			Responses := AppendChange(LastCheck.GI.g, NewGuild, AllChange)
			for _, N := range notifiers {
				sh.dg.ChannelMessageSend(N, "-"+h+":"+mi+";"+s+"-")
				for _, R := range Responses {
					sh.dg.ChannelMessageSend(N, R)
				}
			}
		}
		LastCheck = SetLC(LastCheck, NewGuild)
		err := WriteGLDfile(LastCheck.GI, false)
		if err != nil {
			fmt.Println("Error writing GLD file after check: " + err.Error())
			sh.StopLoop = true
		}
	}

}

func ResumeCheck(Gid string) {
	G, err := GetGLDfile(Gid)
	if err != nil || G.NeedRestall {
		Initnewguild(Gid)
		return
	}
	LastCheck.GI = G
	if G.BotUP == true {
		LastCheck.SiLaChBisTerm = true
		fmt.Println("Warning: Sherlock hasn't been closed properly since last boot!")
	} else {
		LastCheck.SiLaChBisTerm = false
	}
	IsEqual, NewGuild, AllChange := CheckChange(G, Gid)
	for _, N := range notifiers {
		sh.dg.ChannelMessageSend(N, "Sherlock's out of 221b, ready for some investigation!")
	}
	if !IsEqual {
		Responses := AppendChange(LastCheck.GI.g, NewGuild, AllChange)
		for _, N := range notifiers {
			sh.dg.ChannelMessageSend(N, "New Changes in since last taxi ride:")
			for _, R := range Responses {
				sh.dg.ChannelMessageSend(N, R)
			}
		}
	}
	LastCheck = SetLC(LastCheck, NewGuild)
	_ = WriteGLDfile(LastCheck.GI, false)
}

func Initnewguild(GID string) {
	fmt.Println("No guild file for '" + GID + "' found!\nCreating new guild files...")
	BGuild, err := sh.dg.State.Guild(GID)
	if err == nil {
		Y, Mo, D, H, Mi, S := GetTime()
		GLD := &GuildInfo{
			Lastcheck:   TimeFormat{Year: Y, Month: Mo, Day: D, Hour: H, Min: Mi, Sec: S},
			g:           BGuild,
			BotUP:       true,
			NeedRestall: false,
		}
		WriteGLDfile(GLD, true)
		WriteGLDfile(GLD, false)
		fmt.Println("New guild file for " + BGuild.Name + " made, ready for logging!")
	} else {
		fmt.Println("Error getting guild status: " + err.Error())
	}
}

// handlers

func BBReady(s *discordgo.Session, r *discordgo.Ready) {
	sh.OwnID = r.User.ID
	sh.OwnAV = r.User.Avatar
	sh.OwnName = r.User.Username
	fmt.Println("Discord: Ready message received\nSH: I am '" + sh.OwnName + "'!\nSH: My User ID: " + sh.OwnID)

	StartCheckLoop()
}

func BBCreateMessage(Ses *discordgo.Session, MesC *discordgo.MessageCreate) {
	Mes := MesC.Message
	if Mes.Content != "" {
		if Mes.Content[0] == '!' {
			ProcessCMD(Mes.Content[1:], Mes)
		}
	}
}

// misc funx

func ProcessCMD(CMD string, M *discordgo.Message) {
	if CMD[:8] == "SetPrime" {
		PID := CMD[9:]
		data := []byte(PID)
		ioutil.WriteFile("PrimeGuild", data, 9000)
		fmt.Println("Set new Prime Guild to '" + PID + "'")
		sh.dg.ChannelMessageSend(M.ChannelID, "`Set new prime guild to "+PID+"`")
	}
}

func DeepEqual(a *discordgo.Guild, b *discordgo.Guild) (bool, *FullChangeStruct) {
	var Equal = true
	var TotC *FullChangeStruct
	Equal, TotC = CompareGuild(a, b, TotC, Equal)
	return Equal, TotC
}

// init

func Initialize(Token string) {
	isdebug, err := ioutil.ReadFile("debugtoggle")
	sh = &Sherlock{
		version:  Version{0, 1},
		Debug:    (err == nil && len(isdebug) > 0),
		Stop:     false,
		StopLoop: false,
	}
	sh.dg, err = discordgo.New(Token)
	if err != nil {
		fmt.Println("Discord Session error, check token, error message: " + err.Error())
		return
	}
	// handlers
	sh.dg.AddHandler(BBReady)
	sh.dg.AddHandler(BBCreateMessage)

	fmt.Println("SH: Handlers installed")

	notifiers = append(notifiers, "269441095862190082")

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
