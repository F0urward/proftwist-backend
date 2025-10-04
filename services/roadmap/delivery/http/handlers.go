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

		// Пробрасываем оригинальную ошибку с нижнего уровня
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
		errorMsg := err.Error()

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Debug("successfully retrieved roadmap")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) GetByAuthorID(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.GetByAuthorID"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	authorIDStr := vars["author_id"]
	if authorIDStr == "" {
		logger.Warn("author_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "author_id parameter is required")
		return
	}

	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		logger.WithError(err).WithField("author_id", authorIDStr).Warn("invalid author_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid author_id format")
		return
	}

	logger = logger.WithField("author_id", authorID.String())

	res, err := h.uc.GetByAuthorID(r.Context(), authorID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmaps by author ID")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	logger.WithField("count", len(res)).Debug("successfully retrieved roadmaps by author")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) Create(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.Create"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var req dto.CreateRoadmapRequest

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"title":       req.Title,
		"is_public":   req.IsPublic,
		"nodes_count": len(req.Nodes),
	})

	// r.Context().FIXME: id автора вытаскиваем из контекста (контекст иниц. в мидлваре и там же достаем id юзера из токена)

	roadmapEntity := dto.CreateRequestToEntity(&req)

	res, err := h.uc.Create(r.Context(), roadmapEntity)
	// r.Context().FIXME: создать roadmapInfo
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")

		statusCode := http.StatusInternalServerError
		if errs.IsBusinessLogicError(err) || errs.IsAlreadyExistsError(err) {
			statusCode = http.StatusBadRequest
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	logger.WithField("roadmap_id", res.ID.Hex()).Info("successfully created roadmap")
	utils.JSONResponse(r.Context(), w, http.StatusCreated, res)
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

	logger.Info("successfully updated roadmap")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.Delete"
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

	err = h.uc.Delete(r.Context(), roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	logger.Info("successfully deleted roadmap")
	w.WriteHeader(http.StatusNoContent)
}

func (h *RoadmapHandlers) SearchByTitle(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.SearchByTitle"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	query := r.URL.Query().Get("title")
	if query == "" {
		logger.Warn("title query parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "title query parameter is required")
		return
	}

	logger = logger.WithField("query", query)
	res, err := h.uc.SearchByTitle(r.Context(), query)
	if err != nil {
		logger.WithError(err).Error("failed to search roadmaps by title")
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("count", len(res)).Debug("successfully searched roadmaps")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapHandlers) UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapHandlers.UpdatePrivacy"
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

	var req dto.UpdatePrivacyRequest
	if err2 := easyjson.UnmarshalFromReader(r.Body, &req); err2 != nil {
		logger.WithError(err2).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger = logger.WithField("is_public", req.IsPublic)

	// r.Context().FIXME: проверить авторство (только автор может менять приватность)

	err = h.uc.UpdatePrivacy(r.Context(), roadmapID, req.IsPublic)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap privacy")

		statusCode := http.StatusInternalServerError
		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
		}

		utils.JSONError(r.Context(), w, statusCode, err.Error())
		return
	}

	response := dto.NewUpdatePrivacyResponse(roadmapID, req.IsPublic)
	logger.Info("successfully updated roadmap privacy")
	utils.JSONResponse(r.Context(), w, http.StatusOK, response)
}
