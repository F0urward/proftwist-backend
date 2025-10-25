package entities

import (
	"time"

	"github.com/google/uuid"
)

type ChatType string

const (
	ChatTypeDirect  ChatType = "direct"
	ChatTypeGroup   ChatType = "group"
	ChatTypeChannel ChatType = "channel"
)

type MemberRole string

const (
	MemberRoleMember MemberRole = "member"
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleOwner  MemberRole = "owner"
)

type MessageType string

const (
	MessageTypeChat     MessageType = "chat"
	MessageTypeSystem   MessageType = "system"
	MessageTypePresence MessageType = "presence"
	MessageTypeError    MessageType = "error"
	MessageTypeTyping   MessageType = "typing"
	MessageTypeRead     MessageType = "read"
)

type Chat struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Type        ChatType  `json:"type" db:"type"`
	Title       string    `json:"title,omitempty" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	AvatarURL   string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ChatMember struct {
	ID       uuid.UUID  `json:"id" db:"id"`
	ChatID   uuid.UUID  `json:"chat_id" db:"chat_id"`
	UserID   uuid.UUID  `json:"user_id" db:"user_id"`
	Role     MemberRole `json:"role" db:"role"`
	JoinedAt time.Time  `json:"joined_at" db:"joined_at"`
	LastRead time.Time  `json:"last_read" db:"last_read"`
}

type Message struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	ChatID    uuid.UUID              `json:"chat_id" db:"chat_id"`
	UserID    uuid.UUID              `json:"user_id" db:"user_id"`
	Content   string                 `json:"content" db:"content"`
	Type      MessageType            `json:"type" db:"type"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
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

// Добавляем методы проверки прав для разных типов чатов
func (c *Chat) CanAddMembers(userID uuid.UUID, members []*ChatMember, newMemberIDs []uuid.UUID) bool {
	switch c.Type {
	case ChatTypeChannel:
		// В канал можно добавлять только владельца при создании
		return len(newMemberIDs) == 0
	case ChatTypeDirect:
		// В директ можно добавить только одного участника кроме создателя
		return len(newMemberIDs) == 1
	case ChatTypeGroup:
		// В группу можно добавить любое количество участников
		return len(newMemberIDs) >= 0
	default:
		return false
	}
}

func (c *Chat) CanRemoveMember(requestedBy uuid.UUID, targetUserID uuid.UUID, members []*ChatMember) bool {
	role := c.GetMemberRole(requestedBy, members)
	//targetRole := c.GetMemberRole(targetUserID, members)

	switch c.Type {
	case ChatTypeChannel:
		// В канале удалить может владелец или сам пользователь
		return role == MemberRoleOwner || requestedBy == targetUserID
	case ChatTypeDirect:
		// Из директа нельзя выйти
		return false
	case ChatTypeGroup:
		// В группе удалить может владелец или сам пользователь
		return role == MemberRoleOwner || requestedBy == targetUserID
	default:
		return false
	}
}

func (c *Chat) CanSendMessage(userID uuid.UUID, members []*ChatMember) bool {
	if !c.IsUserMember(userID, members) {
		return false
	}

	switch c.Type {
	case ChatTypeChannel:
		// В канале писать может только владелец
		role := c.GetMemberRole(userID, members)
		return role == MemberRoleOwner
	case ChatTypeDirect, ChatTypeGroup:
		// В директ и группу могут писать все участники
		return true
	default:
		return false
	}
}

func (c *Chat) CanDeleteChat(userID uuid.UUID, members []*ChatMember) bool {
	role := c.GetMemberRole(userID, members)

	switch c.Type {
	case ChatTypeChannel, ChatTypeGroup:
		// Канал и группу может удалить только владелец
		return role == MemberRoleOwner
	case ChatTypeDirect:
		// Директ могут удалить оба участника
		return c.IsUserMember(userID, members)
	default:
		return false
	}
}

func (c *Chat) CanViewChat(userID uuid.UUID, members []*ChatMember) bool {
	// Все чаты могут смотреть только участники
	return c.IsUserMember(userID, members)
}
