package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChatMessageResponseDTO struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	ChatID    uuid.UUID              `json:"chat_id" db:"chat_id"`
	UserID    uuid.UUID              `json:"user_id" db:"user_id"`
	Content   string                 `json:"content" db:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

type ChatResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UnreadCount int       `json:"unread_count,omitempty"`
}

type CreateChatRequestDTO struct {
	Type          string      `json:"type" validate:"required,oneof=direct group"`
	Title         string      `json:"title,omitempty"`
	Description   string      `json:"description,omitempty"`
	AvatarURL     string      `json:"avatar_url,omitempty"`
	CreatedByID   uuid.UUID   `json:"-"`
	CreatedByRole string      `json:"-"`
	MemberIDs     []uuid.UUID `json:"member_ids,omitempty"`
}

type SendMessageRequestDTO struct {
	ChatID   uuid.UUID              `json:"chat_id" validate:"required"`
	UserID   uuid.UUID              `json:"-"`
	Content  string                 `json:"content" validate:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type AddMemberRequestDTO struct {
	ChatID      uuid.UUID `json:"chat_id" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Role        string    `json:"role" validate:"required,oneof=member admin owner"`
	RequestedBy uuid.UUID `json:"-"`
}

type RemoveMemberRequestDTO struct {
	ChatID      uuid.UUID `json:"chat_id" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	RequestedBy uuid.UUID `json:"-"`
}

type GetMessagesRequestDTO struct {
	ChatID uuid.UUID `json:"chat_id" validate:"required"`
	Limit  int       `json:"limit" validate:"min=1,max=100"`
	Offset int       `json:"offset" validate:"min=0"`
}

type ChatMemberResponseDTO struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	LastRead time.Time `json:"last_read"`
}

type GetChatMessagesResponseDTO struct {
	ChatMessages []ChatMessageResponseDTO `json:"chat_messages"`
}

type DeleteChatRequestDTO struct {
	ChatID      uuid.UUID `json:"chat_id" validate:"required"`
	RequestedBy uuid.UUID `json:"-"`
}

func (r *CreateChatRequestDTO) Validate() error {
	switch r.Type {
	case "direct":
		if len(r.MemberIDs) != 1 {
			return fmt.Errorf("direct chat must have exactly one member")
		}
	case "group":
	}
	return nil
}
