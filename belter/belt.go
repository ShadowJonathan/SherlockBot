package Belt

import (
	fmt "fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"log"

	"encoding/json"

	"../versions"
	"github.com/bwmarrin/discordgo"
)

type Version struct {
	Major byte
	Minor byte
}

type Sherlock struct {
	dg            *discordgo.Session
	Debug         bool
	version       versions.Version
	OwnID         string
	OwnAV         string
	OwnName       string
	Stop          bool
	StopLoop      bool
	Notifiers     []string
	PrimeSuspects []string
	LoopCooldown  time.Duration
	cl            *ChatLog
}

// Vars after this

var sh *Sherlock

var PER *PermissionBit // do not change this

var restart bool

var upgrade bool

var CheckerCount int

// Functions after this

var PrimeGuild string

var notifiers []string

var timetoswitch int

func CheckChange(G *GuildInfo, GID string) (bool, *discordgo.Guild, *FullChangeStruct) {
	var err error
	var ChGl *discordgo.Guild
	if timetoswitch > 2 {
		ChGl, err = GetGuildState(GID)
		timetoswitch = 0
	} else {
		ChGl, err = GetGuild(GID)
		timetoswitch++
	}
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

var count int

func SaveDebug(data []byte) {
	for {
		_, err := ioutil.ReadFile("debug/data" + strconv.Itoa(count) + ".json")
		if err == nil {
			count++
			continue
		} else {
			ioutil.WriteFile("debug/data"+strconv.Itoa(count)+".json", data, 0777)
			count++
			return
		}
	}
}
func AppendChange(Gold *discordgo.Guild, Gnew *discordgo.Guild, TotC *FullChangeStruct) []string {
	var ChangeString []string
	b, _ := json.Marshal(TotC)
	fmt.Println(string(b))
	SaveDebug(b)
	sh.dg.ChannelMessageSend("269441095862190082", "```json\n"+string(b)+"\n```")
	if TotC.Guild.Name {
		ChangeString = append(ChangeString, "Server **"+Gold.Name+"** changed it's name to **"+Gnew.Name+"**!")
		Log("S " + Gold.ID + ": NEWNAME: " + Gold.Name + "->" + Gnew.Name)
	}
	if TotC.Guild.OwnerID {
		ChangeString = append(ChangeString, "`ALERT!`\n`ALERT!`\n`ALERT!`\nServer **"+Gold.Name+"'s** owner just changed from **"+GetUserName(Gold.OwnerID, Gold)+"** to **"+GetUserName(Gnew.OwnerID, Gnew)+"**!!!")
		Log("S " + Gold.ID + ": NEWOWNER: " + Gold.OwnerID + "->" + Gnew.OwnerID)
	}
	if TotC.Guild.Icon {
		ChangeString = append(ChangeString, "Server **"+Gold.Name+"** changed it's icon!")
		Log("S " + Gold.ID + ": NEWICON: " + Gold.Icon + "->" + Gnew.Icon)
	}
	if TotC.Guild.Region {
		ChangeString = append(ChangeString, "Server **"+Gold.Name+"** changed it's region!\nPreviously: **"+Gold.Region+"**\nNow: **"+Gnew.Region+"**")
		Log("S " + Gold.ID + ": NEWREGION: " + Gold.Region + "->" + Gnew.Region)
	}
	if TotC.Guild.Channels {
		for _, CCh := range TotC.Channels {
			if !CCh.ExistCrisis {
				Cold := GetChannel(CCh.ID, Gold)
				Cnew := GetChannel(CCh.ID, Gnew)
				if CCh.Name {
					ChangeString = append(ChangeString, "Channel **"+Cold.Name+"** changed it's name to "+Cnew.Name+"!")
					Log("CH " + Cold.ID + ": NEWNAME: " + Cold.Name + "->" + Gnew.Name)
				}
				if CCh.Topic {
					ChangeString = append(ChangeString, "Channel **"+Cold.Name+"** changed topics from "+SanCode(Cold.Topic)+" to "+SanCode(Cnew.Topic)+"!")
					Log("CH " + Cold.ID + ": NEWTOPIC:\n" + Cold.Topic + "\n->\n" + Cnew.Topic)
				}
				if CCh.perms {
					for _, P := range CCh.Perms {
						if !P.ExistCrisis {
							Oor := GetOR(P.ID, Cold)
							Nor := GetOR(P.ID, Cnew)
							if Oor.Type == "role" {
								Ror := GetRole(P.ID, Gold)

								if P.Allow {
									ChangeString = append(ChangeString, "Channel **"+Cnew.Name+"** has changed it's allowed overwrite-permissions for the **"+Ror.Name+"** role changed from `"+intToString(Oor.Allow)+"` to `"+intToString(Nor.Allow)+"`")
									Log("CH " + Cold.ID + ": OR-ROLE" + Oor.ID + ": ALLOW: " + intToString(Oor.Allow) + "->" + intToString(Nor.Allow))
								}
								if P.Deny {
									ChangeString = append(ChangeString, "Channel **"+Cnew.Name+"** has changed it's denied overwrite-permissions for the "+Ror.Name+" role changed from `"+intToString(Oor.Deny)+"` to `"+intToString(Nor.Deny)+"`")
									Log("CH " + Cold.ID + ": OR-ROLE" + Oor.ID + ": DENY: " + intToString(Oor.Allow) + "->" + intToString(Nor.Deny))
								}
							} else {
								Mor := GetUser(P.ID, Gold)

								if P.Allow {
									ChangeString = append(ChangeString, "Channel **"+Cnew.Name+"** has changed it's allowed overwrite-permissions for "+Mor.Username+" changed from `"+intToString(Oor.Allow)+"` to `"+intToString(Nor.Allow)+"`")
									Log("CH " + Cold.ID + ": OR-USER " + Mor.ID + ": ALLOW: " + intToString(Oor.Allow) + "->" + intToString(Nor.Allow))
								}
								if P.Deny {
									ChangeString = append(ChangeString, "Channel **"+Cnew.Name+"** has changed it's denied overwrite-permissions for "+Mor.Username+" changed from `"+intToString(Oor.Deny)+"` to `"+intToString(Nor.Deny)+"`")
									Log("CH " + Cold.ID + ": OR-USER " + Mor.ID + ": DENY: " + intToString(Oor.Deny) + "->" + intToString(Nor.Deny))
								}
							}
						} else {
							if P.Mk {
								Pmk := GetOR(P.ID, Cnew)
								if Pmk.Type == "role" {
									PR := GetRole(Pmk.ID, Gnew)
									ChangeString = append(ChangeString, "New overwrite permission for role **"+PR.Name+"** has been granted for channel **"+Cnew.Name+"**!")
									Log("CH " + Cold.ID + ": OR ROLE " + PR.ID + ": NEW")
								} else {
									PM := GetUser(Pmk.ID, Gnew)
									ChangeString = append(ChangeString, "New overwrite permission for member **"+PM.Username+"** has been granted for channel **"+Cnew.Name+"**!")
									Log("CH " + Cold.ID + ": OR USER " + PM.ID + ": NEW")
								}
							} else {
								Pdel := GetOR(P.ID, Cold)
								if Pdel.Type == "role" {
									PR := GetRole(Pdel.ID, Gold)
									ChangeString = append(ChangeString, "The overwritten permission for role **"+PR.Name+"** has been revoked for channel **"+Cold.Name+"**!")
									Log("CH " + Cold.ID + ": OR ROLE " + PR.ID + ": NEW")
								} else {
									PM := GetUser(Pdel.ID, Gold)
									ChangeString = append(ChangeString, "The overwritten permission for member **"+PM.Username+"** has been revoked for channel **"+Cold.Name+"**!")
									Log("CH " + Cold.ID + ": OR USER " + PM.ID + ": DELETE")
								}
							}
						}
					}
				}
			} else {
				if CCh.Mk {
					Cmk := GetChannel(CCh.ID, Gnew)
					ChangeString = append(ChangeString, "New channel created!\nName: **"+Cmk.Name+"**")
					Log("S " + Gold.ID + ": NEWCHANNEL: " + Cmk.ID)
				}
				if CCh.Del {
					Cdel := GetChannel(CCh.ID, Gold)
					ChangeString = append(ChangeString, "Channel **"+Cdel.Name+"** has been deleted.")
					Log("S " + Gold.ID + ": DELETECHANNEL: " + Cdel.ID)
				}
			}
		}
	}
	if TotC.Guild.Roles {
		for _, Rch := range TotC.Roles {
			if !Rch.ExistCrisis {
				Ro := GetRole(Rch.ID, Gold)
				Rn := GetRole(Rch.ID, Gnew)
				if Rch.Name {
					ChangeString = append(ChangeString, "Role **"+Ro.Name+"** has changed it's name to **"+Rn.Name+"**")
					Log("R " + Ro.ID + ": NEWNAME: " + Ro.Name + "->" + Rn.Name)
				}
				if Rch.Color {
					ChangeString = append(ChangeString, "Role **"+Rn.Name+"** has changed it's color from `"+strconv.FormatInt(int64(Ro.Color), 16)+"` to `"+strconv.FormatInt(int64(Rn.Color), 16)+"`")
					Log("R " + Ro.ID + ": NEWCOLOR: " + strconv.FormatInt(int64(Ro.Color), 16) + "->" + strconv.FormatInt(int64(Rn.Color), 16))
				}
				if Rch.Perms {
					ChangeString = append(ChangeString, "Role **"+Rn.Name+"** has changed it's permissions from `"+strconv.FormatInt(int64(Ro.Permissions), 10)+"` to `"+strconv.FormatInt(int64(Rn.Permissions), 10)+"`")
					Log("R " + Ro.ID + ": NEWPERM: " + intToString(Ro.Permissions) + "->" + intToString(Rn.Permissions))
				}
				if Rch.Position {
					ChangeString = append(ChangeString, "Role **"+Rn.Name+"'s** position has been changed, previously: "+strconv.FormatInt(int64(Ro.Position), 10)+", now: "+strconv.FormatInt(int64(Rn.Position), 10))
					Log("R " + Ro.ID + ": NEWPOS: " + intToString(Ro.Position) + "->" + intToString(Rn.Position))
				}
			} else {
				if Rch.Mk {
					Nr := GetRole(Rch.ID, Gnew)
					ChangeString = append(ChangeString, "New role added! ID: `"+Nr.ID+"`")
					Log("S " + Gold.ID + ": NEWROLE: " + Nr.ID)
				}
				if Rch.Del {
					Dr := GetRole(Rch.ID, Gold)
					ChangeString = append(ChangeString, "Role **"+Dr.Name+"** has been deleted")
					Log("S " + Gold.ID + ": DELETEROLE: " + Dr.ID)
				}
			}
		}
	}
	if TotC.Guild.Members {
		for _, Mch := range TotC.Members {
			if !Mch.ExistCrisis {
				Mold := GetMember(Mch.User.ID, Gold)
				Mnew := GetMember(Mch.User.ID, Gnew)
				if Mch.Nick {
					if Mold.Nick == "" {
						ChangeString = append(ChangeString, "Member **"+Mold.User.Username+"** got a nickname! **\""+Mnew.Nick+"\"**!")
					} else if Mnew.Nick == "" {
						ChangeString = append(ChangeString, "Member **"+Mold.User.Username+"** removed his/her nickname, previously: **\""+Mold.Nick+"\"**")
					} else {
						ChangeString = append(ChangeString, "Member **"+Mold.User.Username+"** changed his/her nickname from **\""+Mold.Nick+"\"** to **\""+Mnew.Nick+"\"**")
					}
					Log("M " + Mold.User.ID + ": NICK: " + Mold.Nick + "->" + Mnew.Nick)
				}
				if Mch.User.Username {
					ChangeString = append(ChangeString, "Member **"+Mold.User.Username+"** changed his/her username to **"+Mnew.User.Username+"**")
					Log("M " + Mold.User.ID + ": USERNAME: " + Mold.User.Username + "->" + Mnew.User.Username)
				}
				if Mch.User.Avatar {
					if Mold.User.Avatar == "" {
						ChangeString = append(ChangeString, "Member **"+Mold.User.Username+"** got an avatar!\n"+discordgo.EndpointUserAvatar(Mnew.User.ID, Mnew.User.Avatar))
					} else if Mnew.User.Avatar == "" {
						ChangeString = append(ChangeString, "Member **"+Mold.User.Username+"** removed his/her avatar.\nPrevious:"+discordgo.EndpointUserAvatar(Mold.User.ID, Mold.User.Avatar))
					} else {
						ChangeString = append(ChangeString, "Member **"+Mold.User.Username+"** changed his avatar from "+discordgo.EndpointUserAvatar(Mold.User.ID, Mold.User.Avatar)+"\nTo "+discordgo.EndpointUserAvatar(Mnew.User.ID, Mnew.User.Avatar))
					}
					Log("M " + Mold.User.ID + ": AVATAR: " + Mold.User.Avatar + "->" + Mnew.User.Avatar)
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
									ChangeString = append(ChangeString, "Member **"+Mnew.User.Username+"** got the **"+NR.Name+"** role!")
									Log("M " + Mold.User.ID + ": NEWROLE: " + NR.ID)
								}
							} else {
								NR := GetRole(Rn, Gnew)
								ChangeString = append(ChangeString, "Member **"+Mnew.User.Username+"** got the **"+NR.Name+"** role!")
								Log("M " + Mold.User.ID + ": NEWROLE: " + NR.ID)
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
									ChangeString = append(ChangeString, "Member **"+Mnew.User.Username+"** doesnt have the **"+OR.Name+"** role anymore!")
									Log("M " + Mold.User.ID + ": DELETEROLE: " + OR.ID)
								}
							} else {
								OR := GetRole(Ro, Gold)
								ChangeString = append(ChangeString, "Member **"+Mnew.User.Username+"** doesnt have the **"+OR.Name+"** role anymore!")
								Log("M " + Mold.User.ID + ": DELETEROLE: " + OR.ID)
							}
						}
					}
				}
			} else {
				if Mch.Join {
					NM := GetUser(Mch.User.ID, Gnew)
					ChangeString = append(ChangeString, "Member **"+NM.Username+"** has joined **"+Gnew.Name+"**!")
					Log("S " + Gold.ID + ": JOIN: " + NM.ID)
				}
				if Mch.Leave {
					OM := GetUser(Mch.User.ID, Gold)
					ChangeString = append(ChangeString, "Member **"+OM.Username+"** has left **"+Gnew.Name+"**!")
					Log("S " + Gold.ID + ": LEAVE: " + OM.ID)
				}
			}
		}
	}
	return ChangeString
}

