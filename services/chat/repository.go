package chat

import (
	"context"
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

type Repository interface {
	CreateChat(ctx context.Context, chat *entities.Chat) error
	GetChat(ctx context.Context, chatID uuid.UUID) (*entities.Chat, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*entities.Chat, error)
	SaveMessage(ctx context.Context, message *entities.Message) error
	GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	AddChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, role entities.MemberRole) error
	RemoveChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*entities.ChatMember, error)
	IsChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
}
