package Belt

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var empty string
var emptyB []byte

// overwrite

func CompareORRoles(a *discordgo.Channel, b *discordgo.Channel, ChChange *ChannelChange, Equal bool, Guildchannels bool) (bool, *ChannelChange, bool) {
	var Chanbool bool = Guildchannels
	if len(a.PermissionOverwrites) == 0 {
		return Equal, ChChange, Guildchannels
	}
	for _, Pa := range a.PermissionOverwrites {
		for _, Pb := range b.PermissionOverwrites {
			if Pa.ID == Pb.ID {
				if Pa.Allow != Pb.Allow || Pa.Deny != Pb.Deny {
					Chanbool = true
					Equal = false
					var ChPerm = &PermORChange{}
					ChPerm.ID = Pa.ID
					if Pa.Allow != Pb.Allow {
						ChPerm.Allow = true
					}
					if Pa.Deny != Pb.Deny {
						ChPerm.Deny = true
					}
					Perms := ChChange.Perms
					Perms = append(Perms, ChPerm)
					ChChange.ID = a.ID
					ChChange.Perms = Perms
					ChChange.perms = true
				}
			}
		}
	}
	return Equal, ChChange, Chanbool
}

func NoteORDelete(a *discordgo.Channel, b *discordgo.Channel, ChChange *ChannelChange) *ChannelChange {
	var IsStillThere = false
	var DelOR = &PermORChange{}
	for _, ORA := range a.PermissionOverwrites {
		IsStillThere = false
		for _, ORB := range b.PermissionOverwrites {
			if ORA.ID == ORB.ID {
				IsStillThere = true
				_, ChChange, _ = CompareORRoles(a, b, ChChange, false, false)
			}
		}
		if !IsStillThere {
			DelOR.ID = ORA.ID
			DelOR.ExistCrisis = true
			DelOR.Del = true
			AllOR := ChChange.Perms
			AllOR = append(AllOR, DelOR)
			ChChange.Perms = AllOR
			ChChange.perms = true
			DelOR = &PermORChange{}
		}
	}
	return ChChange
}

func NoteORCreate(a *discordgo.Channel, b *discordgo.Channel, ChChange *ChannelChange) *ChannelChange {
	var IsOld = false
	var CreateOR = &PermORChange{}
	for _, ORA := range a.PermissionOverwrites {
		IsOld = false
		for _, ORB := range b.PermissionOverwrites {
			if ORA.ID == ORB.ID {
				IsOld = true
				_, ChChange, _ = CompareORRoles(a, b, ChChange, false, false)
			}
		}
		if !IsOld {
			CreateOR.ID = ORA.ID
			CreateOR.ExistCrisis = true
			CreateOR.Mk = true
			AllOR := ChChange.Perms
			AllOR = append(AllOR, CreateOR)
			ChChange.Perms = AllOR
			ChChange.perms = true
			CreateOR = &PermORChange{}
		}
	}
	return ChChange
}

// channel

func CompareChannelstruct(a *discordgo.Channel, b *discordgo.Channel, ToTC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	var ChChange = &ChannelChange{}
	TotC := ToTC
	if a.Name != b.Name {
		Equal = false
		TotC.Guild.channels = true
		ChChange.ID = a.ID
		ChChange.Name = true
	}
	if a.Topic != b.Topic {
		Equal = false
		TotC.Guild.channels = true
		ChChange.ID = a.ID
		ChChange.Topic = true
	}
	if len(a.PermissionOverwrites) == len(b.PermissionOverwrites) {
		Equal, ChChange, TotC.Guild.channels = CompareORRoles(a, b, ChChange, Equal, TotC.Guild.channels)
	}
	if len(a.PermissionOverwrites) < len(b.PermissionOverwrites) {
		Equal = false
		TotC.Guild.channels = true
		ChChange.ID = a.ID
		ChChange = NoteORDelete(a, b, ChChange)
	}
	if len(a.PermissionOverwrites) > len(b.PermissionOverwrites) {
		Equal = false
		TotC.Guild.channels = true
		ChChange.ID = a.ID
		ChChange = NoteORCreate(a, b, ChChange)
	}
	if ChChange.ID != "" {
		AllCh := TotC.Channels
		AllCh = append(AllCh, ChChange)
		TotC.Channels = AllCh
	}
	return Equal, TotC
}

