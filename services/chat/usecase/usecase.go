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

func (uc *ChatUsecase) CreateChat(ctx context.Context, req dto.CreateChatRequestDTO) (*dto.ChatResponseDTO, error) {
	const op = "ChatUsecase.CreateChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// if req.Type == "group" {
	// 	isAdmin := req.CreatedByRole == "admin"
	// 	if !isAdmin {
	// 		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	// 	}
	// }

	chat, err := dto.CreateChatRequestToEntity(&req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := uc.repo.CreateChat(ctx, chat); err != nil {
		logger.WithError(err).Error("failed to create chat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	switch chat.Type {
	case entities.ChatTypeDirect:
		participants := append([]uuid.UUID{req.CreatedByID}, req.MemberIDs...)
		for _, participantID := range participants {
			if err := uc.repo.AddChatMember(ctx, chat.ID, participantID, entities.MemberRoleMember); err != nil {
				logger.WithError(err).Error("failed to add participant to direct chat")
				if delErr := uc.repo.DeleteChat(ctx, chat.ID); delErr != nil {
					logger.WithError(delErr).Error("failed to cleanup chat after member addition failure")
				}
				return nil, fmt.Errorf("%s: failed to add participants to direct chat: %w", op, err)
			}
		}

	case entities.ChatTypeGroup:
		if err := uc.repo.AddChatMember(ctx, chat.ID, req.CreatedByID, entities.MemberRoleAdmin); err != nil {
			logger.WithError(err).Error("failed to add chat creator as admin")
			if delErr := uc.repo.DeleteChat(ctx, chat.ID); delErr != nil {
				logger.WithError(delErr).Error("failed to cleanup chat after admin addition failure")
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		for _, memberID := range req.MemberIDs {
			if memberID != req.CreatedByID {
				if err := uc.repo.AddChatMember(ctx, chat.ID, memberID, entities.MemberRoleMember); err != nil {
					logger.WithError(err).Warn("failed to add member to group chat")
				}
			}
		}
	}

	logger.WithField("chat_id", chat.ID.String()).Info("chat created successfully")

	chatDTO := dto.ChatToDTO(chat)
	return &chatDTO, nil
}

func (uc *ChatUsecase) SendMessage(ctx context.Context, req dto.SendMessageRequestDTO) (*dto.ChatMessageResponseDTO, error) {
	const op = "ChatUsecase.SendMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetChat(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if chat == nil {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	members, err := uc.repo.GetChatMembers(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !chat.CanSendMessage(req.UserID, members) {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	message := &entities.Message{
		ChatID:   req.ChatID,
		UserID:   req.UserID,
		Content:  req.Content,
		Metadata: req.Metadata,
	}

	if err := uc.repo.SaveMessage(ctx, message); err != nil {
		logger.WithError(err).Error("failed to save message")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	messageDTO := dto.MessageToDTO(message)

	if err := uc.BroadcastMessageSent(ctx, req.ChatID, messageDTO, members); err != nil {
		logger.WithError(err).Warn("failed to broadcast message")
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID.String(),
		"chat_id":    message.ChatID.String(),
		"user_id":    message.UserID.String(),
	}).Info("message sent and broadcasted successfully")

	return &messageDTO, nil
}

func (uc *ChatUsecase) GetChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error) {
	const op = "ChatUsecase.GetChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetChat(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if chat == nil {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !chat.CanViewChat(userID, members) {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	messages, err := uc.repo.GetChatMessages(ctx, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get chat messages")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.GetChatMessagesResponseToDTO(messages)
	return &response, nil
}

func (uc *ChatUsecase) AddMember(ctx context.Context, req dto.AddMemberRequestDTO) error {
	const op = "ChatUsecase.AddMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetChat(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat")
		return fmt.Errorf("%s: %w", op, err)
	}
	if chat == nil {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	members, err := uc.repo.GetChatMembers(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return fmt.Errorf("%s: %w", op, err)
	}

	if req.UserID == req.RequestedBy {
		return fmt.Errorf("%s: %s", op, "cannot add yourself to chat")
	}

	if !chat.CanManageChat(req.RequestedBy, members) {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	if !chat.CanAddMember(req.RequestedBy, members) {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	for _, member := range members {
		if member.UserID == req.UserID {
			return fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
		}
	}

	if err := uc.repo.AddChatMember(ctx, req.ChatID, req.UserID, entities.MemberRole(req.Role)); err != nil {
		logger.WithError(err).Error("failed to add chat member")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": req.ChatID.String(),
		"user_id": req.UserID.String(),
	}).Info("member added to chat")

	return nil
}

func (uc *ChatUsecase) RemoveMember(ctx context.Context, req dto.RemoveMemberRequestDTO) error {
	const op = "ChatUsecase.RemoveMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetChat(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat")
		return fmt.Errorf("%s: %w", op, err)
	}
	if chat == nil {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	members, err := uc.repo.GetChatMembers(ctx, req.ChatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return fmt.Errorf("%s: %w", op, err)
	}

	if req.UserID == req.RequestedBy {
		return fmt.Errorf("%s: %s", op, "cannot remove yourself from chat")
	}

	if !chat.CanManageChat(req.RequestedBy, members) {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	if !chat.CanRemoveMember(req.RequestedBy, members) {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	if !chat.IsUserMember(req.UserID, members) {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err := uc.repo.RemoveChatMember(ctx, req.ChatID, req.UserID); err != nil {
		logger.WithError(err).Error("failed to remove chat member")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": req.ChatID.String(),
		"user_id": req.UserID.String(),
	}).Info("member removed from chat")

	return nil
}

func (uc *ChatUsecase) GetUserChats(ctx context.Context, userID uuid.UUID) ([]dto.ChatResponseDTO, error) {
	const op = "ChatUsecase.GetUserChats"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chats, err := uc.repo.GetUserChats(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user chats")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatDTOs := dto.ChatListToDTO(chats)
	return chatDTOs, nil
}

func (uc *ChatUsecase) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]dto.ChatMemberResponseDTO, error) {
	const op = "ChatUsecase.GetChatMembers"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	memberDTOs := dto.ChatMemberListToDTO(members)
	return memberDTOs, nil
}

func (uc *ChatUsecase) JoinGroupChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatUsecase.JoinGroupChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetChat(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if chat == nil {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if chat.Type != entities.ChatTypeGroup {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	isMember, err := uc.repo.IsChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check chat membership")
		return fmt.Errorf("%s: %w", op, err)
	}
	if isMember {
		return fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
	}

	if err := uc.repo.AddChatMember(ctx, chatID, userID, entities.MemberRoleMember); err != nil {
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

	chat, err := uc.repo.GetChat(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if chat == nil {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if chat.Type != entities.ChatTypeGroup {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return fmt.Errorf("%s: %w", op, err)
	}

	if !chat.IsUserMember(userID, members) {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	if chat.CanManageChat(userID, members) {
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	if err := uc.repo.RemoveChatMember(ctx, chatID, userID); err != nil {
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

func (uc *ChatUsecase) BroadcastTyping(ctx context.Context, chatID, userID uuid.UUID, typing bool) error {
	const op = "ChatUsecase.BroadcastTyping"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userIDs := uc.extractUserIDsExcept(members, userID.String())

	if err := uc.notifier.NotifyTyping(ctx, userIDs, chatID.String(), userID.String(), typing); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
		"typing":  typing,
	}).Info("typing notification broadcasted")

	return nil
}

func (uc *ChatUsecase) BroadcastUserJoined(ctx context.Context, chatID, userID uuid.UUID) error {
	const op = "ChatUsecase.BroadcastUserJoined"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userIDs := uc.extractUserIDs(members)

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

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userIDs := uc.extractUserIDsExcept(members, userID.String())

	if err := uc.notifier.NotifyUserLeft(ctx, userIDs, chatID.String(), userID.String()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID.String(),
		"user_id": userID.String(),
	}).Info("user left notification broadcasted")

	return nil
}

func (uc *ChatUsecase) BroadcastMessageSent(ctx context.Context, chatID uuid.UUID, message dto.ChatMessageResponseDTO, members []*entities.ChatMember) error {
	userIDs := uc.extractUserIDs(members)

	return uc.notifier.NotifyMessageSent(
		ctx,
		userIDs,
		chatID.String(),
		message.ID.String(),
		message.UserID.String(),
		message.Content,
	)
}

func (uc *ChatUsecase) extractUserIDs(members []*entities.ChatMember) []string {
	var userIDs []string
	for _, member := range members {
		userIDs = append(userIDs, member.UserID.String())
	}
	return userIDs
}

func (uc *ChatUsecase) extractUserIDsExcept(members []*entities.ChatMember, excludeUserID string) []string {
	var userIDs []string
	for _, member := range members {
		userIDStr := member.UserID.String()
		if userIDStr != excludeUserID {
			userIDs = append(userIDs, userIDStr)
		}
	}
	return userIDs
}
