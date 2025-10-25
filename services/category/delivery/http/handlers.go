package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/category"
	"github.com/F0urward/proftwist-backend/services/category/dto"
)

type CategoryHandlers struct {
	uc category.Usecase
}

func NewCategoryHandlers(categoryUC category.Usecase) category.Handlers {
	return &CategoryHandlers{
		uc: categoryUC,
	}
}

func (h *CategoryHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandlers.GetAll"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	res, err := h.uc.GetAll(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get all categories")
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to get all categories")
		return
	}

	logger.WithField("count", len(res.Categories)).Info("successfully retrieved categories")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *CategoryHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandlers.GetByID"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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

	res, err := h.uc.GetByID(r.Context(), categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to get category by ID")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get category"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "category not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully retrieved category")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *CategoryHandlers) Create(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandlers.Create"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var req dto.CreateCategoryRequestDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger = logger.WithFields(map[string]interface{}{
		"name": req.Name,
	})

	res, err := h.uc.Create(r.Context(), &req)
	if err != nil {
		logger.WithError(err).Error("failed to create category")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to create category"

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
		"category_id": res.Category.CategoryID,
	}).Info("successfully created category")

	utils.JSONResponse(r.Context(), w, http.StatusCreated, res)
}

func (h *CategoryHandlers) Update(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandlers.Update"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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

	var req dto.UpdateCategoryRequestDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.uc.Update(r.Context(), categoryID, &req)
	if err != nil {
		logger.WithError(err).Error("failed to update category")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to update category"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "category not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied: you are not the author of this category"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully updated category")
	w.WriteHeader(http.StatusOK)
}

func (h *CategoryHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandlers.Delete"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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

	err = h.uc.Delete(r.Context(), categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to delete category")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to delete category"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "category not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied: you are not the author of this category"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully deleted category")
	w.WriteHeader(http.StatusNoContent)
}
