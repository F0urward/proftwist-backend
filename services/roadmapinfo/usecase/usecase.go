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

func (uc *RoadmapInfoUsecase) GetAllPublic(ctx context.Context) (*dto.GetAllRoadmapsInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetAllPublic"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	roadmaps, err := uc.repo.GetAllPublic(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to get public roadmaps")
		return nil, fmt.Errorf("failed to get public roadmaps: %w", err)
	}

	if roadmaps == nil {
		roadmaps = []*entities.RoadmapInfo{}
	}

	if len(roadmaps) == 0 {
		logger.Debug("no public roadmaps found")
		return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoDTO{}}, nil
	}

	roadmapDTOs := dto.RoadmapInfoListToDTO(roadmaps)

	logger.WithField("count", len(roadmapDTOs)).Info("successfully retrieved public roadmaps")
	return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: roadmapDTOs}, nil
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
		return nil, fmt.Errorf("failed to get roadmap info by ID: %w", err)
	}

	if roadmapInfo == nil {
		logger.Warn("roadmap info not found")
		return nil, errs.ErrNotFound
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

	if roadmapID == "" {
		logger.Warn("roadmap ID is empty")
		return nil, fmt.Errorf("roadmap ID is empty")
	}

	roadmapInfo, err := uc.repo.GetByRoadmapID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info by roadmap ID")
		return nil, fmt.Errorf("failed to get roadmap info by roadmap ID: %w", err)
	}

	if roadmapInfo == nil {
		logger.Warn("roadmap info not found")
		return nil, errs.ErrNotFound
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
		return nil, fmt.Errorf("failed to get roadmaps by category ID: %w", err)
	}

	if roadmaps == nil {
		roadmaps = []*entities.RoadmapInfo{}
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
		return nil, fmt.Errorf("failed to get roadmaps by user ID: %w", err)
	}

	if roadmaps == nil {
		roadmaps = []*entities.RoadmapInfo{}
	}

	if len(roadmaps) == 0 {
		logger.Debug("no roadmaps found for user")
		return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoDTO{}}, nil
	}

	roadmapDTOs := dto.RoadmapInfoListToDTO(roadmaps)

	logger.WithField("count", len(roadmapDTOs)).Info("successfully retrieved roadmaps by user ID")
	return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: roadmapDTOs}, nil
}

func (uc *RoadmapInfoUsecase) CreatePrivate(ctx context.Context, request *dto.CreatePrivateRoadmapInfoRequestDTO) (*dto.CreatePrivateRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"author_id": request.AuthorID,
		"name":      request.Name,
	})

	roadmapCreateRequest := &roadmapclient.CreateRequest{
		Id:       primitive.NewObjectID().Hex(),
		IsPublic: false,
		AuthorId: request.AuthorID,
		Nodes:    []*roadmapclient.NodeWithMaterials{},
		Edges:    []*roadmapclient.Edge{},
	}

	roadmap, err := uc.roadmapClient.Create(ctx, roadmapCreateRequest)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("failed to create roadmap: %w", err)
	}

	if roadmap == nil || roadmap.Roadmap == nil {
		logger.Error("created roadmap is nil")
		return nil, fmt.Errorf("failed to create roadmap")
	}

	roadmapInfo, err := dto.CreatePrivateRequestToEntity(request)
	if err != nil {
		logger.WithError(err).Error("failed to convert create request to entity")
		if _, deleteErr := uc.roadmapClient.Delete(ctx, &roadmapclient.DeleteRequest{
			Id: roadmapCreateRequest.Id,
		}); deleteErr != nil {
			logger.WithError(deleteErr).Error("failed to rollback roadmap creation")
		}
		return nil, fmt.Errorf("failed to convert create request: %w", err)
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
		return nil, fmt.Errorf("failed to create roadmap info: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"roadmap_info_id": createdRoadmapInfo.ID.String(),
		"roadmap_id":      createdRoadmapInfo.RoadmapID,
	}).Info("successfully created roadmap info with roadmap")

	roadmapInfoDTO := dto.RoadmapInfoToDTO(createdRoadmapInfo)
	return &dto.CreatePrivateRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}

