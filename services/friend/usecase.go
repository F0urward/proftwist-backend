package friend

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/services/friend/dto"
)

type Usecase interface {
	GetFriends(ctx context.Context, userID uuid.UUID) (*dto.GetFriendsResponseDTO, error)
	DeleteFriend(ctx context.Context, userID, friendID uuid.UUID) error
	GetFriendRequests(ctx context.Context, userID uuid.UUID) (*dto.GetFriendRequestsResponseDTO, error)
	AcceptFriendRequest(ctx context.Context, userID, requestID uuid.UUID) (*dto.FriendResponseDTO, error)
	DeleteFriendRequest(ctx context.Context, userID, requestID uuid.UUID) error
	CreateFriendRequest(ctx context.Context, userID uuid.UUID, request *dto.CreateFriendRequestDTO) error
}
