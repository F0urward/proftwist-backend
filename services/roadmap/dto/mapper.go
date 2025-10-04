package dto

import (
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EntityToDTO(entity *entities.Roadmap) *RoadmapDTO {
	if entity == nil {
		return nil
	}

	return &RoadmapDTO{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		IsPublic:    entity.IsPublic,
		SubCount:    entity.SubCount,
		CategoryID:  entity.CategoryID,
		AuthorID:    entity.AuthorID,
		Nodes:       nodesToDTO(entity.Nodes),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func DTOToEntity(dto *RoadmapDTO) *entities.Roadmap {
	if dto == nil {
		return nil
	}

	return &entities.Roadmap{
		ID:          dto.ID,
		Title:       dto.Title,
		Description: dto.Description,
		IsPublic:    dto.IsPublic,
		SubCount:    dto.SubCount,
		CategoryID:  dto.CategoryID,
		AuthorID:    dto.AuthorID,
		Nodes:       dtoToNodes(dto.Nodes),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func CreateRequestToEntity(request *CreateRoadmapRequest) *entities.Roadmap {
	if request == nil {
		return nil
	}

	return &entities.Roadmap{
		ID:          primitive.NewObjectID(),
		Title:       request.Title,
		Description: request.Description,
		Nodes:       dtoToNodes(request.Nodes),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func UpdateRequestToEntity(existing *entities.Roadmap, request *UpdateRoadmapRequest) *entities.Roadmap {
	if existing == nil || request == nil {
		return existing
	}

	updated := *existing

	if request.Title != "" {
		updated.Title = request.Title
	}
	if request.Description != "" {
		updated.Description = request.Description
	}
	if request.Nodes != nil {
		updated.Nodes = dtoToNodes(request.Nodes)
	}

	updated.UpdatedAt = time.Now()

	return &updated
}

func nodesToDTO(nodes []entities.RoadmapNode) []RoadmapNodeDTO {
	if nodes == nil {
		return nil
	}

	result := make([]RoadmapNodeDTO, len(nodes))
	for i, node := range nodes {
		result[i] = RoadmapNodeDTO{
			ID:          node.ID,
			Title:       node.Title,
			Description: node.Description,
			Position: Position{
				X: node.Position.X,
				Y: node.Position.Y,
			},
			AuthorID: node.AuthorID,
			Level:    node.Level,
			Links:    linksToDTO(node.Links),
			Children: nodesToDTO(node.Children),
		}
	}
	return result
}

func dtoToNodes(nodesDTO []RoadmapNodeDTO) []entities.RoadmapNode {
	if nodesDTO == nil {
		return nil
	}

	result := make([]entities.RoadmapNode, len(nodesDTO))
	for i, nodeDTO := range nodesDTO {
		result[i] = entities.RoadmapNode{
			ID:          nodeDTO.ID,
			Title:       nodeDTO.Title,
			Description: nodeDTO.Description,
			Position: entities.Position{
				X: nodeDTO.Position.X,
				Y: nodeDTO.Position.Y,
			},
			AuthorID: nodeDTO.AuthorID,
			Level:    int(nodeDTO.Level),
			Links:    dtoToLinks(nodeDTO.Links),
			Children: dtoToNodes(nodeDTO.Children),
		}
	}
	return result
}

func linksToDTO(links []entities.Link) []Link {
	if links == nil {
		return nil
	}

	result := make([]Link, len(links))
	for i, link := range links {
		result[i] = Link{
			URL:      link.URL,
			Title:    link.Title,
			Type:     link.Type,
			AuthorID: link.AuthorID,
		}
	}
	return result
}

func dtoToLinks(linksDTO []Link) []entities.Link {
	if linksDTO == nil {
		return nil
	}

	result := make([]entities.Link, len(linksDTO))
	for i, linkDTO := range linksDTO {
		result[i] = entities.Link{
			URL:      linkDTO.URL,
			Title:    linkDTO.Title,
			Type:     linkDTO.Type,
			AuthorID: linkDTO.AuthorID,
		}
	}
	return result
}

func NewUpdatePrivacyResponse(roadmapID primitive.ObjectID, isPublic bool) UpdatePrivacyResponse {
	return UpdatePrivacyResponse{
		Message:   "Roadmap privacy updated successfully",
		IsPublic:  isPublic,
		RoadmapID: roadmapID,
	}
}
