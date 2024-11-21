package enum

type UserRole string

const (
	UserRoleCreator UserRole = "CREATOR"
)

type UserState string

const (
	UserStateDefault    UserState = "DEFAULT"
	UserStateWaitlisted UserState = "WAITLISTED"
)

type ChannelType string

const (
	ChannelTypeInstagram ChannelType = "INSTAGRAM"
	ChannelTypeTiktok    ChannelType = "TIKTOK"
)

type ChannelMemberRole string

const (
	ChannelMemberRoleOwner ChannelMemberRole = "OWNER"
)

type ChannelMemberVisibility string

const (
	ChannelMemberVisibilityDefault ChannelMemberVisibility = "DEFAULT"
)
