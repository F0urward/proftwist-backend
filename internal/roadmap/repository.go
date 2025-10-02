package roadmap

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*entities.Roadmap, error)
	GetByID(ctx context.Context, roadmapID uuid.UUID) (*entities.Roadmap, error)
	Create(ctx context.Context, roadmap *entities.Roadmap) error
	Update(ctx context.Context, roadmap *entities.Roadmap) error
	Delete(ctx context.Context, roadmapID uuid.UUID) error
}
