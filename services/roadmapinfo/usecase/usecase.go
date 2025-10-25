package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoadmapInfoUsecase struct {
	repo          roadmapinfo.Repository
	roadmapClient roadmapclient.RoadmapServiceClient
}

func NewRoadmapInfoUsecase(
	repo roadmapinfo.Repository,
	roadmapClient roadmapclient.RoadmapServiceClient,
) roadmapinfo.Usecase {
	return &RoadmapInfoUsecase{
		repo:          repo,
		roadmapClient: roadmapClient,
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

func (uc *RoadmapInfoUsecase) GetAllByCategoryID(ctx context.Context, categoryID uuid.UUID) (*dto.GetAllByCategoryIDRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetAllByCategoryID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	roadmaps, err := uc.repo.GetAllByCategoryID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmaps by category ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(roadmaps) == 0 {
		logger.Debug("no roadmaps found for category")
		return &dto.GetAllByCategoryIDRoadmapInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoResponseDTO{}}, nil
	}

	response := dto.RoadmapInfoListToDTO(roadmaps)

	logger.WithField("count", len(response.RoadmapsInfo)).Info("successfully retrieved roadmaps by category")
	return &dto.GetAllByCategoryIDRoadmapInfoResponseDTO{RoadmapsInfo: response.RoadmapsInfo}, nil
}

func (uc *RoadmapInfoUsecase) GetByRoadmapID(ctx context.Context, roadmapID string) (*dto.GetByIDRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetByRoadmapID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID,
	})

	roadmap, err := uc.repo.GetByRoadmapID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by roadmap ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if roadmap == nil {
		logger.Warn("roadmap not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	roadmapDTO := dto.RoadmapInfoToDTO(roadmap)
	return &dto.GetByIDRoadmapInfoResponseDTO{RoadmapInfo: roadmapDTO}, nil
}

func (uc *RoadmapInfoUsecase) Create(ctx context.Context, request *dto.CreateRoadmapInfoRequestDTO) (*dto.CreateRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.Create"

	userIDStr, ok := ctx.Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrUnauthorized)
	}

	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"author_id": userIDStr,
		"name":      request.Name,
	})

	roadmapCreateRequest := &roadmapclient.CreateRequest{
		Id:    primitive.NewObjectID().Hex(),
		Nodes: []*roadmapclient.Node{},
		Edges: []*roadmapclient.Edge{},
	}

	roadmap, err := uc.roadmapClient.Create(ctx, roadmapCreateRequest)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	roadmapInfo, err := dto.CreateRequestToEntity(request)
	if err != nil {
		logger.WithError(err).Error("failed to convert create request to entity")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	roadmapInfo.RoadmapID = roadmap.Roadmap.Id

	createdRoadmapInfo, err := uc.repo.Create(ctx, roadmapInfo)
	if err != nil {
		if _, deleteErr := uc.roadmapClient.Delete(ctx, &roadmapclient.DeleteRequest{
			Id: roadmap.Roadmap.Id,
		}); deleteErr != nil {
			logger.WithError(deleteErr).Error("failed to rollback roadmap creation")
		}
		logger.WithError(err).Error("failed to create roadmap info")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfo.ID.String(),
		"roadmap_id":      roadmapInfo.RoadmapID,
	}).Info("successfully created roadmap info with roadmap")

	roadmapInfoDTO := dto.RoadmapInfoToDTO(createdRoadmapInfo)
	return &dto.CreateRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}

func (uc *RoadmapInfoUsecase) Update(ctx context.Context, roadmapID uuid.UUID, request *dto.UpdateRoadmapInfoRequestDTO) error {
	const op = "RoadmapInfoUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.String(),
	})

	userIDStr, ok := ctx.Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context or invalid")
		return fmt.Errorf("%s: %w", op, errs.ErrUnauthorized)
	}

	logger.WithField("user_id_from_context", userIDStr).Debug("user ID from context")

	existing, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}
	if existing == nil {
		logger.Warn("roadmap not found for update")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.WithFields(map[string]interface{}{
		"request_user_id": userIDStr,
		"author_id":       existing.AuthorID.String(),
		"author_id_raw":   existing.AuthorID,
	}).Debug("comparing user IDs")

	if existing.AuthorID.String() != userIDStr {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userIDStr,
			"author_id":       existing.AuthorID.String(),
		}).Warn("user is not author of the roadmap")
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
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

func (uc *RoadmapInfoUsecase) Delete(ctx context.Context, roadmapInfoID uuid.UUID) error {
	const op = "RoadmapInfoUsecase.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
	})

	userIDStr, ok := ctx.Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context or invalid")
		return fmt.Errorf("%s: %w", op, errs.ErrUnauthorized)
	}

	logger.WithField("user_id_from_context", userIDStr).Debug("user ID from context")

	roadmapInfo, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info")
		return fmt.Errorf("%s: %w", op, err)
	}
	if roadmapInfo == nil {
		logger.Warn("roadmap info not found")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.WithFields(map[string]interface{}{
		"request_user_id": userIDStr,
		"author_id":       roadmapInfo.AuthorID.String(),
		"author_id_raw":   roadmapInfo.AuthorID,
	}).Debug("comparing user IDs")

	if roadmapInfo.AuthorID.String() != userIDStr {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userIDStr,
			"author_id":       roadmapInfo.AuthorID.String(),
		}).Warn("user is not author of the roadmap")
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	roadmapID, err := primitive.ObjectIDFromHex(roadmapInfo.RoadmapID)
	if err != nil {
		logger.WithError(err).WithField("roadmap_id", roadmapInfo.RoadmapID).Error("invalid roadmap ID format")
		return fmt.Errorf("%s: %w", op, err)
	}

	err = uc.repo.Delete(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap info")
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, deleteErr := uc.roadmapClient.Delete(ctx, &roadmapclient.DeleteRequest{
		Id: roadmapInfo.RoadmapID,
	}); deleteErr != nil {
		logger.WithError(deleteErr).Error("failed to rollback roadmap creation")
	}

	logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"roadmap_id":      roadmapID.Hex(),
	}).Info("successfully deleted roadmap info and roadmap")
	return nil
}
