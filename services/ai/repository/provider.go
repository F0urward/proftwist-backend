package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type ErrorProvider struct {
	err error
}

func NewProvider(cfg *config.Config, providerOverride, modelOverride string) ai.Provider {
	provider := strings.ToLower(strings.TrimSpace(providerOverride))
	if provider == "" {
		provider = strings.ToLower(strings.TrimSpace(cfg.AI.Provider))
	}
	if provider == "" {
		provider = "ollama"
	}

	model := modelOverride
	if model == "" {
		switch provider {
		case "ollama", "local":
			model = cfg.AI.Ollama.Model
		case "openai", "openai-compatible", "compatible":
			model = cfg.AI.OpenAI.Model
		}
	}

	switch provider {
	case "ollama", "local":
		return NewOllamaProviderWithCredentials(cfg.AI.Ollama.BaseURL, cfg.AI.Ollama.APIKey, model)
	case "openai", "openai-compatible", "compatible":
		return NewOpenAICompatibleProviderWithCredentials(cfg.AI.OpenAI.BaseURL, cfg.AI.OpenAI.APIKey, model)
	case "mock":
		return NewMockProvider()
	default:
		return &ErrorProvider{err: fmt.Errorf("unsupported AI provider %q", provider)}
	}
}

func (p *ErrorProvider) GenerateRoadmapNodeDescription(_ context.Context, _ dto.GenerateRoadmapNodeDescriptionRequestDTO) (string, error) {
	return "", p.err
}

func (p *ErrorProvider) GenerateRoadmap(_ context.Context, _ dto.GenerateRoadmapRequestDTO) (string, error) {
	return "", p.err
}
