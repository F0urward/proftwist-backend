package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/friendclient"
	"github.com/F0urward/proftwist-backend/services/friend"
	"github.com/F0urward/proftwist-backend/services/friend/dto"
)

type FriendServer struct {
	uc friend.Usecase
	friendclient.UnimplementedFriendServiceServer
}

func NewFriendServer(usecase friend.Usecase) friendclient.FriendServiceServer {
	return &FriendServer{uc: usecase}
}

func (s *FriendServer) GetFriendshipStatus(ctx context.Context, req *friendclient.GetFriendshipStatusRequest) (*friendclient.GetFriendshipStatusResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &friendclient.GetFriendshipStatusResponse{
			Error: "invalid user id format",
		}, nil
	}

	targetUserID, err := uuid.Parse(req.TargetUserId)
	if err != nil {
		return &friendclient.GetFriendshipStatusResponse{
			Error: "invalid target user id format",
		}, nil
	}

	friendshipStatus, err := s.uc.GetFriendshipStatus(ctx, userID, targetUserID)
	if err != nil {
		return &friendclient.GetFriendshipStatusResponse{
			Error: err.Error(),
		}, nil
	}

	protoResponse := s.convertFriendshipStatusToProto(friendshipStatus)

	return protoResponse, nil
}

func (s *FriendServer) convertFriendshipStatusToProto(dto *dto.FriendshipStatusResponseDTO) *friendclient.GetFriendshipStatusResponse {
	response := &friendclient.GetFriendshipStatusResponse{
		Status:   dto.Status,
		IsSender: dto.IsSender,
	}

	if dto.RequestID != nil {
		response.RequestId = dto.RequestID.String()
	}

	return response
}
