package roadmap

import (
	"github.com/google/wire"

	aiHttp "github.com/F0urward/proftwist-backend/services/ai/delivery/http"
	aiUsecase "github.com/F0urward/proftwist-backend/services/ai/usecase"
	roadmapGrpc "github.com/F0urward/proftwist-backend/services/roadmap/delivery/grpc"
	roadmapHttp "github.com/F0urward/proftwist-backend/services/roadmap/delivery/http"
	roadmapRepository "github.com/F0urward/proftwist-backend/services/roadmap/repository"
	roadmapUsecase "github.com/F0urward/proftwist-backend/services/roadmap/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	chatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	gigachatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	moderationClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/moderationclient"
	roadmapClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
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

var AISet = wire.NewSet(
	aiUsecase.NewAIUsecase,
	aiHttp.NewAIHandlers,
	aiHttp.NewAIHttpRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	mongoClient.NewClient,
	mongoClient.NewDatabase,
	gigachatClient.NewGigaChatClient,
	chatClient.NewChatClient,
	roadmapClient.NewRoadmapClient,
	roadmapInfoClient.NewRoadmapInfoClient,
	authClient.NewAuthClient,
	moderationClient.NewModerationClient,
)
