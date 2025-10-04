package sets

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/services/roadmap"
	roadmapHttp "github.com/F0urward/proftwist-backend/services/roadmap/delivery/http"
	roadmapRepository "github.com/F0urward/proftwist-backend/services/roadmap/repository"
	roadmapUsecase "github.com/F0urward/proftwist-backend/services/roadmap/usecase"
)

var (
	RoadmapSet = wire.NewSet(
		roadmapRepository.NewRoadmapRepository,
		wire.Bind(new(roadmap.Repository), new(*roadmapRepository.RoadmapRepository)),
		roadmapUsecase.NewRoadmapUsecase,
		roadmapHttp.NewRoadmapHandlers,
	)
)
