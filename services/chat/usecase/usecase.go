package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/F0urward/proftwist-backend/services/chat/dto"
)

type ChatUsecase struct {
	repo       chat.Repository
	notifier   chat.Notifier
	authClient authclient.AuthServiceClient
}

func NewChatUsecase(repo chat.Repository, notifier chat.Notifier, authClient authclient.AuthServiceClient) chat.Usecase {
	return &ChatUsecase{
		repo:       repo,
		notifier:   notifier,
		authClient: authClient,
	}
}

func (uc *ChatUsecase) CreateGroupChat(ctx context.Context, userID uuid.UUID, req *dto.CreateGroupChatRequestDTO) (*dto.CreateGroupChatResponseDTO, error) {
	const op = "ChatUsecase.CreateGroupChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	groupChat := dto.CreateGroupChatRequestToEntity(req)

	createdChat, err := uc.repo.CreateGroupChat(ctx, groupChat)
	if err != nil {
		logger.WithError(err).Error("failed to create group chat")
		return nil, fmt.Errorf("failed to create group chat: %w", err)
	}

	if err := uc.repo.AddGroupChatMembers(ctx, createdChat.ID, req.MemberIDs); err != nil {
		logger.WithError(err).Error("failed to add group chat members")
		return nil, fmt.Errorf("failed to add group chat members: %w", err)
	}

	response := dto.CreateGroupChatResponseFromEntity(createdChat)
	logger.WithField("chat_id", createdChat.ID.String()).Info("successfully created group chat")
	return &response, nil
}

func (uc *ChatUsecase) DeleteGroupChat(ctx context.Context, chatID uuid.UUID) error {
	const op = "ChatUsecase.DeleteGroupChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if err := uc.repo.DeleteGroupChat(ctx, chatID); err != nil {
		logger.WithError(err).Error("failed to delete group chat")
		return fmt.Errorf("failed to delete group chat: %w", err)
	}

	logger.WithField("chat_id", chatID.String()).Info("successfully deleted group chat")
	return nil
}

func (uc *ChatUsecase) GetGroupChatByNode(ctx context.Context, nodeID string) (*dto.GroupChatResponseDTO, error) {
	const op = "ChatUsecase.GetGroupChatByNode"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetGroupChatByNode(ctx, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to get group chat by node")
		return nil, fmt.Errorf("failed to get group chat by node: %w", err)
	}

	if chat == nil {
		logger.WithField("node_id", nodeID).Warn("group chat not found")
		return nil, errs.ErrNotFound
	}

	response := dto.GroupChatToDTO(chat)
	logger.Info("successfully retrieved group chat by node")
	return &response, nil
}

func (uc *ChatUsecase) GetGroupChatsByUser(ctx context.Context, userID uuid.UUID) (*dto.GroupChatListResponseDTO, error) {
	const op = "ChatUsecase.GetGroupChatsByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chats, err := uc.repo.GetGroupChatsByUser(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user group chats")
		return nil, fmt.Errorf("failed to get user group chats: %w", err)
	}

	if chats == nil {
		chats = []*entities.GroupChat{}
	}

	response := dto.GroupChatListToDTO(chats)
	logger.WithField("count", len(response.GroupChats)).Info("successfully retrieved user group chats")
	return &response, nil
}

func (uc *ChatUsecase) GetGroupChatMembers(ctx context.Context, chatID uuid.UUID) (*dto.ChatMemberListResponseDTO, error) {
	const op = "ChatUsecase.GetGroupChatMembers"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	members, err := uc.repo.GetGroupChatMembers(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get group chat members")
		return nil, fmt.Errorf("failed to get group chat members: %w", err)
	}

	if members == nil {
		members = []*entities.GroupChatMember{}
	}

	userIDs := make([]uuid.UUID, len(members))
	for i, member := range members {
		userIDs[i] = member.UserID
	}

	userData := uc.fetchUserData(ctx, userIDs)
	response := dto.GroupChatMemberListToDTO(members, userData)

	logger.WithField("count", len(response.Members)).Info("successfully retrieved group chat members")
	return &response, nil
}

