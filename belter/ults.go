package Belt

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SetLC(LC *LastChangeStatus, g *discordgo.Guild) *LastChangeStatus {
	GI := LC.GI
	GI.g = g
	GI.Lastcheck = AltTimeForm()
	LC.GI = GI
	return LC
}

func GetOR(ID string, Channel *discordgo.Channel) *discordgo.PermissionOverwrite {
	for _, OR := range Channel.PermissionOverwrites {
		if OR.ID == ID {
			return OR
		}
	}
	var OOR *discordgo.PermissionOverwrite
	return OOR
}

func GetUser(ID string, Guild *discordgo.Guild) *discordgo.User {
	for _, M := range Guild.Members {
		if M.User.ID == ID {
			return M.User
		}
	}
	var UU *discordgo.User
	return UU
}

func GetMember(ID string, Guild *discordgo.Guild) *discordgo.Member {
	for _, M := range Guild.Members {
		if M.User.ID == ID {
			return M
		}
	}
	var MM *discordgo.Member
	return MM
}

func GetUserName(ID string, Guild *discordgo.Guild) string {
	User := GetUser(ID, Guild)
	return User.Username
}

func GetChannel(ID string, Guild *discordgo.Guild) *discordgo.Channel {
	for _, C := range Guild.Channels {
		if C.ID == ID {
			return C
		}
	}
	var CC *discordgo.Channel
	return CC
}

func GetRole(RID string, g *discordgo.Guild) *discordgo.Role {
	for _, R := range g.Roles {
		if RID == R.ID {
			return R
		}
	}
	var Defaultrole *discordgo.Role
	return Defaultrole
}

func AppendSSlices(BeginSlice []string, MergeSlice []string) []string {
	var ProcessSlice = BeginSlice
	for _, MergeString := range MergeSlice {
		ProcessSlice = append(ProcessSlice, MergeString)
	}
	return ProcessSlice
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

func GetTimeForm() *TimeFormat {
	Year, Month, Day, Hour, Min, Sec := GetTime()
	var T = &TimeFormat{
		Year:  Year,
		Month: Month,
		Day:   Day,
		Hour:  Hour,
		Min:   Min,
		Sec:   Sec,
	}
	return T
}
func AltTimeForm() TimeFormat {
	var TM TimeFormat
	TM.Year, TM.Month, TM.Day, TM.Hour, TM.Min, TM.Sec = GetTime()
	return TM
}

func GetGuild(Gid string) (*discordgo.Guild, error) {
	return sh.dg.State.Guild(Gid)
}

func GetGLDfile(GID string) (*GuildInfo, error) {
	DATA, err := ioutil.ReadFile(GID + ".GLD")
	var LLG = &GuildInfo{}
	// Y, Mo, D, H, Mi, S := GetTime()
	if err == nil {
		err := json.Unmarshal(DATA, LLG)
		if err == nil {
			LLG.g = &discordgo.Guild{}
			GU, err := ioutil.ReadFile(GID + "/main.GLD")
			if err == nil {
				err = json.Unmarshal(GU, LLG.g)
				if err == nil {
					return LLG, nil
				} else {
					fmt.Println("Error Unmarshal-ing GLD: " + err.Error())
					return LLG, err
				}
			} else {
				fmt.Println("Error Unmarshal-ing GLD: " + err.Error())
				return LLG, err
			}
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

func intToString(I int) string {
	return strconv.FormatInt(int64(I), 10)
}

func GetNotifiers() ([]string, error) {
	DATA, err := ioutil.ReadFile("Notifiers")
	if DATA == nil {
		ioutil.WriteFile("Notifiers", DATA, 0777)
		var S []string
		fmt.Println("No notifier file found, creating sample file...")
		return S, err
	}
	TotalString := string(DATA)
	return strings.Split(TotalString, "+"), nil
}

func WriteGLDfile(G *GuildInfo, Isb bool) error {
	if G.g == nil {
		panic(G.g)
	}
	GID, GIDerr := json.Marshal(G)
	g := G.g
	GU, GUerr := json.Marshal(g)
	if GU == nil {
		panic(GU)
	}
	if GIDerr == nil && GUerr == nil {
		if Isb == false {
			ioutil.WriteFile(G.g.ID+".GLD", GID, 0777)
			err := ioutil.WriteFile(G.g.ID+"/main.GLD", GU, 0777)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			T := GetTimeForm()
			ioutil.WriteFile(G.g.ID+"/"+intToString(T.Year)+"."+T.Month.String()+"."+intToString(T.Day)+"-"+intToString(T.Hour)+"."+intToString(T.Min)+".GLD", GU, 0777)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
		if Isb == true {
			ioutil.WriteFile("B-"+G.g.ID+".GLD", GID, 0777)
			os.Mkdir(G.g.ID, 0777)
			err := ioutil.WriteFile(G.g.ID+"/MASTER.GLD", GU, 0777)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
		return nil
	} else {
		fmt.Println("GLD writing error: " + GIDerr.Error())
		return GIDerr
	}
}
