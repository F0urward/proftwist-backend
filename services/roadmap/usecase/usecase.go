package roadmap

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type RoadmapUsecase struct {
	mongoRepo         roadmap.MongoRepository
	gigachatWebapi    roadmap.GigachatWebapi
	roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient
}

func NewRoadmapUsecase(
	mongoRepo roadmap.MongoRepository,
	gigichatWebapi roadmap.GigachatWebapi,
	roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient,
) roadmap.Usecase {
	return &RoadmapUsecase{
		mongoRepo:         mongoRepo,
		gigachatWebapi:    gigichatWebapi,
		roadmapInfoClient: roadmapInfoClient,
	}
}

func (uc *RoadmapUsecase) GetByID(ctx context.Context, roadmapID primitive.ObjectID) (*dto.GetByIDRoadmapResponseDTO, error) {
	const op = "RoadmapUsecase.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
	})

	roadmap, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if roadmap == nil {
		logger.Warn("roadmap not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if roadmapInfo == nil {
		logger.Error("roadmap info connected with roadmap doesn't exist")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID := uc.getUserIDFromContext(ctx)
	if !uc.canUserAccessRoadmap(roadmapInfo.RoadmapInfo, userID) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.RoadmapInfo.AuthorId,
			"is_public":       roadmapInfo.RoadmapInfo.IsPublic,
		}).Warn("access denied to roadmap")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	roadmapDTO := dto.EntityToDTO(roadmap)

	logger.WithField("roadmap_id", roadmapDTO.ID.Hex()).Info("successfully retrieved roadmap")
	return &dto.GetByIDRoadmapResponseDTO{Roadmap: roadmapDTO}, nil
}

func (uc *RoadmapUsecase) Create(ctx context.Context, req *dto.RoadmapDTO) (*dto.RoadmapDTO, error) {
	const op = "RoadmapUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"nodes_count": len(req.Nodes),
		"edges_count": len(req.Edges),
	})

	roadmapEntity := dto.DTOToEntity(req)
	if roadmapEntity == nil {
		logger.Warn("failed to convert request to entity")
		return nil, fmt.Errorf("%s: invalid request data", op)
	}

	err := uc.mongoRepo.Create(ctx, roadmapEntity)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	roadmapDTO := dto.EntityToDTO(roadmapEntity)

	logger.WithField("roadmap_id", roadmapDTO.ID.Hex()).Info("successfully created roadmap and roadmap info")
	return &roadmapDTO, nil
}

func (uc *RoadmapUsecase) Update(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.UpdateRoadmapRequestDTO) error {
	const op = "RoadmapUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"user_id":     userID,
		"roadmap_id":  roadmapID.Hex(),
		"nodes_count": len(req.Nodes),
		"edges_count": len(req.Edges),
	})

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return fmt.Errorf("%s: %w", op, err)
	}
	if roadmapInfo == nil {
		logger.Error("roadmap info connected with roadmap doesn't exist")
		return fmt.Errorf("%s: %w", op, err)
	}

	if !uc.isUserOwner(roadmapInfo.RoadmapInfo, userID.String()) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.RoadmapInfo.AuthorId,
		}).Warn("user is not author of the roadmap")
		return fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	existingEntity, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}
	if existingEntity == nil {
		logger.Warn("roadmap not found for update")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	updatedEntity := dto.UpdateRequestToEntity(existingEntity, req)
	if updatedEntity == nil {
		logger.Warn("failed to apply updates to roadmap")
		return fmt.Errorf("%s: invalid update data", op)
	}

	err = uc.mongoRepo.Update(ctx, updatedEntity)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully updated roadmap")
	return nil
}

func (uc *RoadmapUsecase) Delete(ctx context.Context, roadmapID primitive.ObjectID) error {
	const op = "RoadmapUsecase.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
	})

	existing, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap for deletion")
		return fmt.Errorf("%s: %w", op, err)
	}
	if existing == nil {
		logger.Warn("roadmap not found for deletion")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	err = uc.mongoRepo.Delete(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully deleted roadmap")
	return nil
}

func (uc *RoadmapUsecase) Generate(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.GenerateRoadmapRequestDTO) (*dto.GenerateRoadmapResponseDTO, error) {
	const op = "RoadmapUsecase.GenerateRoadmap"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
		"complexity": req.Complexity,
	})

	logger.Info("starting roadmap generation")

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if roadmapInfo == nil {
		logger.Error("roadmap info connected with roadmap doesn't exist")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !uc.isUserOwner(roadmapInfo.RoadmapInfo, userID.String()) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.RoadmapInfo.AuthorId,
		}).Warn("user is not author of the roadmap")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrForbidden)
	}

	roadmapDTO := dto.GenerateRoadmapDTO{
		Topic:       roadmapInfo.RoadmapInfo.Name,
		Description: roadmapInfo.RoadmapInfo.Description,
		Content:     req.Content,
		Complexity:  req.Complexity,
	}

	existingRoadmap, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if existingRoadmap == nil {
		logger.Warn("roadmap not found")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.Info("generating roadmap content with AI")
	generatedRoadmap, err := uc.gigachatWebapi.GenerateRoadmapContent(ctx, &roadmapDTO)
	if err != nil {
		logger.WithError(err).Error("failed to generate roadmap content with AI")
		return nil, fmt.Errorf("%s: failed to generate content: %w", op, err)
	}

	updatedRoadmap := &entities.Roadmap{
		ID:        existingRoadmap.ID,
		Nodes:     generatedRoadmap.Nodes,
		Edges:     generatedRoadmap.Edges,
		CreatedAt: existingRoadmap.CreatedAt,
		UpdatedAt: time.Now(),
	}

	logger.Info("saving generated roadmap to database")
	err = uc.mongoRepo.Update(ctx, updatedRoadmap)
	if err != nil {
		logger.WithError(err).Error("failed to save generated roadmap")
		return nil, fmt.Errorf("%s: failed to save roadmap: %w", op, err)
	}

	response := &dto.GenerateRoadmapResponseDTO{
		RoadmapID: updatedRoadmap.ID,
	}

	logger.WithFields(map[string]interface{}{
		"nodes_count": len(updatedRoadmap.Nodes),
		"edges_count": len(updatedRoadmap.Edges),
	}).Info("successfully generated and saved roadmap")

	return response, nil
}

func (uc *RoadmapUsecase) canUserAccessRoadmap(roadmapInfo *roadmapinfoclient.RoadmapInfo, userID string) bool {
	if roadmapInfo.IsPublic {
		return true
	}
	return roadmapInfo.AuthorId == userID
}

func (uc *RoadmapUsecase) isUserOwner(roadmapInfo *roadmapinfoclient.RoadmapInfo, userID string) bool {
	return roadmapInfo.AuthorId == userID
}

func (uc *RoadmapUsecase) getUserIDFromContext(ctx context.Context) string {
	type userIDKey struct{}
	if userID, ok := ctx.Value(userIDKey{}).(string); ok && userID != "" {
		return userID
	}
	return ""
}
