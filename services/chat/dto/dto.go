package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateGroupChatRequestDTO struct {
	Title         *string     `json:"title,omitempty"`
	AvatarURL     *string     `json:"avatar_url,omitempty"`
	RoadmapNodeID *string     `json:"roadmap_node_id,omitempty"`
	MemberIDs     []uuid.UUID `json:"member_ids"`
}

type CreateGroupChatResponseDTO struct {
	GroupChat GroupChatResponseDTO `json:"group_chat"`
}

type CreateDirectChatRequestDTO struct {
	OtherUserID uuid.UUID `json:"other_user_id"`
}

type CreateDirectChatResponseDTO struct {
	DirectChat DirectChatResponseDTO `json:"direct_chat"`
}

type GroupChatResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	Title         *string   `json:"title,omitempty"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	RoadmapNodeID *string   `json:"roadmap_node_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GroupChatListResponseDTO struct {
	GroupChats []GroupChatResponseDTO `json:"group_chats"`
}

type DirectChatResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	User1ID   uuid.UUID `json:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DirectChatListResponseDTO struct {
	DirectChats []DirectChatResponseDTO `json:"direct_chats"`
}

type MemberResponseDTO struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
}

type ChatMessageResponseDTO struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	ChatID    uuid.UUID              `json:"chat_id" db:"chat_id"`
	User      MemberResponseDTO      `json:"user" db:"user"`
	Content   string                 `json:"content" db:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

type ChatMemberListResponseDTO struct {
	Members []MemberResponseDTO `json:"members"`
}

type SendMessageRequestDTO struct {
	ChatID   uuid.UUID              `json:"chat_id" validate:"required"`
	UserID   uuid.UUID              `json:"-"`
	Content  string                 `json:"content" validate:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type GetChatMessagesResponseDTO struct {
	ChatMessages []ChatMessageResponseDTO `json:"chat_messages"`
}
