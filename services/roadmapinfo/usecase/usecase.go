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

func (uc *RoadmapInfoUsecase) Create(ctx context.Context, request *dto.CreateRoadmapInfoRequestDTO) (*dto.CreateRoadmapInfoResponseDTO, error) {
	const op = "RoadmapInfoUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"author_id": request.AuthorID,
		"name":      request.Name,
	})

	roadmapCreateRequest := &roadmapclient.CreateRequest{
		Id:       primitive.NewObjectID().Hex(),
		IsPublic: request.IsPublic,
		AuthorId: request.AuthorID,
		Nodes:    []*roadmapclient.Node{},
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

	roadmapInfo, err := dto.CreateRequestToEntity(request)
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
	return &dto.CreateRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}

func (uc *RoadmapInfoUsecase) Update(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID, request *dto.UpdateRoadmapInfoRequestDTO) error {
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

	if !uc.isUserOwner(existing, userID) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID.String(),
			"author_id":       existing.AuthorID.String(),
		}).Warn("user is not author of the roadmap info")
		return errs.ErrForbidden
	}

	updated, err := dto.UpdateRequestToEntity(existing, request)
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

func (uc *RoadmapInfoUsecase) Fork(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID) (*dto.CreateRoadmapInfoResponseDTO, error) {
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

	originalRoadmap, err := uc.roadmapClient.GetByID(ctx, &roadmapclient.GetByIDRequest{Id: originalRoadmapInfo.RoadmapID})
	if err != nil {
		logger.WithError(err).Error("failed to get original roadmap")
		return nil, fmt.Errorf("failed to get original roadmap: %w", err)
	}

	if originalRoadmap == nil || originalRoadmap.Roadmap == nil {
		logger.Error("original roadmap is nil")
		return nil, fmt.Errorf("failed to get original roadmap")
	}

	forkRoadmapRequest := &roadmapclient.CreateRequest{
		Id:       primitive.NewObjectID().Hex(),
		IsPublic: false,
		AuthorId: userID.String(),
		Nodes:    originalRoadmap.Roadmap.Nodes,
		Edges:    originalRoadmap.Roadmap.Edges,
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
		ID:          uuid.New(),
		RoadmapID:   forkedRoadmap.Roadmap.Id,
		Name:        originalRoadmapInfo.Name,
		Description: originalRoadmapInfo.Description,
		CategoryID:  originalRoadmapInfo.CategoryID,
		AuthorID:    userID,
		IsPublic:    false,
		CreatedAt:   originalRoadmapInfo.CreatedAt,
		UpdatedAt:   originalRoadmapInfo.UpdatedAt,
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
	return &dto.CreateRoadmapInfoResponseDTO{RoadmapInfo: roadmapInfoDTO}, nil
}