func NoteChannelDelete(cha *discordgo.Channel, TotC *FullChangeStruct) *FullChangeStruct {
	var DelCh = &ChannelChange{}
	DelCh.ID = cha.ID
	DelCh.ExistCrisis = true
	DelCh.Del = true
	AllCh := TotC.Channels
	AllCh = append(AllCh, DelCh)
	TotC.Channels = AllCh
	return TotC
}

func NoteChannelCreate(cha *discordgo.Channel, TotC *FullChangeStruct) *FullChangeStruct {
	var NewCh = &ChannelChange{}
	NewCh.ID = cha.ID
	NewCh.ExistCrisis = true
	NewCh.Mk = true
	AllCh := TotC.Channels
	AllCh = append(AllCh, NewCh)
	TotC.Channels = AllCh
	return TotC
}

func CompareChannels(a []*discordgo.Channel, b []*discordgo.Channel, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	var Aisthere bool
	var Bisthere bool
	for _, Ca := range a {
		Aisthere = false
		for _, Cb := range b {
			if Ca.ID == Cb.ID {
				Equal, TotC = CompareChannelstruct(Ca, Cb, TotC, Equal)
				Aisthere = true
			}
		}
		if !Aisthere {
			Equal = false
			TotC.Guild.channels = true
			TotC = NoteChannelDelete(Ca, TotC)
		}
	}
	for _, Cb := range b {
		Bisthere = false
		for _, Ca := range a {
			if Ca.ID == Cb.ID {
				Bisthere = true
			}
		}
		if !Bisthere {
			Equal = false
			TotC.Guild.channels = true
			TotC = NoteChannelCreate(Cb, TotC)
		}
	}
	return Equal, TotC
}

// members

func CompareMemberstruct(a *discordgo.Member, b *discordgo.Member, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	var MemCh = &MemberChange{}
	if a.Nick != b.Nick {
		Equal = false
		TotC.Guild.members = true

		MemCh.User.ID = a.User.ID
		MemCh.Nick = true
	}
	if a.User.Avatar != b.User.Avatar {
		Equal = false
		TotC.Guild.members = true

		MemCh.User.ID = a.User.ID
		MemCh.User.Avatar = true
	}
	if a.User.Username != b.User.Username {
		Equal = false
		TotC.Guild.members = true

		MemCh.User.ID = a.User.ID
		MemCh.User.Username = true
	}
	if len(a.Roles) > len(b.Roles) {
		Equal = false
		TotC.Guild.members = true

		MemCh.User.ID = a.User.ID
		MemCh.Roles = true
		MemCh.RoleRem = true
	}
	if len(a.Roles) < len(b.Roles) {
		Equal = false
		TotC.Guild.members = true

		MemCh.User.ID = a.User.ID
		MemCh.Roles = true
		MemCh.RoleNew = true
	}
	if MemCh.User.ID != "" {
		AllM := TotC.Members
		AllM = append(AllM, MemCh)
		TotC.Members = AllM
	}
	return Equal, TotC
}

func NoteMemberLeave(mem *discordgo.Member, TotC *FullChangeStruct) *FullChangeStruct {
	var LeaveMem = &MemberChange{}
	LeaveMem.User.ID = mem.User.ID
	LeaveMem.ExistCrisis = true
	LeaveMem.Leave = true
	AllM := TotC.Members
	AllM = append(AllM, LeaveMem)
	TotC.Members = AllM
	return TotC
}

func NoteMemberJoin(mem *discordgo.Member, TotC *FullChangeStruct) *FullChangeStruct {
	var NewM = &MemberChange{}
	NewM.User.ID = mem.User.ID
	NewM.ExistCrisis = true
	NewM.Join = true
	AllM := TotC.Members
	AllM = append(AllM, NewM)
	TotC.Members = AllM
	return TotC
}

func CompareMembers(a []*discordgo.Member, b []*discordgo.Member, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	var Aisthere bool
	var Bisthere bool
	for _, Ca := range a {
		Aisthere = false
		for _, Cb := range b {
			if Ca.User.ID == Cb.User.ID {
				Equal, TotC = CompareMemberstruct(Ca, Cb, TotC, Equal)
				Aisthere = true
			}
		}
		if !Aisthere {
			Equal = false
			TotC.Guild.members = true
			TotC = NoteMemberLeave(Ca, TotC)
		}
	}
	for _, Cb := range b {
		Bisthere = false
		for _, Ca := range a {
			if Ca.User.ID == Cb.User.ID {
				Bisthere = true
			}
		}
		if !Bisthere {
			Equal = false
			TotC.Guild.members = true
			TotC = NoteMemberJoin(Cb, TotC)
		}
	}
	return Equal, TotC
}

