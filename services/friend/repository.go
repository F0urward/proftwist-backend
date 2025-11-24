package friend

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

type Repository interface {
	CreateFriendship(ctx context.Context, userID, friendID, chatID uuid.UUID) error
	DeleteFriendship(ctx context.Context, userID, friendID uuid.UUID) error
	GetFriendIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	IsFriends(ctx context.Context, userID, friendID uuid.UUID) (bool, error)
	GetFriendshipChatID(ctx context.Context, userID, friendID uuid.UUID) (*uuid.UUID, error)

	CreateFriendRequest(ctx context.Context, request *entities.FriendRequest) error
	DeleteFriendRequest(ctx context.Context, requestID uuid.UUID) error
	GetFriendRequestByID(ctx context.Context, requestID uuid.UUID) (*entities.FriendRequest, error)
	GetFriendRequestsForUserByStatus(ctx context.Context, userID uuid.UUID, statuses []entities.FriendStatus) ([]*entities.FriendRequest, error)
	GetSentFriendRequestsByStatus(ctx context.Context, userID uuid.UUID, statuses []entities.FriendStatus) ([]*entities.FriendRequest, error)
	UpdateFriendRequestStatus(ctx context.Context, requestID uuid.UUID, status entities.FriendStatus) error
	UpdateFriendRequest(ctx context.Context, requestID uuid.UUID, fromUserID, toUserID uuid.UUID, status entities.FriendStatus) error
	GetFriendRequestBetweenUsers(ctx context.Context, fromUserID, toUserID uuid.UUID) (*entities.FriendRequest, error)
	GetPendingFriendRequestBetweenUsers(ctx context.Context, fromUserID, toUserID uuid.UUID) (*entities.FriendRequest, error)
}
