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

func Log(i string) {
	logf.Println(i)
	fmt.Println(i)
}

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
	return GetMember(ID, Guild).User
}

func GetUS(ID string) *discordgo.User {
	for _, g := range sh.dg.State.Guilds {
		for _, m := range g.Members {
			if ID == m.User.ID {
				return m.User
			}
		}
	}
	return &discordgo.User{}
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

func GetMemberN(ID string, Guild *discordgo.Guild) (*discordgo.Member, bool) {
	for _, M := range Guild.Members {
		if M.User.ID == ID {
			return M, false
		}
	}
	var MM *discordgo.Member
	return MM, true
}

func GetMemID(name string, disc string, guild *discordgo.Guild, e bool) ([]string, bool) {
	var all []string
	var yus = false
	if disc != "" && !e { //given username+disc, not exact
		for _, m := range guild.Members {
			if strings.ToLower(name) == strings.ToLower(m.User.Username) && disc == m.User.Discriminator {
				all = append(all, m.User.ID)
				yus = true
			}
		}
	} else if disc != "" && e { //given username+disc, match exactly
		for _, m := range guild.Members {
			if name == m.User.Username && disc == m.User.Discriminator {
				all = append(all, m.User.ID)
				yus = true
			}
		}
	} else if disc == "" && !e { //given nickname/username, no disc, not exact
		for _, m := range guild.Members {
			if (strings.ToLower(name) == strings.ToLower(m.User.Username)) || (strings.ToLower(name) == strings.ToLower(m.Nick)) || (name == m.User.ID) {
				all = append(all, m.User.ID)
				yus = true
			}
		}
	} else if disc == "" && e { //given username, nickname, no disc, match exactly (y tho?)
		for _, m := range guild.Members {
			if (name == m.User.Username) || (name == m.Nick) || (name == m.User.ID) {
				all = append(all, m.User.ID)
				yus = true
			}
		}
	}
	return all, yus
}

func GetMemRaw(Raw []string) (*discordgo.Member, int, []string) { //0: no err, 1: nil value, 2: no match (ID), 3: no match (username), 4: no match (username+discriminator), 5: err reading primeguild, 6: multiple users (no disc), 7: multiple users (disc), 8: Unknown error
	PG, err := ioutil.ReadFile("PrimeGuild")
	var MU []string
	var NoMem = &discordgo.Member{}
	if err != nil {
		fmt.Println("Error reading PG file: " + err.Error())
		return NoMem, 5, MU
	}
	RawS, disc, e := ParseRaw(Raw)
	if RawS == "" {
		return NoMem, 1, MU
	}
	Guild, err := GetGuild(string(PG))
	if err != nil {
		logf.Fatal(err)
		return NoMem, 5, MU
	}
	if isnumberstring(RawS) {
		mem, yus := GetMemID(RawS, "", Guild, e)
		if !yus {
			return NoMem, 2, MU
		}
		return GetMember(mem[0], Guild), 0, MU
	}
	IDs, yus := GetMemID(RawS, disc, Guild, e)
	if yus {
		if len(IDs) > 1 {
			if disc == "" {
				return NoMem, 6, IDs
			} else {
				return NoMem, 7, IDs
			}
		}
		ID := IDs[0]
		return GetMember(ID, Guild), 0, MU
	}
	if disc == "" {
		return NoMem, 3, MU
	} else if disc != "" {
		return NoMem, 4, MU
	}
	return NoMem, 8, MU
}

func ParseRaw(Raw []string) (string, string, bool) {
	if len(Raw) == 0 {
		return "", "", false
	}
	if isnumberstring(Raw[0]) {
		if len(Raw) > 1 {
			if strings.ToLower(string(Raw[1][1])) == "e" && string(Raw[1][0]) == "-" {
				return Raw[0], "", true
			}
		} else {
			return Raw[0], "", false
		}
	}
	var exact = false
	var RawS string
	if len(Raw) > 1 {
		if string(Raw[len(Raw)-1][0]) == "-" && strings.ToLower(string(Raw[len(Raw)-1][1])) == "e" {
			exact = true
			Raw = Raw[:len(Raw)-1]
		}
		if len(Raw) > 1 {
			RawS = strings.Join(Raw, " ")
		} else {
			RawS = Raw[0]
		}
	} else {
		RawS = Raw[0]
	}
	if strings.ContainsAny(RawS, "#") {
		FullUN := strings.Split(RawS, "#")
		var UN string
		var disc string
		if len(FullUN) > 2 {
			UN = strings.Join(FullUN[:len(FullUN)-1], " ")
		} else {
			UN = FullUN[0]
		}
		disc = FullUN[len(FullUN)-1]
		return UN, disc, exact
	}
	return RawS, "", exact
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

func GetPG() string {
	PG, err := ioutil.ReadFile("PrimeGuild")
	if err != nil {
		fmt.Println("Error reading PG file: " + err.Error())
		return ""
	}
	return string(PG)
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

var memchan chan []*discordgo.Member

func GetGuildState(Gid string) (*discordgo.Guild, error) {
	state, err := sh.dg.State.Guild(Gid)
	memchan = make(chan []*discordgo.Member)
	sh.dg.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.GuildMembersChunk) { memchan <- m.Members })
	err2 := sh.dg.RequestGuildMembers(Gid, "", 0)
	if err2 != nil {
		return state, err2
	}
	data := <-memchan
	state.Members = data
	return state, err
}

func GetGuild(Gid string) (*discordgo.Guild, error) {
	guild, err := sh.dg.State.Guild(Gid)
	memchan = make(chan []*discordgo.Member)
	sh.dg.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.GuildMembersChunk) { memchan <- m.Members })
	err2 := sh.dg.RequestGuildMembers(Gid, "", 0)
	if err2 != nil {
		return guild, err2
	}
	data := <-memchan
	guild.Members = data
	return guild, err
}

