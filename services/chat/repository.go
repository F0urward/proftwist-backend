package chat

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	CreateGroupChat(ctx context.Context, chat *entities.GroupChat) (*entities.GroupChat, error)
	GetGroupChatByNode(ctx context.Context, nodeID string) (*entities.GroupChat, error)
	GetGroupChatsByUser(ctx context.Context, userID uuid.UUID) ([]*entities.GroupChat, error)
	GetGroupChatMembers(ctx context.Context, chatID uuid.UUID) ([]*entities.GroupChatMember, error)
	IsGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
	AddGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	AddGroupChatMembers(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID) error
	RemoveGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	SaveGroupMessage(ctx context.Context, message *entities.Message) error
	GetGroupChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	DeleteGroupChat(ctx context.Context, chatID uuid.UUID) error

	CreateDirectChat(ctx context.Context, chat *entities.DirectChat) (*entities.DirectChat, error)
	GetDirectChatByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*entities.DirectChat, error)
	GetDirectChatsByUser(ctx context.Context, userID uuid.UUID) ([]*entities.DirectChat, error)
	GetDirectChat(ctx context.Context, chatID uuid.UUID) (*entities.DirectChat, error)
	IsDirectChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
	SaveDirectMessage(ctx context.Context, message *entities.Message) error
	GetDirectChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	DeleteDirectChat(ctx context.Context, chatID uuid.UUID) error
}
