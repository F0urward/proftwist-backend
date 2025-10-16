package sets

import (
	"github.com/google/wire"

	authHandlers "github.com/F0urward/proftwist-backend/services/auth/delivery/http"
	authRepository "github.com/F0urward/proftwist-backend/services/auth/repository"
	authUsecase "github.com/F0urward/proftwist-backend/services/auth/usecase"
)

var AuthSet = wire.NewSet(
	authRepository.NewAuthPostgresRepository,
	authRepository.NewAuthRedisRepository,
	authRepository.NewVKAuthWebapi,
	authUsecase.NewAuthUsecase,
	authHandlers.NewAuthHandlers,
)