func (uc *ChatUsecase) GetGroupChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error) {
	const op = "ChatUsecase.GetGroupChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	messages, err := uc.repo.GetGroupChatMessages(ctx, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get group chat messages")
		return nil, fmt.Errorf("failed to get group chat messages: %w", err)
	}

	if messages == nil {
		messages = []*entities.Message{}
	}

	userIDs := make([]uuid.UUID, len(messages))
	for i, message := range messages {
		userIDs[i] = message.UserID
	}

	userData := uc.fetchUserData(ctx, userIDs)
	response := dto.GetChatMessagesResponseToDTO(messages, userData)

	logger.WithField("count", len(response.ChatMessages)).Info("successfully retrieved group chat messages")
	return &response, nil
}

func (uc *ChatUsecase) SendGroupMessage(ctx context.Context, req *dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error) {
	const op = "ChatUsecase.SendGroupMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsGroupChatMember(ctx, req.ChatID, req.UserID)
	if err != nil {
		logger.WithError(err).Error("failed to check group chat membership")
		return nil, fmt.Errorf("failed to check group chat membership: %w", err)
	}

	if !isMember {
		return nil, errs.ErrForbidden
	}

	message := dto.SendMessageRequestToEntity(req)
	if err := uc.repo.SaveGroupMessage(ctx, message); err != nil {
		logger.WithError(err).Error("failed to save message")
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	userData := uc.fetchSingleUserData(ctx, req.UserID)
	messageDTO := dto.ChatMessageResponseDTO{
		ID:        message.ID,
		ChatID:    message.ChatID,
		User:      userData,
		Content:   message.Content,
		Metadata:  message.Metadata,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}

	members, err := uc.repo.GetGroupChatMembers(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Warn("failed to get group chat members for broadcast")
	} else {
		if members == nil {
			members = []*entities.GroupChatMember{}
		}

		if err := uc.BroadcastGroupMessageSent(ctx, req.ChatID, messageDTO, members); err != nil {
			logger.WithError(err).Warn("failed to broadcast group message")
		}
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID.String(),
		"chat_id":    message.ChatID.String(),
		"user_id":    message.UserID.String(),
	}).Info("group message sent successfully")

	return &messageDTO, nil
}

func (uc *ChatUsecase) JoinGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatUsecase.JoinGroupChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsGroupChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check group chat membership")
		return fmt.Errorf("failed to check group chat membership: %w", err)
	}

	if isMember {
		return errs.ErrAlreadyExists
	}

	if err := uc.repo.AddGroupChatMember(ctx, chatID, userID); err != nil {
		logger.WithError(err).Error("failed to join group chat")
		return fmt.Errorf("failed to join group chat: %w", err)
	}

	if err := uc.BroadcastUserJoined(ctx, chatID, userID); err != nil {
		logger.WithError(err).Warn("failed to broadcast user joined notification")
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user joined group chat")

	return nil
}

func (uc *ChatUsecase) LeaveGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatUsecase.LeaveGroupChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsGroupChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check group chat membership")
		return fmt.Errorf("failed to check group chat membership: %w", err)
	}

	if !isMember {
		return errs.ErrNotFound
	}

	if err := uc.repo.RemoveGroupChatMember(ctx, chatID, userID); err != nil {
		logger.WithError(err).Error("failed to leave group chat")
		return fmt.Errorf("failed to leave group chat: %w", err)
	}

	if err := uc.BroadcastUserLeft(ctx, chatID, userID); err != nil {
		logger.WithError(err).Warn("failed to broadcast user left notification")
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user left group chat")

	return nil
}

func (uc *ChatUsecase) CreateDirectChat(ctx context.Context, userID uuid.UUID, req *dto.CreateDirectChatRequestDTO) (*dto.CreateDirectChatResponseDTO, error) {
	const op = "ChatUsecase.CreateDirectChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	existingChat, err := uc.repo.GetDirectChatByUsers(ctx, userID, req.OtherUserID)
	if err != nil {
		logger.WithError(err).Error("failed to check existing direct chat")
		return nil, fmt.Errorf("failed to check existing direct chat: %w", err)
	}

	if existingChat != nil {
		return nil, errs.ErrAlreadyExists
	}

	directChat := dto.CreateDirectChatRequestToEntity(req, userID)

	createdChat, err := uc.repo.CreateDirectChat(ctx, directChat)
	if err != nil {
		logger.WithError(err).Error("failed to create direct chat")
		return nil, fmt.Errorf("failed to create direct chat: %w", err)
	}

	response := dto.CreateDirectChatResponseFromEntity(createdChat)
	logger.WithField("chat_id", createdChat.ID.String()).Info("successfully created direct chat")
	return &response, nil
}

