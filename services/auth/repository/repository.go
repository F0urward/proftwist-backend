package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/google/uuid"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auth.Repository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	const op = "AuthRepository.CreateUser"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":       op,
		"email":    user.Email,
		"username": user.Username,
	})

	err := r.db.QueryRowContext(ctx, queryCreateUser,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
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

	logger.WithField("user_id", user.ID.String()).Debug("successfully created user")
	return user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	const op = "AuthRepository.GetUserByEmail"
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
		logger.Debug("user not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err != nil {
		logger.WithError(err).Error("failed to get user by email")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("user_id", user.ID.String()).Debug("successfully retrieved user")
	return user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	const op = "AuthRepository.GetUserByID"
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
		logger.Debug("user not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err != nil {
		logger.WithError(err).Error("failed to get user by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("successfully retrieved user")
	return user, nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	const op = "AuthRepository.UpdateUser"
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
		user.Role,
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

	logger.Debug("successfully updated user")
	return nil
}

func (r *AuthRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	const op = "AuthRepository.DeleteUser"
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

	logger.Debug("successfully deleted user")
	return nil
}
