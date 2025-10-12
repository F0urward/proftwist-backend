package roadmapinfo

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
)

type Usecase interface {
	GetAll(context.Context) (*dto.GetAllRoadmapsInfoResponseDTO, error)
	GetByID(context.Context, uuid.UUID) (*dto.GetByIDRoadmapInfoResponseDTO, error)
	GetByRoadmapID(ctx context.Context, roadmapID string) (*dto.GetByIDRoadmapInfoResponseDTO, error)
	Create(context.Context, *dto.CreateRoadmapInfoRequestDTO) error
	Update(context.Context, uuid.UUID, *dto.UpdateRoadmapInfoRequestDTO) error
	Delete(context.Context, uuid.UUID) error
}
