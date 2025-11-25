package roadmap

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type RoadmapUsecase struct {
	mongoRepo         roadmap.MongoRepository
	gigachatWebapi    roadmap.GigachatWebapi
	roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient
	chatClient        chatclient.ChatServiceClient
}

func NewRoadmapUsecase(
	mongoRepo roadmap.MongoRepository,
	gigichatWebapi roadmap.GigachatWebapi,
	roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient,
	chatClient chatclient.ChatServiceClient,
) roadmap.Usecase {
	return &RoadmapUsecase{
		mongoRepo:         mongoRepo,
		gigachatWebapi:    gigichatWebapi,
		roadmapInfoClient: roadmapInfoClient,
		chatClient:        chatClient,
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
		return nil, fmt.Errorf("failed to get roadmap by ID: %w", err)
	}

	if roadmap == nil {
		logger.Warn("roadmap not found")
		return nil, errs.ErrNotFound
	}

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return nil, fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmapInfo == nil || roadmapInfo.RoadmapInfo == nil {
		logger.Error("roadmap info connected with roadmap doesn't exist")
		return nil, errs.ErrNotFound
	}

	roadmapDTO := dto.EntityToDTO(roadmap)

	logger.WithField("roadmap_id", roadmapDTO.ID.Hex()).Info("successfully retrieved roadmap")
	return &dto.GetByIDRoadmapResponseDTO{Roadmap: roadmapDTO}, nil
}

func (uc *RoadmapUsecase) Create(ctx context.Context, req *dto.CreateRoamapRequest) (*dto.RoadmapDTO, error) {
	const op = "RoadmapUsecase.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"nodes_count": len(req.Roadmap.Nodes),
		"edges_count": len(req.Roadmap.Edges),
		"is_public":   req.IsPublic,
	})

	roadmapEntity := dto.DTOToEntity(&req.Roadmap)
	if roadmapEntity == nil {
		logger.Warn("failed to convert request to entity")
		return nil, fmt.Errorf("invalid request data")
	}

	err := uc.mongoRepo.Create(ctx, roadmapEntity)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("failed to create roadmap: %w", err)
	}

	roadmapDTO := dto.EntityToDTO(roadmapEntity)
	if roadmapDTO.ID.IsZero() {
		logger.Error("created roadmap has invalid ID")
		return nil, fmt.Errorf("failed to create roadmap")
	}

	if req.IsPublic {
		go uc.createNodeChats(context.Background(), req.AuthorID, dto.DtoToNodes(req.Roadmap.Nodes))
	}

	logger.WithField("roadmap_id", roadmapDTO.ID.Hex()).Info("successfully created roadmap")
	return &roadmapDTO, nil
}

func (uc *RoadmapUsecase) Update(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.UpdateRoadmapRequestDTO) error {
	const op = "RoadmapUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"user_id":    userID,
		"roadmap_id": roadmapID.Hex(),
	})

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmapInfo == nil || roadmapInfo.RoadmapInfo == nil {
		logger.Error("roadmap info connected with roadmap doesn't exist")
		return errs.ErrNotFound
	}

	if roadmapInfo.RoadmapInfo.IsPublic {
		logger.Warn("attempt to update public roadmap")
		return fmt.Errorf("attempt to update public roadmap")
	}

	if !uc.isUserOwner(roadmapInfo.RoadmapInfo, userID.String()) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.RoadmapInfo.AuthorId,
		}).Warn("user is not author of the roadmap")
		return errs.ErrForbidden
	}

	existingEntity, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap")
		return fmt.Errorf("failed to get existing roadmap: %w", err)
	}

	if existingEntity == nil {
		logger.Warn("roadmap not found for update")
		return errs.ErrNotFound
	}

	updatedEntity := dto.UpdateRequestToEntity(existingEntity, req)
	if updatedEntity == nil {
		logger.Warn("failed to apply updates to roadmap")
		return fmt.Errorf("invalid update data")
	}

	err = uc.mongoRepo.Update(ctx, updatedEntity)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")
		return fmt.Errorf("failed to update roadmap: %w", err)
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
		return fmt.Errorf("failed to get roadmap: %w", err)
	}

	if existing == nil {
		logger.Warn("roadmap not found for deletion")
		return errs.ErrNotFound
	}

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmapInfo == nil || roadmapInfo.RoadmapInfo == nil {
		logger.Error("roadmap info connected with roadmap doesn't exist")
		return errs.ErrNotFound
	}

	if roadmapInfo.RoadmapInfo.IsPublic {
		go uc.deleteNodeChats(context.Background(), existing.Nodes)
	}

	err = uc.mongoRepo.Delete(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap")
		return fmt.Errorf("failed to delete roadmap: %w", err)
	}

	logger.Info("successfully deleted roadmap")
	return nil
}

