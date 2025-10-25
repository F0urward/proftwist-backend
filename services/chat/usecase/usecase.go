package usecase

import (
	"context"
	"fmt"
	"github.com/F0urward/proftwist-backend/services/chat/dto"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/google/uuid"
)

type ChatUseCase struct {
	repo chat.Repository
}

func NewChatUseCase(repo chat.Repository) *ChatUseCase {
	return &ChatUseCase{
		repo: repo,
	}
}

func (uc *ChatUseCase) CreateChat(ctx context.Context, req dto.CreateChatRequest) (*entities.Chat, error) {
	const op = "ChatUseCase.CreateChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	// Валидация в зависимости от типа чата
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chat := &entities.Chat{
		ID:          uuid.New(),
		Type:        entities.ChatType(req.Type),
		Title:       req.Title,
		Description: req.Description,
		AvatarURL:   req.AvatarURL,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.CreateChat(ctx, chat); err != nil {
		logger.WithError(err).Error("failed to create chat")
		return nil, err
	}

	// Добавляем создателя как владельца
	if err := uc.repo.AddChatMember(ctx, chat.ID, req.CreatedBy, entities.MemberRoleOwner); err != nil {
		logger.WithError(err).Error("failed to add chat creator as member")
		return nil, err
	}

	// Добавляем участников в зависимости от типа чата
	switch chat.Type {
	case entities.ChatTypeDirect:
		// Для директа добавляем второго участника
		if len(req.MemberIDs) == 1 {
			memberID := req.MemberIDs[0]
			if err := uc.repo.AddChatMember(ctx, chat.ID, memberID, entities.MemberRoleMember); err != nil {
				logger.WithError(err).Error("failed to add member to direct chat")
				return nil, err
			}
		}
	case entities.ChatTypeGroup:
		// Для группы добавляем всех указанных участников
		for _, memberID := range req.MemberIDs {
			if memberID != req.CreatedBy {
				if err := uc.repo.AddChatMember(ctx, chat.ID, memberID, entities.MemberRoleMember); err != nil {
					logger.WithError(err).Warn("failed to add member to group chat")
				}
			}
		}
	case entities.ChatTypeChannel:
		// Для канала не добавляем участников при создании
		logger.Debug("channel created without additional members")
	}

	logger.WithField("chat_id", chat.ID).Info("chat created successfully")
	return chat, nil
}

func (uc *ChatUseCase) SendMessage(ctx context.Context, req dto.SendMessageRequest) (*entities.Message, error) {
	const op = "ChatUseCase.SendMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	// Получаем информацию о чате и участниках
	chatWithMembers, err := uc.GetChatWithMembers(ctx, req.ChatID, req.UserID)
	if err != nil {
		return nil, err
	}

	// Проверяем права на отправку сообщения
	if !chatWithMembers.Chat.CanSendMessage(req.UserID, chatWithMembers.Members) {
		return nil, errs.ErrForbidden
	}

	message := &entities.Message{
		ID:        uuid.New(),
		ChatID:    req.ChatID,
		UserID:    req.UserID,
		Content:   req.Content,
		Type:      entities.MessageTypeChat,
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.SaveMessage(ctx, message); err != nil {
		logger.WithError(err).Error("failed to save message")
		return nil, err
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID,
		"chat_id":    message.ChatID,
		"user_id":    message.UserID,
	}).Info("message sent successfully")

	return message, nil
}

func (uc *ChatUseCase) GetChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	const op = "ChatUseCase.GetChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check chat membership")
		return nil, err
	}
	if !isMember {
		return nil, errs.ErrForbidden
	}

	messages, err := uc.repo.GetChatMessages(ctx, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get chat messages")
		return nil, err
	}

	return messages, nil
}

