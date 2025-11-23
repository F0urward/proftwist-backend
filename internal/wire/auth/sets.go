package auth

import (
	"github.com/google/wire"

	authGrpc "github.com/F0urward/proftwist-backend/services/auth/delivery/grpc"
	authHandlers "github.com/F0urward/proftwist-backend/services/auth/delivery/http"
	authRepository "github.com/F0urward/proftwist-backend/services/auth/repository"
	authUsecase "github.com/F0urward/proftwist-backend/services/auth/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	vkClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/vkclient"
	awsClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/aws"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
	redisClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/redis"
)

var AuthSet = wire.NewSet(
	authRepository.NewAuthPostgresRepository,
	authRepository.NewAuthRedisRepository,
	authRepository.NewVKAuthWebapi,
	authRepository.NewAuthAWSRepository,
	authUsecase.NewAuthUsecase,
	authHandlers.NewAuthHandlers,
	authHandlers.NewAuthHttpRegistrar,
	authGrpc.NewAuthServer,
	authGrpc.NewAuthGrpcRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	redisClient.NewClient,
	awsClient.NewClient,
	vkClient.NewVKClient,
	authClient.NewAuthClient,
)
