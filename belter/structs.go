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

type FullMention struct {
	ChannelOR []*ChannelChange
	Members   []*MemberChange
	Roles     []*RoleChange
	OwnerID   string
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

type PrimeSuspectChange struct {
	ID       string
	member   bool
	MC       PSMemberchange
	roles    bool
	RC       PSRolechange
	channels bool
	CC       PSchannelchange
}

type PSMemberchange struct {
	name bool // username
	Name struct {
		NName string
		OName string
	}
	nick bool
	Nick struct {
		NNick string
		OBick string
	}
	Owner      bool
	Isownernow bool
}
type PSRolechange struct { // Role that the PS has that is changed or added
	ID    string
	allow bool
	Allow struct {
		Last string
		New  string
	}
	deny bool
	Deny struct {
		Last string
		New  string
	}
	pos         bool
	PosOld      int
	PosNew      int
	Existcrisis bool
	Deleted     bool
	Dperm       string //permissions before the delete
	Made        bool
	Mper        string //permissions after the creation, ignore if 0
}

type PSchannelchange struct { // Channel with OR permission change related to the PS
	ID    string
	Type  string
	allow bool
	Allow struct {
		Last string
		New  string
	}
	deny bool
	Deny struct {
		Last string
		New  string
	}
	Existentcrisis bool
	Deleted        bool
	Dperm          string //permissions before the delete
	Made           bool
	Mper           string //permissions after the creation, ignore if 0
}

type Permissions struct {
	CREATE_INSTANT_INVITE bool //0x00000001	1	Allows creation of instant invites
	KICK_MEMBERS          bool //0x00000002	10	Allows kicking members
	BAN_MEMBERS           bool //0x00000004	100	Allows banning members
	ADMINISTRATOR         bool //0x00000008	1000	Allows all permissions and bypasses channel permission overwrites
	MANAGE_CHANNELS       bool //0x00000010	10000	Allows management and editing of channels
	MANAGE_GUILD          bool //0x00000020	100000	Allows management and editing of the guild
	ADD_REACTIONS         bool //0x00000040	1000000	Allows for the addition of reactions to messages
	READ_MESSAGES         bool //0x00000400	10000000000	Allows reading messages in a channel. The channel will not appear for users without this permission
	SEND_MESSAGES         bool //0x00000800	100000000000	Allows for sending messages in a channel.
	SEND_TTS_MESSAGES     bool //0x00001000	1000000000000	Allows for sending of /tts messages
	MANAGE_MESSAGES       bool //0x00002000	10000000000000	Allows for deletion of other users messages
	EMBED_LINKS           bool //0x00004000	100000000000000	Links sent by this user will be auto-embedded
	ATTACH_FILES          bool //0x00008000	1000000000000000	Allows for uploading images and files
	READ_MESSAGE_HISTORY  bool //0x00010000	10000000000000000	Allows for reading of message history
	MENTION_EVERYONE      bool //0x00020000	100000000000000000	Allows for using the @everyone tag to notify all users in a channel, and the @here tag to notify all online users in a channel
	USE_EXTERNAL_EMOJIS   bool //0x00040000	1000000000000000000	Allows the usage of custom emojis from other servers
	CONNECT               bool //0x00100000	100000000000000000000	Allows for joining of a voice channel
	SPEAK                 bool //0x00200000	1000000000000000000000	Allows for speaking in a voice channel
	MUTE_MEMBERS          bool //0x00400000	10000000000000000000000	Allows for muting members in a voice channel
	DEAFEN_MEMBERS        bool //0x00800000	100000000000000000000000	Allows for deafening of members in a voice channel
	MOVE_MEMBERS          bool //0x01000000	1000000000000000000000000	Allows for moving of members between voice channels
	USE_VAD               bool //0x02000000	10000000000000000000000000	Allows for using voice-activity-detection in a voice channel
	CHANGE_NICKNAME       bool //0x04000000	100000000000000000000000000	Allows for modification of own nickname
	MANAGE_NICKNAMES      bool //0x08000000	1000000000000000000000000000	Allows for modification of other users nicknames
	MANAGE_ROLES          bool //0x10000000	10000000000000000000000000000	Allows management and editing of roles
	MANAGE_WEBHOOKS       bool //0x20000000	100000000000000000000000000000	Allows management and editing of webhooks
	MANAGE_EMOJIS         bool //0x40000000	1000000000000000000000000000000	Allows management and editing of emojis
}

type PermissionBit struct {
	CREATE_INSTANT_INVITE int //0x00000001	1	Allows creation of instant invites
	KICK_MEMBERS          int //0x00000002	10	Allows kicking members
	BAN_MEMBERS           int //0x00000004	100	Allows banning members
	ADMINISTRATOR         int //0x00000008	1000	Allows all permissions and bypasses channel permission overwrites
	MANAGE_CHANNELS       int //0x00000010	10000	Allows management and editing of channels
	MANAGE_GUILD          int //0x00000020	100000	Allows management and editing of the guild
	ADD_REACTIONS         int //0x00000040	1000000	Allows for the addition of reactions to messages
	READ_MESSAGES         int //0x00000400	10000000000	Allows reading messages in a channel. The channel will not appear for users without this permission
	SEND_MESSAGES         int //0x00000800	100000000000	Allows for sending messages in a channel.
	SEND_TTS_MESSAGES     int //0x00001000	1000000000000	Allows for sending of /tts messages
	MANAGE_MESSAGES       int //0x00002000	10000000000000	Allows for deletion of other users messages
	EMBED_LINKS           int //0x00004000	100000000000000	Links sent by this user will be auto-embedded
	ATTACH_FILES          int //0x00008000	1000000000000000	Allows for uploading images and files
	READ_MESSAGE_HISTORY  int //0x00010000	10000000000000000	Allows for reading of message history
	MENTION_EVERYONE      int //0x00020000	100000000000000000	Allows for using the @everyone tag to notify all users in a channel, and the @here tag to notify all online users in a channel
	USE_EXTERNAL_EMOJIS   int //0x00040000	1000000000000000000	Allows the usage of custom emojis from other servers
	CONNECT               int //0x00100000	100000000000000000000	Allows for joining of a voice channel
	SPEAK                 int //0x00200000	1000000000000000000000	Allows for speaking in a voice channel
	MUTE_MEMBERS          int //0x00400000	10000000000000000000000	Allows for muting members in a voice channel
	DEAFEN_MEMBERS        int //0x00800000	100000000000000000000000	Allows for deafening of members in a voice channel
	MOVE_MEMBERS          int //0x01000000	1000000000000000000000000	Allows for moving of members between voice channels
	USE_VAD               int //0x02000000	10000000000000000000000000	Allows for using voice-activity-detection in a voice channel
	CHANGE_NICKNAME       int //0x04000000	100000000000000000000000000	Allows for modification of own nickname
	MANAGE_NICKNAMES      int //0x08000000	1000000000000000000000000000	Allows for modification of other users nicknames
	MANAGE_ROLES          int //0x10000000	10000000000000000000000000000	Allows management and editing of roles
	MANAGE_WEBHOOKS       int //0x20000000	100000000000000000000000000000	Allows management and editing of webhooks
	MANAGE_EMOJIS         int //0x40000000	1000000000000000000000000000000	Allows management and editing of emojis
}

func installPerms() {
	PER = &PermissionBit{
		CREATE_INSTANT_INVITE: 0,  //0x00000001	1	Allows creation of instant invites
		KICK_MEMBERS:          1,  //0x00000002	10	Allows kicking members
		BAN_MEMBERS:           2,  //0x00000004	100	Allows banning members
		ADMINISTRATOR:         3,  //0x00000008	1000	Allows all permissions and bypasses channel permission overwrites
		MANAGE_CHANNELS:       4,  //0x00000010	10000	Allows management and editing of channels
		MANAGE_GUILD:          5,  //0x00000020	100000	Allows management and editing of the guild
		ADD_REACTIONS:         6,  //0x00000040	1000000	Allows for the addition of reactions to messages
		READ_MESSAGES:         10, //0x00000400	10000000000	Allows reading messages in a channel. The channel will not appear for users without this permission
		SEND_MESSAGES:         11, //0x00000800	100000000000	Allows for sending messages in a channel.
		SEND_TTS_MESSAGES:     12, //0x00001000	1000000000000	Allows for sending of /tts messages
		MANAGE_MESSAGES:       13, //0x00002000	10000000000000	Allows for deletion of other users messages
		EMBED_LINKS:           14, //0x00004000	100000000000000	Links sent by this user will be auto-embedded
		ATTACH_FILES:          15, //0x00008000	1000000000000000	Allows for uploading images and files
		READ_MESSAGE_HISTORY:  16, //0x00010000	10000000000000000	Allows for reading of message history
		MENTION_EVERYONE:      17, //0x00020000	100000000000000000	Allows for using the @everyone tag to notify all users in a channel, and the @here tag to notify all online users in a channel
		USE_EXTERNAL_EMOJIS:   18, //0x00040000	1000000000000000000	Allows the usage of custom emojis from other servers
		CONNECT:               20, //0x00100000	100000000000000000000	Allows for joining of a voice channel
		SPEAK:                 21, //0x00200000	1000000000000000000000	Allows for speaking in a voice channel
		MUTE_MEMBERS:          22, //0x00400000	10000000000000000000000	Allows for muting members in a voice channel
		DEAFEN_MEMBERS:        23, //0x00800000	100000000000000000000000	Allows for deafening of members in a voice channel
		MOVE_MEMBERS:          24, //0x01000000	1000000000000000000000000	Allows for moving of members between voice channels
		USE_VAD:               25, //0x02000000	10000000000000000000000000	Allows for using voice-activity-detection in a voice channel
		CHANGE_NICKNAME:       26, //0x04000000	100000000000000000000000000	Allows for modification of own nickname
		MANAGE_NICKNAMES:      27, //0x08000000	1000000000000000000000000000	Allows for modification of other users nicknames
		MANAGE_ROLES:          28, //0x10000000	10000000000000000000000000000	Allows management and editing of roles
		MANAGE_WEBHOOKS:       29, //0x20000000	100000000000000000000000000000	Allows management and editing of webhooks
		MANAGE_EMOJIS:         30, //0x40000000	1000000000000000000000000000000	Allows management and editing of emojis
	}
}
