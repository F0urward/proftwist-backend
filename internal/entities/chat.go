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

type GroupChat struct {
	ID            uuid.UUID
	Title         *string
	AvatarURL     *string
	RoadmapNodeID *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DirectChat struct {
	ID        uuid.UUID
	User1ID   uuid.UUID
	User2ID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GroupChatMember struct {
	ID          uuid.UUID
	GroupChatID uuid.UUID
	UserID      uuid.UUID
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

// func (c *Chat) IsUserMember(userID uuid.UUID, members []*ChatMember) bool {
// 	for _, member := range members {
// 		if member.UserID == userID {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (c *Chat) GetMemberRole(userID uuid.UUID, members []*ChatMember) MemberRole {
// 	for _, member := range members {
// 		if member.UserID == userID {
// 			return member.Role
// 		}
// 	}
// 	return ""
// }

// func (c *Chat) CanManageChat(userID uuid.UUID, members []*ChatMember) bool {
// 	role := c.GetMemberRole(userID, members)

// 	return role == MemberRoleOwner || role == MemberRoleAdmin
// }

// func (c *Chat) CanAddMember(userID uuid.UUID, members []*ChatMember) bool {
// 	role := c.GetMemberRole(userID, members)

// 	switch c.Type {
// 	case ChatTypeDirect:
// 		return false
// 	case ChatTypeGroup:
// 		return role == MemberRoleOwner || role == MemberRoleAdmin
// 	default:
// 		return false
// 	}
// }

// func (c *Chat) CanRemoveMember(requestedBy uuid.UUID, members []*ChatMember) bool {
// 	role := c.GetMemberRole(requestedBy, members)

// 	switch c.Type {
// 	case ChatTypeDirect:
// 		return false
// 	case ChatTypeGroup:
// 		return role == MemberRoleOwner || role == MemberRoleAdmin
// 	default:
// 		return false
// 	}
// }

// func (c *Chat) CanSendMessage(userID uuid.UUID, members []*ChatMember) bool {
// 	return c.IsUserMember(userID, members)
// }

// func (c *Chat) CanViewChat(userID uuid.UUID, members []*ChatMember) bool {
// 	return c.IsUserMember(userID, members)
// }
