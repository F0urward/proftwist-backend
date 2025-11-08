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
