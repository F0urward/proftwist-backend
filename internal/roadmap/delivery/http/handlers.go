package http

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/roadmap"
)

// type RoadmapHandlers struct {
// 	usecase roadmap.Usecase
// }

// func NewRoadmapHandlers(roadmapUC roadmap.Usecase) roadmap.Handlers {
// 	return &RoadmapHandlers{
// 		usecase: roadmapUC,
// 	}
// }

// func (h *RoadmapHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
// 	roadmaps, err := h.usecase.GetAll(r.Context())
// 	if err != nil {
// 		log.Println()
// 	}
// }
