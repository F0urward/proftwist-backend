package roadmap

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type Usecase interface {
	GetAll(ctx context.Context) (*dto.GetAllRoadmapsResponseDTO, error)
	GetByID(ctx context.Context, roadmapID uuid.UUID) (*dto.GetByIDRoadmapResponseDTO, error)
	Create(ctx context.Context, request *dto.CreateRoadmapRequestDTO) error
	Update(ctx context.Context, request *dto.UpdateRoadmapRequestDTO) error
	Delete(ctx context.Context, roadmapID uuid.UUID) error
}