func (uc *RoadmapInfoUsecase) UpdatePrivate(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID, request *dto.UpdatePrivateRoadmapInfoRequestDTO) error {
	const op = "RoadmapInfoUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	existing, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap info")
		return fmt.Errorf("failed to get existing roadmap info: %w", err)
	}

	if existing == nil {
		logger.Warn("roadmap info not found for update")
		return errs.ErrNotFound
	}

	if existing.IsPublic {
		logger.Warn("attempt to update public roadmap")
		return fmt.Errorf("attempt to update public roadmap")
	}

	if !uc.isUserOwner(existing, userID) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID.String(),
			"author_id":       existing.AuthorID.String(),
		}).Warn("user is not author of the roadmap info")
		return errs.ErrForbidden
	}

	updated, err := dto.UpdatePrivateRequestToEntity(existing, request)
	if err != nil {
		logger.WithError(err).Warn("failed to convert update request to entity")
		return fmt.Errorf("invalid input data: %w", err)
	}

	if updated == nil {
		logger.Warn("updated roadmap info is nil")
		return fmt.Errorf("failed to update roadmap info")
	}

	err = uc.repo.Update(ctx, updated)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap info")
		return fmt.Errorf("failed to update roadmap info: %w", err)
	}

	logger.Info("roadmap info updated successfully")
	return nil
}

func (uc *RoadmapInfoUsecase) Delete(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID) error {
	const op = "RoadmapInfoUsecase.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	roadmapInfo, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info")
		return fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmapInfo == nil {
		logger.Warn("roadmap info not found")
		return errs.ErrNotFound
	}

	if !uc.isUserOwner(roadmapInfo, userID) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID.String(),
			"author_id":       roadmapInfo.AuthorID.String(),
		}).Warn("user is not author of the roadmap info")
		return errs.ErrForbidden
	}

	if roadmapInfo.RoadmapID != "" {
		if _, deleteErr := uc.roadmapClient.Delete(ctx, &roadmapclient.DeleteRequest{
			Id: roadmapInfo.RoadmapID,
		}); deleteErr != nil {
			logger.WithError(deleteErr).Error("failed to delete associated roadmap")
		}
	}

	err = uc.repo.Delete(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap info")
		return fmt.Errorf("failed to delete roadmap info: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"roadmap_info_id": roadmapInfoID.String(),
		"roadmap_id":      roadmapInfo.RoadmapID,
	}).Info("successfully deleted roadmap info and roadmap")
	return nil
}

func (uc *RoadmapInfoUsecase) isUserOwner(roadmapInfo *entities.RoadmapInfo, userID uuid.UUID) bool {
	if roadmapInfo == nil {
		return false
	}
	return roadmapInfo.AuthorID == userID
}

