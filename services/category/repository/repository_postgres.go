package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/category"
)

type CategoryPostgresRepository struct {
	db *sql.DB
}

func NewCategoryPostgresRepository(db *sql.DB) category.Repository {
	return &CategoryPostgresRepository{db: db}
}

func (r *CategoryPostgresRepository) GetAll(ctx context.Context) ([]*entities.Category, error) {
	const op = "CategoryRepository.GetAll"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

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

func (r *CategoryPostgresRepository) GetByID(ctx context.Context, categoryID uuid.UUID) (*entities.Category, error) {
	const op = "CategoryRepository.GetByID"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *CategoryPostgresRepository) GetByName(ctx context.Context, name string) (*entities.Category, error) {
	const op = "CategoryRepository.GetByName"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":   op,
		"name": name,
	})

	category := &entities.Category{}

	err := r.db.QueryRowContext(ctx, queryGetByName, name).Scan(
		&category.CategoryID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("category not found by name")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).Error("failed to get category by name")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully retrieved category by name")
	return category, nil
}

func (r *CategoryPostgresRepository) Create(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	const op = "CategoryRepository.Create"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *CategoryPostgresRepository) Update(ctx context.Context, category *entities.Category) error {
	const op = "CategoryRepository.Update"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *CategoryPostgresRepository) Delete(ctx context.Context, categoryID uuid.UUID) error {
	const op = "CategoryRepository.Delete"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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