func (uc *ChatUsecase) DeleteDirectChat(ctx context.Context, chatID uuid.UUID) error {
	const op = "ChatUsecase.DeleteDirectChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if err := uc.repo.DeleteDirectChat(ctx, chatID); err != nil {
		logger.WithError(err).Error("failed to delete direct chat")
		return fmt.Errorf("failed to delete direct chat: %w", err)
	}

	logger.WithField("chat_id", chatID.String()).Info("successfully deleted direct chat")
	return nil
}

func (uc *ChatUsecase) GetDirectChatsByUser(ctx context.Context, userID uuid.UUID) (*dto.DirectChatListResponseDTO, error) {
	const op = "ChatUsecase.GetDirectChatsByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chats, err := uc.repo.GetDirectChatsByUser(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user direct chats")
		return nil, fmt.Errorf("failed to get user direct chats: %w", err)
	}

	if chats == nil {
		chats = []*entities.DirectChat{}
	}

	response := dto.DirectChatListToDTO(chats)
	logger.WithField("count", len(response.DirectChats)).Info("successfully retrieved user direct chats")
	return &response, nil
}

func (uc *ChatUsecase) GetDirectChatMembers(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (*dto.ChatMemberListResponseDTO, error) {
	const op = "ChatUsecase.GetDirectChatMembers"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsDirectChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check direct chat membership")
		return nil, fmt.Errorf("failed to check direct chat membership: %w", err)
	}

	if !isMember {
		return nil, errs.ErrForbidden
	}

	directChat, err := uc.repo.GetDirectChat(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get direct chat")
		return nil, fmt.Errorf("failed to get direct chat: %w", err)
	}

	if directChat == nil {
		logger.WithField("chat_id", chatID.String()).Warn("direct chat not found")
		return nil, errs.ErrNotFound
	}

	user1Data := uc.fetchSingleUserData(ctx, directChat.User1ID)
	user2Data := uc.fetchSingleUserData(ctx, directChat.User2ID)
	members := dto.DirectChatMembersToDTO(directChat.User1ID, directChat.User2ID, user1Data, user2Data)

	logger.Info("successfully retrieved direct chat members")
	return &members, nil
}

func (uc *ChatUsecase) GetDirectChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error) {
	const op = "ChatUsecase.GetDirectChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsDirectChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check direct chat membership")
		return nil, fmt.Errorf("failed to check direct chat membership: %w", err)
	}

	if !isMember {
		return nil, errs.ErrForbidden
	}

	messages, err := uc.repo.GetDirectChatMessages(ctx, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get direct chat messages")
		return nil, fmt.Errorf("failed to get direct chat messages: %w", err)
	}

	if messages == nil {
		messages = []*entities.Message{}
	}

	userIDs := make([]uuid.UUID, len(messages))
	for i, message := range messages {
		userIDs[i] = message.UserID
	}

	userData := uc.fetchUserData(ctx, userIDs)
	response := dto.GetChatMessagesResponseToDTO(messages, userData)

	logger.WithField("count", len(response.ChatMessages)).Info("successfully retrieved direct chat messages")
	return &response, nil
}

