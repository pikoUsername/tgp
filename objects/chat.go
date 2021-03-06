package objects

import "time"

// Chat type
type Chat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

// https://core.telegram.org/bots/api#chatmember
type ChatMember struct {
	// User ofc
	User        `json:"user"`
	Status      string `json:"status"`
	CustomTitle string `json:"custom_title"`
	IsAnon      bool   `json:"is_anonymous"`
	IsMember    bool   `json:"is_member"`
	UntilDate   int64  `json:"until_date"`

	// Copy paste from perms, yeah broken DRY!
	CanBeEdited           bool `json:"can_be_edited"`
	CanManageChat         bool `json:"can_manage_chat"`
	CanPostMessages       bool `json:"can_post_message"`
	CanEditMessages       bool `json:"can_edit_messages"`
	CanDeleteMessage      bool `json:"can_delete_message"`
	CanManageVideochats   bool `json:"can_manage_video_chats"`
	CanRestrictMembers    bool `json:"can_restrict_members"`
	CanPromoteMembers     bool `json:"can_promote_members"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
	CanSendMessage        bool `json:"can_send_message"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessage   bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
}

type ChatMemberMember struct {
	Status string `json:"status"`
	User   *User  `json:"user"`
}

// ChatInviteLink represents ChatInvite object
// https://core.telegram.org/bots/api#chatinvitelink
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 *User  `json:"creator"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
	Name                    string `json:"name"`
	ExpireDate              int64  `json:"expire_date"`
	MemberLimit             uint   `json:"member_limit"`
	PendingJoinRequestCount int    `json:"pending_join_request_count"`
}

type ChatMemberRestricted struct {
	ChatMemberMember
	ChatMemberPermissions
	UntilDate int `json:"until_date"`
}

type ChatMemberLeft struct {
	ChatMemberMember
}

type ChatMemberBanned struct {
	ChatMemberMember
	UntilDate time.Duration `json:"until_date"`
}

type ChatMemberOwner struct {
	Status      string `json:"status"`
	User        *User  `json:"user"`
	IsAnonymous bool   `json:"is_anonymous"`
	CustomTitle string `json:"custom_title"`
}

// ChatMemberUpdated object represents changes in the status of a chat member.
// https://core.telegram.org/bots/api#chatmemberupdated
type ChatMemberUpdated struct {
	Chat          *Chat           `json:"chat"`
	From          *User           `json:"user"`
	Date          uint64          `json:"date"`
	OldChatMember *ChatMember     `json:"old_chat_member"`
	NewChatMember *ChatMember     `json:"new_chat_member"`
	InviteLink    *ChatInviteLink `json:"invite_link"`
}

type ChatMemberAdministrator struct {
	Status              string `json:"status"`
	User                *User  `json:"user"`
	IsAnonymous         bool   `json:"is_anonymous"`
	CanBeEdited         bool   `json:"can_be_edited"`
	CanManageChat       bool   `json:"can_manage_chat"`
	CanPostMessages     bool   `json:"can_post_message"`
	CanEditMessages     bool   `json:"can_edit_messages"`
	CanDeleteMessage    bool   `json:"can_delete_message"`
	CanManageVideochats bool   `json:"can_manage_video_chats"`
	CanRestrictMembers  bool   `json:"can_restrict_members"`
	CanPromoteMembers   bool   `json:"can_promote_members"`
	CanChangeInfo       bool   `json:"can_change_info"`
	CanInviteUsers      bool   `json:"can_invite_users"`
	CanPinMessages      bool   `json:"can_pin_messages"`
	CustomTitle         string `json:"custom_title"`
}

// ChatPhoto object represents photo of chat
// https://core.telegram.org/bots/api#chatphoto
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

type ChatMemberPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
}

type ChatJoinRequest struct {
	Chat       *Chat          `json:"chat"`
	From       *User          `json:"from"`
	Date       int64          `json:"date"`
	Bio        string         `json:"bio"`
	InviteLink ChatInviteLink `json:"invite_link"`
}
