package category

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*entities.Category, error)
	GetByID(ctx context.Context, categoryID uuid.UUID) (*entities.Category, error)
	GetByName(ctx context.Context, name string) (*entities.Category, error)
	Create(ctx context.Context, category *entities.Category) (*entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, categoryID uuid.UUID) error
}
