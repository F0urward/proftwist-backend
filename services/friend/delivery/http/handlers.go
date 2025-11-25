package http

import (
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/friend"
	"github.com/F0urward/proftwist-backend/services/friend/dto"
)

type FriendHandlers struct {
	friendUC friend.Usecase
}

func NewFriendHandlers(friendUC friend.Usecase) friend.Handlers {
	return &FriendHandlers{
		friendUC: friendUC,
	}
}

func (h *FriendHandlers) GetFriends(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.GetFriends"
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

	friends, err := h.friendUC.GetFriends(ctx, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get friends")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get friends")
		return
	}

	logger.WithField("count", len(friends.Friends)).Info("successfully retrieved friends")
	utils.JSONResponse(ctx, w, http.StatusOK, friends)
}

func (h *FriendHandlers) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.DeleteFriend"
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

	vars := mux.Vars(r)
	friendIDStr := vars["friend_id"]
	friendID, err := uuid.Parse(friendIDStr)
	if err != nil {
		logger.WithError(err).WithField("friend_id", friendIDStr).Warn("invalid friend ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid friend ID")
		return
	}

	if err := h.friendUC.DeleteFriend(ctx, userUUID, friendID); err != nil {
		logger.WithError(err).Error("failed to delete friend")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to delete friend"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "friend not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"user_id":   userUUID,
		"friend_id": friendID,
	}).Info("successfully deleted friend")
	utils.JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "friend removed successfully"})
}

func (h *FriendHandlers) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.GetFriendRequests"
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

	requests, err := h.friendUC.GetFriendRequests(ctx, userUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get friend requests")
		utils.JSONError(ctx, w, http.StatusInternalServerError, "failed to get friend requests")
		return
	}

	logger.WithFields(map[string]interface{}{
		"received_count": len(requests.Received),
		"sent_count":     len(requests.Sent),
	}).Info("successfully retrieved friend requests")
	utils.JSONResponse(ctx, w, http.StatusOK, requests)
}

func (h *FriendHandlers) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.AcceptFriendRequest"
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

	vars := mux.Vars(r)
	requestIDStr := vars["request_id"]
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		logger.WithError(err).WithField("request_id", requestIDStr).Warn("invalid request ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request ID")
		return
	}

	friend, err := h.friendUC.AcceptFriendRequest(ctx, userUUID, requestID)
	if err != nil {
		logger.WithError(err).Error("failed to accept friend request")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to accept friend request"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "friend request not found"
		} else if errs.IsForbiddenError(err) {
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
		"user_id":    userUUID,
		"request_id": requestID,
		"friend_id":  friend.UserID,
	}).Info("successfully accepted friend request")
	utils.JSONResponse(ctx, w, http.StatusOK, friend)
}

func (h *FriendHandlers) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.RejectFriendRequest"
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

	vars := mux.Vars(r)
	requestIDStr := vars["request_id"]
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		logger.WithError(err).WithField("request_id", requestIDStr).Warn("invalid request ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request ID")
		return
	}

	if err := h.friendUC.RejectFriendRequest(ctx, userUUID, requestID); err != nil {
		logger.WithError(err).Error("failed to reject friend request")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to reject friend request"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "friend request not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"user_id":    userUUID,
		"request_id": requestID,
	}).Info("successfully rejected friend request")
	utils.JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "friend request rejected successfully"})
}

func (h *FriendHandlers) CreateFriendRequest(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.CreateFriendRequest"
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

	var req dto.CreateFriendRequestDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.friendUC.CreateFriendRequest(ctx, userUUID, &req); err != nil {
		logger.WithError(err).Error("failed to create friend request")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to create friend request"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "target user not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsAlreadyExistsError(err) {
			statusCode = http.StatusConflict
			errorMsg = "friend request already exists or users are already friends"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"user_id":        userUUID,
		"target_user_id": req.TargetUserID,
	}).Info("successfully created friend request")
	utils.JSONResponse(ctx, w, http.StatusCreated, map[string]string{"message": "friend request sent successfully"})
}

func (h *FriendHandlers) DeleteFriendRequest(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.DeleteFriendRequest"
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

	vars := mux.Vars(r)
	requestIDStr := vars["request_id"]
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		logger.WithError(err).WithField("request_id", requestIDStr).Warn("invalid request ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid request ID")
		return
	}

	if err := h.friendUC.DeleteFriendRequest(ctx, userUUID, requestID); err != nil {
		logger.WithError(err).Error("failed to delete friend request")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to delete friend request"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "friend request not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"user_id":    userUUID,
		"request_id": requestID,
	}).Info("successfully deleted friend request")
	utils.JSONResponse(ctx, w, http.StatusOK, map[string]string{"message": "friend request deleted successfully"})
}

func (h *FriendHandlers) GetFriendshipStatus(w http.ResponseWriter, r *http.Request) {
	const op = "FriendHandler.GetFriendshipStatus"
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

	vars := mux.Vars(r)
	targetUserIDStr := vars["target_user_id"]
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		logger.WithError(err).WithField("target_user_id", targetUserIDStr).Warn("invalid target user ID format")
		utils.JSONError(ctx, w, http.StatusBadRequest, "invalid target user ID")
		return
	}

	status, err := h.friendUC.GetFriendshipStatus(ctx, userUUID, targetUserID)
	if err != nil {
		logger.WithError(err).Error("failed to get friendship status")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get friendship status"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "user not found"
		}

		utils.JSONError(ctx, w, statusCode, errorMsg)
		return
	}

	logger.WithFields(map[string]interface{}{
		"user_id":        userUUID,
		"target_user_id": targetUserID,
		"status":         status.Status,
	}).Info("successfully retrieved friendship status")
	utils.JSONResponse(ctx, w, http.StatusOK, status)
}
