package material

import (
	"github.com/google/wire"

	materialHandlers "github.com/F0urward/proftwist-backend/services/material/delivery/http"
	materialRepository "github.com/F0urward/proftwist-backend/services/material/repository"
	materialUsecase "github.com/F0urward/proftwist-backend/services/material/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
)

var MaterialSet = wire.NewSet(
	materialRepository.NewMaterialPostgresRepository,
	materialUsecase.NewMaterialUsecase,
	materialHandlers.NewMaterialHandlers,
	materialHandlers.NewMaterialHttpRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	authClient.NewAuthClient,
)
