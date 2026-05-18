package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type AIHandlers struct {
	uc ai.Usecase
}

func NewAIHandlers(aiUC ai.Usecase) ai.Handlers {
	return &AIHandlers{
		uc: aiUC,
	}
}

func (h *AIHandlers) GenerateRoadmapNodeDescription(w http.ResponseWriter, r *http.Request) {
	const op = "AIHandlers.GenerateRoadmapNodeDescription"
	ctx := r.Context()
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	var req dto.GenerateRoadmapNodeDescriptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.NodeID = strings.TrimSpace(req.NodeID)
	req.NodeLabel = strings.TrimSpace(req.NodeLabel)
	if req.NodeID == "" {
		logger.Warn("node_id is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "node_id is required")
		return
	}
	if req.NodeLabel == "" {
		logger.Warn("node_label is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "node_label is required")
		return
	}

	res, err := h.uc.GenerateRoadmapNodeDescription(ctx, req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.WithError(err).Warn("roadmap node description request canceled by client")
			return
		}
		logger.WithError(err).Error("failed to generate roadmap node description")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to generate roadmap node description")
		return
	}

	utils.JSONResponse(ctx, w, http.StatusOK, res)
}

func (h *AIHandlers) GenerateRoadmap(w http.ResponseWriter, r *http.Request) {
	const op = "AIHandlers.GenerateRoadmap"
	ctx := r.Context()
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	var req dto.GenerateRoadmapRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Prompt = strings.TrimSpace(req.Prompt)
	if req.Prompt == "" {
		logger.Warn("prompt is required")
		utils.JSONError(ctx, w, http.StatusBadRequest, "prompt is required")
		return
	}

	res, err := h.uc.GenerateRoadmap(ctx, req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.WithError(err).Warn("roadmap generation request canceled by client")
			return
		}
		logger.WithError(err).Error("failed to generate roadmap")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to generate roadmap")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
