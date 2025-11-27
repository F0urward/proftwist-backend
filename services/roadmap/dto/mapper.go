package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

func EntityToDTO(entity *entities.Roadmap) RoadmapDTO {
	return RoadmapDTO{
		ID:        entity.ID,
		Nodes:     NodesToDTO(entity.Nodes),
		Edges:     EdgesToDTO(entity.Edges),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func EntityListToDTO(roadmaps []*entities.Roadmap) []RoadmapDTO {
	var roadmapDTOs []RoadmapDTO

	for _, roadmap := range roadmaps {
		roadmapDTOs = append(roadmapDTOs, EntityToDTO(roadmap))
	}

	return roadmapDTOs
}

func DTOToEntity(dto *RoadmapDTO) *entities.Roadmap {
	if dto == nil {
		return nil
	}

	return &entities.Roadmap{
		ID:        dto.ID,
		Nodes:     DtoToNodes(dto.Nodes),
		Edges:     DtoToEdges(dto.Edges),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func UpdateRequestToEntity(existing *entities.Roadmap, request *UpdateRoadmapRequestDTO) *entities.Roadmap {
	if existing == nil || request == nil {
		return existing
	}

	updated := *existing

	if request.Nodes != nil {
		updated.Nodes = DtoToNodes(request.Nodes)
	}
	if request.Edges != nil {
		updated.Edges = DtoToEdges(request.Edges)
	}

	updated.UpdatedAt = time.Now()

	return &updated
}

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

func DtoToNodes(nodesDTO []NodeDTO) []entities.RoadmapNode {
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
			Selected: nodeDTO.Selected,
			Dragging: nodeDTO.Dragging,
		}
	}
	return result
}

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

func DtoToEdges(edgesDTO []EdgeDTO) []entities.RoadmapEdge {
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

func MaterialToDTO(material *entities.Material, author MaterialAuthorDTO) MaterialResponseDTO {
	if material == nil {
		return MaterialResponseDTO{}
	}

	return MaterialResponseDTO{
		ID:        material.ID,
		Name:      material.Name,
		URL:       material.URL,
		Author:    author,
		CreatedAt: material.CreatedAt,
		UpdatedAt: material.UpdatedAt,
	}
}

func MaterialListToDTO(materials []*entities.Material, authorData map[uuid.UUID]MaterialAuthorDTO) MaterialListResponseDTO {
	materialDTOs := make([]MaterialResponseDTO, 0, len(materials))

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

		materialDTOs = append(materialDTOs, MaterialToDTO(material, author))
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
