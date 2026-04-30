package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	gigachatClientDTO "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient/dto"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type GigaChatProvider struct {
	client *gigachatclient.Client
}

func NewGigaChatProvider(client *gigachatclient.Client) ai.Provider {
	return &GigaChatProvider{client: client}
}

func (p *GigaChatProvider) GenerateRoadmapNodeDescription(ctx context.Context, req dto.GenerateRoadmapNodeDescriptionRequestDTO) (string, error) {
	const op = "GigaChatProvider.GenerateRoadmapNodeDescription"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	logger.WithFields(map[string]interface{}{
		"roadmap_id": req.RoadmapID,
		"node_id":    req.NodeID,
		"node_label": req.NodeLabel,
		"node_type":  req.NodeType,
	}).Info("generating roadmap node description")

	chatReq := &gigachatClientDTO.ChatRequest{
		Model: "GigaChat",
		Messages: []gigachatClientDTO.Message{
			{
				Role:    "system",
				Content: "Ты помогаешь создавать образовательные roadmap. Пиши только готовое описание узла на русском языке, без Markdown, списков, кавычек и пояснений. Описание должно быть коротким, практичным и полезным: 1-3 предложения.",
			},
			{
				Role:    "user",
				Content: buildRoadmapNodeDescriptionPrompt(req),
			},
		},
		Temperature:       float64Ptr(0.4),
		MaxTokens:         int64Ptr(350),
		RepetitionPenalty: float64Ptr(1.05),
	}

	chatResp, err := p.client.Chat(ctx, chatReq)
	if err != nil {
		logger.WithError(err).Error("failed to get response from GigaChat")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if len(chatResp.Choices) == 0 {
		logger.Error("empty response from GigaChat")
		return "", fmt.Errorf("%s: empty response from GigaChat", op)
	}

	description := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	if description == "" {
		logger.Error("empty description from GigaChat")
		return "", fmt.Errorf("%s: empty description from GigaChat", op)
	}

	logger.WithField("description_length", len(description)).Info("successfully generated roadmap node description")
	return description, nil
}

func buildRoadmapNodeDescriptionPrompt(req dto.GenerateRoadmapNodeDescriptionRequestDTO) string {
	var b strings.Builder

	b.WriteString("Сгенерируй описание для узла roadmap.\n")
	if req.RoadmapID != "" {
		b.WriteString("Roadmap ID: ")
		b.WriteString(req.RoadmapID)
		b.WriteString("\n")
	}
	b.WriteString("Node ID: ")
	b.WriteString(req.NodeID)
	b.WriteString("\n")
	b.WriteString("Название узла: ")
	b.WriteString(req.NodeLabel)
	b.WriteString("\n")
	if req.NodeType != "" {
		b.WriteString("Тип узла: ")
		b.WriteString(req.NodeType)
		b.WriteString("\n")
	}
	if req.CurrentDescription != "" {
		b.WriteString("Текущее описание: ")
		b.WriteString(req.CurrentDescription)
		b.WriteString("\n")
	}
	b.WriteString("Верни улучшенное или новое описание, которое объясняет, что изучить и на чем потренироваться.")

	return b.String()
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
