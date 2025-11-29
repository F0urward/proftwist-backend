package gigachatclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/baseclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient/dto"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/google/uuid"
)

const (
	AuthUrl    = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	BaseUrl    = "https://gigachat.devices.sberbank.ru/api/v1/"
	ModelsPath = "models"
	ChatPath   = "chat/completions"
	FilesPath  = "files"
)

type Client struct {
	BaseClient *baseclient.BaseClient
	Token      *Token
	cfg        *config.Config
}

func NewGigaChatClient(cfg *config.Config) *Client {
	const op = "gigachat.NewGigaChatClient"
	logger := logctx.GetLogger(context.Background()).WithField("op", op)

	gigachatClient := &Client{
		BaseClient: baseclient.NewInsecureBaseClient(),
		Token:      new(Token),
		cfg:        cfg,
	}

	if err := gigachatClient.Auth(context.Background()); err != nil {
		logger.WithError(err).Error("failed to authenticate during client creation")
	}

	return gigachatClient
}

func (c *Client) Auth(ctx context.Context) error {
	const op = "GigaChatClient.Auth"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if c.Token.Active() {
		logger.Info("token is still active, skipping auth")
		return nil
	}

	logger.Info("starting authentication")

	payload := strings.NewReader("scope=" + c.cfg.GigaChat.Scope)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, AuthUrl, payload)
	if err != nil {
		logger.WithError(err).Error("failed to create auth request")
		return fmt.Errorf("%s: failed to create auth request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("RqUID", uuid.NewString())
	// authKey := base64.StdEncoding.EncodeToString([]byte(c.cfg.GigaChat.ClientId + ":" + c.cfg.GigaChat.SecretKey))
	req.Header.Set("Authorization", "Basic "+c.cfg.GigaChat.AuthKey)

	var oauth dto.OAuthResponse
	if err := c.BaseClient.DoRequest(ctx, req, &oauth); err != nil {
		logger.WithError(err).Error("failed to execute auth request")
		return fmt.Errorf("%s: failed to execute auth request: %w", op, err)
	}

	c.Token.Set(oauth.AccessToken, time.UnixMilli(oauth.ExpiresAt))
	logger.Info("successfully authenticated")
	return nil
}

func (c *Client) Chat(ctx context.Context, in *dto.ChatRequest) (*dto.ChatResponse, error) {
	const op = "GigaChatClient.ChatWithContext"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	logger.Info("sending chat request")

	reqBytes, err := json.Marshal(in)
	if err != nil {
		logger.WithError(err).Error("failed to marshal chat request")
		return nil, fmt.Errorf("%s: failed to marshal chat request: %w", op, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, BaseUrl+ChatPath, bytes.NewReader(reqBytes))
	if err != nil {
		logger.WithError(err).Error("failed to create chat request")
		return nil, fmt.Errorf("%s: failed to create chat request: %w", op, err)
	}

	var chatResponse dto.ChatResponse
	if err := c.sendRequest(ctx, req, &chatResponse); err != nil {
		logger.WithError(err).Error("failed to send chat request")
		return nil, fmt.Errorf("%s: failed to send chat request: %w", op, err)
	}

	logger.Info("successfully received chat response")
	return &chatResponse, nil
}

func (c *Client) sendRequest(ctx context.Context, req *http.Request, result interface{}) error {
	const op = "GigaChatClient.sendRequest"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if !c.Token.Active() {
		logger.Info("token expired, refreshing")
		if err := c.Auth(ctx); err != nil {
			logger.WithError(err).Error("failed to refresh token")
			return fmt.Errorf("%s: failed to refresh token: %w", op, err)
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.Get()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if err := c.BaseClient.DoRequest(ctx, req, result); err != nil {
		logger.WithError(err).Error("request failed")
		return fmt.Errorf("%s: request failed: %w", op, err)
	}

	logger.Info("successfully processed GigaChat request")
	return nil
}
