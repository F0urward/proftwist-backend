package http

import (
	"encoding/json"
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
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
	const op = "RoadmapInfoHandlers.GetAll"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	res, err := h.uc.GetAll(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get all roadmapInfos")
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to get all roadmapInfos")
		return
	}

	logger.WithField("count", len(res.RoadmapsInfo)).Debug("successfully retrieved roadmapInfos")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.GetByID"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	roadmapInfoIDStr := vars["roadmap_info_id"]
	if roadmapInfoIDStr == "" {
		logger.Warn("roadmap_info_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_info_id parameter is required")
		return
	}

	roadmapInfoID, err := uuid.Parse(roadmapInfoIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_info_id", roadmapInfoIDStr).Warn("invalid roadmap_info_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_info_id format")
		return
	}

	logger = logger.WithField("roadmap_info_id", roadmapInfoID.String())

	res, err := h.uc.GetByID(r.Context(), roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmapInfo by ID")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get roadmapInfo"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfo not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Debug("successfully retrieved roadmapInfo")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) Create(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Create"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var req dto.CreateRoadmapInfoRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"author_id": req.AuthorID,
	})

	err := h.uc.Create(r.Context(), &req)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmapInfo")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to create roadmapInfo"

		if errs.IsBusinessLogicError(err) || errs.IsAlreadyExistsError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully created roadmapInfo")
	w.WriteHeader(http.StatusCreated)
}

func (h *RoadmapInfoHandlers) Update(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Update"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	roadmapInfoIDStr := vars["roadmap_info_id"]
	if roadmapInfoIDStr == "" {
		logger.Warn("roadmap_info_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_info_id parameter is required")
		return
	}

	roadmapInfoID, err := uuid.Parse(roadmapInfoIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_info_id", roadmapInfoIDStr).Warn("invalid roadmap_info_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_info_id format")
		return
	}

	logger = logger.WithField("roadmap_info_id", roadmapInfoID.String())

	var req dto.UpdateRoadmapInfoRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.uc.Update(r.Context(), roadmapInfoID, &req)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmapInfo")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to update roadmapInfo"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfo not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully updated roadmapInfo")
	w.WriteHeader(http.StatusOK)
}

func (h *RoadmapInfoHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Delete"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	roadmapInfoIDStr := vars["roadmap_info_id"]
	if roadmapInfoIDStr == "" {
		logger.Warn("roadmap_info_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_info_id parameter is required")
		return
	}

	roadmapInfoID, err := uuid.Parse(roadmapInfoIDStr)
	if err != nil {
		logger.WithError(err).WithField("roadmap_info_id", roadmapInfoIDStr).Warn("invalid roadmap_info_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid roadmap_info_id format")
		return
	}

	logger = logger.WithField("roadmap_info_id", roadmapInfoID.String())

	err = h.uc.Delete(r.Context(), roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmapInfo")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to delete roadmapInfo"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfo not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully deleted roadmapInfo")
	w.WriteHeader(http.StatusNoContent)
}
