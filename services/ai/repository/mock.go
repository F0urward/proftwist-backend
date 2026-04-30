package repository

import (
	"context"
	"fmt"

	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type MockProvider struct{}

func NewMockProvider() ai.Provider {
	return &MockProvider{}
}

func (p *MockProvider) GenerateRoadmapNodeDescription(_ context.Context, req dto.GenerateRoadmapNodeDescriptionRequestDTO) (string, error) {
	return fmt.Sprintf("Изучите тему %s: разберите ключевые понятия, типичные сценарии применения и закрепите материал на небольшом практическом задании.", req.NodeLabel), nil
}
