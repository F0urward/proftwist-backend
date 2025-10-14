package dto

import (
	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

func RoadmapInfoToDTO(roadmap *entities.RoadmapInfo) RoadmapInfoResponseDTO {
	dto := RoadmapInfoResponseDTO{
		ID:        roadmap.ID.String(),
		RoadmapID: roadmap.RoadmapID,
		AuthorID:  roadmap.AuthorID.String(),
		//CategoryID:      roadmap.CategoryID.String(),
		Name:            roadmap.Name,
		Description:     roadmap.Description,
		IsPublic:        roadmap.IsPublic,
		SubscriberCount: roadmap.SubscriberCount,
		CreatedAt:       roadmap.CreatedAt,
		UpdatedAt:       roadmap.UpdatedAt,
	}

	if roadmap.ReferencedRoadmapInfoID != nil {
		dto.ReferencedRoadmapInfoID = roadmap.ReferencedRoadmapInfoID.String()
	}

	return dto
}

func RoadmapInfoListToDTO(roadmaps []*entities.RoadmapInfo) GetAllRoadmapsInfoResponseDTO {
	var roadmapDTOs []RoadmapInfoResponseDTO

	for _, roadmap := range roadmaps {
		roadmapDTOs = append(roadmapDTOs, RoadmapInfoToDTO(roadmap))
	}

	return GetAllRoadmapsInfoResponseDTO{
		RoadmapsInfo: roadmapDTOs,
	}
}

func CreateRequestToEntity(request *CreateRoadmapInfoRequestDTO) (*entities.RoadmapInfo, error) {
	//categoryID, err := uuid.Parse(request.CategoryID)
	//if err != nil {
	//	return nil, err
	//}

	var referencedRoadmapInfoID *uuid.UUID
	if request.ReferencedRoadmapInfoID != nil {
		refID, err := uuid.Parse(*request.ReferencedRoadmapInfoID)
		if err != nil {
			return nil, err
		}
		referencedRoadmapInfoID = &refID
	}

	return &entities.RoadmapInfo{
		//CategoryID:              categoryID,
		Name:                    request.Name,
		Description:             request.Description,
		IsPublic:                request.IsPublic,
		ReferencedRoadmapInfoID: referencedRoadmapInfoID,
	}, nil
}

func UpdateRequestToEntity(existing *entities.RoadmapInfo, request *UpdateRoadmapInfoRequestDTO) (*entities.RoadmapInfo, error) {
	updated := *existing

	//if request.CategoryID != nil {
	//	categoryID, err := uuid.Parse(*request.CategoryID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	updated.CategoryID = categoryID
	//}

	if request.Name != nil {
		updated.Name = *request.Name
	}

	if request.Description != nil {
		updated.Description = *request.Description
	}

	if request.IsPublic != nil {
		updated.IsPublic = *request.IsPublic
	}

	if request.ReferencedRoadmapInfoID != nil {
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
