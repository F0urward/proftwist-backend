package ai

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type Usecase interface {
	GenerateRoadmapNodeDescription(ctx context.Context, req dto.GenerateRoadmapNodeDescriptionRequestDTO) (*dto.GenerateRoadmapNodeDescriptionResponseDTO, error)
	GenerateRoadmap(ctx context.Context, req dto.GenerateRoadmapRequestDTO) (string, error)
}