var currentlyrunningloop chan bool
var running bool

func init() {
	currentlyrunningloop = make(chan bool, 1)
}

func StartCheckLoop() {
	if running {
		currentlyrunningloop <- true
	}
	GLDB, err := ioutil.ReadFile("PrimeGuild")
	if err != nil {
		fmt.Println("Error reading PG file: " + err.Error())
		if GLDB == nil {
			ioutil.WriteFile("PrimeGuild", nil, 9000)
		}
		return
	}
	PrimeGuild = string(GLDB)
	LastCheck := ResumeCheck(PrimeGuild)
	running = true
	go CheckLoop(PrimeGuild, LastCheck, currentlyrunningloop)
}

func CheckLoop(Gid string, LastCheck *LastChangeStatus, closehandle chan bool) {
	for {
		time.Sleep(sh.LoopCooldown)
		select {
		case <-closehandle:
			return
		default:
		}
		GI, err := GetGLDfile(LastCheck.GI.g.ID)
		if err != nil {
			fmt.Println("Error getting GLD file: " + err.Error())
			return
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

var isloaded bool

func BBReady(s *discordgo.Session, r *discordgo.Ready) {
	sh.StopLoop = true
	sh.OwnID = r.User.ID
	sh.OwnAV = r.User.Avatar
	sh.OwnName = r.User.Username
	go sh.cl.work()
	if isloaded {
		fmt.Println("Reconnected!")
		Log("RECONNECT")
	} else {
		Log("CONNECT")
		fmt.Println("Discord: Ready message received\nSH: I am '" + sh.OwnName + "'!\nSH: My User ID: " + sh.OwnID)
		isloaded = true
	}
	StartCheckLoop()
}

func BBCreateMessage(Ses *discordgo.Session, MesC *discordgo.MessageCreate) {
	Mes := MesC.Message
	if Mes.Content != "" {
		if Mes.Content[0] == '!' {
			ProcessCMD(Mes.Content[1:], MesC.Message, sh.Notifiers)
		}
	}
	sh.cl.Mess <- Mes
}

func MESEDIT(Ses *discordgo.Session, MesE *discordgo.MessageUpdate) {
	sh.cl.Edits <- MesE.Message
}

func MESDEL(Ses *discordgo.Session, MesD *discordgo.MessageDelete) {
	sh.cl.Deletes <- MesD.ID
}

// misc funx

func TL(s string) string {
	return strings.ToLower(s)
}

func ProcessCMD(CMD string, M *discordgo.Message, Notifiers []string) {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println("Panicked at ProcessCMD")
		}
	}()
	Commands := getCMD(CMD)
	var SecArg string = ""
	if len(Commands) > 1 {
		SecArg = Commands[1]
	}
	var thirdArg string = ""
	if len(Commands) > 2 {
		thirdArg = Commands[2]
	}
	switch TL(Commands[0]) {
	case "primeguild":
		{
			if SecArg == "" {
				SendMessage(M.ChannelID, "`You gave me a nil ID!`", Notifiers)
			} else {
				data := []byte(SecArg)
				ioutil.WriteFile("PrimeGuild", data, 9000)
				fmt.Println("Set new Prime Guild to '" + SecArg + "'")
				SendMessage(M.ChannelID, "`Set new prime guild to "+SecArg+"`", sh.Notifiers)
			}
		}
	case "stopcheck":
		{
			if !sh.StopLoop {
				sh.StopLoop = true
				fmt.Println("Checking loop stopped")
				SendMessage(M.ChannelID, "`Stopped checking loop`", sh.Notifiers)
			}
			if sh.StopLoop {
				SendMessage(M.ChannelID, "`Checking loop isn't running!`", sh.Notifiers)
			}
		}
	case "kickstart":
		{
			StartCheckLoop()
			SendMessage(M.ChannelID, "`Checking loop started`", sh.Notifiers)
		}
	case "getuser":
		{
			//0: no err, 1: nil value, 2: no match (ID), 3: no match (username), 4: no match (username+discriminator), 5: err reading primeguild, 6: multiple users (no disc), 7: multiple users (disc), 8: Unknown error
			U, errT, IDs := GetMemRaw(Commands[1:])
			var e = false
			if len(Commands[len(Commands)-1]) > 1 {
				if string(Commands[len(Commands)-1][0]) == "-" && strings.ToLower(string(Commands[len(Commands)-1][1])) == "e" {
					e = true
				}
			}
			var TotS string
			switch errT {
			case 0:
				SendMessage(M.ChannelID, GMstring(U.User.ID), sh.Notifiers)
			case 1:
				SendMessage(M.ChannelID, "`You gave me no user to check! Use it like this: !getuser <ID>/<Name>(#<discriminator>) (-e)`\n`Use -e to make lower and uppercase count in the search`", sh.Notifiers)
			case 2:
				SendMessage(M.ChannelID, "`I couldnt find '"+Commands[1]+"'! Note; this user might've left, which means i can't look him up anymore.`\n`You can also have mispelled it, try <Name>(#<discriminator>) instead of an ID, or double-check what you (probably) pasted.`", sh.Notifiers)
			case 3:
				if e {
					TotS = strings.Join(Commands[1:len(Commands)-1], " ")
				} else {
					TotS = strings.Join(Commands[1:], " ")
				}
				SendMessage(M.ChannelID, "`I couldnt find '"+TotS+"'! Note; this user might've left, which means i can't look him up anymore.`\n`You can also have mispelled it, try <Name>(#<discriminator>) instead of just a name or nickname, or double-check what you typed.`", sh.Notifiers)
			case 4:
				if e {
					TotS = strings.Join(Commands[1:len(Commands)-2], " ")
				} else {
					TotS = strings.Join(Commands[1:], " ")
				}
				SendMessage(M.ChannelID, "`I couldnt find '"+TotS+"'! Note; this user might've left, which means i can't look him up anymore.`\n`You can also have mispelled it, try <Name>(#<discriminator>) instead of just a name or nickname, or double-check what you typed.`", sh.Notifiers)
			case 5:
				SendMessage(M.ChannelID, "`Error 5, this is outside your hands, contact the owner of this bot immidiatly.`", sh.Notifiers)
			case 6:
				if len(IDs) > 10 {
					SendMessage(M.ChannelID, "`More than 10 users have this user/nickname... `~~`(Are they having a party over there?)`~~\n`Try to define your search with a discriminator (<Name>#<discriminator>), or count upper and lower case by putting a \"-e\" behind your input.`", sh.Notifiers)
				}
				var NamesA []string
				G, err := GetGuild(GetPG())
				if err != nil {
					log.Fatal("Guild load: " + err.Error())
					return
				}
				for _, ID := range IDs {
					NamesA = append(NamesA, GetUserName(ID, G))
				}
				Names := strings.Join(NamesA[:len(NamesA)-1], ", ")
				Names = Names + " and " + NamesA[len(NamesA)-1]
				SendMessage(M.ChannelID, "`I found more than one user; "+Names+".`\n`Try to define your search with a discriminator (<Name>#<discriminator>), or count upper and lower case by putting a \"-e\" behind your input.`", sh.Notifiers)
			case 7:
				SendMessage(M.ChannelID, "`Error 7, i found more than one user that have this username AND discriminator, use \"-e\" behind your input, please.`", sh.Notifiers)
			case 8:
				SendMessage(M.ChannelID, "`Error 8, please contact the bot's owner.`", sh.Notifiers)
			}
		}
	case "getchannel":
		{
			SendMessage(M.ChannelID, GCstring(SecArg), sh.Notifiers)
		}
	case "getguild":
		{
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
	case "getrole":
		{
			if SecArg == "" {
				SendMessage(M.ChannelID, "`You gave me a nil role!`", sh.Notifiers)
				return
			}
			data, users := GRstring(SecArg)
			SendMessage(M.ChannelID, data, sh.Notifiers)
			if strings.ToLower(thirdArg) == "users" || strings.ToLower(thirdArg) == "user" {

				if len(users) > 1999 {
					SendMessage(M.ChannelID, "`(Cannot send users, too many have this role!)`", sh.Notifiers)
					return
				}
				SendMessage(M.ChannelID, users, sh.Notifiers)
			}
		}
	case "getpermissions":
		{
			var perm string
			if SecArg != "" {
				perm = SecArg
			} else {
				SendMessage(M.ChannelID, "`You gave me a nil permission number!`", sh.Notifiers)
				return
			}
			PerMap := ParsePermissions(GetB10(perm))
			Perms := GetPermissions(PerMap)
			PerS := strings.Join(Perms, ", ")
			SendMessage(M.ChannelID, "`Permissions for: "+SecArg+"`\n`Binary: "+GetBit(GetB10(perm))+"`\n`Perms: "+PerS+"`", sh.Notifiers)
		}
	case "stop":
		{
			SendMessage(M.ChannelID, "`Stopping bot...`", sh.Notifiers)
			sh.Stop = true
		}
	case "restart":
		{
			SendMessage(M.ChannelID, "`Restarting bot...`", sh.Notifiers)
			restart = true
			sh.Stop = true
		}
	case "upgrade":
		{
			SendMessage(M.ChannelID, "`Upgrading bot...`", sh.Notifiers)
			upgrade = true
			sh.Stop = true
		}
	case "getlog":
		{
			for _, n := range sh.Notifiers {
				if M.ChannelID == n {
					file, err := os.Open("log.txt")
					if err != nil {
						fmt.Println(err)
					}
					_, err = sh.dg.ChannelFileSend(M.ChannelID, "log.txt", file)
					if err != nil {
						fmt.Println(err)
					}
					return
				}
			}
		}
	case "broadcast":
		{
			var up bool
			var ver bool = true
			if TL(SecArg) == "up" {
				up = true
			} else if TL(SecArg) == "down" {
				up = false
			} else {
				SendMessage(M.ChannelID, "`You gave me no instruction!`", Notifiers)
				return
			}
			if TL(thirdArg) == "false" {
				ver = false
			}
			swap(up, ver)
			if up && !ver {
				SendMessage(M.ChannelID, "`Broadcasting sherlock identifier...`", Notifiers)
			} else if !up {
				SendMessage(M.ChannelID, "`Stopped broadcasting...`", Notifiers)
			} else if ver && up {
				SendMessage(M.ChannelID, "`Broadcasting sherlock identifier with version...`", Notifiers)
			}
		}
	case "sweep":
		SendMessage(M.ChannelID, "`Sweeping...`", Notifiers)
		shs := sweep()
		if len(shs) != 0 {
			if len(shs) == 1 {
				SendMessage(M.ChannelID, "`I found 1 Sherlock:`", Notifiers)
			} else {
				SendMessage(M.ChannelID, "`I found "+intToString(len(shs))+" Sherlock:`", Notifiers)
			}
			for id := range shs {
				SendMessage(M.ChannelID, "`"+id+": "+GetUS(id).Username+"`", Notifiers)
			}
		} else {
			SendMessage(M.ChannelID, "`I found no Sherlocks...`", Notifiers)
		}
	case "identify":
		U, errT, _ := GetMemRaw(Commands[1:])
		if errT != 0 {
			SendMessage(M.ChannelID, "`Error parsing ID or name, try again, or use !getuser with these same parameters to see what's wrong.`", Notifiers)
			return
		}
		ok, _ := identify(U)
		if ok {
			SendMessage(M.ChannelID, "`"+U.User.Username+" is a Sherlock.`", Notifiers)
		} else {
			SendMessage(M.ChannelID, "`"+U.User.Username+" is not a (visible) Sherlock.`", Notifiers)
		}
	}
}

func DeepEqual(a *discordgo.Guild, b *discordgo.Guild) (bool, *FullChangeStruct) {
	var Equal = true
	var TotC = &FullChangeStruct{}
	Equal, TotC, _ = CompareGuild(a, b, TotC, Equal) // replace _ -> Men
	//var Comp = &CompiledChange{
	//	Old: a,
	//	New: b,
	//}
	//PSchanges, PSyes := HandleMention(Men, Comp) // mentioning chain
	return Equal, TotC
}

// prime suspect dealing

func HandleMention(Men *FullMention, biChange *CompiledChange) (*PrimeSuspectChange, bool) {
	PS := sh.PrimeSuspects
	//var PSC = &PrimeSuspectChange{}
	//Old := biChange.Old
	//New := biChange.New
	for _, Ps := range PS {
		for _, C := range Men.ChannelOR {
			if C.perms {
				for _, Or := range C.Perms {
					if Or.ID == Ps {

					}
				}
			}
		}
	}
	var PSS = &PrimeSuspectChange{}
	return PSS, false
}

// init

var l *os.File

var logf *log.Logger

func Initialize(Token string) (bool, bool) {
	CheckerCount = 0
	isdebug, err := ioutil.ReadFile("debugtoggle")
	restart = false
	upgrade = false
	sh = &Sherlock{
		version:      versions.Version{0, 1, 0, 0},
		Debug:        (err == nil && len(isdebug) > 0),
		Stop:         false,
		StopLoop:     false,
		LoopCooldown: 150 * time.Second, // 150 seconds, normally
		cl:           new(ChatLog),
	}
	sh.cl.init()
	sh.dg, err = discordgo.New(Token)
	if err != nil {
		fmt.Println("Discord Session error, check token, error message: " + err.Error())
		return false, false
	}
	// handlers
	sh.dg.AddHandler(BBReady)
	sh.dg.AddHandler(BBCreateMessage)
	sh.dg.AddHandler(MESEDIT)
	sh.dg.AddHandler(MESDEL)

	fmt.Println("SH: Handlers installed")

	SUS, err := GetSuspects()

	installPerms()

	if err == nil {
		sh.PrimeSuspects = SUS
	} else {
		var sus []string
		sh.PrimeSuspects = sus
	}
	notifiers, err = GetNotifiers()
	if err != nil {
		fmt.Println("Error getting Notifier file: " + err.Error())
	} else {
		sh.Notifiers = notifiers
	}

	l, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logf = log.New(l, "SH: ", log.LstdFlags)
	defer finish()
	Log("BOOT")
	err = sh.dg.Open()
	if err == nil {
		fmt.Println("Discord: Connection established")
		for !sh.Stop {
			time.Sleep(400 * time.Millisecond)
		}
	} else {
		fmt.Println("Error opening websocket connection: ", err.Error())
	}
	sh.StopLoop = true
	fmt.Println("SH: Sherlock stopping...")
	sh.dg.Close()
	sh.cl.Save()
	return restart, upgrade
}

func finish() {
	if restart {
		Log("RESTART")
	} else if upgrade {
		Log("UPGRADE")
	} else if !restart && !upgrade {
		Log("SHUTDOWN")
		l.Close()
	}
}
