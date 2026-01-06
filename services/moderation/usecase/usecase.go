package usecase

import (
	"context"
	"fmt"

	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"

	"github.com/F0urward/proftwist-backend/services/moderation"
	"github.com/F0urward/proftwist-backend/services/moderation/dto"
)

type ModerationUsecase struct {
	gigachatWebapi moderation.GigachatWebapi
}

func NewModerationUsecase(
	gigichatWebapi moderation.GigachatWebapi,
) moderation.Usecase {
	return &ModerationUsecase{
		gigachatWebapi: gigichatWebapi,
	}
}

func (uc *ModerationUsecase) ModerateContent(ctx context.Context, content string) (*dto.ModerationResult, error) {
	const op = "ModerationUsecase.ModerateContent"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if len(content) == 0 {
		logger.Warn("empty content provided for moderation")
		emptyResult := dto.EmptyModerationResult()
		return &emptyResult, nil
	}

	moderationResult, err := uc.gigachatWebapi.GetModerationResult(ctx, content)
	if err != nil {
		logger.WithError(err).Error("failed to get moderation result")
		return nil, fmt.Errorf("failed to moderate content: %w", err)
	}

	result := dto.ModerationResultToDTO(moderationResult.Allowed, moderationResult.Categories)

	logger.WithFields(map[string]interface{}{
		"allowed":    result.Allowed,
		"categories": result.Categories,
	}).Debug("moderation completed")

	return &result, nil
}
