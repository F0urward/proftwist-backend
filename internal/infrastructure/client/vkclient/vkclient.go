package vkclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/baseclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/vkclient/dto"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

const (
	tokenURL = "https://id.vk.com/oauth2/auth"
	infoURL  = "https://id.vk.com/oauth2/user_info"
)

type VKClient struct {
	BaseClient *baseclient.BaseClient
	cfg        *config.Config
}

func NewVKClient(cfg *config.Config) *VKClient {
	return &VKClient{
		BaseClient: baseclient.NewBaseClient(),
		cfg:        cfg,
	}
}

func (c *VKClient) ExchangeCodeForTokens(ctx context.Context, code, codeVerifier, deviceID string) (*dto.VKTokenExchangeResponse, error) {
	const op = "VKClient.ExchangeCodeForTokens"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	logger.WithFields(map[string]interface{}{
		"client_id": c.cfg.Auth.VK.IntegrationID,
		"device_id": deviceID,
	}).Info("exchanging code for tokens")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code_verifier", codeVerifier)
	data.Set("redirect_uri", c.cfg.Auth.VK.RedirectURL)
	data.Set("code", code)
	data.Set("client_id", c.cfg.Auth.VK.IntegrationID)
	data.Set("device_id", deviceID)
	data.Set("client_secret", c.cfg.Auth.VK.SecretKey)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		logger.WithError(err).Error("failed to create request")
		return nil, fmt.Errorf("%s: failed to create request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response dto.VKTokenExchangeResponse
	if err := c.BaseClient.DoRequest(ctx, req, &response); err != nil {
		logger.WithError(err).Error("failed to exchange code for tokens")
		return nil, fmt.Errorf("%s: failed to exchange code for tokens: %w", op, err)
	}

	logger.WithField("user_id", response.UserID).Info("successfully exchanged code for tokens")
	return &response, nil
}

func (c *VKClient) GetUserInfo(ctx context.Context, accessToken string) (*dto.VKUserInfoResponse, error) {
	const op = "VKClient.GetUserInfo"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	logger.Info("getting user info from VK")

	data := url.Values{}
	data.Set("client_id", c.cfg.Auth.VK.IntegrationID)
	data.Set("access_token", accessToken)

	req, err := http.NewRequestWithContext(ctx, "POST", infoURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		logger.WithError(err).Error("failed to create request")
		return nil, fmt.Errorf("%s: failed to create request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response dto.VKUserInfoResponse
	if err := c.BaseClient.DoRequest(ctx, req, &response); err != nil {
		logger.WithError(err).Error("failed to get user info")
		return nil, fmt.Errorf("%s: failed to get user info: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"vk_user_id": response.User.UserID,
		"email":      response.User.Email,
	}).Info("successfully retrieved user info")

	return &response, nil
}
