package dto

import (
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateChatRequest struct {
	Type        string      `json:"type" validate:"required,oneof=direct group channel"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	AvatarURL   string      `json:"avatar_url,omitempty"`
	CreatedBy   uuid.UUID   `json:"-"`
	MemberIDs   []uuid.UUID `json:"member_ids,omitempty"`
}

type SendMessageRequest struct {
	ChatID   uuid.UUID              `json:"chat_id" validate:"required"`
	UserID   uuid.UUID              `json:"-"`
	Content  string                 `json:"content" validate:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type AddMemberRequest struct {
	ChatID      uuid.UUID `json:"chat_id" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Role        string    `json:"role" validate:"required,oneof=member admin owner"`
	RequestedBy uuid.UUID `json:"-"`
}

type RemoveMemberRequest struct {
	ChatID      uuid.UUID `json:"chat_id" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	RequestedBy uuid.UUID `json:"-"`
}

type GetMessagesRequest struct {
	ChatID uuid.UUID `json:"chat_id" validate:"required"`
	Limit  int       `json:"limit" validate:"min=1,max=100"`
	Offset int       `json:"offset" validate:"min=0"`
}

type ChatResponse struct {
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

type MessageResponse struct {
	ID        uuid.UUID              `json:"id"`
	ChatID    uuid.UUID              `json:"chat_id"`
	UserID    uuid.UUID              `json:"user_id"`
	Content   string                 `json:"content"`
	Type      string                 `json:"type"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type ChatMemberResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	LastRead time.Time `json:"last_read"`
}

type ChatWithMembersResponse struct {
	Chat    *ChatResponse        `json:"chat"`
	Members []ChatMemberResponse `json:"members"`
}

func ToChatResponse(chat *entities.Chat) ChatResponse {
	return ChatResponse{
		ID:          chat.ID,
		Type:        string(chat.Type),
		Title:       chat.Title,
		Description: chat.Description,
		AvatarURL:   chat.AvatarURL,
		CreatedBy:   chat.CreatedBy,
		CreatedAt:   chat.CreatedAt,
		UpdatedAt:   chat.UpdatedAt,
	}
}

func ToMessageResponse(message *entities.Message) MessageResponse {
	return MessageResponse{
		ID:        message.ID,
		ChatID:    message.ChatID,
		UserID:    message.UserID,
		Content:   message.Content,
		Type:      string(message.Type),
		Metadata:  message.Metadata,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func ToChatMemberResponse(member *entities.ChatMember) ChatMemberResponse {
	return ChatMemberResponse{
		ID:       member.ID,
		UserID:   member.UserID,
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
		LastRead: member.LastRead,
	}
}

var validate = validator.New()

func (r *CreateChatRequest) Validate() error {
	// Дополнительная валидация в зависимости от типа чата
	switch r.Type {
	case "direct":
		if len(r.MemberIDs) != 1 {
			return fmt.Errorf("direct chat must have exactly one member")
		}
	case "channel":
		if len(r.MemberIDs) > 0 {
			return fmt.Errorf("channel cannot have initial members")
		}
	case "group":
		// Группа может иметь любое количество участников
	}

	return validate.Struct(r)
}

type DeleteChatRequest struct {
	ChatID      uuid.UUID `json:"chat_id" validate:"required"`
	RequestedBy uuid.UUID `json:"-"`
}
