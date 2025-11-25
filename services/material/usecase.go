package material

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/services/material/dto"
)

type Usecase interface {
	CreateMaterial(ctx context.Context, userID uuid.UUID, req dto.CreateMaterialRequestDTO) (*dto.MaterialResponseDTO, error)
	DeleteMaterial(ctx context.Context, materialID uuid.UUID, userID uuid.UUID) error
	GetMaterialsByNode(ctx context.Context, nodeID string) (*dto.MaterialListResponseDTO, error)
	GetUserMaterials(ctx context.Context, userID uuid.UUID) (*dto.MaterialListResponseDTO, error)
}
