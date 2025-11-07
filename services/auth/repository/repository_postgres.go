package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/auth"
)

type AuthPostgresRepository struct {
	db *sql.DB
}

func NewAuthPostgresRepository(db *sql.DB) auth.PostgresRepository {
	return &AuthPostgresRepository{db: db}
}

func (r *AuthPostgresRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	const op = "AuthPostgresRepository.CreateUser"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":       op,
		"email":    user.Email,
		"username": user.Username,
	})

	err := r.db.QueryRowContext(ctx, queryCreateUser,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.AvatarUrl,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to create user")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("user_id", user.ID.String()).Info("successfully created user")
	return user, nil
}

func (r *AuthPostgresRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	const op = "AuthPostgresRepository.GetUserByEmail"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":    op,
		"email": email,
	})

	user := &entities.User{}

	err := r.db.QueryRowContext(ctx, queryGetUserByEmail, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("user not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err != nil {
		logger.WithError(err).Error("failed to get user by email")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("user_id", user.ID.String()).Info("successfully retrieved user")
	return user, nil
}

func (r *AuthPostgresRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	const op = "AuthPostgresRepository.GetUserByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	user := &entities.User{}

	err := r.db.QueryRowContext(ctx, queryGetUserByID, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("user not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err != nil {
		logger.WithError(err).Error("failed to get user by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully retrieved user")
	return user, nil
}

func (r *AuthPostgresRepository) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entities.User, error) {
	const op = "AuthPostgresRepository.GetUsersByIDs"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"user_ids":   userIDs,
		"user_count": len(userIDs),
	})

	if len(userIDs) == 0 {
		logger.Info("empty user IDs list provided")
		return []*entities.User{}, nil
	}

	rows, err := r.db.QueryContext(ctx, queryGetUsersByIDs, pq.Array(userIDs))
	if err != nil {
		logger.WithError(err).Error("failed to query users by IDs")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.WithError(closeErr).Warn("failed to close rows")
		}
	}()

	users := []*entities.User{}

	for rows.Next() {
		user := &entities.User{}

		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.AvatarUrl,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			logger.WithError(err).Error("failed to scan user row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("found_count", len(users)).Info("successfully retrieved users by IDs")
	return users, nil
}

func (r *AuthPostgresRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	const op = "AuthPostgresRepository.UpdateUser"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": user.ID.String(),
	})

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, queryUpdateUser,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.AvatarUrl,
		user.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to update user")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("user not found for update")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.Info("successfully updated user")
	return nil
}

func (r *AuthPostgresRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	const op = "AuthPostgresRepository.DeleteUser"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	result, err := r.db.ExecContext(ctx, queryDeleteUser, userID)
	if err != nil {
		logger.WithError(err).Error("failed to delete user")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("user not found for deletion")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.Info("successfully deleted user")
	return nil
}

func (r *AuthPostgresRepository) CreateVKUser(ctx context.Context, vkUser *entities.VKUser) error {
	const op = "AuthPostgresRepository.CreateVKUser"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"user_id":    vkUser.UserID.String(),
		"vk_user_id": vkUser.VKUserID,
	})

	err := r.db.QueryRowContext(ctx, queryCreateVKUser,
		vkUser.UserID,
		vkUser.VKUserID,
		vkUser.AccessToken,
		vkUser.RefreshToken,
		vkUser.ExpiresAt,
		vkUser.DeviceID,
	).Scan(
		&vkUser.ID,
		&vkUser.CreatedAt,
		&vkUser.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to create vk user")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully created vk user")
	return nil
}

func (r *AuthPostgresRepository) GetVKUserByUserID(ctx context.Context, userID uuid.UUID) (*entities.VKUser, error) {
	const op = "AuthPostgresRepository.GetVKUserByUserID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	vkUser := &entities.VKUser{}

	err := r.db.QueryRowContext(ctx, queryGetVKUserByUserID, userID).Scan(
		&vkUser.ID,
		&vkUser.UserID,
		&vkUser.VKUserID,
		&vkUser.AccessToken,
		&vkUser.RefreshToken,
		&vkUser.ExpiresAt,
		&vkUser.DeviceID,
		&vkUser.CreatedAt,
		&vkUser.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("vk user not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err != nil {
		logger.WithError(err).Error("failed to get vk user by user id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully retrieved vk user")
	return vkUser, nil
}

func (r *AuthPostgresRepository) GetVKUserByID(ctx context.Context, vkUserID int64) (*entities.VKUser, error) {
	const op = "AuthPostgresRepository.GetVKUserByVKUserID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"vk_user_id": vkUserID,
	})

	vkUser := &entities.VKUser{}

	err := r.db.QueryRowContext(ctx, queryGetVKUserByVKUserID, vkUserID).Scan(
		&vkUser.ID,
		&vkUser.UserID,
		&vkUser.VKUserID,
		&vkUser.AccessToken,
		&vkUser.RefreshToken,
		&vkUser.ExpiresAt,
		&vkUser.DeviceID,
		&vkUser.CreatedAt,
		&vkUser.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("vk user not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err != nil {
		logger.WithError(err).Error("failed to get vk user by vk user id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully retrieved vk user")
	return vkUser, nil
}

func (r *AuthPostgresRepository) UpdateVKUser(ctx context.Context, vkUser *entities.VKUser) error {
	const op = "AuthPostgresRepository.UpdateVKUser"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": vkUser.UserID.String(),
	})

	result, err := r.db.ExecContext(ctx, queryUpdateVKUser,
		vkUser.ID,
		vkUser.VKUserID,
		vkUser.AccessToken,
		vkUser.RefreshToken,
		vkUser.ExpiresAt,
		vkUser.DeviceID,
	)
	if err != nil {
		logger.WithError(err).Error("failed to update vk user")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("vk user not found for update")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.Info("successfully updated vk user")
	return nil
}

func (r *AuthPostgresRepository) DeleteVKUser(ctx context.Context, userID uuid.UUID) error {
	const op = "AuthPostgresRepository.DeleteVKUser"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	result, err := r.db.ExecContext(ctx, queryDeleteVKUser, userID)
	if err != nil {
		logger.WithError(err).Error("failed to delete vk user")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("vk user not found for deletion")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.Info("successfully deleted vk user")
	return nil
}
