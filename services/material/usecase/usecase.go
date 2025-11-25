package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/material"
	"github.com/F0urward/proftwist-backend/services/material/dto"
)

type MaterialUsecase struct {
	repo       material.Repository
	authClient authclient.AuthServiceClient
}

func NewMaterialUsecase(repo material.Repository, authClient authclient.AuthServiceClient) material.Usecase {
	return &MaterialUsecase{
		repo:       repo,
		authClient: authClient,
	}
}

func (uc *MaterialUsecase) CreateMaterial(ctx context.Context, userID uuid.UUID, req dto.CreateMaterialRequestDTO) (*dto.MaterialResponseDTO, error) {
	const op = "MaterialUsecase.CreateMaterial"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	materialEntity := dto.CreateMaterialRequestToEntity(req, userID)

	createdMaterial, err := uc.repo.CreateMaterial(ctx, materialEntity)
	if err != nil {
		logger.WithError(err).Error("failed to create material")
		return nil, fmt.Errorf("failed to create material: %w", err)
	}

	authorData, err := uc.fetchAuthorData(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("failed to fetch author data, using fallback")
		authorData = uc.createFallbackAuthorData(userID)
	}

	response := dto.MaterialToDTO(createdMaterial, authorData)

	logger.WithFields(map[string]interface{}{
		"material_id":     createdMaterial.ID,
		"roadmap_node_id": req.RoadmapNodeID,
		"author_id":       userID,
	}).Info("successfully created material")

	return &response, nil
}

func (uc *MaterialUsecase) DeleteMaterial(ctx context.Context, materialID uuid.UUID, userID uuid.UUID) error {
	const op = "MaterialUsecase.DeleteMaterial"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	material, err := uc.repo.GetMaterialByID(ctx, materialID)
	if err != nil {
		logger.WithError(err).Error("failed to get material")
		return fmt.Errorf("failed to get material: %w", err)
	}

	if material == nil {
		logger.WithField("material_id", materialID).Warn("material not found")
		return errs.ErrNotFound
	}

	if material.AuthorID != userID {
		logger.WithFields(map[string]interface{}{
			"material_id": materialID,
			"author_id":   material.AuthorID,
			"user_id":     userID,
		}).Warn("user tried to delete material they don't own")
		return errs.ErrForbidden
	}

	if err := uc.repo.DeleteMaterial(ctx, materialID); err != nil {
		logger.WithError(err).Error("failed to delete material")
		return fmt.Errorf("failed to delete material: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"material_id": materialID,
		"user_id":     userID,
	}).Info("successfully deleted material")

	return nil
}

func (uc *MaterialUsecase) GetMaterialsByNode(ctx context.Context, nodeID string) (*dto.MaterialListResponseDTO, error) {
	const op = "MaterialUsecase.GetMaterialsByNode"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	materials, err := uc.repo.GetMaterialsByNode(ctx, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to get materials by roadmap node")
		return nil, fmt.Errorf("failed to get materials by roadmap node: %w", err)
	}

	if materials == nil {
		materials = []*entities.Material{}
	}

	authorData, err := uc.fetchAuthorsData(ctx, materials)
	if err != nil {
		logger.WithError(err).Warn("failed to fetch some author data, using fallback")
	}

	response := dto.MaterialListToDTO(materials, authorData)

	logger.WithFields(map[string]interface{}{
		"node_id": nodeID,
		"count":   len(response.Materials),
	}).Info("successfully retrieved materials by roadmap node")

	return &response, nil
}

func (uc *MaterialUsecase) GetUserMaterials(ctx context.Context, userID uuid.UUID) (*dto.MaterialListResponseDTO, error) {
	const op = "MaterialUsecase.GetUserMaterials"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	materials, err := uc.repo.GetMaterialsByAuthor(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user materials")
		return nil, fmt.Errorf("failed to get user materials: %w", err)
	}

	if materials == nil {
		materials = []*entities.Material{}
	}

	authorData, err := uc.fetchAuthorData(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("failed to fetch author data, using fallback")
		authorData = uc.createFallbackAuthorData(userID)
	}

	authorDataMap := make(map[uuid.UUID]dto.MaterialAuthorDTO, len(materials))
	for range materials {
		authorDataMap[userID] = authorData
	}

	response := dto.MaterialListToDTO(materials, authorDataMap)

	logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"count":   len(response.Materials),
	}).Info("successfully retrieved user materials")

	return &response, nil
}

func (uc *MaterialUsecase) fetchAuthorData(ctx context.Context, userID uuid.UUID) (dto.MaterialAuthorDTO, error) {
	resp, err := uc.authClient.GetUserByID(ctx, &authclient.GetUserByIDRequest{UserId: userID.String()})
	if err != nil || resp == nil || resp.User == nil {
		return dto.MaterialAuthorDTO{}, fmt.Errorf("failed to fetch user data: %w", err)
	}

	return dto.MaterialAuthorDTO{
		ID:        userID,
		Username:  resp.User.Username,
		AvatarURL: resp.User.AvatarUrl,
	}, nil
}

func (uc *MaterialUsecase) fetchAuthorsData(ctx context.Context, materials []*entities.Material) (map[uuid.UUID]dto.MaterialAuthorDTO, error) {
	if len(materials) == 0 {
		return make(map[uuid.UUID]dto.MaterialAuthorDTO), nil
	}

	authorIDs := make(map[uuid.UUID]bool)
	for _, material := range materials {
		if material != nil {
			authorIDs[material.AuthorID] = true
		}
	}

	if len(authorIDs) == 0 {
		return make(map[uuid.UUID]dto.MaterialAuthorDTO), nil
	}

	authorIDStrings := make([]string, 0, len(authorIDs))
	for authorID := range authorIDs {
		authorIDStrings = append(authorIDStrings, authorID.String())
	}

	resp, err := uc.authClient.GetUsersByIDs(ctx, &authclient.GetUsersByIDsRequest{UserIds: authorIDStrings})
	if err != nil || resp == nil {
		return uc.createFallbackAuthorsData(authorIDs), fmt.Errorf("failed to fetch users data: %w", err)
	}

	authorData := make(map[uuid.UUID]dto.MaterialAuthorDTO, len(resp.Users))
	for _, user := range resp.Users {
		if user == nil {
			continue
		}
		userID, err := uuid.Parse(user.Id)
		if err != nil {
			continue
		}
		authorData[userID] = dto.MaterialAuthorDTO{
			ID:        userID,
			Username:  user.Username,
			AvatarURL: user.AvatarUrl,
		}
	}

	return authorData, nil
}

func (uc *MaterialUsecase) createFallbackAuthorData(userID uuid.UUID) dto.MaterialAuthorDTO {
	return dto.MaterialAuthorDTO{
		ID:        userID,
		Username:  "Unknown User",
		AvatarURL: "",
	}
}

func (uc *MaterialUsecase) createFallbackAuthorsData(authorIDs map[uuid.UUID]bool) map[uuid.UUID]dto.MaterialAuthorDTO {
	authorData := make(map[uuid.UUID]dto.MaterialAuthorDTO, len(authorIDs))
	for authorID := range authorIDs {
		authorData[authorID] = dto.MaterialAuthorDTO{
			ID:        authorID,
			Username:  "Unknown User",
			AvatarURL: "",
		}
	}
	return authorData
}