func (uc *RoadmapInfoUsecase) Fork(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID) (*dto.CreatePrivateRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.Fork"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	originalRoadmapInfo, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get original roadmap info")
		return nil, fmt.Errorf("failed to get original roadmap info: %w", err)
	}

	if originalRoadmapInfo == nil {
		logger.Warn("original roadmap info not found")
		return nil, errs.ErrNotFound
	}

	if !originalRoadmapInfo.IsPublic {
		logger.Warn("attempt to fork private roadmap")
		return nil, errs.ErrForbidden
	}

	originalRoadmap, err := uc.roadmapClient.GetByIDWithMaterials(ctx, &roadmapclient.GetByIDWithMaterialsRequest{Id: originalRoadmapInfo.RoadmapID})
	if err != nil {
		logger.WithError(err).Error("failed to get original roadmap")
		return nil, fmt.Errorf("failed to get original roadmap: %w", err)
	}

	if originalRoadmap == nil || originalRoadmap.Roadmap == nil {
		logger.Error("original roadmap is nil")
		return nil, fmt.Errorf("failed to get original roadmap")
	}

	regeneratedRoadmap, err := uc.roadmapClient.RegenerateNodeIDs(ctx, &roadmapclient.RegenerateNodeIDsRequest{
		Roadmap: originalRoadmap.Roadmap,
	})
	if err != nil {
		logger.WithError(err).Error("failed to regenerate node IDs for forked roadmap")
		return nil, fmt.Errorf("failed to regenerate node IDs: %w", err)
	}

	if regeneratedRoadmap == nil || regeneratedRoadmap.Roadmap == nil {
		logger.Error("regenerated roadmap is nil")
		return nil, fmt.Errorf("failed to regenerate node IDs")
	}

	forkRoadmapRequest := &roadmapclient.CreateRequest{
		Id:       primitive.NewObjectID().Hex(),
		IsPublic: false,
		AuthorId: userID.String(),
		Nodes:    regeneratedRoadmap.Roadmap.Nodes,
		Edges:    regeneratedRoadmap.Roadmap.Edges,
	}

	forkedRoadmap, err := uc.roadmapClient.Create(ctx, forkRoadmapRequest)
	if err != nil {
		logger.WithError(err).Error("failed to create forked roadmap")
		return nil, fmt.Errorf("failed to create forked roadmap: %w", err)
	}

	if forkedRoadmap == nil || forkedRoadmap.Roadmap == nil {
		logger.Error("forked roadmap is nil")
		return nil, fmt.Errorf("failed to create forked roadmap")
	}

	forkedRoadmapInfo := &entities.RoadmapInfo{
		ID:                      uuid.New(),
		RoadmapID:               forkedRoadmap.Roadmap.Id,
		Name:                    originalRoadmapInfo.Name,
		Description:             originalRoadmapInfo.Description,
		CategoryID:              originalRoadmapInfo.CategoryID,
		AuthorID:                userID,
		IsPublic:                false,
		ReferencedRoadmapInfoID: &originalRoadmapInfo.ID,
		CreatedAt:               originalRoadmapInfo.CreatedAt,
		UpdatedAt:               originalRoadmapInfo.UpdatedAt,
	}

	createdRoadmapInfo, err := uc.repo.Create(ctx, forkedRoadmapInfo)
	if err != nil {
		if _, deleteErr := uc.roadmapClient.Delete(ctx, &roadmapclient.DeleteRequest{
			Id: forkedRoadmap.Roadmap.Id,
		}); deleteErr != nil {
			logger.WithError(deleteErr).Error("failed to rollback roadmap creation")
		}
		logger.WithError(err).Error("failed to create forked roadmap info")
		return nil, fmt.Errorf("failed to create forked roadmap info: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"forked_roadmap_info_id": createdRoadmapInfo.ID.String(),
		"forked_roadmap_id":      createdRoadmapInfo.RoadmapID,
	}).Info("successfully forked roadmap info")

	roadmapInfoDTO := dto.RoadmapInfoToDTO(createdRoadmapInfo)
	return &dto.CreatePrivateRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}

