package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/google/uuid"
)

type RoadmapHandlers struct {
	uc roadmap.Usecase
}

func NewRoadmapHandlers(roadmapUC roadmap.Usecase) roadmap.Handlers {
	return &RoadmapHandlers{
		uc: roadmapUC,
	}
}

func (h *RoadmapHandlers) GetByIDWithProgress(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.GetByID"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	var userID uuid.UUID
	if userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string); ok && userIDStr != "" {
		if parsedID, err := uuid.Parse(userIDStr); err == nil {
			userID = parsedID
			logger = logger.WithField("user_id", userID.String())
		}
	}

	res, err := h.uc.GetByIDWithProgress(r.Context(), roadmapID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by ID")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get roadmap"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmap not found"
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"nodes_count": len(res.Roadmap.Nodes),
		"edges_count": len(res.Roadmap.Edges),
		"has_user":    userID != uuid.Nil,
	}).Info("successfully retrieved roadmap")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) Update(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.Update"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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
		errorMsg := "failed to update roadmap"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmap not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied: you are not the author of this roadmap"
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully updated roadmap")
	w.WriteHeader(http.StatusOK)
}

func (h *RoadmapHandlers) Generate(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.GenerateRoadmap"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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
		errorMsg := "failed to generate roadmap"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmap not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	response := dto.GenerateRoadmapResponseDTO{
		RoadmapID: roadmapID,
	}

	logger.Info("successfully generated roadmap")

	utils.JSONResponse(r.Context(), w, http.StatusOK, response)
}

func (h *RoadmapHandlers) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.CreateMaterial"
	ctx := r.Context()
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		logger.Warn("user ID not found in context")
		utils.JSONError(ctx, w, http.StatusUnauthorized, "authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.WithError(err).Warn("invalid user ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid user ID")
		return
	}

	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		logger.Warn("roadmap_id parameter is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := primitive.ObjectIDFromHex(roadmapIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_id", roadmapIDStr).Warn("invalid roadmap_id format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	nodeIDStr := vars["node_id"]
	if nodeIDStr == "" {
		logger.Warn("node_id parameter is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "node_id parameter is required")
		return
	}

	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		logger.WithError(err).WithField("node_id", nodeIDStr).Warn("invalid node ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid node ID format")
		return
	}

	var req dto.CreateMaterialRequestDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request body")
		return
	}

	material, err := h.uc.CreateMaterial(ctx, userUUID, roadmapID, nodeID, req)
	if err != nil {
		logger.WithError(err).Error("failed to create material")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to create material"

		if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmap or node not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"material_id": material.ID,
		"roadmap_id":  roadmapID.Hex(),
		"node_id":     nodeID,
	}).Info("successfully created material")
	utils.JSONResponse(ctx, w, http.StatusCreated, material)
}

func (h *RoadmapHandlers) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.DeleteMaterial"
	ctx := r.Context()
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		logger.Warn("roadmap_id parameter is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := primitive.ObjectIDFromHex(roadmapIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_id", roadmapIDStr).Warn("invalid roadmap_id format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	nodeIDStr := vars["node_id"]
	if nodeIDStr == "" {
		logger.Warn("node_id parameter is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "node_id parameter is required")
		return
	}

	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		logger.WithError(err).WithField("node_id", nodeIDStr).Warn("invalid node ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid node ID format")
		return
	}

	materialIDStr := vars["material_id"]
	materialID, err := uuid.Parse(materialIDStr)
	if err != nil {
		logger.WithError(err).WithField("material_id", materialIDStr).Warn("invalid material ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid material ID")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		logger.Warn("user ID not found in context")
		utils.JSONError(ctx, w, http.StatusUnauthorized, "authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.WithError(err).Warn("invalid user ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.uc.DeleteMaterial(ctx, roadmapID, nodeID, materialID, userUUID); err != nil {
		logger.WithError(err).Error("failed to delete material")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to delete material"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied - you can only delete your own materials"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "material not found"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"material_id": materialID,
		"roadmap_id":  roadmapID.Hex(),
		"node_id":     nodeID,
		"user_id":     userUUID,
	}).Info("successfully deleted material")

	response := dto.DeleteMaterialResponseDTO{
		Message: "material successfully deleted",
	}
	utils.JSONResponse(ctx, w, http.StatusOK, response)
}

func (h *RoadmapHandlers) GetMaterialsByNode(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.GetMaterialsByNode"
	ctx := r.Context()
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		logger.Warn("roadmap_id parameter is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	roadmapID, err := primitive.ObjectIDFromHex(roadmapIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_id", roadmapIDStr).Warn("invalid roadmap_id format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid roadmap_id format")
		return
	}

	nodeIDStr := vars["node_id"]
	if nodeIDStr == "" {
		logger.Warn("node_id parameter is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "node_id parameter is required")
		return
	}

	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		logger.WithError(err).WithField("node_id", nodeIDStr).Warn("invalid node ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid node ID format")
		return
	}

	materials, err := h.uc.GetMaterialsByNode(ctx, roadmapID, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to get materials by node")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get materials")
		return
	}

	logger.WithFields(map[string]interface{}{
		"roadmap_id": roadmapID.Hex(),
		"node_id":    nodeID,
		"count":      len(materials.Materials),
	}).Info("successfully retrieved materials by node")
	utils.JSONResponse(ctx, w, http.StatusOK, materials)
}

func (h *RoadmapHandlers) UpdateNodeProgress(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.UpdateNodeProgress"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	logger = logger.WithField("user_id", userID.String())

	var req dto.UpdateNodeProgressRequestDTO
	if err = easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger = logger.WithField("node_id", req.NodeID.String())

	err = h.uc.UpdateNodeProgress(r.Context(), userID, roadmapID, &req)
	if err != nil {
		logger.WithError(err).Error("failed to update node progress")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to update node progress"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmap or node not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("status", req.Status).Info("successfully updated node progress")
	w.WriteHeader(http.StatusOK)
}
