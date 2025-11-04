package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/google/uuid"
)

type ChatPostgresRepository struct {
	db *sql.DB
}

func NewChatPostgresRepository(db *sql.DB) chat.Repository {
	return &ChatPostgresRepository{
		db: db,
	}
}

func (r *ChatPostgresRepository) CreateChat(ctx context.Context, chat *entities.Chat) error {
	const op = "ChatPostgresRepository.CreateChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if chat.ID == uuid.Nil {
		chat.ID = uuid.New()
	}
	if chat.CreatedAt.IsZero() {
		chat.CreatedAt = time.Now()
	}
	if chat.UpdatedAt.IsZero() {
		chat.UpdatedAt = time.Now()
	}

	_, err := r.db.ExecContext(ctx, queryCreateChat,
		chat.ID,
		chat.Type,
		chat.Title,
		chat.Description,
		chat.AvatarURL,
		chat.CreatedBy,
		chat.CreatedAt,
		chat.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to create chat")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("chat_id", chat.ID).Info("chat created successfully")
	return nil
}

func (r *ChatPostgresRepository) GetChat(ctx context.Context, chatID uuid.UUID) (*entities.Chat, error) {
	const op = "ChatPostgresRepository.GetChat"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", chatID)

	var chat entities.Chat

	err := r.db.QueryRowContext(ctx, queryGetChat, chatID).Scan(
		&chat.ID,
		&chat.Type,
		&chat.Title,
		&chat.Description,
		&chat.AvatarURL,
		&chat.CreatedBy,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Error("failed to get chat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &chat, nil
}

func (r *ChatPostgresRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*entities.Chat, error) {
	const op = "ChatPostgresRepository.GetUserChats"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	rows, err := r.db.QueryContext(ctx, queryGetUserChats, userID)
	if err != nil {
		logger.WithError(err).Error("failed to query user chats")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var chats []*entities.Chat
	for rows.Next() {
		var chat entities.Chat

		err := rows.Scan(
			&chat.ID,
			&chat.Type,
			&chat.Title,
			&chat.Description,
			&chat.AvatarURL,
			&chat.CreatedBy,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan chat row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		chats = append(chats, &chat)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("chats_count", len(chats)).Debug("user chats retrieved")
	return chats, nil
}

func (r *ChatPostgresRepository) SaveMessage(ctx context.Context, message *entities.Message) error {
	const op = "ChatPostgresRepository.SaveMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}
	if message.UpdatedAt.IsZero() {
		message.UpdatedAt = time.Now()
	}

	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		logger.WithError(err).Error("failed to marshal metadata")
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.ExecContext(ctx, querySaveMessage,
		message.ID,
		message.ChatID,
		message.UserID,
		message.Content,
		metadataJSON,
		message.CreatedAt,
		message.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to save message")
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.ExecContext(ctx, queryUpdateChatTimestamp, time.Now(), message.ChatID)
	if err != nil {
		logger.WithError(err).Warn("failed to update chat timestamp")
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID,
		"chat_id":    message.ChatID,
		"user_id":    message.UserID,
	}).Debug("message saved successfully")
	return nil
}

func (r *ChatPostgresRepository) GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	const op = "ChatPostgresRepository.GetChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", chatID)

	rows, err := r.db.QueryContext(ctx, queryGetChatMessages, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to query chat messages")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var messages []*entities.Message
	for rows.Next() {
		var message entities.Message
		var metadataJSON []byte

		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.UserID,
			&message.Content,
			&metadataJSON,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan message row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
				logger.WithError(err).Warn("failed to unmarshal message metadata")
				message.Metadata = make(map[string]interface{})
			}
		} else {
			message.Metadata = make(map[string]interface{})
		}

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("messages_count", len(messages)).Debug("chat messages retrieved")
	return messages, nil
}

func (r *ChatPostgresRepository) AddChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, role entities.MemberRole) error {
	const op = "ChatPostgresRepository.AddChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
		"role":    role,
	})

	_, err := r.db.ExecContext(ctx, queryAddChatMember,
		chatID,
		userID,
		role,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		logger.WithError(err).Error("failed to add chat member")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("chat member added successfully")
	return nil
}

func (r *ChatPostgresRepository) RemoveChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatPostgresRepository.RemoveChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})

	result, err := r.db.ExecContext(ctx, queryRemoveChatMember, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to remove chat member")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.Debug("chat member removed successfully")
	return nil
}

func (r *ChatPostgresRepository) IsChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	const op = "ChatPostgresRepository.IsChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var exists int
	err := r.db.QueryRowContext(ctx, queryIsChatMember, chatID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.WithError(err).Error("failed to check chat membership")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (r *ChatPostgresRepository) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*entities.ChatMember, error) {
	const op = "ChatPostgresRepository.GetChatMembers"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", chatID)

	rows, err := r.db.QueryContext(ctx, queryGetChatMembers, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to query chat members")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var members []*entities.ChatMember
	for rows.Next() {
		var member entities.ChatMember

		err := rows.Scan(
			&member.ID,
			&member.ChatID,
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
			&member.LastRead,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan chat member row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		members = append(members, &member)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return members, nil
}

func (r *ChatPostgresRepository) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	const op = "ChatPostgresRepository.DeleteChat"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	_, err := r.db.ExecContext(ctx, queryDeleteChatMesseges, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to delete chat messages")
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.ExecContext(ctx, queryDeleteChatMembers, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to delete chat members")
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.db.ExecContext(ctx, queryDeleteChat, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to delete chat")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.WithField("chat_id", chatID).Info("chat deleted successfully")
	return nil
}
