package sets

import (
	"github.com/google/wire"

	categoryHandlers "github.com/F0urward/proftwist-backend/services/category/delivery/http"
	categoryRepository "github.com/F0urward/proftwist-backend/services/category/repository"
	categoryUsecase "github.com/F0urward/proftwist-backend/services/category/usecase"
)

var CategorySet = wire.NewSet(
	categoryRepository.NewCategoryRepository,
	categoryUsecase.NewCategoryUsecase,
	categoryHandlers.NewCategoryHandlers,
)
