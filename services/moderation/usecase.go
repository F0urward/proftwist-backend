package moderation

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/moderation/dto"
)

type Usecase interface {
	ModerateContent(ctx context.Context, content string) (*dto.ModerationResult, error)
}
