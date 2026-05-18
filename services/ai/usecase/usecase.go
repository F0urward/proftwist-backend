package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
	"github.com/F0urward/proftwist-backend/services/ai/repository"
)

type AIUsecase struct {
	cfg *config.Config
}

func NewAIUsecase(cfg *config.Config) ai.Usecase {
	return &AIUsecase{
		cfg: cfg,
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

	openaiProvider := repository.NewOpenAICompatibleProviderWithCredentials(
		uc.cfg.AI.OpenAI.BaseURL,
		uc.cfg.AI.OpenAI.APIKey,
		uc.cfg.AI.OpenAI.Model,
	)
	logger.Info("trying OpenAI provider for roadmap node description generation")
	description, err := openaiProvider.GenerateRoadmapNodeDescription(ctx, req)
	if err != nil {
		logger.WithError(err).Warn("OpenAI failed for roadmap node description, trying Ollama")
	}

	if err != nil {
		ollamaProvider := repository.NewOllamaProviderWithCredentials(
			uc.cfg.AI.Ollama.BaseURL,
			uc.cfg.AI.Ollama.APIKey,
			uc.cfg.AI.Ollama.Model,
		)
		logger.Info("trying Ollama provider for roadmap node description generation")
		description, err = ollamaProvider.GenerateRoadmapNodeDescription(ctx, req)
		if err != nil {
			logger.WithError(err).Warn("Ollama failed for roadmap node description, using mock")
		}
	}

	if err != nil {
		logger.Info("falling back to mock provider for roadmap node description")
		mockProvider := repository.NewMockProvider()
		description, err = mockProvider.GenerateRoadmapNodeDescription(ctx, req)
		if err != nil {
			logger.WithError(err).Error("mock provider also failed for roadmap node description")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		logger.Warn("successfully generated roadmap node description using mock provider")
	} else {
		logger.Info("successfully generated roadmap node description using OpenAI")
	}

	return &dto.GenerateRoadmapNodeDescriptionResponseDTO{
		Description: strings.TrimSpace(description),
	}, nil
}

func (uc *AIUsecase) GenerateRoadmap(ctx context.Context, req dto.GenerateRoadmapRequestDTO) (string, error) {
	const op = "AIUsecase.GenerateRoadmap"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	req.RoadmapID = strings.TrimSpace(req.RoadmapID)
	req.Prompt = strings.TrimSpace(req.Prompt)

	openaiProvider := repository.NewOpenAICompatibleProviderWithCredentials(
		uc.cfg.AI.OpenAI.BaseURL,
		uc.cfg.AI.OpenAI.APIKey,
		uc.cfg.AI.OpenAI.Model,
	)
	logger.Info("trying OpenAI provider for roadmap generation")
	response, err := openaiProvider.GenerateRoadmap(ctx, req)
	if err != nil {
		logger.WithError(err).Warn("OpenAI failed for roadmap generation, trying Ollama")
	}

	if err != nil {
		ollamaProvider := repository.NewOllamaProviderWithCredentials(
			uc.cfg.AI.Ollama.BaseURL,
			uc.cfg.AI.Ollama.APIKey,
			uc.cfg.AI.Ollama.Model,
		)
		logger.Info("trying Ollama provider for roadmap generation")
		response, err = ollamaProvider.GenerateRoadmap(ctx, req)
		if err != nil {
			logger.WithError(err).Warn("Ollama failed for roadmap generation, using mock")
		}
	}

	if err != nil {
		logger.Info("falling back to mock provider for roadmap generation")
		mockProvider := repository.NewMockProvider()
		response, err = mockProvider.GenerateRoadmap(ctx, req)
		if err != nil {
			logger.WithError(err).Error("mock provider also failed for roadmap generation")
			return "", fmt.Errorf("%s: %w", op, err)
		}
		logger.Warn("successfully generated roadmap using mock provider")
	} else {
		logger.Info("successfully generated roadmap using OpenAI")
	}

	return response, nil
}