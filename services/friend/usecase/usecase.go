package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/friend"
	"github.com/F0urward/proftwist-backend/services/friend/dto"
)

type FriendUsecase struct {
	repo       friend.Repository
	authClient authclient.AuthServiceClient
	chatClient chatclient.ChatServiceClient
}

func NewFriendUsecase(repo friend.Repository, authClient authclient.AuthServiceClient, chatClient chatclient.ChatServiceClient) friend.Usecase {
	return &FriendUsecase{
		repo:       repo,
		authClient: authClient,
		chatClient: chatClient,
	}
}

func (uc *FriendUsecase) GetFriends(ctx context.Context, userID uuid.UUID) (*dto.GetFriendsResponseDTO, error) {
	const op = "FriendUsecase.GetFriends"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	friendIDs, err := uc.repo.GetFriendIDs(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get friend IDs")
		return nil, fmt.Errorf("failed to get friend IDs: %w", err)
	}

	if len(friendIDs) == 0 {
		return &dto.GetFriendsResponseDTO{Friends: []dto.FriendResponseDTO{}}, nil
	}

	userData, err := uc.fetchUserData(ctx, friendIDs)
	if err != nil {
		logger.WithError(err).Error("failed to fetch user data")
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	chatIDs := make(map[uuid.UUID]*uuid.UUID)
	for _, friendID := range friendIDs {
		chatID, err := uc.repo.GetFriendshipChatID(ctx, userID, friendID)
		if err != nil {
			logger.WithError(err).Warn("failed to get chat ID for friend")
			chatIDs[friendID] = nil
		} else {
			chatIDs[friendID] = chatID
		}
	}

	sharedRoadmaps := make(map[uuid.UUID]int)
	for _, friendID := range friendIDs {
		sharedRoadmaps[friendID] = 0
	}

	response := dto.FriendsToDTO(friendIDs, userData, sharedRoadmaps, chatIDs)
	logger.WithField("count", len(response.Friends)).Info("successfully retrieved friends")
	return &response, nil
}

func (uc *FriendUsecase) DeleteFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	const op = "FriendUsecase.DeleteFriend"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isFriends, err := uc.repo.IsFriends(ctx, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("failed to check friendship")
		return fmt.Errorf("failed to check friendship: %w", err)
	}

	if !isFriends {
		return errs.ErrNotFound
	}

	chatID, err := uc.repo.GetFriendshipChatID(ctx, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("failed to get friendship chat ID")
		return fmt.Errorf("failed to get friendship chat ID: %w", err)
	}

	if err := uc.repo.DeleteFriendship(ctx, userID, friendID); err != nil {
		logger.WithError(err).Error("failed to delete friendship")
		return fmt.Errorf("failed to delete friendship: %w", err)
	}

	if err := uc.repo.DeleteFriendship(ctx, friendID, userID); err != nil {
		logger.WithError(err).Error("failed to delete reverse friendship")
		return fmt.Errorf("failed to delete reverse friendship: %w", err)
	}

	if chatID != nil {
		if err := uc.deleteDirectChat(ctx, *chatID); err != nil {
			logger.WithError(err).Warn("failed to delete direct chat")
		} else {
			logger.WithField("chat_id", *chatID).Debug("direct chat deleted")
		}
	}

	logger.WithFields(map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	}).Info("successfully deleted friend")
	return nil
}

