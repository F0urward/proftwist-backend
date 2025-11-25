package category

import (
	"github.com/google/wire"

	categoryHandlers "github.com/F0urward/proftwist-backend/services/category/delivery/http"
	categoryRepository "github.com/F0urward/proftwist-backend/services/category/repository"
	categoryUsecase "github.com/F0urward/proftwist-backend/services/category/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
	redisClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/redis"
)

var CategorySet = wire.NewSet(
	categoryRepository.NewCategoryPostgresRepository,
	categoryUsecase.NewCategoryUsecase,
	categoryHandlers.NewCategoryHandlers,
	categoryHandlers.NewCategoryHttpRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	redisClient.NewClient,
	authClient.NewAuthClient,
)
