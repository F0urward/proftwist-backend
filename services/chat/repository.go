package chat

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

type Repository interface {
	GetGroupChatByNode(ctx context.Context, nodeID string) (*entities.GroupChat, error)
	GetGroupChatsByUser(ctx context.Context, userID uuid.UUID) ([]*entities.GroupChat, error)
	GetGroupChatMembers(ctx context.Context, chatID uuid.UUID) ([]*entities.GroupChatMember, error)
	IsGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
	AddGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	RemoveGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	SaveGroupMessage(ctx context.Context, message *entities.Message) error

	GeDirectChatsByUser(ctx context.Context, userID uuid.UUID) ([]*entities.DirectChat, error)
	GetDirectChat(ctx context.Context, chatID uuid.UUID) (*entities.DirectChat, error)
	IsDirectChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
	SaveDirectMessage(ctx context.Context, message *entities.Message) error

	GetGroupChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	GetDirectChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error)
}
