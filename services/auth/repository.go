package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type PostgresRepository interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	CreateVKUser(ctx context.Context, vkUser *entities.VKUser) error
	GetVKUserByUserID(ctx context.Context, userID uuid.UUID) (*entities.VKUser, error)
	GetVKUserByID(ctx context.Context, vkUserID int64) (*entities.VKUser, error)
	UpdateVKUser(ctx context.Context, vkUser *entities.VKUser) error
	DeleteVKUser(ctx context.Context, userID uuid.UUID) error
}

type RedisRepository interface {
	AddToBlacklist(ctx context.Context, userID, token string) error
	IsInBlacklist(ctx context.Context, userID, token string) (bool, error)

	StoreState(ctx context.Context, state, codeVerifier string) error
	GetCodeVerifierByState(ctx context.Context, state string) (string, error)
	DeleteState(ctx context.Context, state string) error
}

type VKWebapi interface {
	ExchangeCodeForTokens(ctx context.Context, code, codeVerifier, deviceID string) (*entities.VKTokens, error)
	GetUserInfo(ctx context.Context, accessToken string) (*entities.VKUserInfo, error)
}

type AWSRepository interface {
	PutObject(ctx context.Context, input entities.UploadInput) (*minio.UploadInfo, error)
	GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucket string, fileName string) error
}
