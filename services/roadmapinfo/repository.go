package roadmapinfo

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*entities.RoadmapInfo, error)
	GetByID(ctx context.Context, roadmapID uuid.UUID) (*entities.RoadmapInfo, error)
	Create(ctx context.Context, roadmap *entities.RoadmapInfo) error
	Update(ctx context.Context, roadmap *entities.RoadmapInfo) error
	Delete(ctx context.Context, roadmapID uuid.UUID) error
}
