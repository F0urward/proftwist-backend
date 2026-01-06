package moderation

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

type GigachatWebapi interface {
	GetModerationResult(ctx context.Context, content string) (*entities.ModerationResult, error)
}