func (uc *FriendUsecase) GetFriendRequests(ctx context.Context, userID uuid.UUID) (*dto.GetFriendRequestsResponseDTO, error) {
	const op = "FriendUsecase.GetFriendRequests"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	receivedRequests, err := uc.repo.GetFriendRequestsForUser(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get received friend requests")
		return nil, fmt.Errorf("failed to get received friend requests: %w", err)
	}

	sentRequests, err := uc.repo.GetSentFriendRequests(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get sent friend requests")
		return nil, fmt.Errorf("failed to get sent friend requests: %w", err)
	}

	allUserIDs := uc.collectUserIDsFromRequests(receivedRequests, sentRequests)
	userData, err := uc.fetchUserData(ctx, allUserIDs)
	if err != nil {
		logger.WithError(err).Error("failed to fetch user data")
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	response := dto.GetFriendRequestsResponseToDTO(receivedRequests, sentRequests, userData)
	logger.WithFields(map[string]interface{}{
		"received_count": len(response.Received),
		"sent_count":     len(response.Sent),
	}).Info("successfully retrieved friend requests")
	return &response, nil
}

func (uc *FriendUsecase) AcceptFriendRequest(ctx context.Context, userID, requestID uuid.UUID) (*dto.FriendResponseDTO, error) {
	const op = "FriendUsecase.AcceptFriendRequest"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	request, err := uc.repo.GetFriendRequestByID(ctx, requestID)
	if err != nil {
		logger.WithError(err).Error("failed to get friend request")
		return nil, fmt.Errorf("failed to get friend request: %w", err)
	}

	if request == nil {
		return nil, errs.ErrNotFound
	}

	if request.ToUserID != userID {
		return nil, errs.ErrForbidden
	}

	if request.Status != entities.FriendStatusPending {
		return nil, errs.ErrBusinessLogic
	}

	chatID, err := uc.createDirectChat(ctx, request.FromUserID, request.ToUserID)
	if err != nil {
		logger.WithError(err).Error("failed to create direct chat")
		return nil, fmt.Errorf("failed to create direct chat: %w", err)
	}

	if err := uc.repo.DeleteFriendRequest(ctx, requestID); err != nil {
		logger.WithError(err).Error("failed to delete friend request")
		return nil, fmt.Errorf("failed to delete friend request: %w", err)
	}

	if err := uc.repo.CreateFriendship(ctx, request.FromUserID, request.ToUserID, chatID); err != nil {
		logger.WithError(err).Error("failed to create friendship")
		return nil, fmt.Errorf("failed to create friendship: %w", err)
	}

	if err := uc.repo.CreateFriendship(ctx, request.ToUserID, request.FromUserID, chatID); err != nil {
		logger.WithError(err).Error("failed to create reverse friendship")
		return nil, fmt.Errorf("failed to create reverse friendship: %w", err)
	}

	friendData, err := uc.fetchSingleUserData(ctx, request.FromUserID)
	if err != nil {
		logger.WithError(err).Error("failed to fetch friend data")
		return nil, fmt.Errorf("failed to fetch friend data: %w", err)
	}

	response := dto.FriendToDTO(request.FromUserID, friendData, 0, &chatID)
	logger.WithFields(map[string]interface{}{
		"user_id":    userID,
		"request_id": requestID,
		"friend_id":  request.FromUserID,
		"chat_id":    chatID,
	}).Info("successfully accepted friend request")
	return &response, nil
}

func (uc *FriendUsecase) DeleteFriendRequest(ctx context.Context, userID, requestID uuid.UUID) error {
	const op = "FriendUsecase.DeleteFriendRequest"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	request, err := uc.repo.GetFriendRequestByID(ctx, requestID)
	if err != nil {
		logger.WithError(err).Error("failed to get friend request")
		return fmt.Errorf("failed to get friend request: %w", err)
	}

	if request == nil {
		return errs.ErrNotFound
	}

	if request.FromUserID != userID && request.ToUserID != userID {
		return errs.ErrForbidden
	}

	if err := uc.repo.DeleteFriendRequest(ctx, requestID); err != nil {
		logger.WithError(err).Error("failed to delete friend request")
		return fmt.Errorf("failed to delete friend request: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"user_id":    userID,
		"request_id": requestID,
	}).Info("successfully deleted friend request")
	return nil
}

func (uc *FriendUsecase) CreateFriendRequest(ctx context.Context, userID uuid.UUID, req *dto.CreateFriendRequestDTO) error {
	const op = "FriendUsecase.CreateFriendRequest"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if userID == req.TargetUserID {
		return errs.ErrBusinessLogic
	}

	existingRequest, err := uc.repo.GetPendingFriendRequestBetweenUsers(ctx, userID, req.TargetUserID)
	if err != nil {
		logger.WithError(err).Error("failed to check existing friend request")
		return fmt.Errorf("failed to check existing friend request: %w", err)
	}

	if existingRequest != nil {
		return errs.ErrAlreadyExists
	}

	isFriends, err := uc.repo.IsFriends(ctx, userID, req.TargetUserID)
	if err != nil {
		logger.WithError(err).Error("failed to check friendship")
		return fmt.Errorf("failed to check friendship: %w", err)
	}

	if isFriends {
		return errs.ErrAlreadyExists
	}

	friendRequest := dto.CreateFriendRequestToEntity(userID, req)
	if err := uc.repo.CreateFriendRequest(ctx, friendRequest); err != nil {
		logger.WithError(err).Error("failed to create friend request")
		return fmt.Errorf("failed to create friend request: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"user_id":        userID,
		"target_user_id": req.TargetUserID,
	}).Info("successfully created friend request")
	return nil
}

func (uc *FriendUsecase) createDirectChat(ctx context.Context, user1ID, user2ID uuid.UUID) (uuid.UUID, error) {
	const op = "FriendUsecase.createDirectChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	req := &chatclient.CreateDirectChatRequest{
		UserId:      user1ID.String(),
		OtherUserId: user2ID.String(),
	}

	resp, err := uc.chatClient.CreateDirectChat(ctx, req)
	if err != nil {
		logger.WithError(err).Error("failed to call chat client")
		return uuid.Nil, fmt.Errorf("failed to call chat client: %w", err)
	}

	if resp.Error != "" {
		logger.WithField("error", resp.Error).Error("chat service returned error")
		return uuid.Nil, fmt.Errorf("chat service error: %s", resp.Error)
	}

	if resp.DirectChat == nil {
		logger.Error("chat service returned nil direct chat")
		return uuid.Nil, fmt.Errorf("chat service returned nil direct chat")
	}

	chatID, err := uuid.Parse(resp.DirectChat.Id)
	if err != nil {
		logger.WithError(err).Error("failed to parse chat ID")
		return uuid.Nil, fmt.Errorf("failed to parse chat ID: %w", err)
	}

	logger.WithField("chat_id", chatID).Debug("direct chat created")
	return chatID, nil
}

