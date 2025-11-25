package dto

import (
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

func MaterialToDTO(material *entities.Material, author MaterialAuthorDTO) MaterialResponseDTO {
	if material == nil {
		return MaterialResponseDTO{}
	}

	return MaterialResponseDTO{
		ID:            material.ID,
		Name:          material.Name,
		URL:           material.URL,
		RoadmapNodeID: material.RoadmapNodeID,
		Author:        author,
		CreatedAt:     material.CreatedAt,
		UpdatedAt:     material.UpdatedAt,
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
		ID:            uuid.New(),
		Name:          req.Name,
		URL:           req.URL,
		RoadmapNodeID: req.RoadmapNodeID,
		AuthorID:      authorID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