func (uc *RoadmapInfoUsecase) Publish(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID) (*dto.CreatePrivateRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.Publish"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	originalRoadmapInfo, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get original roadmap info")
		return nil, fmt.Errorf("failed to get original roadmap info: %w", err)
	}

	if originalRoadmapInfo == nil {
		logger.Warn("original roadmap info not found")
		return nil, errs.ErrNotFound
	}

	if originalRoadmapInfo.IsPublic {
		logger.Warn("attempt to publish already public roadmap")
		return nil, fmt.Errorf("attempt to publish public roadmap")
	}

	if !uc.isUserOwner(originalRoadmapInfo, userID) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID.String(),
			"author_id":       originalRoadmapInfo.AuthorID.String(),
		}).Warn("user is not author of the roadmap info")
		return nil, errs.ErrForbidden
	}

	originalRoadmap, err := uc.roadmapClient.GetByIDWithMaterials(ctx, &roadmapclient.GetByIDWithMaterialsRequest{Id: originalRoadmapInfo.RoadmapID})
	if err != nil {
		logger.WithError(err).Error("failed to get original roadmap")
		return nil, fmt.Errorf("failed to get original roadmap: %w", err)
	}

	if originalRoadmap == nil || originalRoadmap.Roadmap == nil {
		logger.Error("original roadmap is nil")
		return nil, fmt.Errorf("failed to get original roadmap")
	}

	regeneratedRoadmap, err := uc.roadmapClient.RegenerateNodeIDs(ctx, &roadmapclient.RegenerateNodeIDsRequest{
		Roadmap: originalRoadmap.Roadmap,
	})
	if err != nil {
		logger.WithError(err).Error("failed to regenerate node IDs for published roadmap")
		return nil, fmt.Errorf("failed to regenerate node IDs: %w", err)
	}

	if regeneratedRoadmap == nil || regeneratedRoadmap.Roadmap == nil {
		logger.Error("regenerated roadmap is nil")
		return nil, fmt.Errorf("failed to regenerate node IDs")
	}

	publishRoadmapRequest := &roadmapclient.CreateRequest{
		Id:       primitive.NewObjectID().Hex(),
		IsPublic: true,
		AuthorId: userID.String(),
		Nodes:    regeneratedRoadmap.Roadmap.Nodes,
		Edges:    regeneratedRoadmap.Roadmap.Edges,
	}

	publishedRoadmap, err := uc.roadmapClient.Create(ctx, publishRoadmapRequest)
	if err != nil {
		logger.WithError(err).Error("failed to create published roadmap")
		return nil, fmt.Errorf("failed to create published roadmap: %w", err)
	}

	if publishedRoadmap == nil || publishedRoadmap.Roadmap == nil {
		logger.Error("published roadmap is nil")
		return nil, fmt.Errorf("failed to create published roadmap")
	}

	publishedRoadmapInfo := &entities.RoadmapInfo{
		ID:                      uuid.New(),
		RoadmapID:               publishedRoadmap.Roadmap.Id,
		Name:                    originalRoadmapInfo.Name,
		Description:             originalRoadmapInfo.Description,
		CategoryID:              originalRoadmapInfo.CategoryID,
		AuthorID:                userID,
		IsPublic:                true,
		ReferencedRoadmapInfoID: originalRoadmapInfo.ReferencedRoadmapInfoID,
		CreatedAt:               originalRoadmapInfo.CreatedAt,
		UpdatedAt:               originalRoadmapInfo.UpdatedAt,
	}

	createdRoadmapInfo, err := uc.repo.Create(ctx, publishedRoadmapInfo)
	if err != nil {
		if _, deleteErr := uc.roadmapClient.Delete(ctx, &roadmapclient.DeleteRequest{
			Id: publishedRoadmap.Roadmap.Id,
		}); deleteErr != nil {
			logger.WithError(deleteErr).Error("failed to rollback roadmap creation")
		}
		logger.WithError(err).Error("failed to create published roadmap info")
		return nil, fmt.Errorf("failed to create published roadmap info: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"published_roadmap_info_id": createdRoadmapInfo.ID.String(),
		"published_roadmap_id":      createdRoadmapInfo.RoadmapID,
	}).Info("successfully published roadmap info")

	roadmapInfoDTO := dto.RoadmapInfoToDTO(createdRoadmapInfo)
	return &dto.CreatePrivateRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}

