package sets

import (
	"github.com/google/wire"

	roadmapHttp "github.com/F0urward/proftwist-backend/services/roadmap/delivery/http"
	roadmapRepository "github.com/F0urward/proftwist-backend/services/roadmap/repository"
	roadmapUsecase "github.com/F0urward/proftwist-backend/services/roadmap/usecase"
)

var RoadmapSet = wire.NewSet(
	roadmapRepository.NewRoadmapMongoRepository,
	roadmapRepository.NewRoadmapGigaChatWebapi,
	roadmapUsecase.NewRoadmapUsecase,
	roadmapHttp.NewRoadmapHandlers,
)
