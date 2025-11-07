package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/auth/dto"
)

type AuthServer struct {
	uc auth.Usecase
	authclient.UnimplementedAuthServiceServer
}

func NewAuthServer(usecase auth.Usecase) *AuthServer {
	return &AuthServer{uc: usecase}
}

func (s *AuthServer) GetUserByID(ctx context.Context, req *authclient.GetUserByIDRequest) (*authclient.GetUserByIDResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &authclient.GetUserByIDResponse{
			Error: "invalid user ID format",
		}, nil
	}

	user, err := s.uc.GetByID(ctx, userID)
	if err != nil {
		return &authclient.GetUserByIDResponse{
			Error: err.Error(),
		}, nil
	}

	protoUser := convertUserToProto(&user.User)

	return &authclient.GetUserByIDResponse{
		User: protoUser,
	}, nil
}

func (s *AuthServer) GetUsersByIDs(ctx context.Context, req *authclient.GetUsersByIDsRequest) (*authclient.GetUsersByIDsResponse, error) {
	if len(req.UserIds) == 0 {
		return &authclient.GetUsersByIDsResponse{
			Users: []*authclient.User{},
		}, nil
	}

	userIDs := make([]uuid.UUID, 0, len(req.UserIds))
	for _, idStr := range req.UserIds {
		userID, err := uuid.Parse(idStr)
		if err != nil {
			return &authclient.GetUsersByIDsResponse{
				Error: "invalid user ID format: " + idStr,
			}, nil
		}
		userIDs = append(userIDs, userID)
	}

	users, err := s.uc.GetByIDs(ctx, userIDs)
	if err != nil {
		return &authclient.GetUsersByIDsResponse{
			Error: err.Error(),
		}, nil
	}

	protoUsers := make([]*authclient.User, 0, len(users.Users))
	for _, user := range users.Users {
		protoUsers = append(protoUsers, convertUserToProto(&user))
	}

	return &authclient.GetUsersByIDsResponse{
		Users: protoUsers,
	}, nil
}

func convertUserToProto(dto *dto.UserDTO) *authclient.User {
	if dto == nil {
		return nil
	}

	return &authclient.User{
		Id:        dto.ID.String(),
		Username:  dto.Username,
		Email:     dto.Email,
		AvatarUrl: dto.AvatarUrl,
	}
}
