package roadmap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/moderationclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type RoadmapUsecase struct {
	mongoRepo         roadmap.MongoRepository
	gigachatWebapi    roadmap.GigachatWebapi
	roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient
	chatClient        chatclient.ChatServiceClient
	authClient        authclient.AuthServiceClient
	moderationClient  moderationclient.ModerationServiceClient
}

func NewRoadmapUsecase(
	mongoRepo roadmap.MongoRepository,
	gigichatWebapi roadmap.GigachatWebapi,
	roadmapInfoClient roadmapinfoclient.RoadmapInfoServiceClient,
	chatClient chatclient.ChatServiceClient,
	authClient authclient.AuthServiceClient,
	moderationClient moderationclient.ModerationServiceClient,
) roadmap.Usecase {
	return &RoadmapUsecase{
		mongoRepo:         mongoRepo,
		gigachatWebapi:    gigichatWebapi,
		roadmapInfoClient: roadmapInfoClient,
		chatClient:        chatClient,
		authClient:        authClient,
		moderationClient:  moderationClient,
	}
}

func (uc *RoadmapUsecase) GetByID(ctx context.Context, roadmapID primitive.ObjectID) (*dto.GetByIDRoadmapResponseDTO, error) {
	const op = "RoadmapUsecase.GetByID"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (uc *RoadmapUsecase) GetByIDWithMaterials(ctx context.Context, roadmapID primitive.ObjectID) (*dto.GetByIDRoadmapWithMaterialsResponseDTO, error) {
	const op = "RoadmapUsecase.GetByIDWithMaterials"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
	})

	roadmap, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by ID with materials")
		return nil, fmt.Errorf("failed to get roadmap by ID with materials: %w", err)
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

	roadmapWithMaterialsDTO := dto.EntityToWithMaterialsDTO(roadmap)

	logger.WithFields(map[string]interface{}{
		"roadmap_id":  roadmapWithMaterialsDTO.ID.Hex(),
		"nodes_count": len(roadmapWithMaterialsDTO.NodesWithMaterials),
	}).Info("successfully retrieved roadmap with materials")

	return &dto.GetByIDRoadmapWithMaterialsResponseDTO{RoadmapWithMaterials: roadmapWithMaterialsDTO}, nil
}

func (uc *RoadmapUsecase) Create(ctx context.Context, req *dto.CreateRoadmapRequestDTO) (*dto.CreateRoadmapResponseDTO, error) {
	const op = "RoadmapUsecase.Create"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"nodes_count": len(req.Roadmap.NodesWithMaterials),
		"edges_count": len(req.Roadmap.Edges),
		"is_public":   req.IsPublic,
	})

	roadmapEntity := dto.DTOWithMaterialsToEntity(&req.Roadmap)
	if roadmapEntity == nil {
		logger.Warn("failed to convert request to entity")
		return nil, fmt.Errorf("invalid request data")
	}

	err := uc.mongoRepo.Create(ctx, roadmapEntity)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("failed to create roadmap: %w", err)
	}

	roadmapDTO := dto.EntityToWithMaterialsDTO(roadmapEntity)
	if roadmapDTO.ID.IsZero() {
		logger.Error("created roadmap has invalid ID")
		return nil, fmt.Errorf("failed to create roadmap")
	}

	if req.IsPublic {
		nodes := make([]dto.NodeDTO, len(req.Roadmap.NodesWithMaterials))
		for i, node := range req.Roadmap.NodesWithMaterials {
			nodes[i] = dto.NodeDTO{
				ID:          node.ID,
				Type:        node.Type,
				Description: node.Description,
				Position:    node.Position,
				Data:        node.Data,
				Measured:    node.Measured,
				Selected:    node.Selected,
				Dragging:    node.Dragging,
			}
		}
		go uc.createNodeChats(context.Background(), req.AuthorID, dto.DTOToNodes(nodes))
	}

	logger.WithField("roadmap_id", roadmapDTO.ID.Hex()).Info("successfully created roadmap")
	return &dto.CreateRoadmapResponseDTO{RoadmapWithMaterials: roadmapDTO}, nil
}