func (uc *ChatUseCase) AddMember(ctx context.Context, req dto.AddMemberRequest) error {
	const op = "ChatUseCase.AddMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	// Получаем информацию о чате
	chatWithMembers, err := uc.GetChatWithMembers(ctx, req.ChatID, req.RequestedBy)
	if err != nil {
		return err
	}

	// Проверяем права на добавление участника
	if !chatWithMembers.Chat.CanManageChat(req.RequestedBy, chatWithMembers.Members) {
		return errs.ErrForbidden
	}

	// Проверяем возможность добавления в зависимости от типа чата
	newMemberIDs := []uuid.UUID{req.UserID}
	if !chatWithMembers.Chat.CanAddMembers(req.RequestedBy, chatWithMembers.Members, newMemberIDs) {
		return errs.ErrForbidden
	}

	isMember, err := uc.repo.IsChatMember(ctx, req.ChatID, req.UserID)
	if err != nil {
		logger.WithError(err).Error("failed to check existing membership")
		return err
	}
	if isMember {
		return errs.ErrAlreadyExists
	}

	if err := uc.repo.AddChatMember(ctx, req.ChatID, req.UserID, entities.MemberRole(req.Role)); err != nil {
		logger.WithError(err).Error("failed to add chat member")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": req.ChatID,
		"user_id": req.UserID,
	}).Info("member added to chat")

	return nil
}

func (uc *ChatUseCase) RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error {
	const op = "ChatUseCase.RemoveMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	// Получаем информацию о чате
	chatWithMembers, err := uc.GetChatWithMembers(ctx, req.ChatID, req.RequestedBy)
	if err != nil {
		return err
	}

	// Проверяем права на удаление участника
	if !chatWithMembers.Chat.CanRemoveMember(req.RequestedBy, req.UserID, chatWithMembers.Members) {
		return errs.ErrForbidden
	}

	if err := uc.repo.RemoveChatMember(ctx, req.ChatID, req.UserID); err != nil {
		logger.WithError(err).Error("failed to remove chat member")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": req.ChatID,
		"user_id": req.UserID,
	}).Info("member removed from chat")

	return nil
}

func (uc *ChatUseCase) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*entities.Chat, error) {
	const op = "ChatUseCase.GetUserChats"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chats, err := uc.repo.GetUserChats(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user chats")
		return nil, err
	}

	return chats, nil
}

func (uc *ChatUseCase) GetChatWithMembers(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (*entities.ChatWithMembers, error) {
	const op = "ChatUseCase.GetChatWithMembers"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check chat membership")
		return nil, err
	}
	if !isMember {
		return nil, errs.ErrForbidden
	}

	chat, err := uc.repo.GetChat(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat")
		return nil, err
	}
	if chat == nil {
		return nil, errs.ErrNotFound
	}

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return nil, err
	}

	return &entities.ChatWithMembers{
		Chat:    chat,
		Members: members,
	}, nil
}

func (uc *ChatUseCase) canManageChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	chat, err := uc.repo.GetChat(ctx, chatID)
	if err != nil {
		return false, err
	}
	if chat == nil {
		return false, errs.ErrNotFound
	}

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		return false, err
	}

	return chat.CanManageChat(userID, members), nil
}

func (uc *ChatUseCase) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*entities.ChatMember, error) {
	const op = "ChatUseCase.GetChatMembers"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	members, err := uc.repo.GetChatMembers(ctx, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to get chat members")
		return nil, err
	}

	return members, nil
}

func (uc *ChatUseCase) IsChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	const op = "ChatUseCase.IsChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	isMember, err := uc.repo.IsChatMember(ctx, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check chat membership")
		return false, err
	}

	return isMember, nil
}

func (uc *ChatUseCase) JoinChannel(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatUseCase.JoinChannel"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chat, err := uc.repo.GetChat(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return errs.ErrNotFound
	}

	if chat.Type != entities.ChatTypeChannel {
		return errs.ErrForbidden
	}

	if err := uc.repo.AddChatMember(ctx, chatID, userID, entities.MemberRoleMember); err != nil {
		logger.WithError(err).Error("failed to join channel")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("user joined channel")

	return nil
}

func (uc *ChatUseCase) DeleteChat(ctx context.Context, req dto.DeleteChatRequest) error {
	const op = "ChatUseCase.DeleteChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	// Получаем информацию о чате
	chatWithMembers, err := uc.GetChatWithMembers(ctx, req.ChatID, req.RequestedBy)
	if err != nil {
		return err
	}

	// Проверяем права на удаление чата
	if !chatWithMembers.Chat.CanDeleteChat(req.RequestedBy, chatWithMembers.Members) {
		return errs.ErrForbidden
	}

	if err := uc.repo.DeleteChat(ctx, req.ChatID); err != nil {
		logger.WithError(err).Error("failed to delete chat")
		return err
	}

	logger.WithField("chat_id", req.ChatID).Info("chat deleted successfully")
	return nil
}
