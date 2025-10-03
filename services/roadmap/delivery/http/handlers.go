package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RoadmapHandlers struct {
	uc roadmap.Usecase
}

func NewRoadmapHandlers(roadmapUC roadmap.Usecase) roadmap.Handlers {
	return &RoadmapHandlers{
		uc: roadmapUC,
	}
}

func (h *RoadmapHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.GetAll(r.Context())
	if err != nil {
		log.Printf("Failed to get all roadmaps: %v", err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to get all roadmaps")
		return
	}

	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := uuid.Parse(roadmapIDStr)
	if err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	res, err := h.uc.GetByID(r.Context(), roadmapID)
	if err != nil {
		log.Printf("Failed to get roadmap by ID %s: %v", roadmapID, err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to get roadmap")
		return
	}

	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoadmapRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.uc.Create(r.Context(), &req)
	if err != nil {
		log.Printf("Failed to create roadmap: %v", err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to create roadmap")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RoadmapHandlers) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := uuid.Parse(roadmapIDStr)
	if err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	var req dto.UpdateRoadmapRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.ID = roadmapIDStr

	err = h.uc.Update(r.Context(), &req)
	if err != nil {
		log.Printf("Failed to update roadmap with ID %s: %v", roadmapID, err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to update roadmap")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RoadmapHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := uuid.Parse(roadmapIDStr)
	if err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	err = h.uc.Delete(r.Context(), roadmapID)
	if err != nil {
		log.Printf("Failed to delete roadmap with ID %s: %v", roadmapID, err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to delete roadmap")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
