package usecase

import (
	"context"
	"fmt"

	"github.com/F0urward/proftwist-backend/services/chat/dto"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/google/uuid"
)

type ChatUsecase struct {
	repo     chat.Repository
	notifier chat.Notifier
}

func NewChatUsecase(repo chat.Repository, notifier chat.Notifier) chat.Usecase {
	return &ChatUsecase{
		repo:     repo,
		notifier: notifier,
	}
}

func (uc *ChatUsecase) GetGroupChatByNode(ctx context.Context, nodeID string) (*dto.GroupChatResponseDTO, error) {
	const op = "ChatUsecase.GetGroupChatByNode"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetGroupChatByNode(ctx, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to get group chats by node")
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.GroupChatMemberListToDTO(members)
	logger.WithField("count", len(response.Members)).Info("successfully retrieved group chat members")
	return &response, nil
}

func (uc *ChatUsecase) GetGroupChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error) {
	const op = "ChatUsecase.GetGroupChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	messages, err := uc.repo.GetGroupChatMessages(ctx, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get group chat messages")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.GetChatMessagesResponseToDTO(messages)
	logger.WithField("count", len(response.ChatMessages)).Info("successfully retrieved group chat messages")
	return &response, nil
}

func (uc *ChatUsecase) SendGroupMessage(ctx context.Context, req *dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error) {
	const op = "ChatUsecase.SendGroupMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsGroupChatMember(ctx, req.ChatID, req.UserID)
	if err != nil {
		logger.WithError(err).Error("failed to check group chat membership")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !isMember {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	message := dto.SendMessageRequestToEntity(req)

	if err := uc.repo.SaveGroupMessage(ctx, message); err != nil {
		logger.WithError(err).Error("failed to save message")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	messageDTO := dto.MessageToDTO(message)

	members, err := uc.repo.GetGroupChatMembers(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Warn("failed to get group chat members for broadcast")
	} else {
		if err := uc.BroadcastGroupMessageSent(ctx, req.ChatID, messageDTO, members); err != nil {
			logger.WithError(err).Warn("failed to broadcast group message")
		}
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID.String(),
		"chat_id":    message.ChatID.String(),
		"user_id":    message.UserID.String(),
	}).Info("group message sent and broadcasted successfully")

	return &messageDTO, nil
}

func (uc *ChatUsecase) JoinGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatUsecase.JoinGroupChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsGroupChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check group chat membership")
		return fmt.Errorf("%s: %w", op, err)
	}
	if isMember {
		return fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
	}

	if err := uc.repo.AddGroupChatMember(ctx, chatID, userID); err != nil {
		logger.WithError(err).Error("failed to join group chat")
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := uc.BroadcastUserJoined(ctx, chatID, userID); err != nil {
		logger.WithError(err).Warn("failed to broadcast user joined notification")
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user joined group chat and notification broadcasted")

	return nil
}

func (uc *ChatUsecase) LeaveGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatUsecase.LeaveGroupChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsGroupChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check group chat membership")
		return fmt.Errorf("%s: %w", op, err)
	}
	if !isMember {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err := uc.repo.RemoveGroupChatMember(ctx, chatID, userID); err != nil {
		logger.WithError(err).Error("failed to leave group chat")
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := uc.BroadcastUserLeft(ctx, chatID, userID); err != nil {
		logger.WithError(err).Warn("failed to broadcast user left notification")
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user left group chat and notification broadcasted")

	return nil
}

func (uc *ChatUsecase) GetDirectChatsByUser(ctx context.Context, userID uuid.UUID) (*dto.DirectChatListResponseDTO, error) {
	const op = "ChatUsecase.GeDirectChatsByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chats, err := uc.repo.GeDirectChatsByUser(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user direct chats")
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !isMember {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	directChat, err := uc.repo.GetDirectChat(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get direct chat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if directChat == nil {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	members := dto.DirectChatMembersToDTO(directChat.User1ID, directChat.User2ID)
	logger.Info("successfully retrieved direct chat members")
	return &members, nil
}

func (uc *ChatUsecase) GetDirectChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error) {
	const op = "ChatUsecase.GetDirectChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsDirectChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check direct chat membership")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !isMember {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	messages, err := uc.repo.GetDirectChatMessages(ctx, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get direct chat messages")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.GetChatMessagesResponseToDTO(messages)
	logger.WithField("count", len(response.ChatMessages)).Info("successfully retrieved direct chat messages")
	return &response, nil
}

func (uc *ChatUsecase) SendDirectMessage(ctx context.Context, req *dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error) {
	const op = "ChatUsecase.SendDirectMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsDirectChatMember(ctx, req.ChatID, req.UserID)
	if err != nil {
		logger.WithError(err).Error("failed to check direct chat membership")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !isMember {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	message := dto.SendMessageRequestToEntity(req)

	if err := uc.repo.SaveDirectMessage(ctx, message); err != nil {
		logger.WithError(err).Error("failed to save message")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	messageDTO := dto.MessageToDTO(message)

	directChat, err := uc.repo.GetDirectChat(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Warn("failed to get direct chat for broadcast")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	memberDTOs := []dto.ChatMemberResponseDTO{
		{UserID: directChat.User1ID},
		{UserID: directChat.User2ID},
	}

	if err := uc.BroadcastDirectMessageSent(ctx, req.ChatID, messageDTO, memberDTOs); err != nil {
		logger.WithError(err).Warn("failed to broadcast direct message")
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID.String(),
		"chat_id":    message.ChatID.String(),
		"user_id":    message.UserID.String(),
	}).Info("direct message sent and broadcasted successfully")

	return &messageDTO, nil
}

func (uc *ChatUsecase) BroadcastGroupMessageSent(ctx context.Context, chatID uuid.UUID, message dto.ChatMessageResponseDTO, members []*entities.GroupChatMember) error {
	userIDs := uc.extractUserIDsFromGroupMembers(members)
	return uc.notifier.NotifyMessageSent(
		ctx,
		userIDs,
		chatID.String(),
		message.ID.String(),
		message.UserID.String(),
		message.Content,
	)
}

func (uc *ChatUsecase) BroadcastDirectMessageSent(ctx context.Context, chatID uuid.UUID, message dto.ChatMessageResponseDTO, members []dto.ChatMemberResponseDTO) error {
	userIDs := uc.extractUserIDsFromMemberDTOs(members)
	return uc.notifier.NotifyMessageSent(
		ctx,
		userIDs,
		chatID.String(),
		message.ID.String(),
		message.UserID.String(),
		message.Content,
	)
}

func (uc *ChatUsecase) BroadcastTyping(ctx context.Context, chatID, userID uuid.UUID, typing bool, isGroup bool) error {
	const op = "ChatUsecase.BroadcastTyping"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var userIDs []string

	if isGroup {
		isMember, err := uc.repo.IsGroupChatMember(ctx, chatID, userID)
		if err != nil || !isMember {
			return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
		}
		members, err := uc.repo.GetGroupChatMembers(ctx, chatID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		userIDs = uc.extractUserIDsFromGroupMembersExcept(members, userID.String())
	} else {
		isMember, err := uc.repo.IsDirectChatMember(ctx, chatID, userID)
		if err != nil || !isMember {
			return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
		}
		directChat, err := uc.repo.GetDirectChat(ctx, chatID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		var otherUserID uuid.UUID
		if directChat.User1ID == userID {
			otherUserID = directChat.User2ID
		} else {
			otherUserID = directChat.User1ID
		}
		userIDs = []string{otherUserID.String()}
	}

	if err := uc.notifier.NotifyTyping(ctx, userIDs, chatID.String(), userID.String(), typing); err != nil {
		return fmt.Errorf("%s: %w", op, err)
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
		return fmt.Errorf("%s: %w", op, err)
	}

	userIDs := uc.extractUserIDsFromGroupMembers(members)

	if err := uc.notifier.NotifyUserJoined(ctx, userIDs, chatID.String(), userID.String()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
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
		return fmt.Errorf("%s: %w", op, err)
	}

	userIDs := uc.extractUserIDsFromGroupMembersExcept(members, userID.String())

	if err := uc.notifier.NotifyUserLeft(ctx, userIDs, chatID.String(), userID.String()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user left notification broadcasted")

	return nil
}

func (uc *ChatUsecase) extractUserIDsFromGroupMembers(members []*entities.GroupChatMember) []string {
	var userIDs []string
	for _, member := range members {
		userIDs = append(userIDs, member.UserID.String())
	}
	return userIDs
}

func (uc *ChatUsecase) extractUserIDsFromGroupMembersExcept(members []*entities.GroupChatMember, excludeUserID string) []string {
	var userIDs []string
	for _, member := range members {
		userIDStr := member.UserID.String()
		if userIDStr != excludeUserID {
			userIDs = append(userIDs, userIDStr)
		}
	}
	return userIDs
}

func (uc *ChatUsecase) extractUserIDsFromMemberDTOs(members []dto.ChatMemberResponseDTO) []string {
	var userIDs []string
	for _, member := range members {
		userIDs = append(userIDs, member.UserID.String())
	}
	return userIDs
}
