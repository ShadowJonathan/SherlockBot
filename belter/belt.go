package Belt

import (
	fmt "fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Version struct {
	Major byte
	Minor byte
}

type Sherlock struct {
	dg        *discordgo.Session
	Debug     bool
	version   Version
	OwnID     string
	OwnAV     string
	OwnName   string
	Stop      bool
	StopLoop  bool
	Notifiers []string
}

// Vars after this

var sh *Sherlock

// Functions after this

var PrimeGuild string

var notifiers []string

func CheckChange(G *GuildInfo, GID string) (bool, *discordgo.Guild, *FullChangeStruct) {
	ChGl, err := GetGuild(GID)
	if err != nil {
		fmt.Println("Error getting guild: " + err.Error())
	}
	if G.g == nil || ChGl == nil {
		fmt.Println("!!CRITICAL ERROR!!\nOne of the two guild values given to check is nil.")
		if G.g == nil {
			fmt.Println("Guild value that is nil: G.g")
			panic(G.g)
		}
		if ChGl == nil {
			fmt.Println("Guild value that is nil: ChGl")
			panic(ChGl)
		}
	}
	Equal, FCS := DeepEqual(G.g, ChGl)
	return Equal, ChGl, FCS
}

// rest funcs

func AppendChange(Gold *discordgo.Guild, Gnew *discordgo.Guild, TotC *FullChangeStruct) []string {
	var ChangeString []string
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
						} else {
							if P.Mk {
								Pmk := GetOR(P.ID, Cnew)
								if Pmk.Type == "role" {
									PR := GetRole(Pmk.ID, Gnew)
									ChangeString = append(ChangeString, "New overwrite permission for role "+PR.Name+" has been granted for channel "+Cnew.Name+"!")
								} else {
									PM := GetUser(Pmk.ID, Gnew)
									ChangeString = append(ChangeString, "New overwrite permission for member "+PM.Username+" has been granted for channel "+Cnew.Name+"!")
								}
							} else {
								Pdel := GetOR(P.ID, Cold)
								if Pdel.Type == "role" {
									PR := GetRole(Pdel.ID, Gold)
									ChangeString = append(ChangeString, "The overwritten permission for role "+PR.Name+" has been revoked for channel "+Cold.Name+"!")
								} else {
									PM := GetUser(Pdel.ID, Gold)
									ChangeString = append(ChangeString, "The overwritten permission for member "+PM.Username+" has been revoked for channel "+Cold.Name+"!")
								}
							}
						}
					}
				}
			} else {
				if CCh.Mk {
					Cmk := GetChannel(CCh.ID, Gnew)
					ChangeString = append(ChangeString, "New channel created!\nName: "+Cmk.Name)
				}
				if CCh.Del {
					Cdel := GetChannel(CCh.ID, Gold)
					ChangeString = append(ChangeString, "Channel "+Cdel.Name+" has been deleted.")
				}
			}
		}
	}
	if TotC.Guild.roles {
		for _, Rch := range TotC.Roles {
			if !Rch.ExistCrisis {
				Ro := GetRole(Rch.ID, Gold)
				Rn := GetRole(Rch.ID, Gnew)
				if Rch.Name {
					ChangeString = append(ChangeString, "Role "+Ro.Name+" has changed it's name to "+Rn.Name)
				}
				if Rch.Color {
					ChangeString = append(ChangeString, "Role "+Rn.Name+" has changed it's color from `"+strconv.FormatInt(int64(Ro.Color), 16)+"` to `"+strconv.FormatInt(int64(Rn.Color), 16)+"`")
				}
				if Rch.Perms {
					ChangeString = append(ChangeString, "Role "+Rn.Name+" has changed it's permissions from `"+strconv.FormatInt(int64(Ro.Permissions), 10)+"` to `"+strconv.FormatInt(int64(Rn.Permissions), 10)+"`")
				}
				if Rch.Position {
					ChangeString = append(ChangeString, "Role "+Rn.Name+"'s position has been changed, previously: "+strconv.FormatInt(int64(Ro.Position), 10)+", now: "+strconv.FormatInt(int64(Rn.Position), 10))
				}
			} else {
				if Rch.Mk {
					Nr := GetRole(Rch.ID, Gnew)
					ChangeString = append(ChangeString, "New role added! ID: `"+Nr.ID+"`")
				}
				if Rch.Del {
					Dr := GetRole(Rch.ID, Gold)
					ChangeString = append(ChangeString, "Role "+Dr.Name+" has been deleted")
				}
			}
		}
	}
	if TotC.Guild.members {
		for _, Mch := range TotC.Members {
			if !Mch.ExistCrisis {
				Mold := GetMember(Mch.User.ID, Gold)
				Mnew := GetMember(Mch.User.ID, Gnew)
				if Mch.Nick {
					ChangeString = append(ChangeString, "Member "+Mold.User.Username+" changed his/her nickname from "+Mold.Nick+" to "+Mnew.Nick)
				}
				if Mch.User.Username {
					ChangeString = append(ChangeString, "Member "+Mold.User.Username+" changed his/her username to "+Mnew.User.Username)
				}
				if Mch.User.Avatar {
					ChangeString = append(ChangeString, "Member "+Mold.User.Username+" changed his avatar from "+discordgo.EndpointUserAvatar(Mold.User.ID, Mold.User.Avatar)+"\nTo "+discordgo.EndpointUserAvatar(Mnew.User.ID, Mnew.User.Avatar))
				}
				if Mch.Roles {
					if Mch.RoleNew {
						var isOld bool
						for _, Rn := range Mnew.Roles {
							if len(Mold.Roles) != 0 {
								isOld = false
								for _, Ro := range Mold.Roles {
									if Rn == Ro {
										isOld = true
									}
								}
								if !isOld {
									NR := GetRole(Rn, Gnew)
									ChangeString = append(ChangeString, "Member "+Mnew.User.Username+" got the "+NR.Name+" role!")
								}
							} else {
								NR := GetRole(Rn, Gnew)
								ChangeString = append(ChangeString, "Member "+Mnew.User.Username+" got the "+NR.Name+" role!")
							}
						}
					}
					if Mch.RoleRem {
						var IsThere bool
						for _, Ro := range Mold.Roles {
							if len(Mnew.Roles) != 0 {
								IsThere = false
								for _, Rn := range Mnew.Roles {
									if Rn == Ro {
										IsThere = true
									}
								}
								if !IsThere {
									OR := GetRole(Ro, Gold)
									ChangeString = append(ChangeString, "Member "+Mnew.User.Username+" doesnt have the "+OR.Name+" role anymore!")
								}
							} else {
								OR := GetRole(Ro, Gold)
								ChangeString = append(ChangeString, "Member "+Mnew.User.Username+" doesnt have the "+OR.Name+" role anymore!")
							}
						}
					}
				}
			} else {
				if Mch.Join {
					NM := GetUser(Mch.User.ID, Gnew)
					ChangeString = append(ChangeString, "Member "+NM.Username+" has joined "+Gnew.Name+"!")
				}
				if Mch.Leave {
					OM := GetUser(Mch.User.ID, Gold)
					ChangeString = append(ChangeString, "Member "+OM.Username+" has left "+Gnew.Name+"!")
				}
			}
		}
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
	sh.StopLoop = true
	PrimeGuild = string(GLDB)
	LastCheck := ResumeCheck(PrimeGuild)
	sh.StopLoop = false
	go CheckLoop(PrimeGuild, LastCheck)
}

