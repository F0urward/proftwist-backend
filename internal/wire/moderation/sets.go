package moderation

import (
	"github.com/google/wire"

	moderationGrpc "github.com/F0urward/proftwist-backend/services/moderation/delivery/grpc"
	moderationUsecase "github.com/F0urward/proftwist-backend/services/moderation/usecase"
)

var ModerationSet = wire.NewSet(
	moderationUsecase.NewModerationUsecase,
	moderationGrpc.NewModerationServer,
	moderationGrpc.NewModerationGrpcRegistrar,
)