func (uc *ChatUsecase) SendDirectMessage(ctx context.Context, req *dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error) {
	const op = "ChatUsecase.SendDirectMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsDirectChatMember(ctx, req.ChatID, req.UserID)
	if err != nil {
		logger.WithError(err).Error("failed to check direct chat membership")
		return nil, fmt.Errorf("failed to check direct chat membership: %w", err)
	}

	if !isMember {
		return nil, errs.ErrForbidden
	}

	message := dto.SendMessageRequestToEntity(req)
	if err := uc.repo.SaveDirectMessage(ctx, message); err != nil {
		logger.WithError(err).Error("failed to save message")
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	userData := uc.fetchSingleUserData(ctx, req.UserID)
	messageDTO := dto.ChatMessageResponseDTO{
		ID:        message.ID,
		ChatID:    message.ChatID,
		User:      userData,
		Content:   message.Content,
		Metadata:  message.Metadata,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}

	directChat, err := uc.repo.GetDirectChat(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Warn("failed to get direct chat for broadcast")
	} else if directChat != nil {
		userIDs := []uuid.UUID{directChat.User1ID, directChat.User2ID}
		userDataMap := uc.fetchUserData(ctx, userIDs)
		memberDTOs := []dto.MemberResponseDTO{userDataMap[directChat.User1ID], userDataMap[directChat.User2ID]}
		if err := uc.BroadcastDirectMessageSent(ctx, req.ChatID, messageDTO, memberDTOs); err != nil {
			logger.WithError(err).Warn("failed to broadcast direct message")
		}
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID.String(),
		"chat_id":    message.ChatID.String(),
		"user_id":    message.UserID.String(),
	}).Info("direct message sent successfully")

	return &messageDTO, nil
}

func (uc *ChatUsecase) fetchUserData(ctx context.Context, userIDs []uuid.UUID) map[uuid.UUID]dto.MemberResponseDTO {
	if len(userIDs) == 0 {
		return make(map[uuid.UUID]dto.MemberResponseDTO)
	}

	userIDStrings := make([]string, len(userIDs))
	for i, id := range userIDs {
		userIDStrings[i] = id.String()
	}

	resp, err := uc.authClient.GetUsersByIDs(ctx, &authclient.GetUsersByIDsRequest{UserIds: userIDStrings})
	if err != nil || resp == nil {
		return uc.createFallbackUserData(userIDs)
	}

	userData := make(map[uuid.UUID]dto.MemberResponseDTO, len(resp.Users))
	for _, user := range resp.Users {
		if user == nil {
			continue
		}
		userID, err := uuid.Parse(user.Id)
		if err != nil {
			continue
		}
		userData[userID] = dto.MemberResponseDTO{
			UserID:    userID,
			Username:  user.Username,
			AvatarURL: user.AvatarUrl,
		}
	}

	return userData
}

func (uc *ChatUsecase) fetchSingleUserData(ctx context.Context, userID uuid.UUID) dto.MemberResponseDTO {
	resp, err := uc.authClient.GetUserByID(ctx, &authclient.GetUserByIDRequest{UserId: userID.String()})
	if err != nil || resp == nil || resp.User == nil {
		return dto.MemberResponseDTO{UserID: userID}
	}

	return dto.MemberResponseDTO{
		UserID:    userID,
		Username:  resp.User.Username,
		AvatarURL: resp.User.AvatarUrl,
	}
}

func (uc *ChatUsecase) createFallbackUserData(userIDs []uuid.UUID) map[uuid.UUID]dto.MemberResponseDTO {
	userData := make(map[uuid.UUID]dto.MemberResponseDTO, len(userIDs))
	for _, id := range userIDs {
		userData[id] = dto.MemberResponseDTO{UserID: id}
	}
	return userData
}

func (uc *ChatUsecase) BroadcastGroupMessageSent(ctx context.Context, chatID uuid.UUID, message dto.ChatMessageResponseDTO, members []*entities.GroupChatMember) error {
	userIDs := uc.extractUserIDsFromGroupMembers(members)
	return uc.notifier.NotifyMessageSent(
		ctx,
		userIDs,
		chatID.String(),
		message.ID.String(),
		message.User.UserID.String(),
		message.Content,
		message.User.Username,
		message.User.AvatarURL,
	)
}

func (uc *ChatUsecase) BroadcastDirectMessageSent(ctx context.Context, chatID uuid.UUID, message dto.ChatMessageResponseDTO, members []dto.MemberResponseDTO) error {
	userIDs := uc.extractUserIDsFromMemberDTOs(members)
	return uc.notifier.NotifyMessageSent(
		ctx,
		userIDs,
		chatID.String(),
		message.ID.String(),
		message.User.UserID.String(),
		message.Content,
		message.User.Username,
		message.User.AvatarURL,
	)
}