func (uc *RoadmapUsecase) Update(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.UpdateRoadmapRequestDTO) error {
	const op = "RoadmapUsecase.Update"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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
			"author_id":       roadmapInfo.RoadmapInfo.Author.UserId,
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

	updatedEntity := dto.UpdateRequestToEntityWithMaterials(existingEntity, req)
	if updatedEntity == nil {
		logger.Warn("failed to apply updates to roadmap")
		return fmt.Errorf("invalid update data")
	}

	err = uc.moderateRoadmap(ctx, updatedEntity)
	if err != nil {
		logger.WithError(err).Warn("roadmap update rejected due to moderation")
		return fmt.Errorf("moderation check failed: %w", err)
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
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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
			"author_id":       roadmapInfo.RoadmapInfo.Author.UserId,
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

func (uc *RoadmapUsecase) RegenerateNodeIDs(roadmapDTO *dto.RoadmapWithMaterialsDTO) *dto.RoadmapWithMaterialsDTO {
	if roadmapDTO == nil {
		return nil
	}

	nodeIDMap := make(map[string]string)

	regeneratedNodes := make([]dto.NodeWithMaterialsDTO, 0, len(roadmapDTO.NodesWithMaterials))
	for _, node := range roadmapDTO.NodesWithMaterials {
		oldID := node.ID
		newID := uuid.New()
		nodeIDMap[oldID.String()] = newID.String()

		regeneratedNode := dto.NodeWithMaterialsDTO{
			ID:          newID,
			Type:        node.Type,
			Description: node.Description,
			Position:    node.Position,
			Data:        node.Data,
			Measured:    node.Measured,
			Selected:    node.Selected,
			Dragging:    node.Dragging,
			Materials:   node.Materials,
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

	return &dto.RoadmapWithMaterialsDTO{
		ID:                 roadmapDTO.ID,
		NodesWithMaterials: regeneratedNodes,
		Edges:              regeneratedEdges,
		CreatedAt:          roadmapDTO.CreatedAt,
		UpdatedAt:          roadmapDTO.UpdatedAt,
	}
}

func (uc *RoadmapUsecase) isUserOwner(roadmapInfo *roadmapinfoclient.RoadmapInfo, userID string) bool {
	if roadmapInfo == nil {
		return false
	}
	return roadmapInfo.Author.UserId == userID
}

func (uc *RoadmapUsecase) createNodeChats(ctx context.Context, userID uuid.UUID, nodes []entities.RoadmapNode) {
	logger := ctxutil.GetLogger(ctx)

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
	logger := ctxutil.GetLogger(ctx)

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

func (uc *RoadmapUsecase) CreateMaterial(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, nodeID uuid.UUID, req dto.CreateMaterialRequestDTO) (*dto.EnrichedMaterialResponseDTO, error) {
	const op = "RoadmapUsecase.CreateMaterial"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	roadmapEntity, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap")
		return nil, fmt.Errorf("failed to get roadmap: %w", err)
	}

	if roadmapEntity == nil {
		logger.Warn("roadmap not found")
		return nil, errs.ErrNotFound
	}

	roadmapInfo, err := uc.roadmapInfoClient.GetByRoadmapID(ctx, &roadmapinfoclient.GetByRoadmapIDRequest{RoadmapId: roadmapID.Hex()})
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap info for authorization check")
		return nil, fmt.Errorf("failed to get roadmap info: %w", err)
	}

	if roadmapInfo == nil || roadmapInfo.RoadmapInfo == nil {
		logger.Error("roadmap info not found")
		return nil, errs.ErrNotFound
	}

	if !roadmapInfo.RoadmapInfo.IsPublic && !uc.isUserOwner(roadmapInfo.RoadmapInfo, userID.String()) {
		logger.WithFields(map[string]interface{}{
			"request_user_id": userID,
			"author_id":       roadmapInfo.RoadmapInfo.Author.UserId,
		}).Warn("user is not author of the roadmap")
		return nil, errs.ErrForbidden
	}

	nodeExists := false
	for _, node := range roadmapEntity.Nodes {
		if node.ID == nodeID {
			nodeExists = true
			break
		}
	}

	if !nodeExists {
		logger.WithField("node_id", nodeID).Warn("node not found in roadmap")
		return nil, errs.ErrNotFound
	}

	materialEntity := dto.CreateMaterialRequestToEntity(req, userID)

	createdMaterial, err := uc.mongoRepo.CreateMaterial(ctx, roadmapID, nodeID, materialEntity)
	if err != nil {
		logger.WithError(err).Error("failed to create material")
		return nil, fmt.Errorf("failed to create material: %w", err)
	}

	authorData, err := uc.fetchAuthorData(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("failed to fetch author data, using fallback")
		authorData = uc.createFallbackAuthorData(userID)
	}

	response := dto.MaterialToEnrichedDTO(createdMaterial, authorData)

	logger.WithFields(map[string]interface{}{
		"material_id": createdMaterial.ID,
		"roadmap_id":  roadmapID.Hex(),
		"node_id":     nodeID,
		"author_id":   userID,
	}).Info("successfully created material")

	return &response, nil
}

func (uc *RoadmapUsecase) DeleteMaterial(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, materialID uuid.UUID, userID uuid.UUID) error {
	const op = "RoadmapUsecase.DeleteMaterial"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	roadmapEntity, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap")
		return fmt.Errorf("failed to get roadmap: %w", err)
	}

	if roadmapEntity == nil {
		logger.Warn("roadmap not found")
		return errs.ErrNotFound
	}

	nodeExists := false
	for _, node := range roadmapEntity.Nodes {
		if node.ID == nodeID {
			nodeExists = true
			break
		}
	}

	if !nodeExists {
		logger.WithField("node_id", nodeID).Warn("node not found in roadmap")
		return errs.ErrNotFound
	}

	material, err := uc.mongoRepo.GetMaterialByID(ctx, roadmapID, nodeID, materialID)
	if err != nil {
		logger.WithError(err).Error("failed to get material")
		return fmt.Errorf("failed to get material: %w", err)
	}

	if material == nil {
		logger.WithField("material_id", materialID).Warn("material not found")
		return errs.ErrNotFound
	}

	isMaterialAuthor := material.AuthorID == userID

	if !isMaterialAuthor {
		logger.WithFields(map[string]interface{}{
			"material_id": materialID,
			"author_id":   material.AuthorID,
			"user_id":     userID,
		}).Warn("user tried to delete material they don't own")
		return errs.ErrForbidden
	}

	if err := uc.mongoRepo.DeleteMaterial(ctx, roadmapID, nodeID, materialID); err != nil {
		logger.WithError(err).Error("failed to delete material")
		return fmt.Errorf("failed to delete material: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"material_id": materialID,
		"roadmap_id":  roadmapID.Hex(),
		"node_id":     nodeID,
		"user_id":     userID,
	}).Info("successfully deleted material")

	return nil
}

func (uc *RoadmapUsecase) GetMaterialsByNode(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID) (*dto.MaterialListResponseDTO, error) {
	const op = "RoadmapUsecase.GetMaterialsByNode"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	roadmapEntity, err := uc.mongoRepo.GetByID(ctx, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to get roadmap")
		return nil, fmt.Errorf("failed to get roadmap: %w", err)
	}

	if roadmapEntity == nil {
		logger.Warn("roadmap not found")
		return nil, errs.ErrNotFound
	}

	nodeExists := false
	for _, node := range roadmapEntity.Nodes {
		if node.ID == nodeID {
			nodeExists = true
			break
		}
	}

	if !nodeExists {
		logger.WithField("node_id", nodeID).Warn("node not found in roadmap")
		return nil, errs.ErrNotFound
	}

	materials, err := uc.mongoRepo.GetMaterialsByNode(ctx, roadmapID, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to get materials by node")
		return nil, fmt.Errorf("failed to get materials by node: %w", err)
	}

	if materials == nil {
		materials = []*entities.Material{}
	}

	authorData, err := uc.fetchAuthorsData(ctx, materials)
	if err != nil {
		logger.WithError(err).Warn("failed to fetch some author data, using fallback")
	}

	response := dto.MaterialListToEnrichedDTO(materials, authorData)

	logger.WithFields(map[string]interface{}{
		"roadmap_id": roadmapID.Hex(),
		"node_id":    nodeID,
		"count":      len(response.Materials),
	}).Info("successfully retrieved materials by node")

	return &response, nil
}

func (uc *RoadmapUsecase) extractRoadmapContent(roadmap *entities.Roadmap) string {
	if roadmap == nil || len(roadmap.Nodes) == 0 {
		return ""
	}

	var contentBuilder strings.Builder

	for _, node := range roadmap.Nodes {
		if node.Data.Label != "" {
			contentBuilder.WriteString(node.Data.Label)
			contentBuilder.WriteString(" ")
		}

		if node.Description != "" {
			contentBuilder.WriteString(node.Description)
			contentBuilder.WriteString(" ")
		}
	}

	return strings.TrimSpace(contentBuilder.String())
}

func (uc *RoadmapUsecase) moderateRoadmap(ctx context.Context, roadmap *entities.Roadmap) error {
	const op = "RoadmapUsecase.checkRoadmapModeration"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	content := uc.extractRoadmapContent(roadmap)
	if content == "" {
		logger.Warn("roadmap has no content to moderate")
		return nil
	}

	logger.WithField("content_length", len(content)).Debug("sending content for moderation")

	resp, err := uc.moderationClient.ModerateContent(ctx, &moderationclient.ModerateContentRequest{
		Content: content,
	})
	if err != nil {
		logger.WithError(err).Error("failed to call moderation service")
		return fmt.Errorf("moderation service unavailable: %w", err)
	}

	if resp.Error != "" {
		logger.WithField("error", resp.Error).Error("moderation service returned error")
		return fmt.Errorf("moderation error: %s", resp.Error)
	}

	if resp.Result == nil {
		logger.Error("moderation service returned nil result")
		return fmt.Errorf("invalid moderation response")
	}

	if !resp.Result.Allowed {
		logger.WithFields(map[string]interface{}{
			"categories": resp.Result.Categories,
		}).Warn("roadmap content failed moderation")

		categoriesStr := strings.Join(resp.Result.Categories, ", ")
		return fmt.Errorf("content violates moderation rules: %s", categoriesStr)
	}

	logger.WithFields(map[string]interface{}{
		"allowed":    resp.Result.Allowed,
		"categories": resp.Result.Categories,
	}).Debug("roadmap passed moderation")

	return nil
}

func (uc *RoadmapUsecase) fetchAuthorData(ctx context.Context, userID uuid.UUID) (dto.MaterialAuthorDTO, error) {
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

func (uc *RoadmapUsecase) fetchAuthorsData(ctx context.Context, materials []*entities.Material) (map[uuid.UUID]dto.MaterialAuthorDTO, error) {
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

func (uc *RoadmapUsecase) createFallbackAuthorData(userID uuid.UUID) dto.MaterialAuthorDTO {
	return dto.MaterialAuthorDTO{
		ID:        userID,
		Username:  "Unknown User",
		AvatarURL: "",
	}
}

func (uc *RoadmapUsecase) createFallbackAuthorsData(authorIDs map[uuid.UUID]bool) map[uuid.UUID]dto.MaterialAuthorDTO {
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