// roles

func CompareRolestruct(a *discordgo.Role, b *discordgo.Role, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	var RoleCh = &RoleChange{}
	var detect bool
	if a.Name != b.Name {
		TotC.Guild.roles = true
		RoleCh.ID = a.ID
		RoleCh.Name = true
		detect = true
	}
	if a.Permissions != b.Permissions {
		TotC.Guild.roles = true
		RoleCh.ID = a.ID
		RoleCh.Perms = true
		detect = true
	}
	if a.Position != b.Position {
		TotC.Guild.roles = true
		RoleCh.ID = a.ID
		RoleCh.Position = true
		detect = true
	}
	if a.Color != b.Color {
		TotC.Guild.roles = true
		RoleCh.ID = a.ID
		RoleCh.Color = true
		detect = true
	}
	if detect {
		AllR := TotC.Roles
		AllR = append(AllR, RoleCh)
		TotC.Roles = AllR
		Equal = false
	}
	return Equal, TotC
}

func NoteRoleRemove(rol *discordgo.Role, TotC *FullChangeStruct) *FullChangeStruct {
	var DelR = &RoleChange{}
	DelR.ID = rol.ID
	DelR.ExistCrisis = true
	DelR.Del = true
	AllR := TotC.Roles
	AllR = append(AllR, DelR)
	TotC.Roles = AllR
	DelR = &RoleChange{}
	return TotC
}

func NoteRoleCreate(rol *discordgo.Role, TotC *FullChangeStruct) *FullChangeStruct {
	var NewR = &RoleChange{}
	NewR.ID = rol.ID
	NewR.ExistCrisis = true
	NewR.Mk = true
	AllR := TotC.Roles
	AllR = append(AllR, NewR)
	TotC.Roles = AllR
	return TotC
}

func CompareRoles(a []*discordgo.Role, b []*discordgo.Role, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	var Aisthere bool
	var Bisthere bool
	for _, Ca := range a {
		Aisthere = false
		for _, Cb := range b {
			if Ca.ID == Cb.ID {
				Equal, TotC = CompareRolestruct(Ca, Cb, TotC, Equal)
				Aisthere = true
			}
		}
		if !Aisthere {
			Equal = false
			TotC.Guild.roles = true
			TotC = NoteRoleRemove(Ca, TotC)
		}
	}
	for _, Cb := range b {
		Bisthere = false
		for _, Ca := range a {
			if Ca.ID == Cb.ID {
				Bisthere = true
			}
		}
		if !Bisthere {
			Equal = false
			TotC.Guild.roles = true
			TotC = NoteRoleCreate(Cb, TotC)
		}
	}
	return Equal, TotC
}

func CompareGuild(a *discordgo.Guild, b *discordgo.Guild, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct, *FullMention) {
	if a.Name != b.Name {
		Equal = false
		TotC.Guild.Name = true
	}
	if a.OwnerID != b.OwnerID {
		Equal = false
		TotC.Guild.OwnerID = true
	}
	if a.Icon != b.Icon {
		Equal = false
		TotC.Guild.Icon = true
	}
	if a.Region != b.Region {
		Equal = false
		TotC.Guild.Region = true
	}
	Equal, TotC = CompareChannels(a.Channels, b.Channels, TotC, Equal)
	Equal, TotC = CompareMembers(a.Members, b.Members, TotC, Equal)
	Equal, TotC = CompareRoles(a.Roles, b.Roles, TotC, Equal)
	line, _ := json.Marshal(TotC)
	fmt.Println(GetTimeString(GetTimeForm(), false))
	if !Equal {
		fmt.Println("Change data:")
		fmt.Println(string(line))
		fmt.Println(Equal)
	}
	var MenT = &FullMention{
		ChannelOR: TotC.Channels,
		Members:   TotC.Members,
		Roles:     TotC.Roles,
		OwnerID:   b.OwnerID,
	}
	return Equal, TotC, MenT
}