func (uc *RoadmapUsecase) Generate(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.GenerateRoadmapRequestDTO) (*dto.GenerateRoadmapResponseDTO, error) {
	const op = "RoadmapUsecase.GenerateRoadmap"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
		"user_id":    userID,
	})

	logger.Info("starting roadmap generation")

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return nil, fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmapInfo == nil || roadmapInfo.RoadmapInfo == nil {
		logger.Error("roadmap info connected with roadmap doesn't exist")
		return nil, errs.ErrNotFound
	}

	if !uc.isUserOwner(roadmapInfo.RoadmapInfo, userID.String()) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.RoadmapInfo.AuthorId,
		}).Warn("user is not author of the roadmap")
		return nil, errs.ErrForbidden
	}

	existingRoadmap, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get existing roadmap")
		return nil, fmt.Errorf("failed to get existing roadmap: %w", err)
	}

	if existingRoadmap == nil {
		logger.Warn("roadmap not found")
		return nil, errs.ErrNotFound
	}

	roadmapDTO := dto.GenerateRoadmapDTO{
		Topic:       roadmapInfo.RoadmapInfo.Name,
		Description: roadmapInfo.RoadmapInfo.Description,
		Content:     req.Content,
		Complexity:  req.Complexity,
	}

	logger.Info("generating roadmap content with AI")
	generatedRoadmap, err := uc.gigachatWebapi.GenerateRoadmapContent(ctx, &roadmapDTO)
	if err != nil {
		logger.WithError(err).Error("failed to generate roadmap content with AI")
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if generatedRoadmap == nil {
		logger.Error("AI generated roadmap is nil")
		return nil, fmt.Errorf("failed to generate roadmap content")
	}

	updatedRoadmap := &entities.Roadmap{
		ID:        existingRoadmap.ID,
		Nodes:     generatedRoadmap.Nodes,
		Edges:     generatedRoadmap.Edges,
		CreatedAt: existingRoadmap.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if updatedRoadmap.Nodes == nil {
		updatedRoadmap.Nodes = []entities.RoadmapNode{}
	}
	if updatedRoadmap.Edges == nil {
		updatedRoadmap.Edges = []entities.RoadmapEdge{}
	}

	logger.Info("saving generated roadmap to database")
	err = uc.mongoRepo.Update(ctx, updatedRoadmap)
	if err != nil {
		logger.WithError(err).Error("failed to save generated roadmap")
		return nil, fmt.Errorf("failed to save roadmap: %w", err)
	}

	if roadmapInfo.RoadmapInfo.IsPublic {
		go uc.createNodeChats(context.Background(), userID, generatedRoadmap.Nodes)
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

func (uc *RoadmapUsecase) RegenerateNodeIDs(roadmapDTO *dto.RoadmapDTO) *dto.RoadmapDTO {
	if roadmapDTO == nil {
		return nil
	}

	nodeIDMap := make(map[string]string)

	regeneratedNodes := make([]dto.NodeDTO, 0, len(roadmapDTO.Nodes))
	for _, node := range roadmapDTO.Nodes {
		oldID := node.ID
		newID := uuid.New()
		nodeIDMap[oldID.String()] = newID.String()

		regeneratedNode := dto.NodeDTO{
			ID:       newID,
			Type:     node.Type,
			Position: node.Position,
			Data:     node.Data,
			Measured: node.Measured,
			Selected: node.Selected,
			Dragging: node.Dragging,
		}

		regeneratedNodes = append(regeneratedNodes, regeneratedNode)
	}

	regeneratedEdges := make([]dto.EdgeDTO, 0, len(roadmapDTO.Edges))
	for _, edge := range roadmapDTO.Edges {
		regeneratedEdge := dto.EdgeDTO{
			ID:     uuid.New().String(),
			Source: edge.Source,
			Target: edge.Target,
		}

		if newSourceID, exists := nodeIDMap[edge.Source]; exists {
			regeneratedEdge.Source = newSourceID
		}
		if newTargetID, exists := nodeIDMap[edge.Target]; exists {
			regeneratedEdge.Target = newTargetID
		}

		regeneratedEdges = append(regeneratedEdges, regeneratedEdge)
	}

	return &dto.RoadmapDTO{
		ID:        roadmapDTO.ID,
		Nodes:     regeneratedNodes,
		Edges:     regeneratedEdges,
		CreatedAt: roadmapDTO.CreatedAt,
		UpdatedAt: roadmapDTO.UpdatedAt,
	}
}

func (uc *RoadmapUsecase) isUserOwner(roadmapInfo *roadmapinfoclient.RoadmapInfo, userID string) bool {
	if roadmapInfo == nil {
		return false
	}
	return roadmapInfo.AuthorId == userID
}

func (uc *RoadmapUsecase) createNodeChats(ctx context.Context, userID uuid.UUID, nodes []entities.RoadmapNode) {
	logger := logctx.GetLogger(ctx)

	for _, node := range nodes {
		chatReq := &chatclient.CreateGroupChatRequest{
			UserId:        userID.String(),
			Title:         fmt.Sprintf("Discussion: %s", node.Data.Label),
			RoadmapNodeId: node.ID.String(),
			MemberIds:     []string{},
		}

		_, err := uc.chatClient.CreateGroupChat(ctx, chatReq)
		if err != nil {
			logger.WithError(err).WithField("node_id", node.ID.String()).Warn("failed to create chat for node")
		} else {
			logger.WithField("node_id", node.ID.String()).Info("successfully created chat for node")
		}
	}
}

func (uc *RoadmapUsecase) deleteNodeChats(ctx context.Context, nodes []entities.RoadmapNode) {
	logger := logctx.GetLogger(ctx)

	for _, node := range nodes {
		chatResp, err := uc.chatClient.GetGroupChatByNode(ctx, &chatclient.GetGroupChatByNodeRequest{
			NodeId: node.ID.String(),
		})

		if err == nil && chatResp.GroupChat != nil && chatResp.GroupChat.Id != "" {
			_, err := uc.chatClient.DeleteGroupChat(ctx, &chatclient.DeleteGroupChatRequest{
				ChatId: chatResp.GroupChat.Id,
			})
			if err != nil {
				logger.WithError(err).WithField("node_id", node.ID.String()).Warn("failed to delete chat for node")
			} else {
				logger.WithField("node_id", node.ID.String()).Info("successfully deleted chat for node")
			}
		}
	}
}
