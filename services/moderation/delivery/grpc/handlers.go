package grpc

import (
	"context"
	"strings"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/moderationclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/moderation"
	"github.com/F0urward/proftwist-backend/services/moderation/dto"
)

type ModerationServer struct {
	uc moderation.Usecase
	moderationclient.UnimplementedModerationServiceServer
}

func NewModerationServer(usecase moderation.Usecase) moderationclient.ModerationServiceServer {
	return &ModerationServer{uc: usecase}
}

func (s *ModerationServer) ModerateContent(ctx context.Context, req *moderationclient.ModerateContentRequest) (*moderationclient.ModerateContentResponse, error) {
	const op = "ModerationServer.ModerateContent"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	logger.WithField("content_length", len(req.Content)).Debug("received moderation request")

	content := strings.TrimSpace(req.Content)
	if content == "" {
		logger.Warn("empty content provided")
		return &moderationclient.ModerateContentResponse{
			Result: &moderationclient.ModerationResult{
				Allowed:    true,
				Categories: []string{},
			},
		}, nil
	}

	result, err := s.uc.ModerateContent(ctx, content)
	if err != nil {
		logger.WithError(err).Error("failed to moderate content")
		return &moderationclient.ModerateContentResponse{
			Error: err.Error(),
		}, nil
	}

	protoResult := convertModerationResultToProto(result)

	logger.WithFields(map[string]interface{}{
		"allowed":    protoResult.Allowed,
		"categories": protoResult.Categories,
	}).Debug("moderation completed")

	return &moderationclient.ModerateContentResponse{
		Result: protoResult,
	}, nil
}

func convertModerationResultToProto(dto *dto.ModerationResult) *moderationclient.ModerationResult {
	if dto == nil {
		return &moderationclient.ModerationResult{
			Allowed:    true,
			Categories: []string{},
		}
	}

	return &moderationclient.ModerationResult{
		Allowed:    dto.Allowed,
		Categories: dto.Categories,
	}
}
