package chat

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/chat/dto"

	"github.com/google/uuid"
)

type Usecase interface {
	CreateGroupChat(ctx context.Context, userID uuid.UUID, req *dto.CreateGroupChatRequestDTO) (*dto.CreateGroupChatResponseDTO, error)
	DeleteGroupChat(ctx context.Context, chatID uuid.UUID) error
	GetGroupChatByNode(ctx context.Context, nodeID string) (*dto.GroupChatResponseDTO, error)
	GetGroupChatsByUser(ctx context.Context, userID uuid.UUID) (*dto.GroupChatListResponseDTO, error)
	GetGroupChatMembers(ctx context.Context, chatID uuid.UUID) (*dto.ChatMemberListResponseDTO, error)
	GetGroupChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error)
	SendGroupMessage(ctx context.Context, req *dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error)
	JoinGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	LeaveGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error

	CreateDirectChat(ctx context.Context, userID uuid.UUID, req *dto.CreateDirectChatRequestDTO) (*dto.CreateDirectChatResponseDTO, error)
	DeleteDirectChat(ctx context.Context, chatID uuid.UUID) error
	GetDirectChatsByUser(ctx context.Context, userID uuid.UUID) (*dto.DirectChatListResponseDTO, error)
	GetDirectChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error)
	SendDirectMessage(ctx context.Context, req *dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error)

	BroadcastTyping(ctx context.Context, chatID, userID uuid.UUID, typing bool, isGroup bool) error
	BroadcastUserJoined(ctx context.Context, chatID, userID uuid.UUID) error
	BroadcastUserLeft(ctx context.Context, chatID, userID uuid.UUID) error
}
