package dto

import (
	"time"

	"github.com/google/uuid"
)

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

type ChatMessageResponseDTO struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	ChatID    uuid.UUID              `json:"chat_id" db:"chat_id"`
	UserID    uuid.UUID              `json:"user_id" db:"user_id"`
	Content   string                 `json:"content" db:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

type ChatMemberResponseDTO struct {
	UserID uuid.UUID `json:"user_id"`
}

type ChatMemberListResponseDTO struct {
	Members []ChatMemberResponseDTO `json:"members"`
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
