package roadmapinfo

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	GetAll(context.Context) ([]*entities.RoadmapInfo, error)
	GetAllPublicByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*entities.RoadmapInfo, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.RoadmapInfo, error)
	GetByID(context.Context, uuid.UUID) (*entities.RoadmapInfo, error)
	GetByRoadmapID(ctx context.Context, roadmapID string) (*entities.RoadmapInfo, error)
	Create(ctx context.Context, roadmap *entities.RoadmapInfo) (*entities.RoadmapInfo, error)
	Update(context.Context, *entities.RoadmapInfo) error
	Delete(context.Context, uuid.UUID) error
}
