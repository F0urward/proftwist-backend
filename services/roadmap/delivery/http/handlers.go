package http

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/google/uuid"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	const op = "RoadmapHandlers.GetAll"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	res, err := h.uc.GetAll(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get all roadmaps")
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("count", len(res.Roadmaps)).Debug("successfully retrieved roadmaps")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.GetByID"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]

	if roadmapIDStr == "" {
		logger.Warn("roadmap_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := primitive.ObjectIDFromHex(roadmapIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_id", roadmapIDStr).Warn("invalid roadmap_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	logger = logger.WithField("roadmap_id", roadmapID.Hex())

	res, err := h.uc.GetByID(r.Context(), roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by ID")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	logger.WithFields(map[string]interface{}{
		"nodes_count": len(res.Roadmap.Nodes),
		"edges_count": len(res.Roadmap.Edges),
	}).Debug("successfully retrieved roadmap")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) Update(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.Update"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		logger.Warn("roadmap_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := primitive.ObjectIDFromHex(roadmapIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_id", roadmapIDStr).Warn("invalid roadmap_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	logger = logger.WithField("roadmap_id", roadmapID.Hex())

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	var req dto.UpdateRoadmapRequestDTO
	if err = easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.uc.Update(r.Context(), userID, roadmapID, &req)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")

		statusCode := http.StatusInternalServerError
		switch {
		case errs.IsNotFoundError(err):
			statusCode = http.StatusNotFound
		case errs.IsBusinessLogicError(err):
			statusCode = http.StatusBadRequest
		case errs.IsForbiddenError(err):
			statusCode = http.StatusForbidden
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	logger.Info("successfully updated roadmap")
	w.WriteHeader(http.StatusOK)
}

func (h *RoadmapHandlers) Generate(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.GenerateRoadmap"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		logger.Warn("roadmap_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := primitive.ObjectIDFromHex(roadmapIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_id", roadmapIDStr).Warn("invalid roadmap_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	logger = logger.WithField("roadmap_id", roadmapID.Hex())

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	var generateReq dto.GenerateRoadmapRequestDTO
	if err = easyjson.UnmarshalFromReader(r.Body, &generateReq); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger.WithFields(map[string]interface{}{
		"complexity": generateReq.Complexity,
	}).Info("starting roadmap generation")

	_, err = h.uc.Generate(r.Context(), userID, roadmapID, &generateReq)
	if err != nil {
		logger.WithError(err).Error("failed to generate roadmap")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	response := dto.GenerateRoadmapResponseDTO{
		RoadmapID: roadmapID,
	}

	logger.Info("successfully generated roadmap")

	utils.JSONResponse(r.Context(), w, http.StatusOK, response)
}
