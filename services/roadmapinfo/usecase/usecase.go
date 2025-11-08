package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
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

func (uc *RoadmapInfoUsecase) GetByID(ctx context.Context, roadmapInfoID uuid.UUID) (*dto.GetByIDRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
	})

	roadmapInfo, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if roadmapInfo == nil {
		logger.Warn("roadmap info not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID := uc.getUserIDFromContext(ctx)
	if !uc.canUserAccessRoadmapInfo(roadmapInfo, userID) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.AuthorID.String(),
			"is_public":       roadmapInfo.IsPublic,
		}).Warn("access denied to roadmap info")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	roadmapInfoDTO := dto.RoadmapInfoToDTO(roadmapInfo)

	logger.Info("successfully retrieved roadmap info")
	return &dto.GetByIDRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}

func (uc *RoadmapInfoUsecase) GetByRoadmapID(ctx context.Context, roadmapID string) (*dto.GetByIDRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetByRoadmapID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID,
	})

	roadmapInfo, err := uc.repo.GetByRoadmapID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info by roadmap ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if roadmapInfo == nil {
		logger.Warn("roadmap info not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID := uc.getUserIDFromContext(ctx)
	if !uc.canUserAccessRoadmapInfo(roadmapInfo, userID) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.AuthorID.String(),
			"is_public":       roadmapInfo.IsPublic,
		}).Warn("access denied to roadmap info")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	roadmapInfoDTO := dto.RoadmapInfoToDTO(roadmapInfo)
	return &dto.GetByIDRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}

func (uc *RoadmapInfoUsecase) GetAllPublicByCategoryID(ctx context.Context, categoryID uuid.UUID) (*dto.GetAllRoadmapsInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetAllPublicByCategoryID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	roadmaps, err := uc.repo.GetAllPublicByCategoryID(ctx, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmaps by category ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(roadmaps) == 0 {
		logger.Debug("no roadmaps found for category")
		return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoDTO{}}, nil
	}

	roadmapDTOs := dto.RoadmapInfoListToDTO(roadmaps)

	logger.WithField("count", len(roadmapDTOs)).Info("successfully retrieved roadmaps by category")
	return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: roadmapDTOs}, nil
}

func (uc *RoadmapInfoUsecase) GetAllByUserID(ctx context.Context, userID uuid.UUID) (*dto.GetAllRoadmapsInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetAllByUserID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID,
	})

	roadmaps, err := uc.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmaps by user ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(roadmaps) == 0 {
		logger.Debug("no roadmaps found for user")
		return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoDTO{}}, nil
	}

	roadmapDTOs := dto.RoadmapInfoListToDTO(roadmaps)

	logger.WithField("count", len(roadmapDTOs)).Info("successfully retrieved roadmaps by user ID")
	return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: roadmapDTOs}, nil
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

func (uc *RoadmapInfoUsecase) Update(ctx context.Context, roadmapInfoID uuid.UUID, request *dto.UpdateRoadmapInfoRequestDTO) error {
	const op = "RoadmapInfoUsecase.Update"
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

	existing, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap info")
		return fmt.Errorf("%s: %w", op, err)
	}
	if existing == nil {
		logger.Warn("roadmap info not found for update")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.WithFields(map[string]interface{}{
		"request_user_id": userIDStr,
		"author_id":       existing.AuthorID.String(),
		"author_id_raw":   existing.AuthorID,
	}).Debug("comparing user IDs")

	if !uc.isUserOwner(existing, userIDStr) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userIDStr,
			"author_id":       existing.AuthorID.String(),
		}).Warn("user is not author of the roadmap info")
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	updated, err := dto.UpdateRequestToEntity(existing, request)
	if err != nil {
		logger.WithError(err).Warn("failed to convert update request to entity")
		return fmt.Errorf("%s: invalid input data: %w", op, err)
	}

	err = uc.repo.Update(ctx, updated)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap info")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("roadmap info updated successfully")
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

	if !uc.isUserOwner(roadmapInfo, userIDStr) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userIDStr,
			"author_id":       roadmapInfo.AuthorID.String(),
		}).Warn("user is not author of the roadmap info")
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
		logger.WithError(deleteErr).Error("failed to delete associated roadmap")
	}

	logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"roadmap_id":      roadmapID.Hex(),
	}).Info("successfully deleted roadmap info and roadmap")
	return nil
}

func (uc *RoadmapInfoUsecase) canUserAccessRoadmapInfo(roadmapInfo *entities.RoadmapInfo, userID string) bool {
	if roadmapInfo.IsPublic {
		return true
	}
	return roadmapInfo.AuthorID.String() == userID
}

func (uc *RoadmapInfoUsecase) isUserOwner(roadmapInfo *entities.RoadmapInfo, userID string) bool {
	return roadmapInfo.AuthorID.String() == userID
}

func (uc *RoadmapInfoUsecase) getUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(utils.UserIDKey{}).(string); ok && userID != "" {
		return userID
	}
	return ""
}
