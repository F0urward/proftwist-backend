package roadmap

import (
	"context"
	"fmt"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmap"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoadmapUsecase struct {
	repo roadmap.Repository
}

func NewRoadmapUsecase(repo roadmap.Repository) roadmap.Usecase {
	return &RoadmapUsecase{
		repo: repo,
	}
}

func (uc *RoadmapUsecase) GetAll(ctx context.Context) ([]*entities.Roadmap, error) {
	const op = "RoadmapUsecase.GetAll"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	roadmaps, err := uc.repo.GetAll(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to get all roadmaps")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("count", len(roadmaps)).Debug("successfully retrieved roadmaps")
	return roadmaps, nil
}

func (uc *RoadmapUsecase) GetByID(ctx context.Context, roadmapID primitive.ObjectID) (*entities.Roadmap, error) {
	const op = "RoadmapUsecase.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
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

	logger.WithFields(map[string]interface{}{
		"nodes_count": len(roadmap.Nodes),
		"edges_count": len(roadmap.Edges),
	}).Debug("successfully retrieved roadmap")
	return roadmap, nil
}

func (uc *RoadmapUsecase) Create(ctx context.Context, roadmap *entities.Roadmap) (*entities.Roadmap, error) {
	const op = "RoadmapUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"roadmap_id":  roadmap.ID.Hex(),
		"nodes_count": len(roadmap.Nodes),
		"edges_count": len(roadmap.Edges),
	})

	err := uc.repo.Create(ctx, roadmap)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully created roadmap")
	return roadmap, nil
}

func (uc *RoadmapUsecase) Update(ctx context.Context, roadmap *entities.Roadmap) (*entities.Roadmap, error) {
	const op = "RoadmapUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"roadmap_id":  roadmap.ID.Hex(),
		"nodes_count": len(roadmap.Nodes),
		"edges_count": len(roadmap.Edges),
	})

	existing, err := uc.repo.GetByID(ctx, roadmap.ID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if existing == nil {
		logger.Warn("roadmap not found for update")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	err = uc.repo.Update(ctx, roadmap)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully updated roadmap")
	return roadmap, nil
}

func (uc *RoadmapUsecase) Delete(ctx context.Context, roadmapID primitive.ObjectID) error {
	const op = "RoadmapUsecase.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
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

	logger.Info("successfully deleted roadmap")
	return nil
}
