package roadmapinfo

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
)

type Usecase interface {
	GetAll(ctx context.Context) (*dto.GetAllRoadmapsInfoResponseDTO, error)
	GetByID(ctx context.Context, roadmapID uuid.UUID) (*dto.GetByIDRoadmapInfoResponseDTO, error)
	Create(ctx context.Context, request *dto.CreateRoadmapInfoRequestDTO) error
	Update(ctx context.Context, roadmapID uuid.UUID, request *dto.UpdateRoadmapInfoRequestDTO) error
	Delete(ctx context.Context, roadmapID uuid.UUID) error
}