func CheckLoop(Gid string, LastCheck *LastChangeStatus) {
	for sh.StopLoop == false {
		time.Sleep(45 * time.Second)
		GI, err := GetGLDfile(LastCheck.GI.g.ID)
		fmt.Println("Loop debug")
		if err != nil {
			fmt.Println("Error getting GLD file: " + err.Error())
		}
		LastCheck.GI = GI
		IsEqual, NewGuild, AllChange := CheckChange(LastCheck.GI, Gid)
		if !IsEqual {
			_, _, _, H, Mi, S := GetTime()
			h := strconv.FormatInt(int64(H), 10)
			mi := strconv.FormatInt(int64(Mi), 10)
			s := strconv.FormatInt(int64(S), 10)

			Responses := AppendChange(LastCheck.GI.g, NewGuild, AllChange)
			for _, N := range notifiers {
				SendMessage(N, "-"+h+":"+mi+";"+s+"-", sh.Notifiers)
				for _, R := range Responses {
					SendMessage(N, R, sh.Notifiers)
				}
			}
		}
		if NewGuild == nil {
			panic(NewGuild)
		}
		LastCheck = SetLC(LastCheck, NewGuild)
		err = WriteGLDfile(LastCheck.GI, false)
		if err != nil {
			fmt.Println("Error writing GLD file after check: " + err.Error())
			sh.StopLoop = true
		}
	}
	WriteGLDfile(LastCheck.GI, false)
}

func ResumeCheck(Gid string) *LastChangeStatus {
	var LastCheck = &LastChangeStatus{}
	G, err := GetGLDfile(Gid)
	if err != nil || G.NeedRestall {
		if err != nil {
			fmt.Println(err.Error())
		}
		GI := Initnewguild(Gid)
		LastCheck.GI = GI
		return LastCheck
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
		SendMessage(N, "Sherlock's out of 221b, ready for some investigation!", sh.Notifiers)
	}
	if !IsEqual {
		Responses := AppendChange(LastCheck.GI.g, NewGuild, AllChange)
		if len(Responses) == 0 {
			panic(Responses)
		}
		for _, N := range notifiers {
			SendMessage(N, "New Changes in since last taxi ride:", sh.Notifiers)
			for _, R := range Responses {
				SendMessage(N, R, sh.Notifiers)
			}
		}
	}
	LastCheck = SetLC(LastCheck, NewGuild)
	err = WriteGLDfile(LastCheck.GI, false)
	if err != nil {
		fmt.Println("Error writing GLD file: " + err.Error())
	}
	return LastCheck
}

