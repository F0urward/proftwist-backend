package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

// ==================== Roadmap Mappers ====================

func EntityToDTO(entity *entities.Roadmap) RoadmapDTO {
	return RoadmapDTO{
		ID:        entity.ID,
		Nodes:     NodesToDTO(entity.Nodes),
		Edges:     EdgesToDTO(entity.Edges),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func EntityToWithMaterialsDTO(entity *entities.Roadmap) RoadmapWithMaterialsDTO {
	return RoadmapWithMaterialsDTO{
		ID:                 entity.ID,
		NodesWithMaterials: NodesToWithMaterialsDTO(entity.Nodes),
		Edges:              EdgesToDTO(entity.Edges),
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
	}
}

func EntityToDTOWithProgress(entity *entities.Roadmap, userProgress *entities.UserProgress) RoadmapWithProgressDTO {
	return RoadmapWithProgressDTO{
		ID:        entity.ID,
		Nodes:     NodesToDTOWithProgress(entity.Nodes, userProgress),
		Edges:     EdgesToDTO(entity.Edges),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func DTOToEntity(dto *RoadmapDTO) *entities.Roadmap {
	if dto == nil {
		return nil
	}

	return &entities.Roadmap{
		ID:        dto.ID,
		Nodes:     DTOToNodes(dto.Nodes),
		Edges:     DTOToEdges(dto.Edges),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func DTOWithMaterialsToEntity(dto *RoadmapWithMaterialsDTO) *entities.Roadmap {
	if dto == nil {
		return nil
	}

	return &entities.Roadmap{
		ID:        dto.ID,
		Nodes:     DTOToNodesWithMaterials(dto.NodesWithMaterials),
		Edges:     DTOToEdges(dto.Edges),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func EntityListToDTO(roadmaps []*entities.Roadmap) []RoadmapDTO {
	var roadmapDTOs []RoadmapDTO

	for _, roadmap := range roadmaps {
		roadmapDTOs = append(roadmapDTOs, EntityToDTO(roadmap))
	}

	return roadmapDTOs
}

func UpdateRequestToEntityWithMaterials(existing *entities.Roadmap, request *UpdateRoadmapRequestDTO) *entities.Roadmap {
	if existing == nil || request == nil {
		return existing
	}

	updated := *existing

	if request.Nodes != nil {
		updated.Nodes = mergeNodesWithMaterials(existing.Nodes, request.Nodes)
	}

	if request.Edges != nil {
		updated.Edges = DTOToEdges(request.Edges)
	}

	updated.UpdatedAt = time.Now()

	return &updated
}

func mergeNodesWithMaterials(existingNodes []entities.RoadmapNode, newNodesDTO []NodeDTO) []entities.RoadmapNode {
	existingNodesMap := make(map[uuid.UUID]entities.RoadmapNode)
	for _, node := range existingNodes {
		existingNodesMap[node.ID] = node
	}

	result := make([]entities.RoadmapNode, len(newNodesDTO))

	for i, newNodeDTO := range newNodesDTO {
		newNode := entities.RoadmapNode{
			ID:          newNodeDTO.ID,
			Type:        newNodeDTO.Type,
			Description: newNodeDTO.Description,
			Position: entities.Position{
				X: newNodeDTO.Position.X,
				Y: newNodeDTO.Position.Y,
			},
			Data: entities.NodeData{
				Label: newNodeDTO.Data.Label,
				Type:  newNodeDTO.Data.Type,
			},
			Measured: entities.Measured{
				Width:  newNodeDTO.Measured.Width,
				Height: newNodeDTO.Measured.Height,
			},
			Selected: newNodeDTO.Selected,
			Dragging: newNodeDTO.Dragging,
		}

		if existingNode, exists := existingNodesMap[newNodeDTO.ID]; exists {
			newNode.Materials = existingNode.Materials
		} else {
			newNode.Materials = []entities.Material{}
		}

		result[i] = newNode
	}

	return result
}

// ==================== Node Mappers ====================

func NodesToDTO(nodes []entities.RoadmapNode) []NodeDTO {
	if nodes == nil {
		return nil
	}

	result := make([]NodeDTO, len(nodes))
	for i, node := range nodes {
		result[i] = NodeDTO{
			ID:          node.ID,
			Type:        node.Type,
			Description: node.Description,
			Position: Position{
				X: node.Position.X,
				Y: node.Position.Y,
			},
			Data: NodeData{
				Label: node.Data.Label,
				Type:  node.Data.Type,
			},
			Measured: Measured{
				Width:  node.Measured.Width,
				Height: node.Measured.Height,
			},
			Selected: node.Selected,
			Dragging: node.Dragging,
		}
	}
	return result
}

func NodesToWithMaterialsDTO(nodes []entities.RoadmapNode) []NodeWithMaterialsDTO {
	if nodes == nil {
		return nil
	}

	result := make([]NodeWithMaterialsDTO, len(nodes))
	for i, node := range nodes {
		result[i] = NodeWithMaterialsDTO{
			ID:          node.ID,
			Type:        node.Type,
			Description: node.Description,
			Position: Position{
				X: node.Position.X,
				Y: node.Position.Y,
			},
			Data: NodeData{
				Label: node.Data.Label,
				Type:  node.Data.Type,
			},
			Measured: Measured{
				Width:  node.Measured.Width,
				Height: node.Measured.Height,
			},
			Selected:  node.Selected,
			Dragging:  node.Dragging,
			Materials: MaterialListToDTO(node.Materials),
		}
	}
	return result
}

func NodesToDTOWithProgress(nodes []entities.RoadmapNode, userProgress *entities.UserProgress) []NodeWithProgressDTO {
	if nodes == nil {
		return nil
	}

	result := make([]NodeWithProgressDTO, len(nodes))

	if userProgress == nil || userProgress.Progress == nil {
		for i, node := range nodes {
			result[i] = NodeWithProgressDTO{
				ID:          node.ID,
				Type:        node.Type,
				Description: node.Description,
				Position: Position{
					X: node.Position.X,
					Y: node.Position.Y,
				},
				Data: NodeData{
					Label: node.Data.Label,
					Type:  node.Data.Type,
				},
				Measured: Measured{
					Width:  node.Measured.Width,
					Height: node.Measured.Height,
				},
				Selected: node.Selected,
				Dragging: node.Dragging,
				Progress: nil,
			}
		}
		return result
	}

	for i, node := range nodes {
		nodeDTO := NodeWithProgressDTO{
			ID:          node.ID,
			Type:        node.Type,
			Description: node.Description,
			Position: Position{
				X: node.Position.X,
				Y: node.Position.Y,
			},
			Data: NodeData{
				Label: node.Data.Label,
				Type:  node.Data.Type,
			},
			Measured: Measured{
				Width:  node.Measured.Width,
				Height: node.Measured.Height,
			},
			Selected: node.Selected,
			Dragging: node.Dragging,
		}

		if progress, exists := userProgress.Progress[node.ID]; exists {
			nodeDTO.Progress = &NodeProgress{
				Status: NodeProgressStatus(progress.Status),
			}
		} else {
			nodeDTO.Progress = &NodeProgress{
				Status: NodeProgressPending,
			}
		}

		result[i] = nodeDTO
	}
	return result
}

func DTOToNodes(nodesDTO []NodeDTO) []entities.RoadmapNode {
	if nodesDTO == nil {
		return nil
	}

	result := make([]entities.RoadmapNode, len(nodesDTO))
	for i, nodeDTO := range nodesDTO {
		result[i] = entities.RoadmapNode{
			ID:          nodeDTO.ID,
			Type:        nodeDTO.Type,
			Description: nodeDTO.Description,
			Position: entities.Position{
				X: nodeDTO.Position.X,
				Y: nodeDTO.Position.Y,
			},
			Data: entities.NodeData{
				Label: nodeDTO.Data.Label,
				Type:  nodeDTO.Data.Type,
			},
			Measured: entities.Measured{
				Width:  nodeDTO.Measured.Width,
				Height: nodeDTO.Measured.Height,
			},
			Selected:  nodeDTO.Selected,
			Dragging:  nodeDTO.Dragging,
			Materials: []entities.Material{},
		}
	}
	return result
}

func DTOToNodesWithMaterials(nodesDTO []NodeWithMaterialsDTO) []entities.RoadmapNode {
	if nodesDTO == nil {
		return nil
	}

	result := make([]entities.RoadmapNode, len(nodesDTO))
	for i, nodeDTO := range nodesDTO {
		materials := DTOToMaterials(nodeDTO.Materials)
		if materials == nil {
			materials = []entities.Material{}
		}
		result[i] = entities.RoadmapNode{
			ID:          nodeDTO.ID,
			Type:        nodeDTO.Type,
			Description: nodeDTO.Description,
			Position: entities.Position{
				X: nodeDTO.Position.X,
				Y: nodeDTO.Position.Y,
			},
			Data: entities.NodeData{
				Label: nodeDTO.Data.Label,
				Type:  nodeDTO.Data.Type,
			},
			Measured: entities.Measured{
				Width:  nodeDTO.Measured.Width,
				Height: nodeDTO.Measured.Height,
			},
			Selected:  nodeDTO.Selected,
			Dragging:  nodeDTO.Dragging,
			Materials: materials,
		}
	}
	return result
}

// ==================== Edge Mappers ====================

func EdgesToDTO(edges []entities.RoadmapEdge) []EdgeDTO {
	if edges == nil {
		return nil
	}

	result := make([]EdgeDTO, len(edges))
	for i, edge := range edges {
		result[i] = EdgeDTO{
			Source: edge.Source,
			Target: edge.Target,
			ID:     edge.ID,
		}
	}
	return result
}

func DTOToEdges(edgesDTO []EdgeDTO) []entities.RoadmapEdge {
	if edgesDTO == nil {
		return nil
	}

	result := make([]entities.RoadmapEdge, len(edgesDTO))
	for i, edgeDTO := range edgesDTO {
		result[i] = entities.RoadmapEdge{
			Source: edgeDTO.Source,
			Target: edgeDTO.Target,
			ID:     edgeDTO.ID,
		}
	}
	return result
}

// ==================== Material Mappers ====================

func MaterialListToDTO(materials []entities.Material) []Material {
	if materials == nil {
		return nil
	}

	result := make([]Material, len(materials))
	for i, material := range materials {
		result[i] = Material{
			ID:        material.ID,
			Name:      material.Name,
			URL:       material.URL,
			AuthorID:  material.AuthorID,
			CreatedAt: material.CreatedAt,
			UpdatedAt: material.UpdatedAt,
		}
	}
	return result
}

func DTOToMaterials(materialsDTO []Material) []entities.Material {
	if materialsDTO == nil {
		return nil
	}

	result := make([]entities.Material, len(materialsDTO))
	for i, materialDTO := range materialsDTO {
		result[i] = entities.Material{
			ID:        materialDTO.ID,
			Name:      materialDTO.Name,
			URL:       materialDTO.URL,
			AuthorID:  materialDTO.AuthorID,
			CreatedAt: materialDTO.CreatedAt,
			UpdatedAt: materialDTO.UpdatedAt,
		}
	}
	return result
}

func MaterialToEnrichedDTO(material *entities.Material, author MaterialAuthorDTO) EnrichedMaterialResponseDTO {
	if material == nil {
		return EnrichedMaterialResponseDTO{}
	}

	return EnrichedMaterialResponseDTO{
		ID:        material.ID,
		Name:      material.Name,
		URL:       material.URL,
		Author:    author,
		CreatedAt: material.CreatedAt,
		UpdatedAt: material.UpdatedAt,
	}
}

func MaterialListToEnrichedDTO(materials []*entities.Material, authorData map[uuid.UUID]MaterialAuthorDTO) MaterialListResponseDTO {
	materialDTOs := make([]EnrichedMaterialResponseDTO, 0, len(materials))

	for _, material := range materials {
		if material == nil {
			continue
		}

		author, exists := authorData[material.AuthorID]
		if !exists {
			author = MaterialAuthorDTO{
				ID:        material.AuthorID,
				Username:  "Unknown User",
				AvatarURL: "",
			}
		}

		materialDTOs = append(materialDTOs, MaterialToEnrichedDTO(material, author))
	}

	return MaterialListResponseDTO{
		Materials: materialDTOs,
		Total:     len(materialDTOs),
	}
}

func CreateMaterialRequestToEntity(req CreateMaterialRequestDTO, authorID uuid.UUID) *entities.Material {
	now := time.Now()
	return &entities.Material{
		ID:        uuid.New(),
		Name:      req.Name,
		URL:       req.URL,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ==================== Progress Mappers ====================

func (s NodeProgressStatus) IsValid() bool {
	switch s {
	case NodeProgressPending,
		NodeProgressInProgress,
		NodeProgressDone,
		NodeProgressSkipped:
		return true
	default:
		return false
	}
}

func UpdateNodeProgressRequestToEntity(req UpdateNodeProgressRequestDTO) entities.NodeProgress {
	return entities.NodeProgress{
		Status: entities.NodeProgressStatus(req.Status),
	}
}
