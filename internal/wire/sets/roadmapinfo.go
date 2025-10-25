package sets

import (
	"github.com/google/wire"

	roadmapInfoGrpc "github.com/F0urward/proftwist-backend/services/roadmapinfo/delivery/grpc"
	roadmapInfoHandlers "github.com/F0urward/proftwist-backend/services/roadmapinfo/delivery/http"
	roadmapInfoRepository "github.com/F0urward/proftwist-backend/services/roadmapinfo/repository"
	roadmapInfoUsecase "github.com/F0urward/proftwist-backend/services/roadmapinfo/usecase"
)

var RoadmapInfoSet = wire.NewSet(
	roadmapInfoRepository.NewRoadmapInfoRepository,
	roadmapInfoUsecase.NewRoadmapInfoUsecase,
	roadmapInfoHandlers.NewRoadmapInfoHandlers,
	roadmapInfoGrpc.NewRoadmapInfoServer,
)
