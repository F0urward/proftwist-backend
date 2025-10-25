package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/category"
	"github.com/google/uuid"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) category.Repository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]*entities.Category, error) {
	const op = "CategoryRepository.GetAll"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetAll)
	if err != nil {
		logger.WithError(err).Error("failed to query categories")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.WithError(closeErr).Warn("failed to close rows")
		}
	}()

	categories := []*entities.Category{}

	for rows.Next() {
		category := &entities.Category{}

		if err = rows.Scan(
			&category.CategoryID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			logger.WithError(err).Error("failed to scan category row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("categories_count", len(categories)).Info("successfully retrieved categories")
	return categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, categoryID uuid.UUID) (*entities.Category, error) {
	const op = "CategoryRepository.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	category := &entities.Category{}

	err := r.db.QueryRowContext(ctx, queryGetByID, categoryID).Scan(
		&category.CategoryID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("category not found")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).Error("failed to get category by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully retrieved category")
	return category, nil
}

func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*entities.Category, error) {
	const op = "CategoryRepository.GetByName"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":   op,
		"name": name,
	})

	category := &entities.Category{}

	err := r.db.QueryRowContext(ctx, queryGetByName, name).Scan(
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Debug("category not found by name")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).Error("failed to get category by name")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("successfully retrieved category by name")
	return category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	const op = "CategoryRepository.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": category.CategoryID.String(),
		"name":        category.Name,
	})

	createdCategory := &entities.Category{}

	err := r.db.QueryRowContext(ctx, queryCreate,
		category.Name,
		category.Description,
	).Scan(
		&createdCategory.CategoryID,
		&createdCategory.Name,
		&createdCategory.Description,
		&createdCategory.CreatedAt,
		&createdCategory.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to create category")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully created category")
	return createdCategory, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *entities.Category) error {
	const op = "CategoryRepository.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": category.CategoryID.String(),
	})

	category.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, queryUpdate,
		category.CategoryID,
		category.Name,
		category.Description,
		category.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to update category")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("category not found for update")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("category not found"))
	}

	logger.Info("successfully updated category")
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, categoryID uuid.UUID) error {
	const op = "CategoryRepository.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	result, err := r.db.ExecContext(ctx, queryDelete, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to delete category")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("category not found for deletion")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("category not found"))
	}

	logger.Info("successfully deleted category")
	return nil
}
