package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
	"github.com/F0urward/proftwist-backend/services/ai/repository"
)

type AIUsecase struct {
	cfg              *config.Config
	roadmapClient    roadmapclient.RoadmapServiceClient
	roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient
}

func NewAIUsecase(cfg *config.Config, roadmapClient roadmapclient.RoadmapServiceClient, roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient) ai.Usecase {
	return &AIUsecase{
		cfg:              cfg,
		roadmapClient:    roadmapClient,
		roadmapInfoClient: roadmapInfoClient,
	}
}

func (uc *AIUsecase) GenerateRoadmapNodeDescription(ctx context.Context, req dto.GenerateRoadmapNodeDescriptionRequestDTO) (string, error) {
	const op = "AIUsecase.GenerateRoadmapNodeDescription"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	req.NodeID = strings.TrimSpace(req.NodeID)
	req.NodeLabel = strings.TrimSpace(req.NodeLabel)
	req.NodeType = strings.TrimSpace(req.NodeType)
	req.RoadmapID = strings.TrimSpace(req.RoadmapID)
	req.CurrentDescription = strings.TrimSpace(req.CurrentDescription)

	uc.enrichWithRoadmapContext(ctx, &req)

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
			return "", fmt.Errorf("%s: %w", op, err)
		}
		logger.Warn("successfully generated roadmap node description using mock provider")
	} else {
		logger.Info("successfully generated roadmap node description using OpenAI")
	}

	return strings.TrimSpace(description), nil
}

func (uc *AIUsecase) enrichWithRoadmapContext(ctx context.Context, req *dto.GenerateRoadmapNodeDescriptionRequestDTO) {
	const op = "AIUsecase.enrichWithRoadmapContext"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	if req.RoadmapID == "" {
		return
	}

	roadmapResp, err := uc.roadmapClient.GetByIDWithMaterials(ctx, &roadmapclient.GetByIDWithMaterialsRequest{Id: req.RoadmapID})
	if err != nil || roadmapResp == nil || roadmapResp.Roadmap == nil {
		logger.WithError(err).Warn("failed to fetch roadmap context")
		return
	}

	roadmap := roadmapResp.Roadmap

	infoResp, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: req.RoadmapID})
	if err == nil && infoResp != nil && infoResp.RoadmapInfo != nil {
		req.RoadmapName = infoResp.RoadmapInfo.Name
	}

	req.TotalNodeCount = len(roadmap.Nodes)

	incomingEdges := make(map[string]bool)
	for _, edge := range roadmap.Edges {
		incomingEdges[edge.Target] = true
	}

	for _, node := range roadmap.Nodes {
		if !incomingEdges[node.Id] {
			req.RootNodeLabel = node.Data.GetLabel()
			req.RootNodeType = node.Type
			break
		}
	}

	targetNodeFound := false
	targetNodeID := req.NodeID
	for _, node := range roadmap.Nodes {
		if node.Id == targetNodeID {
			targetNodeFound = true
			break
		}
	}
	if !targetNodeFound {
		return
	}

	parentEdges := make(map[string][]string)
	childEdges := make(map[string][]string)
	for _, edge := range roadmap.Edges {
		parentEdges[edge.Target] = append(parentEdges[edge.Target], edge.Source)
		childEdges[edge.Source] = append(childEdges[edge.Source], edge.Target)
	}

	parentIDs := parentEdges[targetNodeID]

	siblingSet := make(map[string]bool)
	for _, parentID := range parentIDs {
		for _, siblingID := range childEdges[parentID] {
			if siblingID != targetNodeID {
				siblingSet[siblingID] = true
			}
		}
	}

	nodeLabelByID := make(map[string]string)
	for _, node := range roadmap.Nodes {
		nodeLabelByID[node.Id] = node.Data.GetLabel()
	}

	for siblingID := range siblingSet {
		if label, ok := nodeLabelByID[siblingID]; ok {
			req.SiblingLabels = append(req.SiblingLabels, label)
		}
	}

	for _, childID := range childEdges[targetNodeID] {
		if label, ok := nodeLabelByID[childID]; ok {
			req.ChildLabels = append(req.ChildLabels, label)
		}
	}

	logger.WithFields(map[string]interface{}{
		"roadmap_name":    req.RoadmapName,
		"total_nodes":     req.TotalNodeCount,
		"sibling_count":   len(req.SiblingLabels),
		"child_count":     len(req.ChildLabels),
		"has_root":        req.RootNodeLabel != "",
	}).Info("enriched node description request with roadmap context")
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