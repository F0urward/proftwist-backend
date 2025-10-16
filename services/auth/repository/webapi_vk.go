package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/vkclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/auth"
)

type AuthVKWebapi struct {
	client *vkclient.VKClient
}

func NewVKAuthWebapi(client *vkclient.VKClient) auth.VKWebapi {
	return &AuthVKWebapi{client: client}
}

func (r *AuthVKWebapi) ExchangeCodeForTokens(ctx context.Context, code, codeVerifier, deviceID string) (*entities.VKTokens, error) {
	const op = "VKWebapi.ExchangeCodeForTokens"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	tokens, err := r.client.ExchangeCodeForTokens(ctx, code, codeVerifier, deviceID)
	if err != nil {
		logger.WithError(err).Error("failed to exchange code for tokens")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokenResponse := &entities.VKTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IDToken:      tokens.IDToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
		UserID:       tokens.UserID,
		State:        tokens.State,
		Scope:        tokens.Scope,
		ExpiresAt:    time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second),
	}

	logger.WithField("user_id", tokens.UserID).Info("successfully exchanged code for tokens")
	return tokenResponse, nil
}

func (r *AuthVKWebapi) GetUserInfo(ctx context.Context, accessToken string) (*entities.VKUserInfo, error) {
	const op = "VKWebapi.GetUserInfo"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userInfo, err := r.client.GetUserInfo(ctx, accessToken)
	if err != nil {
		logger.WithError(err).Error("failed to get user info from vk")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userInfoEntity := &entities.VKUserInfo{
		VKUserID:  userInfo.User.UserID,
		Email:     userInfo.User.Email,
		FirstName: userInfo.User.FirstName,
		LastName:  userInfo.User.LastName,
		Avatar:    userInfo.User.Avatar,
	}

	logger.WithFields(map[string]interface{}{
		"vk_user_id": userInfoEntity.VKUserID,
		"email":      userInfoEntity.Email,
	}).Info("successfully retrieved user info from vk")

	return userInfoEntity, nil
}
