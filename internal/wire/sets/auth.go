package sets

import (
	"github.com/google/wire"

	authGrpc "github.com/F0urward/proftwist-backend/services/auth/delivery/grpc"
	authHandlers "github.com/F0urward/proftwist-backend/services/auth/delivery/http"
	authRepository "github.com/F0urward/proftwist-backend/services/auth/repository"
	authUsecase "github.com/F0urward/proftwist-backend/services/auth/usecase"
)

var AuthSet = wire.NewSet(
	authRepository.NewAuthPostgresRepository,
	authRepository.NewAuthRedisRepository,
	authRepository.NewVKAuthWebapi,
	authRepository.NewAuthAWSRepository,
	authUsecase.NewAuthUsecase,
	authHandlers.NewAuthHandlers,
	authGrpc.NewAuthServer,
)
