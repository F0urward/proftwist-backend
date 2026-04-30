package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type ErrorProvider struct {
	err error
}

func NewProvider(cfg *config.Config, gigachatClient *gigachatclient.Client) ai.Provider {
	provider := strings.ToLower(strings.TrimSpace(cfg.AI.Provider))
	if provider == "" {
		provider = "gigachat"
	}

	switch provider {
	case "gigachat":
		return NewGigaChatProvider(gigachatClient)
	case "openai", "openai-compatible", "compatible":
		return NewOpenAICompatibleProvider(cfg)
	case "mock":
		return NewMockProvider()
	default:
		return &ErrorProvider{err: fmt.Errorf("unsupported AI provider %q", cfg.AI.Provider)}
	}
}

func (p *ErrorProvider) GenerateRoadmapNodeDescription(_ context.Context, _ dto.GenerateRoadmapNodeDescriptionRequestDTO) (string, error) {
	return "", p.err
}