func (uc *RoadmapInfoUsecase) Subscribe(ctx context.Context, roadmapInfoID, userID uuid.UUID) error {
	const op = "RoadmapInfoUsecase.Subscribe"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	roadmap, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info")
		return fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmap == nil {
		logger.Warn("roadmap info not found")
		return errs.ErrNotFound
	}

	if !roadmap.IsPublic {
		logger.Warn("attempt to subscribe to private roadmap")
		return errs.ErrForbidden
	}

	if roadmap.AuthorID == userID {
		logger.Warn("attempt to subscribe to own roadmap")
		return fmt.Errorf("cannot subscribe to your own roadmap")
	}

	exists, err := uc.repo.SubscriptionExists(ctx, userID, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to check subscription existence")
		return fmt.Errorf("failed to check subscription: %w", err)
	}

	if exists {
		logger.Warn("subscription already exists")
		return errs.ErrAlreadyExists
	}

	err = uc.repo.CreateSubscription(ctx, userID, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to create subscription")
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	logger.Info("successfully subscribed to roadmap")
	return nil
}

func (uc *RoadmapInfoUsecase) Unsubscribe(ctx context.Context, roadmapInfoID, userID uuid.UUID) error {
	const op = "RoadmapInfoUsecase.Unsubscribe"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	exists, err := uc.repo.SubscriptionExists(ctx, userID, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to check subscription existence")
		return fmt.Errorf("failed to check subscription: %w", err)
	}

	if !exists {
		logger.Warn("subscription not found")
		return errs.ErrNotFound
	}

	err = uc.repo.DeleteSubscription(ctx, userID, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to delete subscription")
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	logger.Info("successfully unsubscribed from roadmap")
	return nil
}

func (uc *RoadmapInfoUsecase) GetSubscribed(ctx context.Context, userID uuid.UUID) (*dto.GetSubscribedRoadmapsInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.GetSubscribedRoadmaps"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	subscribedIDs, err := uc.repo.GetSubscribedRoadmapIDs(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get subscribed roadmap IDs")
		return nil, fmt.Errorf("failed to get subscribed roadmap IDs: %w", err)
	}

	if len(subscribedIDs) == 0 {
		logger.Debug("no subscriptions found")
		return &dto.GetSubscribedRoadmapsInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoDTO{}}, nil
	}

	roadmaps, err := uc.repo.GetByIDs(ctx, subscribedIDs)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmaps by IDs")
		return nil, fmt.Errorf("failed to get roadmaps by IDs: %w", err)
	}

	if len(roadmaps) == 0 {
		logger.Debug("no roadmaps found for subscriptions")
		return &dto.GetSubscribedRoadmapsInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoDTO{}}, nil
	}

	roadmapDTOs := dto.RoadmapInfoListToDTO(roadmaps)

	logger.WithField("count", len(roadmapDTOs)).Info("successfully retrieved subscribed roadmaps")
	return &dto.GetSubscribedRoadmapsInfoResponseDTO{RoadmapsInfo: roadmapDTOs}, nil
}

func (uc *RoadmapInfoUsecase) CheckSubscription(ctx context.Context, roadmapInfoID, userID uuid.UUID) (bool, error) {
	const op = "RoadmapInfoUsecase.CheckSubscription"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"roadmap_info_id": roadmapInfoID.String(),
		"user_id":         userID.String(),
	})

	roadmap, err := uc.repo.GetByID(ctx, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info")
		return false, fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmap == nil {
		logger.Warn("roadmap info not found")
		return false, errs.ErrNotFound
	}

	if !roadmap.IsPublic {
		logger.Debug("roadmap is private and user is not author")
		return false, nil
	}

	isSubscribed, err := uc.repo.SubscriptionExists(ctx, userID, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to check subscription existence")
		return false, fmt.Errorf("failed to check subscription: %w", err)
	}

	logger.WithField("is_subscribed", isSubscribed).Debug("subscription check completed")
	return isSubscribed, nil
}

func (uc *RoadmapInfoUsecase) SearchPublic(ctx context.Context, query string, categoryID *uuid.UUID) (*dto.GetAllRoadmapsInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.SearchPublic"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"query":       query,
		"category_id": categoryID,
	})

	if query == "" {
		logger.Warn("search query is empty")
		return nil, fmt.Errorf("search query cannot be empty")
	}

	roadmaps, err := uc.repo.SearchPublic(ctx, query, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to search public roadmaps")
		return nil, fmt.Errorf("failed to search public roadmaps: %w", err)
	}

	if roadmaps == nil {
		roadmaps = []*entities.RoadmapInfo{}
	}

	if len(roadmaps) == 0 {
		logger.Debug("no public roadmaps found for search query")
		return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: []dto.RoadmapInfoDTO{}}, nil
	}

	roadmapDTOs := dto.RoadmapInfoListToDTO(roadmaps)

	logger.WithField("count", len(roadmapDTOs)).Info("successfully searched public roadmaps")
	return &dto.GetAllRoadmapsInfoResponseDTO{RoadmapsInfo: roadmapDTOs}, nil
}
