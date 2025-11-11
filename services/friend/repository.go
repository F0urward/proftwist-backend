package friend

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

type Repository interface {
	CreateFriendship(ctx context.Context, userID, friendID uuid.UUID) error
	DeleteFriendship(ctx context.Context, userID, friendID uuid.UUID) error
	GetFriendIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	IsFriends(ctx context.Context, userID, friendID uuid.UUID) (bool, error)

	CreateFriendRequest(ctx context.Context, request *entities.FriendRequest) error
	GetFriendRequestByID(ctx context.Context, requestID uuid.UUID) (*entities.FriendRequest, error)
	GetFriendRequestsForUser(ctx context.Context, userID uuid.UUID) ([]*entities.FriendRequest, error)
	GetSentFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entities.FriendRequest, error)
	UpdateFriendRequestStatus(ctx context.Context, requestID uuid.UUID, status entities.FriendStatus) error
	DeleteFriendRequest(ctx context.Context, requestID uuid.UUID) error
	GetFriendRequestBetweenUsers(ctx context.Context, fromUserID, toUserID uuid.UUID) (*entities.FriendRequest, error)
	GetPendingFriendRequestBetweenUsers(ctx context.Context, fromUserID, toUserID uuid.UUID) (*entities.FriendRequest, error)
}
