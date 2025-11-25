package material

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type Repository interface {
	CreateMaterial(ctx context.Context, material *entities.Material) (*entities.Material, error)
	GetMaterialByID(ctx context.Context, materialID uuid.UUID) (*entities.Material, error)
	GetMaterialsByNode(ctx context.Context, nodeID string) ([]*entities.Material, error)
	GetMaterialsByAuthor(ctx context.Context, authorID uuid.UUID) ([]*entities.Material, error)
	DeleteMaterial(ctx context.Context, materialID uuid.UUID) error
}
