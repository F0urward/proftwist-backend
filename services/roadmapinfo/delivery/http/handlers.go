package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RoadmapInfoHandlers struct {
	uc roadmapinfo.Usecase
}

func NewRoadmapInfoHandlers(roadmapInfoUC roadmapinfo.Usecase) roadmapinfo.Handlers {
	return &RoadmapInfoHandlers{
		uc: roadmapInfoUC,
	}
}

func (h *RoadmapInfoHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.GetAll(r.Context())
	if err != nil {
		log.Printf("Failed to get all roadmapInfos: %v", err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to get all roadmapInfos")
		return
	}

	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roadmapInfoIDStr := vars["roadmap_info_id"]
	if roadmapInfoIDStr == "" {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_info_id parameter is required")
		return
	}

	roadmapInfoID, err := uuid.Parse(roadmapInfoIDStr)
	if err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_info_id format")
		return
	}

	res, err := h.uc.GetByID(r.Context(), roadmapInfoID)
	if err != nil {
		log.Printf("Failed to get roadmapInfo by ID %s: %v", roadmapInfoID, err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to get roadmapInfo")
		return
	}

	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoadmapInfoRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.uc.Create(r.Context(), &req)
	if err != nil {
		log.Printf("Failed to create roadmapInfo: %v", err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to create roadmapInfo")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RoadmapInfoHandlers) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roadmapInfoIDStr := vars["roadmap_info_id"]
	if roadmapInfoIDStr == "" {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_info_id parameter is required")
		return
	}

	roadmapInfoID, err := uuid.Parse(roadmapInfoIDStr)
	if err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_info_id format")
		return
	}

	var req dto.UpdateRoadmapInfoRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.uc.Update(r.Context(), roadmapInfoID, &req)
	if err != nil {
		log.Printf("Failed to update roadmapInfo with ID %s: %v", roadmapInfoID, err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to update roadmapInfo")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RoadmapInfoHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roadmapInfoIDStr := vars["roadmap_info_id"]
	if roadmapInfoIDStr == "" {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_info_id parameter is required")
		return
	}

	roadmapInfoID, err := uuid.Parse(roadmapInfoIDStr)
	if err != nil {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_info_id format")
		return
	}

	err = h.uc.Delete(r.Context(), roadmapInfoID)
	if err != nil {
		log.Printf("Failed to delete roadmapInfo with ID %s: %v", roadmapInfoID, err)
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to delete roadmapInfo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
