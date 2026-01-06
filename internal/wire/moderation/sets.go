package moderation

import (
	"github.com/google/wire"

	moderationGrpc "github.com/F0urward/proftwist-backend/services/moderation/delivery/grpc"
	moderationRepository "github.com/F0urward/proftwist-backend/services/moderation/repository"
	moderationUsecase "github.com/F0urward/proftwist-backend/services/moderation/usecase"

	gigachatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
)

var ModerationSet = wire.NewSet(
	moderationRepository.NewModerationGigaChatWebapi,
	moderationUsecase.NewModerationUsecase,
	moderationGrpc.NewModerationServer,
	moderationGrpc.NewModerationGrpcRegistrar,
)

var ClientsSet = wire.NewSet(
	gigachatClient.NewGigaChatClient,
)
