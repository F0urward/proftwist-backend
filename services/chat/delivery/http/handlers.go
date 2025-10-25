package http

import (
	"context"
	"encoding/json"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/chat/dto"
	"github.com/F0urward/proftwist-backend/services/chat/usecase"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
	}
}

func JSONResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")

	rawBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to marshal response: %v", err)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(rawBytes)
	if err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}

func JSONError(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(ctx, w, statusCode, NewErrorResponse(message))
}

type ChatHandler struct {
	chatUC *usecase.ChatUseCase
	authMW *auth.AuthMiddleware
}

func NewChatHandler(chatUC *usecase.ChatUseCase, authMW *auth.AuthMiddleware) *ChatHandler {
	return &ChatHandler{
		chatUC: chatUC,
		authMW: authMW,
	}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.CreateChat"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var req dto.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req.CreatedBy = userUUID

	chat, err := h.chatUC.CreateChat(ctx, req)
	if err != nil {
		logger.WithError(err).Error("failed to create chat")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to create chat")
		return
	}

	if len(req.MemberIDs) > 0 {
		for _, memberID := range req.MemberIDs {
			if memberID != userUUID {
				addReq := dto.AddMemberRequest{
					ChatID:      chat.ID,
					UserID:      memberID,
					Role:        "member",
					RequestedBy: userUUID,
				}
				if err := h.chatUC.AddMember(ctx, addReq); err != nil {
					logger.WithError(err).Warn("failed to add member to chat")
				}
			}
		}
	}

	chatWithMembers, err := h.chatUC.GetChatWithMembers(ctx, chat.ID, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat with members")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to get chat details")
		return
	}

	memberResponses := make([]dto.ChatMemberResponse, len(chatWithMembers.Members))
	for i, member := range chatWithMembers.Members {
		memberResponses[i] = dto.ToChatMemberResponse(member)
	}

	chatResponse := dto.ToChatResponse(chatWithMembers.Chat)
	response := dto.ChatWithMembersResponse{
		Chat:    &chatResponse,
		Members: memberResponses,
	}

	JSONResponse(ctx, w, http.StatusCreated, response)
}

func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetUserChats"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	chats, err := h.chatUC.GetUserChats(ctx, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get user chats")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to get chats")
		return
	}

	JSONResponse(ctx, w, http.StatusOK, chats)
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetChat"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	chatWithMembers, err := h.chatUC.GetChatWithMembers(ctx, chatID, userUUID)
	if err != nil {
		if err == errs.ErrForbidden {
			JSONError(ctx, w, http.StatusForbidden, "Access denied")
			return
		}
		if err == errs.ErrNotFound {
			JSONError(ctx, w, http.StatusNotFound, "Chat not found")
			return
		}
		logger.WithError(err).Error("failed to get chat")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to get chat")
		return
	}

	memberResponses := make([]dto.ChatMemberResponse, len(chatWithMembers.Members))
	for i, member := range chatWithMembers.Members {
		memberResponses[i] = dto.ToChatMemberResponse(member)
	}

	chatResponse := dto.ToChatResponse(chatWithMembers.Chat)
	response := dto.ChatWithMembersResponse{
		Chat:    &chatResponse,
		Members: memberResponses,
	}

	JSONResponse(ctx, w, http.StatusOK, response)
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.SendMessage"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req.ChatID = chatID
	req.UserID = userUUID

	message, err := h.chatUC.SendMessage(ctx, req)
	if err != nil {
		if err == errs.ErrForbidden {
			JSONError(ctx, w, http.StatusForbidden, "Access denied")
			return
		}
		logger.WithError(err).Error("failed to send message")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to send message")
		return
	}

	response := dto.ToMessageResponse(message)
	JSONResponse(ctx, w, http.StatusCreated, response)
}

func (h *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.GetChatMessages"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
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

	messages, err := h.chatUC.GetChatMessages(ctx, chatID, userUUID, limit, offset)
	if err != nil {
		if err == errs.ErrForbidden {
			JSONError(ctx, w, http.StatusForbidden, "Access denied")
			return
		}
		logger.WithError(err).Error("failed to get chat messages")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	responses := make([]dto.MessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = dto.ToMessageResponse(message)
	}

	JSONResponse(ctx, w, http.StatusOK, responses)
}

func (h *ChatHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.AddMember"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	var req dto.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req.ChatID = chatID
	req.RequestedBy = userUUID

	if err := h.chatUC.AddMember(ctx, req); err != nil {
		if err == errs.ErrForbidden {
			JSONError(ctx, w, http.StatusForbidden, "Access denied")
			return
		}
		if err == errs.ErrAlreadyExists {
			JSONError(ctx, w, http.StatusConflict, "User already in chat")
			return
		}
		if err == errs.ErrNotFound {
			JSONError(ctx, w, http.StatusNotFound, "Chat not found")
			return
		}
		logger.WithError(err).Error("failed to add member")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to add member")
		return
	}

	JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "Member added successfully"})
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
		JSONError(ctx, w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	requestedBy, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	requestedByUUID, err := uuid.Parse(requestedBy)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req := dto.RemoveMemberRequest{
		ChatID:      chatID,
		UserID:      userID,
		RequestedBy: requestedByUUID,
	}

	if err := h.chatUC.RemoveMember(ctx, req); err != nil {
		if err == errs.ErrForbidden {
			JSONError(ctx, w, http.StatusForbidden, "Access denied")
			return
		}
		if err == errs.ErrNotFound {
			JSONError(ctx, w, http.StatusNotFound, "Chat or member not found")
			return
		}
		logger.WithError(err).Error("failed to remove member")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to remove member")
		return
	}

	JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "Member removed successfully"})
}

func (h *ChatHandler) JoinChannel(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.JoinChannel"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.chatUC.JoinChannel(ctx, chatID, userUUID); err != nil {
		if err == errs.ErrForbidden {
			JSONError(ctx, w, http.StatusForbidden, "Cannot join this chat")
			return
		}
		logger.WithError(err).Error("failed to join channel")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to join channel")
		return
	}

	JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "Joined channel successfully"})
}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	const op = "ChatHandler.DeleteChat"
	ctx := r.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		JSONError(ctx, w, http.StatusUnauthorized, "Authentication required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		JSONError(ctx, w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req := dto.DeleteChatRequest{
		ChatID:      chatID,
		RequestedBy: userUUID,
	}

	if err := h.chatUC.DeleteChat(ctx, req); err != nil {
		if err == errs.ErrForbidden {
			JSONError(ctx, w, http.StatusForbidden, "Access denied")
			return
		}
		logger.WithError(err).Error("failed to delete chat")
		JSONError(ctx, w, http.StatusInternalServerError, "Failed to delete chat")
		return
	}

	JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "Chat deleted successfully"})
}
