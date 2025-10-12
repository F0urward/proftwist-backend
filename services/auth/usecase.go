package auth

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/auth/dto"
)

type Usecase interface {
	Register(context.Context, *dto.RegisterRequestDTO) (*dto.UserTokenDTO, error)
	Login(context.Context, *dto.LoginRequestDTO) (*dto.UserTokenDTO, error)
	Logout(context.Context, string) error
	VKOauthLink(ctx context.Context) (*dto.VKOauthLinkResponse, error)
	VKOAuthCallback(ctx context.Context, request *dto.VKCallbackRequestDTO) (*dto.UserTokenDTO, error)
}
