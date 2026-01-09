package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/friend"
)

type FriendPostgresRepository struct {
	db *sql.DB
}

func NewFriendPostgresRepository(db *sql.DB) friend.Repository {
	return &FriendPostgresRepository{
		db: db,
	}
}

func (r *FriendPostgresRepository) CreateFriendship(ctx context.Context, userID, friendID, chatID uuid.UUID) error {
	const op = "FriendRepository.CreateFriendship"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
		"chat_id":   chatID,
	})

	_, err := r.db.ExecContext(ctx, queryCreateFriendship, userID, friendID, chatID)
	if err != nil {
		logger.WithError(err).Error("failed to create friendship")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("friendship created successfully")
	return nil
}

func (r *FriendPostgresRepository) DeleteFriendship(ctx context.Context, userID, friendID uuid.UUID) error {
	const op = "FriendRepository.DeleteFriendship"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	})

	result, err := r.db.ExecContext(ctx, queryDeleteFriendship, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("failed to delete friendship")
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

	logger.Info("friendship deleted successfully")
	return nil
}

func (r *FriendPostgresRepository) GetFriendIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	const op = "FriendRepository.GetFriendIDs"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	rows, err := r.db.QueryContext(ctx, queryGetFriendIDs, userID)
	if err != nil {
		logger.WithError(err).Error("failed to query friend IDs")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var friendIDs []uuid.UUID
	for rows.Next() {
		var friendID uuid.UUID
		err := rows.Scan(&friendID)
		if err != nil {
			logger.WithError(err).Error("failed to scan friend ID row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		friendIDs = append(friendIDs, friendID)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("friend_count", len(friendIDs)).Info("friend IDs retrieved")
	return friendIDs, nil
}

func (r *FriendPostgresRepository) IsFriends(ctx context.Context, userID, friendID uuid.UUID) (bool, error) {
	const op = "FriendRepository.IsFriends"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	var exists int
	err := r.db.QueryRowContext(ctx, queryIsFriends, userID, friendID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.WithError(err).Error("failed to check friendship")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (r *FriendPostgresRepository) GetFriendshipChatID(ctx context.Context, userID, friendID uuid.UUID) (*uuid.UUID, error) {
	const op = "FriendRepository.GetFriendshipChatID"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	})

	var chatID *uuid.UUID
	err := r.db.QueryRowContext(ctx, queryGetFriendshipChatID, userID, friendID).Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("friendship not found")
			return nil, nil
		}
		logger.WithError(err).Error("failed to get friendship chat ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("chat_id", chatID).Info("friendship chat ID retrieved")
	return chatID, nil
}

func (r *FriendPostgresRepository) CreateFriendRequest(ctx context.Context, request *entities.FriendRequest) error {
	const op = "FriendRepository.CreateFriendRequest"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"from_user_id": request.FromUserID,
		"to_user_id":   request.ToUserID,
	})

	err := r.db.QueryRowContext(ctx, queryCreateFriendRequest,
		request.FromUserID,
		request.ToUserID,
		request.Message,
	).Scan(
		&request.ID,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to create friend request")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("request_id", request.ID.String()).Info("successfully created friend request")
	return nil
}

func (r *FriendPostgresRepository) DeleteFriendRequest(ctx context.Context, requestID uuid.UUID) error {
	const op = "FriendRepository.DeleteFriendRequest"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithField("request_id", requestID)

	result, err := r.db.ExecContext(ctx, queryDeleteFriendRequest, requestID)
	if err != nil {
		logger.WithError(err).Error("failed to delete friend request")
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

	logger.Info("friend request deleted successfully")
	return nil
}

func (r *FriendPostgresRepository) GetFriendRequestByID(ctx context.Context, requestID uuid.UUID) (*entities.FriendRequest, error) {
	const op = "FriendRepository.GetFriendRequestByID"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithField("request_id", requestID)

	var request entities.FriendRequest

	err := r.db.QueryRowContext(ctx, queryGetFriendRequestByID, requestID).Scan(
		&request.ID,
		&request.FromUserID,
		&request.ToUserID,
		&request.Status,
		&request.Message,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Error("failed to get friend request by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &request, nil
}

func (r *FriendPostgresRepository) GetFriendRequestsForUserByStatus(ctx context.Context, userID uuid.UUID, statuses []entities.FriendStatus) ([]*entities.FriendRequest, error) {
	const op = "FriendRepository.GetFriendRequestsForUserByStatus"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	placeholders := make([]string, len(statuses))
	args := make([]interface{}, len(statuses)+1)
	args[0] = userID

	for i, status := range statuses {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = status
	}

	query := fmt.Sprintf(`
		SELECT id, from_user_id, to_user_id, status, message, created_at, updated_at
		FROM friend_requests 
		WHERE to_user_id = $1 AND status IN (%s)`,
		strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.WithError(err).Error("failed to query friend requests for user by status")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	return r.scanFriendRequests(ctx, rows)
}

func (r *FriendPostgresRepository) GetSentFriendRequestsByStatus(ctx context.Context, userID uuid.UUID, statuses []entities.FriendStatus) ([]*entities.FriendRequest, error) {
	const op = "FriendRepository.GetSentFriendRequestsByStatus"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	placeholders := make([]string, len(statuses))
	args := make([]interface{}, len(statuses)+1)
	args[0] = userID

	for i, status := range statuses {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = status
	}

	query := fmt.Sprintf(`
		SELECT id, from_user_id, to_user_id, status, message, created_at, updated_at
		FROM friend_requests 
		WHERE from_user_id = $1 AND status IN (%s)`,
		strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.WithError(err).Error("failed to query sent friend requests by status")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	return r.scanFriendRequests(ctx, rows)
}

func (r *FriendPostgresRepository) UpdateFriendRequestStatus(ctx context.Context, requestID uuid.UUID, status entities.FriendStatus) error {
	const op = "FriendRepository.UpdateFriendRequestStatus"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"request_id": requestID,
		"status":     status,
	})

	result, err := r.db.ExecContext(ctx, queryUpdateFriendRequestStatus, status, requestID)
	if err != nil {
		logger.WithError(err).Error("failed to update friend request status")
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

	logger.Info("friend request status updated successfully")
	return nil
}

func (r *FriendPostgresRepository) UpdateFriendRequest(ctx context.Context, requestID uuid.UUID, fromUserID, toUserID uuid.UUID, status entities.FriendStatus) error {
	const op = "FriendRepository.UpdateFriendRequest"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"request_id":   requestID,
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
		"status":       status,
	})

	result, err := r.db.ExecContext(ctx, queryUpdateFriendRequest, fromUserID, toUserID, status, requestID)
	if err != nil {
		logger.WithError(err).Error("failed to update friend request")
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

	logger.Info("friend request updated successfully")
	return nil
}

func (r *FriendPostgresRepository) GetFriendRequestBetweenUsers(ctx context.Context, fromUserID, toUserID uuid.UUID) (*entities.FriendRequest, error) {
	const op = "FriendRepository.GetFriendRequestBetweenUsers"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
	})

	var request entities.FriendRequest

	err := r.db.QueryRowContext(ctx, queryGetFriendRequestBetweenUsers, fromUserID, toUserID).Scan(
		&request.ID,
		&request.FromUserID,
		&request.ToUserID,
		&request.Status,
		&request.Message,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Error("failed to get friend request between users")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &request, nil
}

func (r *FriendPostgresRepository) GetPendingFriendRequestBetweenUsers(ctx context.Context, fromUserID, toUserID uuid.UUID) (*entities.FriendRequest, error) {
	const op = "FriendRepository.GetPendingFriendRequestBetweenUsers"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithFields(map[string]interface{}{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
	})

	var request entities.FriendRequest

	err := r.db.QueryRowContext(ctx, queryGetPendingFriendRequestBetweenUsers, fromUserID, toUserID).Scan(
		&request.ID,
		&request.FromUserID,
		&request.ToUserID,
		&request.Status,
		&request.Message,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Error("failed to get pending friend request between users")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &request, nil
}

func (r *FriendPostgresRepository) scanFriendRequests(ctx context.Context, rows *sql.Rows) ([]*entities.FriendRequest, error) {
	const op = "FriendRepository.scanFriendRequests"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	var requests []*entities.FriendRequest
	for rows.Next() {
		var request entities.FriendRequest

		err := rows.Scan(
			&request.ID,
			&request.FromUserID,
			&request.ToUserID,
			&request.Status,
			&request.Message,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan friend request row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		requests = append(requests, &request)
	}

	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("requests_count", len(requests)).Info("friend requests retrieved")
	return requests, nil
}
