package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/category"
	"github.com/F0urward/proftwist-backend/services/category/dto"
)

type CategoryUsecase struct {
	repo category.Repository
}

func NewCategoryUsecase(repo category.Repository) category.Usecase {
	return &CategoryUsecase{
		repo: repo,
	}
}

func (uc *CategoryUsecase) GetAll(ctx context.Context) (*dto.GetAllCategoriesResponse, error) {
	const op = "CategoryUsecase.GetAll"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	categories, err := uc.repo.GetAll(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to get all categories")
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}

	response := dto.CategoryListToDTO(categories)

	logger.WithField("count", len(response.Categories)).Info("successfully retrieved categories")
	return &response, nil
}

func (uc *CategoryUsecase) GetByID(ctx context.Context, categoryID uuid.UUID) (*dto.CategoryResponse, error) {
	const op = "CategoryUsecase.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	category, err := uc.repo.GetByID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to get category by ID")
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}
	if category == nil {
		logger.Warn("category not found")
		return nil, errs.ErrNotFound
	}

	categoryDTO := dto.CategoryToDTO(category)

	logger.Info("successfully retrieved category")
	return &categoryDTO, nil
}

func (uc *CategoryUsecase) Create(ctx context.Context, req *dto.CreateCategoryRequestDTO) (*dto.CreateCategoryResponse, error) {
	const op = "CategoryUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":   op,
		"name": req.Name,
	})

	existing, err := uc.repo.GetByName(ctx, req.Name)
	if err != nil {
		logger.WithError(err).Error("failed to check category existence")
		return nil, fmt.Errorf("failed to check category existence: %w", err)
	}
	if existing != nil {
		logger.Warn("category with this name already exists")
		return nil, errs.ErrAlreadyExists
	}

	category := dto.CreateRequestToEntity(req)

	createdCategory, err := uc.repo.Create(ctx, category)
	if err != nil {
		logger.WithError(err).Error("failed to create category")
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"category_id": category.CategoryID.String(),
	}).Info("successfully created category")

	categoryDTO := dto.CategoryToDTO(createdCategory)
	return &dto.CreateCategoryResponse{Category: categoryDTO}, nil
}

func (uc *CategoryUsecase) Update(ctx context.Context, categoryID uuid.UUID, req *dto.UpdateCategoryRequestDTO) error {
	const op = "CategoryUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	existing, err := uc.repo.GetByID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing category")
		return fmt.Errorf("failed to get existing category: %w", err)
	}
	if existing == nil {
		logger.Warn("category not found for update")
		return errs.ErrNotFound
	}

	if req.Name != nil && *req.Name != existing.Name {
		existingWithName, err := uc.repo.GetByName(ctx, *req.Name)
		if err != nil {
			logger.WithError(err).Error("failed to check category name existence")
			return fmt.Errorf("failed to check category name existence: %w", err)
		}
		if existingWithName != nil {
			logger.WithField("new_name", *req.Name).Warn("category with this name already exists")
			return errs.ErrAlreadyExists
		}
	}

	updated := dto.UpdateRequestToEntity(existing, req)

	err = uc.repo.Update(ctx, updated)
	if err != nil {
		logger.WithError(err).Error("failed to update category")
		return fmt.Errorf("failed to update category: %w", err)
	}

	logger.Info("category updated successfully")
	return nil
}

func (uc *CategoryUsecase) Delete(ctx context.Context, categoryID uuid.UUID) error {
	const op = "CategoryUsecase.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	category, err := uc.repo.GetByID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to get category")
		return fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		logger.Warn("category not found")
		return errs.ErrNotFound
	}

	err = uc.repo.Delete(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to delete category")
		return fmt.Errorf("failed to delete category: %w", err)
	}

	logger.Info("successfully deleted category")
	return nil
}
