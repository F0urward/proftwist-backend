package sets

import (
	"github.com/google/wire"

	roadmapInfoHandlers "github.com/F0urward/proftwist-backend/services/auth/delivery/http"
	roadmapInfoRepository "github.com/F0urward/proftwist-backend/services/auth/repository"
	roadmapInfoUsecase "github.com/F0urward/proftwist-backend/services/auth/usecase"
)

var AuthSet = wire.NewSet(
	roadmapInfoRepository.NewAuthRepository,
	roadmapInfoUsecase.NewAuthUsecase,
	roadmapInfoHandlers.NewAuthHandlers,
)