func Initnewguild(GID string) *GuildInfo {
	fmt.Println("No guild file for '" + GID + "' found!\nCreating new guild files...")
	BGuild, err := GetGuild(GID)
	if err == nil {
		Y, Mo, D, H, Mi, S := GetTime()
		GLD := &GuildInfo{
			Lastcheck: TimeFormat{Year: Y, Month: Mo, Day: D, Hour: H, Min: Mi, Sec: S},
		}
		if BGuild == nil {
			panic(BGuild)
		}
		GLD.g = BGuild
		GLD.BotUP = true
		GLD.NeedRestall = false
		WriteGLDfile(GLD, true)
		WriteGLDfile(GLD, false)
		fmt.Println("New guild file for " + BGuild.Name + " made, ready for logging!")
		return GLD
	} else {
		fmt.Println("Error getting guild status: " + err.Error())
	}
	var GLD *GuildInfo
	return GLD
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
			ProcessCMD(Mes.Content[1:], MesC.Message, sh.Notifiers)
		}
	}
}

// misc funx

func ProcessCMD(CMD string, M *discordgo.Message, Notifiers []string) {
	Commands := getCMD(CMD)
	var SecArg string
	if len(Commands) > 1 {
		SecArg = Commands[1]
	}
	if strings.ToLower(Commands[0]) == "primeguild" {
		PID := Commands[1]
		if PID == "" {
			SendMessage(M.ChannelID, "You gave me a nil ID!", Notifiers)
		} else {
			data := []byte(PID)
			ioutil.WriteFile("PrimeGuild", data, 9000)
			fmt.Println("Set new Prime Guild to '" + PID + "'")
			SendMessage(M.ChannelID, "`Set new prime guild to "+PID+"`", sh.Notifiers)
		}
	}
	if strings.ToLower(Commands[0]) == "stopcheck" {
		if !sh.StopLoop {
			sh.StopLoop = true
			fmt.Println("Checking loop stopped")
			SendMessage(M.ChannelID, "`Stopped checking loop`", sh.Notifiers)
		}
		if sh.StopLoop {
			SendMessage(M.ChannelID, "`Checking loop isn't running!`", sh.Notifiers)
		}
	}
	if strings.ToLower(Commands[0]) == "kickstart" {
		StartCheckLoop()
		SendMessage(M.ChannelID, "`Checking loop restarted`", sh.Notifiers)
	}
	if strings.ToLower(Commands[0]) == "getuser" {
		SendMessage(M.ChannelID, GMstring(SecArg), sh.Notifiers)
	}
	if strings.ToLower(Commands[0]) == "getchannel" {
		SendMessage(M.ChannelID, GCstring(SecArg), sh.Notifiers)
	}
	if strings.ToLower(Commands[0]) == "getguild" {
		PG, err := ioutil.ReadFile("PrimeGuild")
		if err != nil {
			fmt.Println("Error reading PG file: " + err.Error())
			return
		}
		GG, err := GetGuild(string(PG))
		if len(Commands) < 1 {
			Commands[1] = GG.ID
		}
		var Channels []string
		for _, Ch := range GG.Channels {
			Channels = append(Channels, GCstring(Ch.ID))
		}
		SendChannel, err := sh.dg.UserChannelCreate(M.Author.ID)
		if err != nil {
			return
		}
		var not []string
		not = append(not, SendChannel.ID)
		Owner := GetUser(GG.OwnerID, GG)
		var own string
		own = Owner.ID
		SendMessage(SendChannel.ID, "`Guild:`\n`ID: "+GG.ID+"`\n`Name: "+GG.Name+"`\n`Region: "+GG.Region+"`\n`Icon: `"+discordgo.EndpointGuildIcon(GG.ID, GG.Icon), not)
		SendMessage(SendChannel.ID, "`Owner`\n"+GMstring(own), not)
		SendMessage(SendChannel.ID, "`Channels:`", not)
		for _, CHS := range Channels {
			SendMessage(SendChannel.ID, CHS, not)
		}
	}
}

func DeepEqual(a *discordgo.Guild, b *discordgo.Guild) (bool, *FullChangeStruct) {
	var Equal = true
	var TotC = &FullChangeStruct{}
	Equal, TotC = CompareGuild(a, b, TotC, Equal)
	return Equal, TotC
}

// init

func Initialize(Token string) {
	isdebug, err := ioutil.ReadFile("debugtoggle")
	sh = &Sherlock{
		version:  Version{1, 0},
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
	// "269441095862190082", "270639478379511809"
	notifiers, err = GetNotifiers()
	if err != nil {
		fmt.Println("Error getting Notifier file: " + err.Error())
	} else {
		sh.Notifiers = notifiers
	}
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
