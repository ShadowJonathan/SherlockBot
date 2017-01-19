package Belt

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type PermORChange struct {
	ID          string `json:"ID"`
	Allow       bool   `json:"A"`
	Deny        bool   `json:"D"`
	ExistCrisis bool   `json:"Ex,omitempty"`
	Mk          bool   `json:"Mk,omitempty"`
	Del         bool   `json:"Del,omitempty"`
}

type ChannelChange struct {
	ID          string          `json:"ID"`
	Name        bool            `json:"Name"`
	Topic       bool            `json:"Topic"`
	perms       bool            `json:"Perms"`
	Perms       []*PermORChange `json:"ORs"`
	ExistCrisis bool            `json:"Ex,omitempty"`
	Del         bool            `json:"Del,omitempty"`
	Mk          bool            `json:"Mk,omitempty"`
}

type RoleChange struct {
	ID          string `json:"ID"`
	Name        bool   `json:"Name"`
	Perms       bool   `json:"Perms"`
	Position    bool   `json:"Pos"`
	Color       bool   `json:"Color"`
	ExistCrisis bool   `json:"Ex,omitempty"`
	Del         bool   `json:"Del,omitempty"`
	Mk          bool   `json:"Mk,omitempty"`
}

type MemberChange struct {
	User struct {
		ID       string `json:"ID"`
		Username bool   `json:"UN"`
		Avatar   bool   `json:"AV"`
	}
	Nick        bool `json:"Nick"`
	Roles       bool `json:"Roles"`
	RoleNew     bool `json:"RN,omitempty"`
	RoleRem     bool `json:"RR,omitempty"`
	ExistCrisis bool `json:"Ex,omitempty"`
	Leave       bool `json:"L,omitempty"`
	Join        bool `json:"J,omitempty"`
}

type FullChangeStruct struct {
	Guild struct {
		Name     bool `json:"Name"`
		OwnerID  bool `json:"OwnerID"`
		Icon     bool `json:"Icon"`
		Region   bool `json:"Region"`
		channels bool `json:"Channels"`
		roles    bool `json:"Roles"`
		members  bool `json:"Members"`
	}
	Channels []*ChannelChange `json:"CHs"`
	Roles    []*RoleChange    `json:"Rs"`
	Members  []*MemberChange  `json:"Ms"`
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
