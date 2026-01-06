package roadmap

import (
	"github.com/google/wire"

	roadmapGrpc "github.com/F0urward/proftwist-backend/services/roadmap/delivery/grpc"
	roadmapHttp "github.com/F0urward/proftwist-backend/services/roadmap/delivery/http"
	roadmapRepository "github.com/F0urward/proftwist-backend/services/roadmap/repository"
	roadmapUsecase "github.com/F0urward/proftwist-backend/services/roadmap/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	chatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	gigachatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	moderationClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/moderationclient"
	roadmapInfoClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
	mongoClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/mongo"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
)

var RoadmapSet = wire.NewSet(
	roadmapRepository.NewRoadmapMongoRepository,
	roadmapRepository.NewRoadmapGigaChatWebapi,
	roadmapUsecase.NewRoadmapUsecase,
	roadmapHttp.NewRoadmapHandlers,
	roadmapHttp.NewRoadmapHttpRegistrar,
	roadmapGrpc.NewRoadmapServer,
	roadmapGrpc.NewRoadmapGrpcRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	mongoClient.NewClient,
	mongoClient.NewDatabase,
	gigachatClient.NewGigaChatClient,
	chatClient.NewChatClient,
	roadmapInfoClient.NewRoadmapInfoClient,
	authClient.NewAuthClient,
	moderationClient.NewModerationClient,
)
