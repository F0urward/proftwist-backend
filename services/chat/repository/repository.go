package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/chat"
)

type ChatPostgresRepository struct {
	db *sql.DB
}

func NewChatPostgresRepository(db *sql.DB) chat.Repository {
	return &ChatPostgresRepository{
		db: db,
	}
}

func (r *ChatPostgresRepository) GetGroupChatByNode(ctx context.Context, nodeID string) (*entities.GroupChat, error) {
	const op = "ChatPostgresRepository.GetGroupChatByNode"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("node_id", nodeID)

	var chat entities.GroupChat

	err := r.db.QueryRowContext(ctx, queryGetGroupChatByNode, nodeID).Scan(
		&chat.ID,
		&chat.Title,
		&chat.AvatarURL,
		&chat.RoadmapNodeID,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Debug("no group chat found for node")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).Error("failed to get group chat by node")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("chat_id", chat.ID).Debug("group chat by node retrieved")
	return &chat, nil
}

func (r *ChatPostgresRepository) GetGroupChatsByUser(ctx context.Context, userID uuid.UUID) ([]*entities.GroupChat, error) {
	const op = "ChatPostgresRepository.GetGroupChatsByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	rows, err := r.db.QueryContext(ctx, queryGetGroupChatsByUser, userID)
	if err != nil {
		logger.WithError(err).Error("failed to query user group chats")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var chats []*entities.GroupChat
	for rows.Next() {
		var chat entities.GroupChat

		err := rows.Scan(
			&chat.ID,
			&chat.Title,
			&chat.AvatarURL,
			&chat.RoadmapNodeID,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan group chat row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		chats = append(chats, &chat)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("chats_count", len(chats)).Debug("user group chats retrieved")
	return chats, nil
}

func (r *ChatPostgresRepository) GetGroupChatMembers(ctx context.Context, chatID uuid.UUID) ([]*entities.GroupChatMember, error) {
	const op = "ChatPostgresRepository.GetGroupChatMembers"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", chatID)

	rows, err := r.db.QueryContext(ctx, queryGetGroupChatMembers, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to query group chat members")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var members []*entities.GroupChatMember
	for rows.Next() {
		var member entities.GroupChatMember

		err := rows.Scan(
			&member.ID,
			&member.GroupChatID,
			&member.UserID,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan group chat member row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		members = append(members, &member)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("members_count", len(members)).Debug("group chat members retrieved")
	return members, nil
}

func (r *ChatPostgresRepository) IsGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	const op = "ChatPostgresRepository.IsGroupChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var exists int
	err := r.db.QueryRowContext(ctx, queryIsGroupChatMember, chatID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.WithError(err).Error("failed to check group chat membership")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (r *ChatPostgresRepository) AddGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatPostgresRepository.AddGroupChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})

	_, err := r.db.ExecContext(ctx, queryAddGroupChatMember, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to add group chat member")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("group chat member added successfully")
	return nil
}

func (r *ChatPostgresRepository) RemoveGroupChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	const op = "ChatPostgresRepository.RemoveGroupChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})

	result, err := r.db.ExecContext(ctx, queryRemoveGroupChatMember, chatID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to remove group chat member")
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

	logger.Debug("group chat member removed successfully")
	return nil
}

func (r *ChatPostgresRepository) GeDirectChatsByUser(ctx context.Context, userID uuid.UUID) ([]*entities.DirectChat, error) {
	const op = "ChatPostgresRepository.GeDirectChatsByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	rows, err := r.db.QueryContext(ctx, queryGeDirectChatsByUser, userID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to query user direct chats")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var chats []*entities.DirectChat
	for rows.Next() {
		var chat entities.DirectChat

		err := rows.Scan(
			&chat.ID,
			&chat.User1ID,
			&chat.User2ID,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan direct chat row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		chats = append(chats, &chat)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("chats_count", len(chats)).Debug("user direct chats retrieved")
	return chats, nil
}

func (r *ChatPostgresRepository) GetDirectChat(ctx context.Context, chatID uuid.UUID) (*entities.DirectChat, error) {
	const op = "ChatPostgresRepository.GetDirectChat"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", chatID)

	var chat entities.DirectChat

	err := r.db.QueryRowContext(ctx, queryGetDirectChat, chatID).Scan(
		&chat.ID,
		&chat.User1ID,
		&chat.User2ID,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Error("failed to get direct chat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &chat, nil
}

func (r *ChatPostgresRepository) IsDirectChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	const op = "ChatPostgresRepository.IsDirectChatMember"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var exists int
	err := r.db.QueryRowContext(ctx, queryIsDirectChatMember, chatID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.WithError(err).Error("failed to check direct chat membership")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (r *ChatPostgresRepository) SaveGroupMessage(ctx context.Context, message *entities.Message) error {
	const op = "ChatPostgresRepository.SaveGroupMessage"
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

	_, err = r.db.ExecContext(ctx, querySaveGroupMessage,
		message.ID,
		message.ChatID,
		message.UserID,
		message.Content,
		metadataJSON,
		message.CreatedAt,
		message.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to save group message")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID,
		"chat_id":    message.ChatID,
		"user_id":    message.UserID,
	}).Debug("group message saved successfully")
	return nil
}

func (r *ChatPostgresRepository) SaveDirectMessage(ctx context.Context, message *entities.Message) error {
	const op = "ChatPostgresRepository.SaveDirectMessage"
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

	_, err = r.db.ExecContext(ctx, querySaveDirectMessage,
		message.ID,
		message.ChatID,
		message.UserID,
		message.Content,
		metadataJSON,
		message.CreatedAt,
		message.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to save direct message")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"message_id": message.ID,
		"chat_id":    message.ChatID,
		"user_id":    message.UserID,
	}).Debug("direct message saved successfully")
	return nil
}

func (r *ChatPostgresRepository) GetGroupChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	const op = "ChatPostgresRepository.GetGroupChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", chatID)

	rows, err := r.db.QueryContext(ctx, queryGetGroupChatMessages, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to query group chat messages")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	return r.scanMessages(ctx, rows)
}

func (r *ChatPostgresRepository) GetDirectChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	const op = "ChatPostgresRepository.GetDirectChatMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", chatID)

	rows, err := r.db.QueryContext(ctx, queryGetDirectChatMessages, chatID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to query direct chat messages")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	return r.scanMessages(ctx, rows)
}

func (r *ChatPostgresRepository) scanMessages(ctx context.Context, rows *sql.Rows) ([]*entities.Message, error) {
	const op = "ChatPostgresRepository.scanMessages"
	logger := logctx.GetLogger(ctx).WithField("op", op)

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

	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("messages_count", len(messages)).Debug("messages retrieved")
	return messages, nil
}
