package http

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/F0urward/proftwist-backend/services/chat/dto"
)

type ChatHandler struct {
	chatUC chat.Usecase
}

func NewChatHandler(chatUC chat.Usecase) chat.Handlers {
	return &ChatHandler{
		chatUC: chatUC,
	}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.CreateChat"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var req dto.CreateChatRequestDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request body")
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

	// role, ok := r.Context().Value(utils.RoleKey{}).(string)
	// if !ok {
	// 	logger.Warn("user role not found in context")
	// 	utils.JSONError(ctx, w, http.StatusUnauthorized, "authentication required")
	// 	return
	// }

	req.CreatedByID = userUUID
	// req.CreatedByRole = role

	chat, err := h.chatUC.CreateChat(ctx, req)
	if err != nil {
		logger.WithError(err).Error("failed to create chat")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to create chat"

		if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithField("chat_id", chat.ID).Info("successfully created chat")
	utils.JSONResponse(ctx, w, http.StatusCreated, chat)
}

func (h *ChatHandler) GetChatsByUser(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetChatsByUser"
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

	chats, err := h.chatUC.GetUserChats(ctx, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get user chats")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get chats")
		return
	}

	logger.WithField("count", len(chats)).Info("successfully retrieved user chats")
	utils.JSONResponse(ctx, w, http.StatusOK, chats)
}

func (h *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetChatMessages"
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

	res, err := h.chatUC.GetChatMessages(ctx, chatID, userUUID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get chat messages")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get messages"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"count":   len(res.ChatMessages),
	}).Info("successfully retrieved chat messages")
	utils.JSONResponse(ctx, w, http.StatusOK, res)
}

func (h *ChatHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.AddMember"
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

	var req dto.AddMemberRequestDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request body")
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

	req.ChatID = chatID
	req.RequestedBy = userUUID

	if err := h.chatUC.AddMember(ctx, req); err != nil {
		logger.WithError(err).Error("failed to add member")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to add member"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsAlreadyExistsError(err) {
			statusCode = http.StatusConflict
			errorMsg = "user already in chat"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "chat not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": req.UserID,
	}).Info("successfully added member to chat")
	utils.JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "member added successfully"})
}

func (h *ChatHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.RemoveMember"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	userIDStr := vars["user_id"]

	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		logger.WithError(err).WithField("chat_id", chatIDStr).Warn("invalid chat ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid chat ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid user ID")
		return
	}

	requestedBy, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		logger.Warn("user ID not found in context")
		utils.JSONError(ctx, w, http.StatusUnauthorized, "authentication required")
		return
	}

	requestedByUUID, err := uuid.Parse(requestedBy)
	if err != nil {
		logger.WithError(err).Warn("invalid user ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid user ID")
		return
	}

	req := dto.RemoveMemberRequestDTO{
		ChatID:      chatID,
		UserID:      userID,
		RequestedBy: requestedByUUID,
	}

	if err := h.chatUC.RemoveMember(ctx, req); err != nil {
		logger.WithError(err).Error("failed to remove member")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to remove member"

		if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "chat or member not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("successfully removed member from chat")
	utils.JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "member removed successfully"})
}
