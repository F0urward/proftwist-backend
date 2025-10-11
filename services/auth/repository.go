package auth

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type PostgresRepository interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}

type RedisRepository interface {
	AddToBlacklist(ctx context.Context, userID, token string) error
	IsInBlacklist(ctx context.Context, userID, token string) (bool, error)
}