func (uc *FriendUsecase) deleteDirectChat(ctx context.Context, chatID uuid.UUID) error {
	const op = "FriendUsecase.deleteDirectChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	req := &chatclient.DeleteDirectChatRequest{
		ChatId: chatID.String(),
	}

	resp, err := uc.chatClient.DeleteDirectChat(ctx, req)
	if err != nil {
		logger.WithError(err).Error("failed to call chat client")
		return fmt.Errorf("failed to call chat client: %w", err)
	}

	if !resp.Success {
		logger.WithField("error", resp.Error).Error("chat service failed to delete chat")
		return fmt.Errorf("chat service failed to delete chat: %s", resp.Error)
	}

	logger.WithField("chat_id", chatID).Debug("direct chat deleted")
	return nil
}

func (uc *FriendUsecase) fetchUserData(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]*dto.UserDTO, error) {
	if len(userIDs) == 0 {
		return make(map[uuid.UUID]*dto.UserDTO), nil
	}

	userIDStrings := make([]string, len(userIDs))
	for i, id := range userIDs {
		userIDStrings[i] = id.String()
	}

	resp, err := uc.authClient.GetUsersByIDs(ctx, &authclient.GetUsersByIDsRequest{UserIds: userIDStrings})
	if err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}

	userData := make(map[uuid.UUID]*dto.UserDTO, len(resp.Users))
	for _, user := range resp.Users {
		if user == nil {
			continue
		}
		userID, err := uuid.Parse(user.Id)
		if err != nil {
			continue
		}
		userData[userID] = &dto.UserDTO{
			ID:        userID,
			Username:  user.Username,
			AvatarURL: &user.AvatarUrl,
		}
	}

	return userData, nil
}

func (uc *FriendUsecase) fetchSingleUserData(ctx context.Context, userID uuid.UUID) (*dto.UserDTO, error) {
	resp, err := uc.authClient.GetUserByID(ctx, &authclient.GetUserByIDRequest{UserId: userID.String()})
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if resp == nil || resp.User == nil {
		return &dto.UserDTO{ID: userID}, nil
	}

	return &dto.UserDTO{
		ID:        userID,
		Username:  resp.User.Username,
		AvatarURL: &resp.User.AvatarUrl,
	}, nil
}

func (uc *FriendUsecase) collectUserIDsFromRequests(receivedRequests, sentRequests []*entities.FriendRequest) []uuid.UUID {
	userIDMap := make(map[uuid.UUID]bool)

	for _, req := range receivedRequests {
		userIDMap[req.FromUserID] = true
		userIDMap[req.ToUserID] = true
	}

	for _, req := range sentRequests {
		userIDMap[req.FromUserID] = true
		userIDMap[req.ToUserID] = true
	}

	userIDs := make([]uuid.UUID, 0, len(userIDMap))
	for userID := range userIDMap {
		userIDs = append(userIDs, userID)
	}

	return userIDs
}
