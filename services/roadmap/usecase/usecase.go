package roadmap

import (
	"context"
	"fmt"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"log"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
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

	// r.Context().FIXME: Добавьте проверку авторства здесь
	// Если roadmap не публичный, проверяем что пользователь - автор

	return roadmap, nil
}

func (uc *RoadmapUsecase) GetByAuthorID(ctx context.Context, authorID uuid.UUID) ([]*entities.Roadmap, error) {
	const op = "RoadmapUsecase.GetByAuthorID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"author_id": authorID.String(),
	})

	roadmaps, err := uc.repo.GetByAuthorID(ctx, authorID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmaps by author ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// r.Context().FIXME: Добавьте проверку авторства здесь
	// Если пользователь запрашивает свои roadmap, показываем все
	// Если чужие - только публичные

	var publicRoadmaps []*entities.Roadmap
	for _, roadmap := range roadmaps {
		if roadmap.IsPublic {
			publicRoadmaps = append(publicRoadmaps, roadmap)
		}
	}

	return publicRoadmaps, nil
}

func (uc *RoadmapUsecase) Create(ctx context.Context, roadmap *entities.Roadmap) (*entities.Roadmap, error) {
	const op = "RoadmapUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.Hex(),
		"author_id":  roadmap.AuthorID.String(),
		"title":      roadmap.Title,
	})

	err := uc.repo.Create(ctx, roadmap)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return roadmap, nil
}

func (uc *RoadmapUsecase) Update(ctx context.Context, roadmap *entities.Roadmap) (*entities.Roadmap, error) {
	const op = "RoadmapUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.Hex(),
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

	// r.Context().FIXME: Добавить проверку авторства здесь

	err = uc.repo.Update(ctx, roadmap)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

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

	// r.Context().FIXME: Добавьте проверку авторства здесь
	// Показываем только публичные roadmap в поиске
	// Если пользователь авторизован, показываем его приватные roadmap + все публичные

	err = uc.repo.Delete(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *RoadmapUsecase) SearchByTitle(ctx context.Context, title string) ([]*entities.Roadmap, error) {
	const op = "RoadmapUsecase.SearchByTitle"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":    op,
		"title": title,
	})

	roadmaps, err := uc.repo.SearchByTitle(ctx, title)
	if err != nil {
		logger.WithError(err).Error("failed to search roadmaps by title")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// r.Context().FIXME: Добавьте проверку авторства здесь
	// Если пользователь авторизован, показываем его приватные roadmap + все публичные

	return roadmaps, nil
}

func (uc *RoadmapUsecase) UpdatePrivacy(ctx context.Context, roadmapID primitive.ObjectID, isPublic bool) error {
	const op = "RoadmapUsecase.UpdatePrivacy"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
		"is_public":  isPublic,
	})

	existing, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap for privacy update")
		return fmt.Errorf("%s: %w", op, err)
	}
	if existing == nil {
		logger.Warn("roadmap not found for privacy update")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	// r.Context().FIXME: Добавить проверку авторства здесь

	err = uc.repo.UpdatePrivacy(ctx, roadmapID, isPublic)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap privacy")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Printf("Successfully updated privacy for roadmap %s to %v", roadmapID.Hex(), isPublic)
	return nil
}
