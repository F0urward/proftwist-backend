package category

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/services/category/dto"
)

type Usecase interface {
	GetAll(ctx context.Context) (*dto.GetAllCategoriesResponse, error)
	GetByID(ctx context.Context, categoryID uuid.UUID) (*dto.CategoryResponse, error)
	Create(ctx context.Context, req *dto.CreateCategoryRequestDTO) (*dto.CreateCategoryResponse, error)
	Update(ctx context.Context, categoryID uuid.UUID, req *dto.UpdateCategoryRequestDTO) error
	Delete(ctx context.Context, categoryID uuid.UUID) error
}