func getCMD(CMD string) []string {
	CMDS := strings.Split(CMD, " ")
	return CMDS
}

func GMstring(U string) string {
	PG, err := ioutil.ReadFile("PrimeGuild")
	if err != nil {
		fmt.Println("Error reading PG file: " + err.Error())
		return ""
	}
	GG, err := GetGuild(string(PG))
	if err != nil {
		panic(PG)
	}
	if U == "" {
		return ""
	}
	User := GetMember(U, GG)
	var Roles []string
	for _, R := range User.Roles {
		var Role = GetRole(R, GG)
		Roles = append(Roles, Role.Name)
	}
	return "`User:`\n`ID: " + User.User.ID + "`\n`Username: " + User.User.Username + "#" + User.User.Discriminator + "`\n`Nickname: " + User.Nick + "`\n`Bot: " + strconv.FormatBool(User.User.Bot) + "`\n`Roles: " + strings.Join(Roles, ", ") + "`\n`Avatar: `" + discordgo.EndpointUserAvatar(User.User.ID, User.User.Avatar)
}

func GCstring(Channel string) string {
	PG, err := ioutil.ReadFile("PrimeGuild")
	if err != nil {
		fmt.Println("Error reading PG file: " + err.Error())
		return ""
	}
	GG, err := GetGuild(string(PG))
	if err != nil {
		panic(PG)
	}
	if !isnumberstring(Channel) {
		for _, ch := range GG.Channels {
			if TL(Channel) == TL(ch.Name) {
				Channel = ch.ID
				break
			}
		}
		if !isnumberstring(Channel) {
			return ""
		}
	}
	CH := GetChannel(Channel, GG)
	if strings.ToLower(CH.Type) == "voice" {
		return "`Channel:`\n`ID: " + CH.ID + "`\n`Name: " + CH.Name + "`\n`Type: " + CH.Type + "`"
	} else {
		return "`Channel:`\n`ID: " + CH.ID + "`\n`Name: " + CH.Name + "`\n`Topic: " + SanCode(CH.Topic) + "`\n`Type: " + CH.Type + "`"
	}
}

