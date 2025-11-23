package auth

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/auth/dto"
	"github.com/google/uuid"
)

type Usecase interface {
	Register(context.Context, *dto.RegisterRequestDTO) (*dto.UserTokenDTO, error)
	Login(context.Context, *dto.LoginRequestDTO) (*dto.UserTokenDTO, error)
	Logout(context.Context, string) error
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserDTO, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*dto.GetUserByIDResponseDTO, error)
	GetByIDs(ctx context.Context, userIDs []uuid.UUID) (*dto.GetUsersByIDsResponseDTO, error)
	Update(ctx context.Context, userID uuid.UUID, request *dto.UpdateUserRequestDTO) error
	UploadAvatar(ctx context.Context, request *dto.UploadAvatarRequestDTO) (*dto.UploadAvatarResponseDTO, error)
	IsInBlacklist(ctx context.Context, userID string, token string) (bool, error)
	VKOauthLink(ctx context.Context) (*dto.VKOauthLinkResponse, error)
	VKOAuthCallback(ctx context.Context, request *dto.VKCallbackRequestDTO) (*dto.UserTokenDTO, error)
}