func (uc *ChatUsecase) BroadcastTyping(ctx context.Context, chatID, userID uuid.UUID, typing bool, isGroup bool) error {
	const op = "ChatUsecase.BroadcastTyping"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userData := uc.fetchSingleUserData(ctx, userID)
	var userIDs []string

	if isGroup {
		isMember, err := uc.repo.IsGroupChatMember(ctx, chatID, userID)
		if err != nil {
			return fmt.Errorf("failed to check group chat membership: %w", err)
		}
		if !isMember {
			return errs.ErrForbidden
		}

		members, err := uc.repo.GetGroupChatMembers(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to get group chat members: %w", err)
		}

		if members == nil {
			members = []*entities.GroupChatMember{}
		}
		userIDs = uc.extractUserIDsFromGroupMembersExcept(members, userID.String())
	} else {
		isMember, err := uc.repo.IsDirectChatMember(ctx, chatID, userID)
		if err != nil {
			return fmt.Errorf("failed to check direct chat membership: %w", err)
		}
		if !isMember {
			return errs.ErrForbidden
		}

		directChat, err := uc.repo.GetDirectChat(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to get direct chat: %w", err)
		}
		if directChat == nil {
			return errs.ErrNotFound
		}

		var otherUserID uuid.UUID
		if directChat.User1ID == userID {
			otherUserID = directChat.User2ID
		} else {
			otherUserID = directChat.User1ID
		}
		userIDs = []string{otherUserID.String()}
	}

	if err := uc.notifier.NotifyTyping(ctx, userIDs, chatID.String(), userID.String(), userData.Username, typing); err != nil {
		return fmt.Errorf("failed to notify typing: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id":  chatID.String(),
		"user_id":  userID.String(),
		"typing":   typing,
		"is_group": isGroup,
	}).Info("typing notification broadcasted")

	return nil
}

func (uc *ChatUsecase) BroadcastUserJoined(ctx context.Context, chatID, userID uuid.UUID) error {
	const op = "ChatUsecase.BroadcastUserJoined"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	members, err := uc.repo.GetGroupChatMembers(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to get group chat members: %w", err)
	}
	if members == nil {
		members = []*entities.GroupChatMember{}
	}

	userIDs := uc.extractUserIDsFromGroupMembers(members)
	userData := uc.fetchSingleUserData(ctx, userID)

	if err := uc.notifier.NotifyUserJoined(ctx, userIDs, chatID.String(), userID.String(), userData.Username); err != nil {
		return fmt.Errorf("failed to notify user joined: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user joined notification broadcasted")

	return nil
}

func (uc *ChatUsecase) BroadcastUserLeft(ctx context.Context, chatID, userID uuid.UUID) error {
	const op = "ChatUsecase.BroadcastUserLeft"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	members, err := uc.repo.GetGroupChatMembers(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to get group chat members: %w", err)
	}
	if members == nil {
		members = []*entities.GroupChatMember{}
	}

	userIDs := uc.extractUserIDsFromGroupMembersExcept(members, userID.String())
	userData := uc.fetchSingleUserData(ctx, userID)

	if err := uc.notifier.NotifyUserLeft(ctx, userIDs, chatID.String(), userID.String(), userData.Username); err != nil {
		return fmt.Errorf("failed to notify user left: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user left notification broadcasted")

	return nil
}

func (uc *ChatUsecase) extractUserIDsFromGroupMembers(members []*entities.GroupChatMember) []string {
	if members == nil {
		return []string{}
	}
	userIDs := make([]string, len(members))
	for i, member := range members {
		if member == nil {
			continue
		}
		userIDs[i] = member.UserID.String()
	}
	return userIDs
}

func (uc *ChatUsecase) extractUserIDsFromGroupMembersExcept(members []*entities.GroupChatMember, excludeUserID string) []string {
	if members == nil {
		return []string{}
	}
	var userIDs []string
	for _, member := range members {
		if member == nil {
			continue
		}
		userIDStr := member.UserID.String()
		if userIDStr != excludeUserID {
			userIDs = append(userIDs, userIDStr)
		}
	}
	return userIDs
}

func (uc *ChatUsecase) extractUserIDsFromMemberDTOs(members []dto.MemberResponseDTO) []string {
	if members == nil {
		return []string{}
	}
	userIDs := make([]string, len(members))
	for i, member := range members {
		userIDs[i] = member.UserID.String()
	}
	return userIDs
}
