package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type AIUsecase struct {
	provider ai.Provider
}

func NewAIUsecase(provider ai.Provider) ai.Usecase {
	return &AIUsecase{
		provider: provider,
	}
}

func (uc *AIUsecase) GenerateRoadmapNodeDescription(ctx context.Context, req dto.GenerateRoadmapNodeDescriptionRequestDTO) (*dto.GenerateRoadmapNodeDescriptionResponseDTO, error) {
	const op = "AIUsecase.GenerateRoadmapNodeDescription"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	req.NodeID = strings.TrimSpace(req.NodeID)
	req.NodeLabel = strings.TrimSpace(req.NodeLabel)
	req.NodeType = strings.TrimSpace(req.NodeType)
	req.RoadmapID = strings.TrimSpace(req.RoadmapID)
	req.CurrentDescription = strings.TrimSpace(req.CurrentDescription)

	description, err := uc.provider.GenerateRoadmapNodeDescription(ctx, req)
	if err != nil {
		logger.WithError(err).Error("failed to generate roadmap node description")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.GenerateRoadmapNodeDescriptionResponseDTO{
		Description: strings.TrimSpace(description),
	}, nil
}
