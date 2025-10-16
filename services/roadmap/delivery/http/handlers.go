package http

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoadmapHandlers struct {
	uc            roadmap.Usecase
	roadmapInfoUC roadmapinfo.Usecase
}

func NewRoadmapHandlers(roadmapUC roadmap.Usecase, roadmapInfoUC roadmapinfo.Usecase) roadmap.Handlers {
	return &RoadmapHandlers{
		uc:            roadmapUC,
		roadmapInfoUC: roadmapInfoUC,
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

	logger.WithField("count", len(res)).Debug("successfully retrieved roadmaps")
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
		"nodes_count": len(res.Nodes),
		"edges_count": len(res.Edges),
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

	// Получаем userID из контекста для проверки авторства
	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	// Получаем RoadmapInfo для проверки авторства
	roadmapInfo, err := h.roadmapInfoUC.GetByRoadmapID(r.Context(), roadmapIDStr)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg := "roadmap not found"
			utils.JSONError(r.Context(), w, statusCode, errorMsg)
			return
		}

		utils.JSONError(r.Context(), w, statusCode, "failed to verify roadmap ownership")
		return
	}

	// Проверяем авторство
	if roadmapInfo.RoadmapInfo.AuthorID != userIDStr {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userIDStr,
			"author_id":       roadmapInfo.RoadmapInfo.AuthorID,
		}).Warn("user is not author of the roadmap")
		utils.JSONError(r.Context(), w, http.StatusForbidden, "access denied: you are not the author of this roadmap")
		return
	}

	var req dto.UpdateRoadmapRequest

	if err = easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing, err := h.uc.GetByID(r.Context(), roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by ID")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}
	if existing == nil {
		logger.Warn("roadmap not found")
		utils.JSONError(r.Context(), w, http.StatusNotFound, "roadmap not found")
		return
	}

	updatedRoadmap := dto.UpdateRequestToEntity(existing, &req)

	res, err := h.uc.Update(r.Context(), updatedRoadmap)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	logger.WithFields(map[string]interface{}{
		"nodes_count": len(res.Nodes),
		"edges_count": len(res.Edges),
	}).Info("successfully updated roadmap")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

//
//func (h *RoadmapHandlers) Delete(w http.ResponseWriter, r *http.Request) {
//	const op = "RoadmapHandlers.Delete"
//	logger := logctx.GetLogger(r.Context()).WithField("op", op)
//
//	vars := mux.Vars(r)
//	roadmapIDStr := vars["roadmap_id"]
//	if roadmapIDStr == "" {
//		logger.Warn("roadmap_id parameter is required")
//		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
//		return
//	}
//
//	roadmapID, err := primitive.ObjectIDFromHex(roadmapIDStr)
//	if err != nil {
//		logger.WithError(err).WithField("roadmap_id", roadmapIDStr).Warn("invalid roadmap_id format")
//		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_id format")
//		return
//	}
//
//	logger = logger.WithField("roadmap_id", roadmapID.Hex())
//
//	// Получаем userID из контекста для проверки авторства
//	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
//	if !ok || userIDStr == "" {
//		logger.Warn("user ID not found in context")
//		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
//		return
//	}
//
//	// Получаем RoadmapInfo для проверки авторства
//	roadmapInfo, err := h.roadmapInfoUC.GetByRoadmapID(r.Context(), roadmapIDStr)
//	if err != nil {
//		logger.WithError(err).Error("failed to get roadmap info for authorization check")
//
//		statusCode := http.StatusInternalServerError
//		if errs.IsNotFoundError(err) {
//			statusCode = http.StatusNotFound
//			errorMsg := "roadmap not found"
//			utils.JSONError(r.Context(), w, statusCode, errorMsg)
//			return
//		}
//
//		utils.JSONError(r.Context(), w, statusCode, "failed to verify roadmap ownership")
//		return
//	}
//
//	// Проверяем авторство
//	if roadmapInfo.RoadmapInfo.AuthorID != userIDStr {
//		logger.WithFields(map[string]interface{}{
//			"request_user_id": userIDStr,
//			"author_id":       roadmapInfo.RoadmapInfo.AuthorID,
//		}).Warn("user is not author of the roadmap")
//		utils.JSONError(r.Context(), w, http.StatusForbidden, "access denied: you are not the author of this roadmap")
//		return
//	}
//
//	err = h.uc.Delete(r.Context(), roadmapID)
//	if err != nil {
//		logger.WithError(err).Error("failed to delete roadmap")
//
//		statusCode := http.StatusInternalServerError
//		if errs.IsNotFoundError(err) {
//			statusCode = http.StatusNotFound
//		}
//
//		utils.JSONError(r.Context(), w, statusCode, err.Error())
//		return
//	}
//
//	logger.Info("successfully deleted roadmap")
//	w.WriteHeader(http.StatusNoContent)
//}
