package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
)

type RoadmapInfoHandlers struct {
	uc roadmapinfo.Usecase
}

func NewRoadmapInfoHandlers(roadmapInfoUC roadmapinfo.Usecase) roadmapinfo.Handlers {
	return &RoadmapInfoHandlers{
		uc: roadmapInfoUC,
	}
}

func (h *RoadmapInfoHandlers) GetAllPublic(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.GetAllPublic"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	res, err := h.uc.GetAllPublic(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get public roadmapInfos")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get public roadmapInfos"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "no public roadmapInfos found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("count", len(res.RoadmapsInfo)).Info("successfully retrieved public roadmapInfos")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) GetAllPublicByCategoryID(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.GetAllPublicByCategoryID"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	categoryIDStr := vars["category_id"]
	if categoryIDStr == "" {
		logger.Warn("category_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "category_id parameter is required")
		return
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		logger.WithError(err).WithField("category_id", categoryIDStr).Warn("invalid category_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid category_id format")
		return
	}

	logger = logger.WithField("category_id", categoryID.String())

	res, err := h.uc.GetAllPublicByCategoryID(r.Context(), categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmapInfos by category ID")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get roadmapInfos by category"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfos not found for this category"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("count", len(res.RoadmapsInfo)).Info("successfully retrieved roadmapInfos by category")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.GetAllByUserID"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Warn("invalid user ID format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user ID")
		return
	}

	res, err := h.uc.GetAllByUserID(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmapInfos by user ID")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get roadmapInfos"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfos not found for this user"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("count", len(res.RoadmapsInfo)).Info("successfully retrieved roadmapInfos by user ID")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.GetByID"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	logger.Info("successfully retrieved roadmapInfo")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) GetByRoadmapID(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.GetByRoadmapID"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	roadmapIDStr := vars["roadmap_id"]
	if roadmapIDStr == "" {
		logger.Warn("roadmap_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "roadmap_id parameter is required")
		return
	}

	logger = logger.WithField("roadmap_id", roadmapIDStr)

	res, err := h.uc.GetByRoadmapID(r.Context(), roadmapIDStr)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmapInfo by roadmap ID")

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

	logger.Info("successfully retrieved roadmapInfo by roadmap ID")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) CreatePrivate(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Create"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req dto.CreatePrivateRoadmapInfoRequestDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.AuthorID = userIDStr

	logger = logger.WithFields(map[string]interface{}{
		"author_id": req.AuthorID,
	})

	res, err := h.uc.CreatePrivate(r.Context(), &req)
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

	logger.WithFields(map[string]interface{}{
		"roadmap_info_id": res.RoadmapInfo.ID,
		"roadmap_id":      res.RoadmapInfo.RoadmapID,
	}).Info("successfully created roadmapInfo")

	utils.JSONResponse(r.Context(), w, http.StatusCreated, res)
}

func (h *RoadmapInfoHandlers) UpdatePrivate(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Update"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	var req dto.UpdatePrivateRoadmapInfoRequestDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.uc.UpdatePrivate(r.Context(), roadmapInfoID, userID, &req)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmapInfo")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to update roadmapInfo"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfo not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied: you are not the author of this roadmap"
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
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	err = h.uc.Delete(r.Context(), roadmapInfoID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmapInfo")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to delete roadmapInfo"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfo not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied: you are not the author of this roadmap"
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

func (h *RoadmapInfoHandlers) Fork(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Fork"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	res, err := h.uc.Fork(r.Context(), roadmapInfoID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to fork roadmapInfo")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to fork roadmapInfo"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfo not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied: cannot fork this roadmap"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("forked_roadmap_info_id", res.RoadmapInfo.ID).Info("successfully forked roadmapInfo")
	utils.JSONResponse(r.Context(), w, http.StatusCreated, res)
}

func (h *RoadmapInfoHandlers) Publish(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Publish"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	res, err := h.uc.Publish(r.Context(), roadmapInfoID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to publish roadmapInfo")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to publish roadmapInfo"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmapInfo not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied: cannot publish this roadmap"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("published_roadmap_info_id", res.RoadmapInfo.ID).Info("successfully published roadmapInfo")
	utils.JSONResponse(r.Context(), w, http.StatusCreated, res)
}

func (h *RoadmapInfoHandlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Subscribe"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	err = h.uc.Subscribe(r.Context(), roadmapInfoID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to subscribe to roadmap")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to subscribe to roadmap"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "roadmap not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "cannot subscribe to private roadmap"
		} else if errs.IsAlreadyExistsError(err) {
			statusCode = http.StatusConflict
			errorMsg = "already subscribed to this roadmap"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully subscribed to roadmap")
	w.WriteHeader(http.StatusCreated)
}

func (h *RoadmapInfoHandlers) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.Unsubscribe"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	err = h.uc.Unsubscribe(r.Context(), roadmapInfoID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to unsubscribe from roadmap")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to unsubscribe from roadmap"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "subscription not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully unsubscribed from roadmap")
	w.WriteHeader(http.StatusNoContent)
}

func (h *RoadmapInfoHandlers) GetSubscribedRoadmaps(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.GetSubscribedRoadmaps"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Warn("invalid user ID format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user ID")
		return
	}

	res, err := h.uc.GetSubscribed(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("failed to get subscribed roadmaps")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get subscribed roadmaps"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "no subscriptions found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("count", len(res.RoadmapsInfo)).Info("successfully retrieved subscribed roadmaps")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *RoadmapInfoHandlers) CheckSubscription(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.CheckSubscription"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

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

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	isSubscribed, err := h.uc.CheckSubscription(r.Context(), roadmapInfoID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check subscription status")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to check subscription status"

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

	response := map[string]bool{
		"is_subscribed": isSubscribed,
	}

	logger.WithField("is_subscribed", isSubscribed).Info("successfully checked subscription status")
	utils.JSONResponse(r.Context(), w, http.StatusOK, response)
}

func (h *RoadmapInfoHandlers) SearchPublic(w http.ResponseWriter, r *http.Request) {
	const op = "RoadmapInfoHandlers.SearchPublic"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	query := r.URL.Query().Get("q")
	if query == "" {
		logger.Warn("search query parameter 'q' is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "search query parameter 'q' is required")
		return
	}

	categoryIDStr := r.URL.Query().Get("category_id")
	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err != nil {
			logger.WithError(err).WithField("category_id", categoryIDStr).Warn("invalid category_id format")
			utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid category_id format")
			return
		}
		categoryID = &id
	}

	logger = logger.WithFields(map[string]interface{}{
		"query":       query,
		"category_id": categoryID,
	})

	res, err := h.uc.SearchPublic(r.Context(), query, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to search public roadmapInfos")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to search roadmapInfos"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "no roadmapInfos found for this search query"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("count", len(res.RoadmapsInfo)).Info("successfully searched public roadmapInfos")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}
