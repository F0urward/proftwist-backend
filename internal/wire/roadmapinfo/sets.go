package roadmapinfo

import (
	"github.com/google/wire"

	roadmapInfoGrpc "github.com/F0urward/proftwist-backend/services/roadmapinfo/delivery/grpc"
	roadmapInfoHandlers "github.com/F0urward/proftwist-backend/services/roadmapinfo/delivery/http"
	roadmapInfoRepository "github.com/F0urward/proftwist-backend/services/roadmapinfo/repository"
	roadmapInfoUsecase "github.com/F0urward/proftwist-backend/services/roadmapinfo/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	moderationClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/moderationclient"
	roadmapClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
)

var RoadmapInfoSet = wire.NewSet(
	roadmapInfoRepository.NewRoadmapInfoPostgresRepository,
	roadmapInfoUsecase.NewRoadmapInfoUsecase,
	roadmapInfoHandlers.NewRoadmapInfoHandlers,
	roadmapInfoHandlers.NewRoadmapInfoHttpRegistrar,
	roadmapInfoGrpc.NewRoadmapInfoServer,
	roadmapInfoGrpc.NewRoadmapInfoGrpcRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	roadmapClient.NewRoadmapClient,
	authClient.NewAuthClient,
	moderationClient.NewModerationClient,
)