func GRstring(role string) (string, string) {
	PG, err := ioutil.ReadFile("PrimeGuild")
	if err != nil {
		fmt.Println("Error reading PG file: " + err.Error())
		return "", ""
	}
	GG, err := GetGuild(string(PG))
	if err != nil {
		panic(PG)
	}
	Role := GetRole(role, GG)
	var returnS string
	var return2 string
	returnS = "`Role:`\n`ID: " + Role.ID + "`\n`Name: " + Role.Name + "`\n`Bot Role: " + strconv.FormatBool(Role.Managed) + "`\n`Mentionable: " + strconv.FormatBool(Role.Mentionable) + "`\n`Special tab in sidebar: " + strconv.FormatBool(Role.Hoist) + "`\n`Color: " + strconv.FormatInt(int64(Role.Color), 16) + "`\n`Position (acending): " + strconv.FormatInt(int64(Role.Position), 10) + "`\n`Permissions: " + strconv.FormatInt(int64(Role.Permissions), 10) + "`"
	var people []string
	for _, P := range GG.Members {
		for _, R := range P.Roles {
			if R == Role.ID {
				people = append(people, P.User.Username)
			}
		}
	}
	return2 = "`Users:`\n`" + strings.Join(people, ", ") + "`"
	return returnS, return2

}

