package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/material"
	"github.com/F0urward/proftwist-backend/services/material/dto"
)

var validate = validator.New()

type MaterialHandlers struct {
	materialUC material.Usecase
}

func NewMaterialHandlers(materialUC material.Usecase) material.Handlers {
	return &MaterialHandlers{
		materialUC: materialUC,
	}
}

func (h *MaterialHandlers) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	const op = "MaterialHandler.CreateMaterial"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

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

	var req dto.CreateMaterialRequestDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		logger.WithError(err).Warn("validation failed")
		utils.JSONError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	material, err := h.materialUC.CreateMaterial(ctx, userUUID, req)
	if err != nil {
		logger.WithError(err).Error("failed to create material")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to create material"

		if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"material_id":     material.ID,
		"roadmap_node_id": req.RoadmapNodeID,
	}).Info("successfully created material")
	utils.JSONResponse(ctx, w, http.StatusCreated, material)
}

func (h *MaterialHandlers) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	const op = "MaterialHandler.DeleteMaterial"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
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

	if err := h.materialUC.DeleteMaterial(ctx, materialID, userUUID); err != nil {
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
		"user_id":     userUUID,
	}).Info("successfully deleted material")

	response := dto.DeleteMaterialResponseDTO{
		Message: "material successfully deleted",
	}
	utils.JSONResponse(ctx, w, http.StatusOK, response)
}

func (h *MaterialHandlers) GetMaterialsByNode(w http.ResponseWriter, r *http.Request) {
	const op = "MaterialHandler.GetMaterialsByNode"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	nodeID := vars["node_id"]

	materials, err := h.materialUC.GetMaterialsByNode(ctx, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to get materials by roadmap node")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get materials")
		return
	}

	logger.WithFields(map[string]interface{}{
		"node_id": nodeID,
		"count":   len(materials.Materials),
	}).Info("successfully retrieved materials by roadmap node")
	utils.JSONResponse(ctx, w, http.StatusOK, materials)
}

func (h *MaterialHandlers) GetUserMaterials(w http.ResponseWriter, r *http.Request) {
	const op = "MaterialHandler.GetUserMaterials"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

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

	materials, err := h.materialUC.GetUserMaterials(ctx, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get user materials")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get materials")
		return
	}

	logger.WithFields(map[string]interface{}{
		"user_id": userUUID,
		"count":   len(materials.Materials),
	}).Info("successfully retrieved user materials")
	utils.JSONResponse(ctx, w, http.StatusOK, materials)
}
