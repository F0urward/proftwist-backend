package entities

import (
	"time"

	"github.com/google/uuid"
)

type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup  ChatType = "group"
)

type MemberRole string

const (
	MemberRoleMember MemberRole = "member"
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleOwner  MemberRole = "owner"
)

type Chat struct {
	ID          uuid.UUID
	Type        ChatType
	Title       string
	Description string
	AvatarURL   string
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ChatMember struct {
	ID       uuid.UUID
	ChatID   uuid.UUID
	UserID   uuid.UUID
	Role     MemberRole
	JoinedAt time.Time
	LastRead time.Time
}

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	Metadata  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChatWithMembers struct {
	Chat    *Chat
	Members []*ChatMember
}

func (c *Chat) IsUserMember(userID uuid.UUID, members []*ChatMember) bool {
	for _, member := range members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}

func (c *Chat) GetMemberRole(userID uuid.UUID, members []*ChatMember) MemberRole {
	for _, member := range members {
		if member.UserID == userID {
			return member.Role
		}
	}
	return ""
}

func (c *Chat) CanManageChat(userID uuid.UUID, members []*ChatMember) bool {
	role := c.GetMemberRole(userID, members)

	return role == MemberRoleOwner || role == MemberRoleAdmin
}

func (c *Chat) CanAddMember(userID uuid.UUID, members []*ChatMember) bool {
	role := c.GetMemberRole(userID, members)

	switch c.Type {
	case ChatTypeDirect:
		return false
	case ChatTypeGroup:
		return role == MemberRoleOwner || role == MemberRoleAdmin
	default:
		return false
	}
}

func (c *Chat) CanRemoveMember(requestedBy uuid.UUID, members []*ChatMember) bool {
	role := c.GetMemberRole(requestedBy, members)

	switch c.Type {
	case ChatTypeDirect:
		return false
	case ChatTypeGroup:
		return role == MemberRoleOwner || role == MemberRoleAdmin
	default:
		return false
	}
}

func (c *Chat) CanSendMessage(userID uuid.UUID, members []*ChatMember) bool {
	return c.IsUserMember(userID, members)
}

func (c *Chat) CanViewChat(userID uuid.UUID, members []*ChatMember) bool {
	return c.IsUserMember(userID, members)
}