func SendMessage(Channel string, Message string, Priviledged []string) (bool, string, string) {
	var ok bool
	ok = false
	for _, AN := range Priviledged {
		if AN == Channel {
			ok = true
		}
	}
	if !ok {
		fmt.Println("Permission denied for " + Channel)
		return false, "", ""
	}
	m, err := sh.dg.ChannelMessageSend(Channel, Message)
	if err != nil {
		fmt.Println("Err parsing message to "+Channel+" from sending:", err)
		return true, "", ""
	} else {
		return true, m.ID, Channel
	}
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

func GetSuspects() ([]string, error) {
	DATA, err := ioutil.ReadFile("PS")
	if DATA == nil {
		var S []string
		return S, err
	}
	TotalString := string(DATA)
	return strings.Split(TotalString, "+"), nil
}

func SanCode(S string) string {
	return strings.Replace(S, "\n", "`\n`", -1)
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
			ioutil.WriteFile(G.g.ID+"/"+intToString(T.Year)+"."+T.Month.String()+"."+intToString(T.Day)+"-"+intToString(T.Hour)+".GLD", GU, 0777)
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

func GSFI(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func isnumberstring(s string) bool {
	var check = true
	for _, l := range s {
		if !strings.ContainsAny(string(l), "0 1 2 3 4 5 6 7 8 9") {
			check = false
		}
	}
	return check
}

func GetBit(base10 int) string {
	return strconv.FormatInt(int64(base10), 2)
}

func GetB10(Bit string) int {
	B10, _ := strconv.ParseInt(Bit, 10, 0)
	return int(B10)
}

func AltB10(Bit string) int {
	B10, _ := strconv.ParseInt(Bit, 0, 0)
	return int(B10)
}

func CompByteArray(a []byte, b []byte) bool {
	var aS = string(a)
	var bS = string(b)
	return aS == bS
}

func GetTimeString(T *TimeFormat, short bool) string {
	if short {
		h := strconv.FormatInt(int64(T.Hour), 10)
		mi := strconv.FormatInt(int64(T.Min), 10)
		s := strconv.FormatInt(int64(T.Sec), 10)
		return "-" + h + ":" + mi + ";" + s + "-"
	} else {
		y := strconv.FormatInt(int64(T.Year), 10)
		mo := strconv.FormatInt(int64(T.Month), 10)
		d := strconv.FormatInt(int64(T.Day), 10)
		h := strconv.FormatInt(int64(T.Hour), 10)
		mi := strconv.FormatInt(int64(T.Min), 10)
		s := strconv.FormatInt(int64(T.Sec), 10)
		return "|" + y + "/" + mo + "/" + d + "|" + h + ":" + mi + ";" + s
	}
}

func ParsePermissions(Perm int) *Permissions { // perm in base 10
	Perms := &Permissions{}
	Bits := getBit(uint64(Perm))
	for b, yes := range Bits[:32] {
		switch b {
		case 0:
			Perms.CREATE_INSTANT_INVITE = yes
		case 1:
			Perms.KICK_MEMBERS = yes
		case 2:
			Perms.BAN_MEMBERS = yes
		case 3:
			Perms.ADMINISTRATOR = yes
		case 4:
			Perms.MANAGE_CHANNELS = yes
		case 5:
			Perms.MANAGE_GUILD = yes
		case 6:
			Perms.ADD_REACTIONS = yes
		case 10:
			Perms.READ_MESSAGES = yes
		case 11:
			Perms.SEND_MESSAGES = yes
		case 12:
			Perms.SEND_TTS_MESSAGES = yes
		case 13:
			Perms.MANAGE_MESSAGES = yes
		case 14:
			Perms.EMBED_LINKS = yes
		case 15:
			Perms.ATTACH_FILES = yes
		case 16:
			Perms.READ_MESSAGE_HISTORY = yes
		case 17:
			Perms.MENTION_EVERYONE = yes
		case 18:
			Perms.USE_EXTERNAL_EMOJIS = yes
		case 20:
			Perms.CONNECT = yes
		case 21:
			Perms.SPEAK = yes
		case 22:
			Perms.MUTE_MEMBERS = yes
		case 23:
			Perms.DEAFEN_MEMBERS = yes
		case 24:
			Perms.MOVE_MEMBERS = yes
		case 25:
			Perms.USE_VAD = yes
		case 26:
			Perms.CHANGE_NICKNAME = yes
		case 27:
			Perms.MANAGE_NICKNAMES = yes
		case 28:
			Perms.MANAGE_ROLES = yes
		case 29:
			Perms.MANAGE_WEBHOOKS = yes
		case 30:
			Perms.MANAGE_EMOJIS = yes
		}
	}
	return Perms
}

func GetPermissions(Perms *Permissions) []string {
	var P []string
	if Perms.CREATE_INSTANT_INVITE == true {
		P = append(P, "CREATE_INSTANT_INVITE")
	}
	if Perms.KICK_MEMBERS == true {
		P = append(P, "KICK_MEMBERS")
	}
	if Perms.BAN_MEMBERS == true {
		P = append(P, "BAN_MEMBERS")
	}
	if Perms.ADMINISTRATOR == true {
		P = append(P, "ADMINISTRATOR")
	}
	if Perms.MANAGE_CHANNELS == true {
		P = append(P, "MANAGE_CHANNELS")
	}
	if Perms.MANAGE_GUILD == true {
		P = append(P, "MANAGE_GUILD")
	}
	if Perms.ADD_REACTIONS == true {
		P = append(P, "ADD_REACTIONS")
	}
	if Perms.READ_MESSAGES == true {
		P = append(P, "READ_MESSAGES")
	}
	if Perms.SEND_MESSAGES == true {
		P = append(P, "SEND_MESSAGES")
	}
	if Perms.SEND_TTS_MESSAGES == true {
		P = append(P, "SEND_TTS_MESSAGES")
	}
	if Perms.MANAGE_MESSAGES == true {
		P = append(P, "MANAGE_MESSAGES")
	}
	if Perms.EMBED_LINKS == true {
		P = append(P, "EMBED_LINKS")
	}
	if Perms.ATTACH_FILES == true {
		P = append(P, "ATTACH_FILES")
	}
	if Perms.READ_MESSAGE_HISTORY == true {
		P = append(P, "READ_MESSAGE_HISTORY")
	}
	if Perms.MENTION_EVERYONE == true {
		P = append(P, "MENTION_EVERYONE")
	}
	if Perms.USE_EXTERNAL_EMOJIS == true {
		P = append(P, "USE_EXTERNAL_EMOJIS")
	}
	if Perms.CONNECT == true {
		P = append(P, "CONNECT")
	}
	if Perms.SPEAK == true {
		P = append(P, "SPEAK")
	}
	if Perms.MUTE_MEMBERS == true {
		P = append(P, "MUTE_MEMBERS")
	}
	if Perms.DEAFEN_MEMBERS == true {
		P = append(P, "DEAFEN_MEMBERS")
	}
	if Perms.MOVE_MEMBERS == true {
		P = append(P, "MOVE_MEMBERS")
	}
	if Perms.USE_VAD == true {
		P = append(P, "USE_VAD")
	}
	if Perms.CHANGE_NICKNAME == true {
		P = append(P, "CHANGE_NICKNAME")
	}
	if Perms.MANAGE_NICKNAMES == true {
		P = append(P, "MANAGE_NICKNAMES")
	}
	if Perms.MANAGE_ROLES == true {
		P = append(P, "MANAGE_ROLES")
	}
	if Perms.MANAGE_WEBHOOKS == true {
		P = append(P, "MANAGE_WEBHOOKS")
	}
	if Perms.MANAGE_EMOJIS == true {
		P = append(P, "MANAGE_EMOJIS")
	}
	return P
}

/*
CREATE_INSTANT_INVITE 0,  //0x00000001
KICK_MEMBERS          1,  //0x00000002
BAN_MEMBERS           2,  //0x00000004
ADMINISTRATOR         3,  //0x00000008
MANAGE_CHANNELS       4,  //0x00000010
MANAGE_GUILD          5,  //0x00000020
ADD_REACTIONS         6,  //0x00000040
READ_MESSAGES         10, //0x00000400
SEND_MESSAGES         11, //0x00000800
SEND_TTS_MESSAGES     12, //0x00001000
MANAGE_MESSAGES       13, //0x00002000
EMBED_LINKS           14, //0x00004000
ATTACH_FILES          15, //0x00008000
READ_MESSAGE_HISTORY  16, //0x00010000
MENTION_EVERYONE      17, //0x00020000
USE_EXTERNAL_EMOJIS   18, //0x00040000
CONNECT               20, //0x00100000
SPEAK                 21, //0x00200000
MUTE_MEMBERS          22, //0x00400000
DEAFEN_MEMBERS        23, //0x00800000
MOVE_MEMBERS          24, //0x01000000
USE_VAD               25, //0x02000000
CHANGE_NICKNAME       26, //0x04000000
MANAGE_NICKNAMES      27, //0x08000000
MANAGE_ROLES          28, //0x10000000
MANAGE_WEBHOOKS       29, //0x20000000
MANAGE_EMOJIS         30, //0x40000000

		Permissions in decending order:

		MANAGE_EMOJIS
		MANAGE_WEBHOOKS
		MANAGE_ROLES
		MANAGE_NICKNAMES
		CHANGE_NICKNAME
		USE_VAD
		MOVE_MEMBERS
		DEAFEN_MEMBERS
		MUTE_MEMBERS
		SPEAK
		CONNECT
		USE_EXTERNAL_EMOJIS
		MENTION_EVERYONE
		READ_MESSAGE_HISTORY
		ATTACH_FILES
		EMBED_LINKS
		MANAGE_MESSAGES
		SEND_TTS_MESSAGES
		SEND_MESSAGES
		READ_MESSAGES
		ADD_REACTIONS
		MANAGE_GUILD
		MANAGE_CHANNELS
		ADMINISTRATOR
		BAN_MEMBERS
		KICK_MEMBERS
		CREATE_INSTANT_INVITE

*/

func getBit(data uint64) []bool {
	var dataBitmap = make([]bool, 64)
	var index uint64 = 0
	for index < 64 {
		dataBitmap[index] = data&(1<<index) > 0
		index++
	}
	return dataBitmap
}

/*
var dataBitmap = make([]bool, 64)
    var index uint64 = 0
    for index < 64 {
        dataBitmap[index] = data & (1 << index) > 0
        index++
    }
    return dataBitmap

    ------

    var dataBitmap = make([]bool, 64)
	Bit := GetBit(data)
	var index int = 0
	for index < 64 {
		if index < len(Bit) {
			switch Bit[len(Bit)-index-1] {
			case '0':
				dataBitmap[index] = false
			case '1':
				dataBitmap[index] = true
			}
		} else {
			dataBitmap[index] = false
		}
		index++
	}
	return dataBitmap
*/
