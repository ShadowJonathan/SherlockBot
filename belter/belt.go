package Belt

import (
	//"bytes"
	"encoding/json"
	fmt "fmt"
	"io/ioutil"
	"strings"
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

var PrimeGuild int
var LastCheck *LastChangeStatus
var Changes map[string]map[string]string

func CheckChange(G *GuildInfo, GID string) *discordgo.Guild {
	ChGl, err := sh.dg.State.Guild(GID)
	if err != nil {
		fmt.Println("Error getting guild: " + err.Error())
	}
	if ChGl.Channels != G.g.Channels {
		LastCheck.Changeyes = true
		LastCheck.Change.channel = true
	}
	if ChGl.Members != G.g.Members || ChGl.MemberCount != G.g.MemberCount {
		LastCheck.Changeyes = true
		LastCheck.Change.members = true
	}
	if ChGl.Roles != G.g.Roles {
		LastCheck.Changeyes = true
		LastCheck.Change.roles = true
	}
	if ChGl.Presences != G.g.Presences {
		LastCheck.Change.yes = true
		LastCheck.Change.presence = true
	}
	if ChGl.Icon != G.g.Icon || ChGl.Large != G.g.Large || ChGl.Name != G.g.Name || ChGl.OwnerID != G.g.OwnerID || ChGl.Region != G.g.Region {
		LastCheck.Changeyes = true
		LastCheck.Change.beginvars = true
	}
}

// compare vals here

func CompareChannel(Last *discordgo.Guild, New *discordgo.Guild) (bool, Changes) {
	var ChS *Changes
	var Ident bool
	Ident = false
	if Last.ID == New.ID {
		Ident = true
		if Last.Name != New.Name {

		}
	}
	return
}

// rest funcs

func AppendChange(G *discordgo.Guild) []string {
	Resume := LastCheck.Firstcheck
	var ChangeString []string
	var NewString string
	if LastCheck.Change.channel == true {
		LastChannels := LastCheck.GI.g.Channels
		NewChannels := G.Channels
		if len(LastChannels) != len(NewChannels) {
			// code for "new channel" here >.>
		} else {
			for _, LCh := range LastChannels {
				for _, NCh := range NewChannels {
					Yes, Changes := CompareChannel(LCh, NCh)
				}
			}
		}
	}
	ChangeString = append(ChangeString, NewString)
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
}

func Initnewguild(GID string) {
	fmt.Println("No guild file for '" + "' found!\nCreating new guild files...")
	BGuild, err := sh.dg.State.Guild(GID)
	if err == nil {
		GLD := &GuildInfo{
			Lastcheck: GetTime(),
			g:         BGuild,
			BotUP:     true,
		}
		WriteGLDfile(GLD, true)
		WriteGLDfile(GLD, false)
	}
	fmt.Println("New guild file for " + BGuild.Name + " made, ready for logging!")
}

func ResumeCheck(Gid string) {
	G := GetGLDfile(Gid)
	if G.g == nil {
		Initnewguild(Gid)
		return
	}
	if G.BotUP == true {
		LastCheck.SiLaChBisTerm = true
		fmt.Println("Warning: Sherlock hasn't been closed properly since last boot!")
	} else {
		LastCheck.SiLaChBisTerm = false
	}
	LastCheck.Firstcheck = true
	NewChange := CheckChange(G, Gid)
	if LastCheck.Changeyes == true {
		AppendChange(NewChange)
	}
}

type GuildInfo struct {
	g         *discordgo.Guild
	Lastcheck struct {
		Year  int
		Month *time.Month
		Day   int
		Hour  int
		Min   int
		Sec   int
	}
	BotUP bool
}

type LastChangeStatus struct {
	GI            *GuildInfo
	SiLaChBisTerm bool
	Firstcheck    bool
	Changeyes     bool
	Change        struct {
		beginvars bool
		roles     bool
		members   bool
		presence  bool
		channel   bool
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
	if Mes.Content[0] == '!' {
		ProcessCMD(Mes.Content[1:], Mes)
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

func GetTime() (int, *time.Month, int, int, int, int) {
	Year := time.Now().Year()
	Month := time.Now().Month()
	Day := time.Now().Day()
	Hour := time.Now().Hour()
	Min := time.Now().Minute()
	Sec := time.Now().Second()
}

func GetGLDfile(GID string) GuildInfo {
	DATA, err := ioutil.ReadFile(GID + ".GLD")
	if err == nil {
		var LLG GuildInfo
		err := json.Unmarshal(DATA, LLG)
		if err == nil {
			return LLG
		} else {
			fmt.Println("Error Unmarshal-ing GLD: " + err.Error())
		}
	}
	var Errorstring string
	var NoGFile bool
	if err != nil {
		Errorstring = err.Error()
		NoGFile = strings.ContainsAny(GID, Errorstring)
	} else {
		NoGFile = false
	}
	if err != nil {
		if NoGFile == true {
			return GuildInfo{}
		} else {
			fmt.Println("Error reading GLD file: " + Errorstring)
		}
	}
	return GuildInfo{}
}

func WriteGLDfile(G *GuildInfo, Isb bool) error {
	GID, GIDerr := json.Marshal(G)
	if GIDerr == nil {
		if Isb == false {
			ioutil.WriteFile(G.g.ID+".GLD", GID, 0777)
		}
		if Isb == true {
			ioutil.WriteFile("B-"+G.g.ID+".GLD", GID, 0777)
		}
		return nil
	} else {
		fmt.Println("GLD writing error: " + GIDerr.Error())
		return GIDerr
	}
}

// init

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
	sh.dg.AddHandler(BBCreateMessage)

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
