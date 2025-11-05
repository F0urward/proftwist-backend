package chat

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/chat/dto"

	"github.com/google/uuid"
)

type Usecase interface {
	CreateChat(ctx context.Context, req dto.CreateChatRequestDTO) (*dto.ChatResponseDTO, error)
	SendMessage(ctx context.Context, req dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error)
	GetChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error)
	AddMember(ctx context.Context, req dto.AddMemberRequestDTO) error
	RemoveMember(ctx context.Context, req dto.RemoveMemberRequestDTO) error
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]dto.ChatResponseDTO, error)
	GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]dto.ChatMemberResponseDTO, error)
	JoinGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	LeaveGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error

	BroadcastTyping(ctx context.Context, chatID, userID uuid.UUID, typing bool) error
	BroadcastUserJoined(ctx context.Context, chatID, userID uuid.UUID) error
	BroadcastUserLeft(ctx context.Context, chatID, userID uuid.UUID) error
}
