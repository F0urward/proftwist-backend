package dto

import (
	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

func RoadmapInfoToDTO(roadmap *entities.RoadmapInfo) RoadmapInfoDTO {
	dto := RoadmapInfoDTO{
		ID:          roadmap.ID.String(),
		RoadmapID:   roadmap.RoadmapID,
		AuthorID:    roadmap.AuthorID.String(),
		CategoryID:  roadmap.CategoryID.String(),
		Name:        roadmap.Name,
		Description: roadmap.Description,
		IsPublic:    roadmap.IsPublic,
		CreatedAt:   roadmap.CreatedAt,
		UpdatedAt:   roadmap.UpdatedAt,
	}

	if roadmap.ReferencedRoadmapInfoID != nil {
		dto.ReferencedRoadmapInfoID = roadmap.ReferencedRoadmapInfoID.String()
	} else {
		dto.ReferencedRoadmapInfoID = ""
	}

	return dto
}

func RoadmapInfoListToDTO(roadmaps []*entities.RoadmapInfo) []RoadmapInfoDTO {
	var roadmapDTOs []RoadmapInfoDTO

	for _, roadmap := range roadmaps {
		roadmapDTOs = append(roadmapDTOs, RoadmapInfoToDTO(roadmap))
	}

	return roadmapDTOs
}

func CreatePrivateRequestToEntity(request *CreatePrivateRoadmapInfoRequestDTO) (*entities.RoadmapInfo, error) {
	authorID, err := uuid.Parse(request.AuthorID)
	if err != nil {
		return nil, err
	}

	categoryID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		return nil, err
	}

	var referencedRoadmapInfoID *uuid.UUID
	if request.ReferencedRoadmapInfoID != nil && *request.ReferencedRoadmapInfoID != "" {
		refID, err := uuid.Parse(*request.ReferencedRoadmapInfoID)
		if err != nil {
			return nil, err
		}
		referencedRoadmapInfoID = &refID
	}

	return &entities.RoadmapInfo{
		AuthorID:                authorID,
		CategoryID:              categoryID,
		Name:                    request.Name,
		Description:             request.Description,
		IsPublic:                request.IsPublic,
		ReferencedRoadmapInfoID: referencedRoadmapInfoID,
	}, nil
}

func UpdatePrivateRequestToEntity(existing *entities.RoadmapInfo, request *UpdatePrivateRoadmapInfoRequestDTO) (*entities.RoadmapInfo, error) {
	updated := *existing

	if request.CategoryID != nil {
		categoryID, err := uuid.Parse(*request.CategoryID)
		if err != nil {
			return nil, err
		}
		updated.CategoryID = categoryID
	}

	if request.Name != nil {
		updated.Name = *request.Name
	}

	if request.Description != nil {
		updated.Description = *request.Description
	}

	if request.ReferencedRoadmapInfoID != nil && *request.ReferencedRoadmapInfoID != "" {
		refID, err := uuid.Parse(*request.ReferencedRoadmapInfoID)
		if err != nil {
			return nil, err
		}
		updated.ReferencedRoadmapInfoID = &refID
	} else {
		updated.ReferencedRoadmapInfoID = nil
	}

	return &updated, nil
}
