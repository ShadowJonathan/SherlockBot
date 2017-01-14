package Belt

import (
	//"bytes"
	"encoding/json"
	fmt "fmt"
	"io/ioutil"
	"reflect"
	"strconv"
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

var PrimeGuild string
var LastCheck = &LastChangeStatus{}

var notifiers []string

type ChangeMap map[string]map[string]map[string]map[string]string

// (changed type) (ID) (type change) (old val/misc):new val/misc

var Changes = make(map[string]map[string]map[string]map[string]string)
var SubChanges = make(map[string]map[string]map[string]string)

func CheckChange(G *GuildInfo, GID string) (*discordgo.Guild, bool) {
	ChGl, err := sh.dg.State.Guild(GID)
	if err != nil {
		fmt.Println("Error getting guild: " + err.Error())
	}
	LastCheck.Changeyes = false
	if !reflect.DeepEqual(ChGl, G.g) {
		fmt.Println("Change detected")
	}
	if !reflect.DeepEqual(ChGl.Channels, G.g.Channels) {
		LastCheck.Changeyes = true
		LastCheck.Change.channel = true
	}
	if !reflect.DeepEqual(ChGl.Members, G.g.Members) || ChGl.MemberCount != G.g.MemberCount {
		LastCheck.Changeyes = true
		LastCheck.Change.members = true
	}
	if !reflect.DeepEqual(ChGl.Roles, G.g.Roles) {
		LastCheck.Changeyes = true
		LastCheck.Change.roles = true
	}
	if !reflect.DeepEqual(ChGl.Presences, G.g.Presences) {
		LastCheck.Changeyes = true
		LastCheck.Change.presence = true
	}
	if ChGl.Icon != G.g.Icon || ChGl.Large != G.g.Large || ChGl.Name != G.g.Name || ChGl.OwnerID != G.g.OwnerID || ChGl.Region != G.g.Region {
		LastCheck.Changeyes = true
		LastCheck.Change.beginvars = true
	}
	return ChGl, LastCheck.Changeyes
}

// compare vals here

func CompareChannel(Last *discordgo.Channel, New *discordgo.Channel) (bool, map[string]map[string]map[string]string) {
	var ChS = make(map[string]map[string]map[string]string)
	var Ident bool
	Ident = false
	if Last.ID == New.ID {
		Ident = true
		if !reflect.DeepEqual(Last.Name, New.Name) {
			ChS[Last.ID]["Name"][Last.Name] = New.Name
		}
		if !reflect.DeepEqual(Last.Topic, New.Topic) {
			ChS[Last.ID]["Topic"][Last.Topic] = New.Topic
		}
		if !reflect.DeepEqual(Last.PermissionOverwrites, New.PermissionOverwrites) {
			for _, LPerm := range Last.PermissionOverwrites {
				for _, NPerm := range New.PermissionOverwrites {
					if reflect.DeepEqual(LPerm.ID, NPerm.ID) {
						if !reflect.DeepEqual(LPerm.Allow, NPerm.Allow) {
							var LAllow int64
							LAllow = int64(LPerm.Allow)
							var NAllow int64
							NAllow = int64(NPerm.Allow)
							ChS[Last.ID]["AllowPerm"][LPerm.ID] = LPerm.Type + " " + LPerm.ID + " " + strconv.FormatInt(LAllow, 10) + " " + strconv.FormatInt(NAllow, 10)

						}
						if reflect.DeepEqual(LPerm.Deny, NPerm.Deny) {
							var LDeny int64
							LDeny = int64(LPerm.Deny)
							var NDeny int64
							NDeny = int64(NPerm.Deny)
							ChS[Last.ID]["DenyPerm"][LPerm.ID] = LPerm.Type + " " + LPerm.ID + " " + strconv.FormatInt(LDeny, 10) + " " + strconv.FormatInt(NDeny, 10)

						}
					}
					// delete
					if len(Last.PermissionOverwrites) < len(New.PermissionOverwrites) {
						var check bool
						check = false
						for _, PO := range Last.PermissionOverwrites {
							for _, PA := range New.PermissionOverwrites {
								if PO.ID == PA.ID {
									check = true
								}
							}
							if check == false {
								var LDeny int64
								LDeny = int64(PO.Deny)
								var LAllow int64
								LAllow = int64(PO.Allow)
								ChS[Last.ID]["PermPoof"][PO.ID] = PO.Type + " " + PO.ID + " " + " " + strconv.FormatInt(LAllow, 10) + " " + strconv.FormatInt(LDeny, 10)
							}
						}
					}
					// create
					if len(Last.PermissionOverwrites) > len(New.PermissionOverwrites) {
						var check bool
						check = false
						var PO *discordgo.PermissionOverwrite
						for _, PA := range Last.PermissionOverwrites {
							for _, PO = range New.PermissionOverwrites {
								if PO.ID == PA.ID {
									check = true
								}
							}
							if check == false {
								var NDeny int64
								NDeny = int64(PO.Deny)
								var NAllow int64
								NAllow = int64(PO.Allow)
								ChS[Last.ID]["PermWoah"][PO.ID] = PO.Type + " " + strconv.FormatInt(NAllow, 10) + " " + strconv.FormatInt(NDeny, 10)
							}
						}
					}
				}
			}
		}
	}
	return Ident, ChS
}

// rest funcs

func AppendChange(G *discordgo.Guild) []string {
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
					Yes, ChChanges := CompareChannel(LCh, NCh)
					if Yes {
						Changes["Channel"] = ChChanges
					}
				}
			}
		}
	}
	for i := 0; i <= len(Changes); i++ {
		for i, ChCh := range Changes {
			if reflect.DeepEqual(Changes[i], Changes["Channel"]) { //(ID) (type change) (old val/misc):new val/misc
				ChannelChange := ParseChannelChange(ChCh, LastCheck.GI.g)
				ChangeString = AppendSSlices(ChangeString, ChannelChange)
			}
		}
	}

	ChangeString = append(ChangeString, NewString)
	return ChangeString
}

