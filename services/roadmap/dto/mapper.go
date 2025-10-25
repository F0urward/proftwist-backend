package dto

import (
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

func EntityToDTO(entity *entities.Roadmap) *RoadmapDTO {
	if entity == nil {
		return nil
	}

	return &RoadmapDTO{
		ID:        entity.ID,
		Nodes:     nodesToDTO(entity.Nodes),
		Edges:     edgesToDTO(entity.Edges),
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
		Nodes:     dtoToNodes(dto.Nodes),
		Edges:     dtoToEdges(dto.Edges),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func UpdateRequestToEntity(existing *entities.Roadmap, request *UpdateRoadmapRequest) *entities.Roadmap {
	if existing == nil || request == nil {
		return existing
	}

	updated := *existing

	if request.Nodes != nil {
		updated.Nodes = dtoToNodes(request.Nodes)
	}
	if request.Edges != nil {
		updated.Edges = dtoToEdges(request.Edges)
	}

	updated.UpdatedAt = time.Now()

	return &updated
}

func nodesToDTO(nodes []entities.RoadmapNode) []NodeDTO {
	if nodes == nil {
		return nil
	}

	result := make([]NodeDTO, len(nodes))
	for i, node := range nodes {
		result[i] = NodeDTO{
			ID:   node.ID,
			Type: node.Type,
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

func dtoToNodes(nodesDTO []NodeDTO) []entities.RoadmapNode {
	if nodesDTO == nil {
		return nil
	}

	result := make([]entities.RoadmapNode, len(nodesDTO))
	for i, nodeDTO := range nodesDTO {
		result[i] = entities.RoadmapNode{
			ID:   nodeDTO.ID,
			Type: nodeDTO.Type,
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

func edgesToDTO(edges []entities.RoadmapEdge) []EdgeDTO {
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

func dtoToEdges(edgesDTO []EdgeDTO) []entities.RoadmapEdge {
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
