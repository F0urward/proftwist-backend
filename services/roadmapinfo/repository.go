package roadmapinfo

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	GetAllPublic(ctx context.Context) ([]*entities.RoadmapInfo, error)
	GetAllPublicByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*entities.RoadmapInfo, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.RoadmapInfo, error)
	GetByID(context.Context, uuid.UUID) (*entities.RoadmapInfo, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.RoadmapInfo, error)
	GetByRoadmapID(ctx context.Context, roadmapID string) (*entities.RoadmapInfo, error)
	Create(ctx context.Context, roadmap *entities.RoadmapInfo) (*entities.RoadmapInfo, error)
	Update(context.Context, *entities.RoadmapInfo) error
	Delete(context.Context, uuid.UUID) error
	CreateSubscription(ctx context.Context, userID, roadmapInfoID uuid.UUID) error
	DeleteSubscription(ctx context.Context, userID, roadmapInfoID uuid.UUID) error
	SubscriptionExists(ctx context.Context, userID, roadmapInfoID uuid.UUID) (bool, error)
	GetSubscribedRoadmapIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	SearchPublic(ctx context.Context, query string, categoryID *uuid.UUID) ([]*entities.RoadmapInfo, error)
}
