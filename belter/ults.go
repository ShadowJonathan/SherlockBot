package Belt

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SetLC(LC *LastChangeStatus, g *discordgo.Guild) *LastChangeStatus {
	GI := LC.GI
	GI.g = g
	GI.Lastcheck.Year, GI.Lastcheck.Month, GI.Lastcheck.Day, GI.Lastcheck.Hour, GI.Lastcheck.Min, GI.Lastcheck.Sec = GetTime()
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
