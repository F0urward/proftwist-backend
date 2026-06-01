package usecase

import (
	"context"

	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/moderation"
	"github.com/F0urward/proftwist-backend/services/moderation/dto"
)

type ModerationUsecase struct{}

func NewModerationUsecase() moderation.Usecase {
	return &ModerationUsecase{}
}

func (uc *ModerationUsecase) ModerateContent(ctx context.Context, content string) (*dto.ModerationResult, error) {
	const op = "ModerationUsecase.ModerateContent"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	if len(content) == 0 {
		logger.Warn("empty content provided for moderation")
		emptyResult := dto.EmptyModerationResult()
		return &emptyResult, nil
	}

	result := dto.EmptyModerationResult()

	logger.WithField("allowed", result.Allowed).Debug("moderation passed (gigachat support removed)")

	return &result, nil
}
