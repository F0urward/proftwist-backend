package http

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/chat"
)

type ChatHandlers struct {
	chatUC chat.Usecase
}

func NewChatHandler(chatUC chat.Usecase) chat.Handlers {
	return &ChatHandlers{
		chatUC: chatUC,
	}
}

func (h *ChatHandlers) GetGroupChatByNode(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetGroupChatByNode"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	nodeID := vars["node_id"]

	chat, err := h.chatUC.GetGroupChatByNode(ctx, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to get group chats by node")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get chats")
		return
	}

	logger.Info("successfully retrieved group chats by node")
	utils.JSONResponse(ctx, w, http.StatusOK, chat)
}

func (h *ChatHandlers) GetGroupChatsByUser(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetGroupChatsByUser"
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

	chats, err := h.chatUC.GetGroupChatsByUser(ctx, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get user group chats")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get chats")
		return
	}

	logger.WithField("count", len(chats.GroupChats)).Info("successfully retrieved user group chats")
	utils.JSONResponse(ctx, w, http.StatusOK, chats)
}

func (h *ChatHandlers) GetGroupChatMembers(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetGroupChatMembers"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		logger.WithError(err).WithField("chat_id", chatIDStr).Warn("invalid chat ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid chat ID")
		return
	}

	members, err := h.chatUC.GetGroupChatMembers(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get group chat members")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get members"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "chat not found"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"count":   len(members.Members),
	}).Info("successfully retrieved group chat members")
	utils.JSONResponse(ctx, w, http.StatusOK, members)
}

func (h *ChatHandlers) GetGroupChatMessages(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetGroupChatMessages"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		logger.WithError(err).WithField("chat_id", chatIDStr).Warn("invalid chat ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid chat ID")
		return
	}

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	res, err := h.chatUC.GetGroupChatMessages(ctx, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get group chat messages")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get messages"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "chat not found"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"count":   len(res.ChatMessages),
	}).Info("successfully retrieved group chat messages")
	utils.JSONResponse(ctx, w, http.StatusOK, res)
}

func (h *ChatHandlers) JoinGroupChat(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.JoinGroupChat"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		logger.WithError(err).WithField("chat_id", chatIDStr).Warn("invalid chat ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid chat ID")
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

	if err := h.chatUC.JoinGroupChat(ctx, chatID, userUUID); err != nil {
		logger.WithError(err).Error("failed to join group chat")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to join chat"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "chat not found"
		} else if errs.IsAlreadyExistsError(err) {
			statusCode = http.StatusConflict
			errorMsg = "already a member of this chat"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userUUID,
	}).Info("successfully joined group chat")
	utils.JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "successfully joined chat"})
}

func (h *ChatHandlers) LeaveGroupChat(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.LeaveGroupChat"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		logger.WithError(err).WithField("chat_id", chatIDStr).Warn("invalid chat ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid chat ID")
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

	if err := h.chatUC.LeaveGroupChat(ctx, chatID, userUUID); err != nil {
		logger.WithError(err).Error("failed to leave group chat")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to leave chat"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "chat not found or not a member"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userUUID,
	}).Info("successfully left group chat")
	utils.JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "successfully left chat"})
}

func (h *ChatHandlers) GetDirectChatsByUser(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetDirectChatsByUser"
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

	chats, err := h.chatUC.GetDirectChatsByUser(ctx, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get user direct chats")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get chats")
		return
	}

	logger.WithField("count", len(chats.DirectChats)).Info("successfully retrieved user direct chats")
	utils.JSONResponse(ctx, w, http.StatusOK, chats)
}

func (h *ChatHandlers) GetDirectChatMessages(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetDirectChatMessages"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		logger.WithError(err).WithField("chat_id", chatIDStr).Warn("invalid chat ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid chat ID")
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

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	res, err := h.chatUC.GetDirectChatMessages(ctx, chatID, userUUID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get direct chat messages")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get messages"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "chat not found"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"count":   len(res.ChatMessages),
	}).Info("successfully retrieved direct chat messages")
	utils.JSONResponse(ctx, w, http.StatusOK, res)
}
