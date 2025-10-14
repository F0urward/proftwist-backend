package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
)

type RoadmapInfoUsecase struct {
	repo roadmapinfo.Repository
}

func NewRoadmapInfoUsecase(repo roadmapinfo.Repository) roadmapinfo.Usecase {
	return &RoadmapInfoUsecase{
		repo: repo,
	}
}

func (uc *RoadmapInfoUsecase) GetAll(ctx context.Context) (*dto.GetAllRoadmapsInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetAll"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	roadmaps, err := uc.repo.GetAll(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to get all roadmaps")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.RoadmapInfoListToDTO(roadmaps)
	return &response, nil
}

func (uc *RoadmapInfoUsecase) GetByID(ctx context.Context, roadmapID uuid.UUID) (*dto.GetByIDRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.String(),
	})

	roadmap, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if roadmap == nil {
		logger.Warn("roadmap not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	roadmapDTO := dto.RoadmapInfoToDTO(roadmap)

	logger.Info("successfully retrieved roadmapinfo")
	return &dto.GetByIDRoadmapInfoResponseDTO{RoadmapInfo: roadmapDTO}, nil
}

func (uc *RoadmapInfoUsecase) Create(ctx context.Context, request *dto.CreateRoadmapInfoRequestDTO) error {
	const op = "RoadmapInfoUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"author_id": request.AuthorID,
		"name":      request.Name,
	})

	newRoadmapInfo, err := dto.CreateRequestToEntity(request)
	if err != nil {
		logger.WithError(err).Warn("failed to convert request to entity")
		return fmt.Errorf("%s: invalid input data: %w", op, err)
	}

	err = uc.repo.Create(ctx, newRoadmapInfo)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("roadmap_id", newRoadmapInfo.ID.String()).Info("roadmap created successfully")
	return nil
}

func (uc *RoadmapInfoUsecase) Update(ctx context.Context, roadmapID uuid.UUID, request *dto.UpdateRoadmapInfoRequestDTO) error {
	const op = "RoadmapInfoUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.String(),
	})

	existing, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}
	if existing == nil {
		logger.Warn("roadmap not found for update")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	updated, err := dto.UpdateRequestToEntity(existing, request)
	if err != nil {
		logger.WithError(err).Warn("failed to convert update request to entity")
		return fmt.Errorf("%s: invalid input data: %w", op, err)
	}

	err = uc.repo.Update(ctx, updated)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("roadmap updated successfully")
	return nil
}

func (uc *RoadmapInfoUsecase) Delete(ctx context.Context, roadmapID uuid.UUID) error {
	const op = "RoadmapInfoUsecase.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.String(),
	})

	existing, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap for deletion")
		return fmt.Errorf("%s: %w", op, err)
	}
	if existing == nil {
		logger.Warn("roadmap not found for deletion")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	err = uc.repo.Delete(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("roadmap deleted successfully")
	return nil
}
