package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/redis/go-redis/v9"
)

const (
	userTokensPrefix = "user_id:"
	oauthStatePrefix = "oauth_state:"
)

type AuthRedisRepository struct {
	client *redis.Client
	cfg    *config.Config
}

func NewAuthRedisRepository(client *redis.Client, cfg *config.Config) auth.RedisRepository {
	return &AuthRedisRepository{
		client: client,
		cfg:    cfg,
	}
}

func (r *AuthRedisRepository) AddToBlacklist(ctx context.Context, userID, token string) error {
	expiration := time.Until(time.Now().Add(time.Duration(r.cfg.Auth.Jwt.Expire) * time.Second))
	userKey := fmt.Sprintf("%s%s", userTokensPrefix, userID)

	if err := r.client.SAdd(ctx, userKey, token).Err(); err != nil {
		return fmt.Errorf("failed to add token to user's blacklist: %w", err)
	}

	if err := r.client.Expire(ctx, userKey, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration for user's blacklist: %w", err)
	}

	return nil
}

func (r *AuthRedisRepository) IsInBlacklist(ctx context.Context, userID, token string) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}

	userKey := fmt.Sprintf("%s%s", userTokensPrefix, userID)
	isMember, err := r.client.SIsMember(ctx, userKey, token).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token in blacklist: %w", err)
	}
	return isMember, nil
}

func (r *AuthRedisRepository) StoreState(ctx context.Context, state, codeVerifier string) error {
	const op = "AuthRedisRepository.StoreOAuthState"

	stateKey := fmt.Sprintf("%s%s", oauthStatePrefix, state)
	expiration := 10 * time.Minute

	err := r.client.Set(ctx, stateKey, codeVerifier, expiration).Err()
	if err != nil {
		return fmt.Errorf("%s: failed to store oauth state: %w", op, err)
	}

	return nil
}

func (r *AuthRedisRepository) GetCodeVerifierByState(ctx context.Context, state string) (string, error) {
	const op = "AuthRedisRepository.GetCodeVerifierByState"

	stateKey := fmt.Sprintf("%s%s", oauthStatePrefix, state)

	codeVerifier, err := r.client.Get(ctx, stateKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("%s: oauth state not found or expired", op)
	}
	if err != nil {
		return "", fmt.Errorf("%s: failed to get code verifier by state: %w", op, err)
	}

	return codeVerifier, nil
}

func (r *AuthRedisRepository) DeleteState(ctx context.Context, state string) error {
	const op = "AuthRedisRepository.DeleteOAuthState"

	stateKey := fmt.Sprintf("%s%s", oauthStatePrefix, state)

	err := r.client.Del(ctx, stateKey).Err()
	if err != nil {
		return fmt.Errorf("%s: failed to delete oauth state: %w", op, err)
	}

	return nil
}
