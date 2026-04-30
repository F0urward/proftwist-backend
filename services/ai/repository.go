package ai

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type Provider interface {
	GenerateRoadmapNodeDescription(ctx context.Context, req dto.GenerateRoadmapNodeDescriptionRequestDTO) (string, error)
}
