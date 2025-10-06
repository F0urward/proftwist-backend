package auth

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}
