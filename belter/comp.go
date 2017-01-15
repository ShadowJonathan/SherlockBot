package Belt

import (
	"github.com/bwmarrin/discordgo"
)

// channel

func CompareORRoles(a *discordgo.Channel, b *discordgo.Channel, ChChange *ChannelChange, Equal bool) (bool, *ChannelChange) {
	if len(a.PermissionOverwrites) == 0 {
		return Equal, ChChange
	}
	for _, Pa := range a.PermissionOverwrites {
		for _, Pb := range b.PermissionOverwrites {
			if Pa.ID == Pb.ID {
				if Pa.Allow != Pb.Allow || Pa.Deny != Pb.Deny {
					Equal = false
					var ChPerm *PermORChange
					ChPerm.ID = Pa.ID
					if Pa.Allow != Pb.Allow {
						ChPerm.Allow = true
					}
					if Pa.Deny != Pb.Deny {
						ChPerm.Deny = true
					}
					Perms := ChChange.Perms
					Perms = append(Perms, ChPerm)
					ChChange.Perms = Perms
				}
			}
		}
	}
	return Equal, ChChange
}

func NoteORDelete(a *discordgo.Channel, b *discordgo.Channel, ChChange *ChannelChange) *ChannelChange {
	var IsStillThere = false
	var DelOR *PermORChange
	for _, ORA := range a.PermissionOverwrites {
		for _, ORB := range b.PermissionOverwrites {
			if ORA.ID == ORB.ID {
				IsStillThere = true
				_, ChChange = CompareORRoles(a, b, ChChange, false)
			}
		}
		if !IsStillThere {
			DelOR.ID = ORA.ID
			DelOR.ExistCrisis = true
			DelOR.Del = true
			AllOR := ChChange.Perms
			AllOR = append(AllOR, DelOR)
			ChChange.Perms = AllOR
		}
	}
	return ChChange
}

func NoteORCreate(a *discordgo.Channel, b *discordgo.Channel, ChChange *ChannelChange) *ChannelChange {
	var IsOld = false
	var CreateOR *PermORChange
	for _, ORA := range a.PermissionOverwrites {
		for _, ORB := range b.PermissionOverwrites {
			if ORA.ID == ORB.ID {
				IsOld = true
				_, ChChange = CompareORRoles(a, b, ChChange, false)
			}
		}
		if !IsOld {
			CreateOR.ID = ORA.ID
			CreateOR.ExistCrisis = true
			CreateOR.Mk = true
			AllOR := ChChange.Perms
			AllOR = append(AllOR, CreateOR)
			ChChange.Perms = AllOR
		}
	}
	return ChChange
}

func CompareChannelstruct(a *discordgo.Channel, b *discordgo.Channel, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	var ChChange *ChannelChange
	var Changeyes = false
	if a.Name != b.Name {
		Equal = false
		ChChange.Name = true
	}
	if a.Topic != b.Topic {
		Equal = false
		ChChange.Topic = true
	}
	if len(a.PermissionOverwrites) == len(b.PermissionOverwrites) {
		Equal, ChChange = CompareORRoles(a, b, ChChange, Equal)
	}
	if len(a.PermissionOverwrites) < len(b.PermissionOverwrites) {
		Equal = false
		ChChange = NoteORDelete(a, b, ChChange)
	}
	if len(a.PermissionOverwrites) > len(b.PermissionOverwrites) {
		Equal = false
		ChChange = NoteORCreate(a, b, ChChange)
	}
	TotC.Channels = append(TotC.Channels, ChChange)
	return Equal, TotC
}

func NoteChannelDelete(a []*discordgo.Channel, b []*discordgo.Channel, TotC *FullChangeStruct) *FullChangeStruct {
	var IsStillThere = false
	var DelCh *ChannelChange
	for _, ChA := range a {
		for _, ChB := range b {
			if ChA.ID == ChB.ID {
				IsStillThere = true
				_, TotC = CompareChannelstruct(ChA, ChB, TotC, false)
			}
		}
		if !IsStillThere {
			DelCh.ID = ChA.ID
			DelCh.ExistCrisis = true
			DelCh.Del = true
			AllCh := TotC.Channels
			AllCh = append(AllCh, DelCh)
			TotC.Channels = AllCh
		}
	}
	return TotC
}

func NoteChannelCreate(a []*discordgo.Channel, b []*discordgo.Channel, TotC *FullChangeStruct) *FullChangeStruct {
	var IsOld = false
	var NewCh *ChannelChange
	for _, ChA := range a {
		for _, ChB := range b {
			if ChA.ID == ChB.ID {
				IsOld = true
				_, TotC = CompareChannelstruct(ChA, ChB, TotC, false)
			}
		}
		if !IsOld {
			NewCh.ID = ChA.ID
			NewCh.ExistCrisis = true
			NewCh.Mk = true
			AllCh := TotC.Channels
			AllCh = append(AllCh, NewCh)
			TotC.Channels = AllCh
		}
	}
	return TotC
}

func CompareChannels(a []*discordgo.Channel, b []*discordgo.Channel, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
	if len(a) > len(b) {
		Equal = false
		TotC = NoteChannelDelete(a, b, TotC)
	}
	if len(a) < len(b) {
		Equal = false
		TotC = NoteChannelCreate(a, b, TotC)
	}
	if len(a) == len(b) {
		for _, Ca := range a {
			for _, Cb := range b {
				if Ca.ID == Cb.ID {
					Equal, TotC = CompareChannelstruct(Ca, Cb, TotC, Equal)
				}
			}
		}
	}
	return Equal, TotC
}

func CompareGuild(a *discordgo.Guild, b *discordgo.Guild, TotC *FullChangeStruct, Equal bool) (bool, *FullChangeStruct) {
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
	if a.MemberCount != b.MemberCount {
		Equal = false
		TotC.Guild.Membercount = true
	}
	if a.Region != b.Region {
		Equal = false
		TotC.Guild.Region = true
	}
	Equal, TotC = CompareChannels(a.Channels, b.Channels, TotC, Equal)
	Equal, TotC = CompareMembers(a.Members, b.Members, TotC, Equal)
	Equal, TotC = CompareRoles(a.Roles, b.Roles, TotC, Equal)
	return Equal, TotC
}