func ParseChannelChange(ChCh map[string]map[string]map[string]string, G *discordgo.Guild) []string {
	//(ID) (type change) (old val/misc):new val/misc
	var ReturnStrings []string
	for _, check := range G.Channels {
		if len(ChCh[check.ID]) > 0 {
			var CheckID = ChCh[check.ID]
			// Name Topic DenyPerm AllowPerm PermPoof PermWoah
			for i, CheckType := range CheckID {
				if reflect.DeepEqual(CheckID[i], CheckID["Name"]) {
					ReturnStrings = append(ReturnStrings, "Channel '"+check.ID+"' just changed it's name from '"+check.Name+"' to '"+CheckType[check.Name]+"'!")
				}
				if reflect.DeepEqual(CheckID[i], CheckID["Topic"]) {
					ReturnStrings = append(ReturnStrings, "Channel '"+check.Name+"' just changed it's topic from '"+check.Topic+"' to '"+CheckType[check.Topic])
				}
				if reflect.DeepEqual(CheckID[i], CheckID["AllowPerm"]) {
					for F, AllowMap := range CheckID["AllowPerm"] {
						for _, ChannelPerm := range check.PermissionOverwrites {
							if reflect.DeepEqual(CheckID["AllowPerm"][F], CheckID["AllowPerm"][ChannelPerm.ID]) {
								var Perm = AllowMap
								// type ID Last New
								Perms := strings.Split(Perm, " ")
								if Perms[0] == "role" {
									Role := GetRole(Perms[1], G)
									ReturnStrings = append(ReturnStrings, "Role '"+Role.Name+"' on channel '"+check.Name+"' just got it's allowed permissions changed from `"+Perms[2]+"` to `"+Perms[3]+"`!")
								}
							}
						}
					}
				}
				if reflect.DeepEqual(CheckID[i], CheckID["DenyPerm"]) {
					for F, DenyMap := range CheckID["DenyPerm"] {
						for _, ChannelPerm := range check.PermissionOverwrites {
							if reflect.DeepEqual(CheckID["DenyPerm"][F], CheckID["DenyPerm"][ChannelPerm.ID]) {
								var Perm = DenyMap
								// type ID Last New
								Perms := strings.Split(Perm, " ")
								if Perms[0] == "role" {
									Role := GetRole(Perms[1], G)
									ReturnStrings = append(ReturnStrings, "Role '"+Role.Name+"' on channel '"+check.Name+"' just got it's denied permissions changed from `"+Perms[2]+"` to `"+Perms[3]+"`!")
								}
							}
						}
					}
				}
				if reflect.DeepEqual(CheckID[i], CheckID["PermPoof"]) {

				}
				if reflect.DeepEqual(CheckID[i], CheckID["PermWoah"]) {

				}
			}
		}
	}
	return ReturnStrings

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

func CheckLoop(LC *LastChangeStatus, Gid string) {
	LastCheck = LC
	var G *GuildInfo
	for sh.Stop == false {
		time.Sleep(30 * time.Second)
		fmt.Println("Loop debug")
		G = LastCheck.GI
		NewChange, change := CheckChange(G, Gid)
		if change {
			_, _, _, H, Mi, S := GetTime()
			h := strconv.FormatInt(int64(H), 10)
			mi := strconv.FormatInt(int64(Mi), 10)
			s := strconv.FormatInt(int64(S), 10)

			Responses := AppendChange(NewChange)
			for _, N := range notifiers {
				sh.dg.ChannelMessageSend(N, "-"+h+":"+mi+";"+s+"-")
				for _, R := range Responses {
					sh.dg.ChannelMessageSend(N, R)
				}
			}
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
	LastCheck.Firstcheck = true
	NewChange, change := CheckChange(G, Gid)
	for _, N := range notifiers {
		sh.dg.ChannelMessageSend(N, "Sherlock's out of 221b, ready for some investigation!")
	}
	if change {
		Responses := AppendChange(NewChange)
		for _, N := range notifiers {
			sh.dg.ChannelMessageSend(N, "New Changes in since last taxi ride:")
			for _, R := range Responses {
				sh.dg.ChannelMessageSend(N, R)
			}
		}
	}
	LastCheck.GI.g = NewChange
	LastCheck.GI.Lastcheck.Year, LastCheck.GI.Lastcheck.Month, LastCheck.GI.Lastcheck.Day, LastCheck.GI.Lastcheck.Hour, LastCheck.GI.Lastcheck.Min, LastCheck.GI.Lastcheck.Sec = GetTime()
	_ = WriteGLDfile(LastCheck.GI, false)
}

type GuildInfo struct {
	g           *discordgo.Guild
	Lastcheck   TimeFormat
	BotUP       bool
	NeedRestall bool
}

type TimeFormat struct {
	Year  int
	Month time.Month
	Day   int
	Hour  int
	Min   int
	Sec   int
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

func GetTime() (int, time.Month, int, int, int, int) {
	Year := time.Now().Year()
	Month := time.Now().Month()
	Day := time.Now().Day()
	Hour := time.Now().Hour()
	Min := time.Now().Minute()
	Sec := time.Now().Second()
	return Year, Month, Day, Hour, Min, Sec
}

func GetGLDfile(GID string) (*GuildInfo, error) {
	DATA, err := ioutil.ReadFile(GID + ".GLD")
	var LLG *GuildInfo
	Y, Mo, D, H, Mi, S := GetTime()
	TG, _ := sh.dg.State.Guild(GID)
	LLG = &GuildInfo{
		Lastcheck:   TimeFormat{Year: Y, Month: Mo, Day: D, Hour: H, Min: Mi, Sec: S},
		g:           TG,
		BotUP:       true,
		NeedRestall: false,
	}
	if err == nil {
		err := json.Unmarshal(DATA, LLG)
		if err == nil {
			return LLG, nil
		} else {
			fmt.Println("Error Unmarshal-ing GLD: " + err.Error())
			return LLG, err
		}
	}
	LLG = &GuildInfo{
		NeedRestall: true,
	}
	return LLG, nil
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

func GetRole(RID string, g *discordgo.Guild) *discordgo.Role {
	for _, R := range g.Roles {
		if RID == R.ID {
			return R
		}
	}
	var Defaultrole = &discordgo.Role{
		ID:          g.ID,
		Name:        "@everyone",
		Managed:     false,
		Mentionable: false,
		Hoist:       false,
		Color:       0,
		Position:    0,
		Permissions: 0,
	}
	return Defaultrole
}

func AppendSSlices(BeginSlice []string, MergeSlice []string) []string {
	var ProcessSlice = BeginSlice
	for _, MergeString := range MergeSlice {
		ProcessSlice = append(ProcessSlice, MergeString)
	}
	return ProcessSlice
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
