package Belt

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type PermORChange struct {
	ID          string
	Allow       bool
	Deny        bool
	ExistCrisis bool
	Mk          bool
	Del         bool
}

type ChannelChange struct {
	ID          string
	Name        bool
	Topic       bool
	perms       bool
	Perms       []*PermORChange
	ExistCrisis bool
	Del         bool
	Mk          bool
}

type RoleChange struct {
	ID          string
	Name        bool
	Perms       bool
	Position    bool
	Color       bool
	ExistCrisis bool
	Del         bool
	Mk          bool
}

type MemberChange struct {
	User struct {
		ID       string
		Username bool
		Avatar   bool
	}
	Nick        bool
	Roles       bool
	RoleNew     bool
	RoleRem     bool
	ExistCrisis bool
	Leave       bool
	Join        bool
}

type FullChangeStruct struct {
	Guild struct {
		Name     bool
		OwnerID  bool
		Icon     bool
		Region   bool
		channels bool
		roles    bool
		members  bool
	}
	Channels []*ChannelChange
	Roles    []*RoleChange
	Members  []*MemberChange
}

type GuildInfo struct {
	g           *discordgo.Guild `json:"g"`
	Lastcheck   TimeFormat       `json:"TimeForm"`
	BotUP       bool             `json:"BU"`
	NeedRestall bool             `json:"RESTALL"`
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
}
