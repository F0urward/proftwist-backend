package chat

import (
	"context"
	"github.com/F0urward/proftwist-backend/services/chat/dto"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

type ChatUseCaseInterface interface {
	CreateChat(ctx context.Context, req dto.CreateChatRequest) (*entities.Chat, error)
	SendMessage(ctx context.Context, req dto.SendMessageRequest) (*entities.Message, error)
	GetChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	AddMember(ctx context.Context, req dto.AddMemberRequest) error
	RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*entities.Chat, error)
	GetChatWithMembers(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (*entities.ChatWithMembers, error)
	GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*entities.ChatMember, error)
	DeleteChat(ctx context.Context, req dto.DeleteChatRequest) error
	JoinChannel(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	IsChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
}
