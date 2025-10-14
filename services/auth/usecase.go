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
	VKOauthLink(ctx context.Context) (*dto.VKOauthLinkResponse, error)
	VKOAuthCallback(ctx context.Context, request *dto.VKCallbackRequestDTO) (*dto.UserTokenDTO, error)
}
