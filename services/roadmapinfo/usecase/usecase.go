package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
)

type RoadmapInfoUsecase struct {
	repo roadmapinfo.Repository
}

func NewRoadmapInfoUsecase(repo roadmapinfo.Repository) roadmapinfo.Usecase {
	return &RoadmapInfoUsecase{
		repo: repo,
	}
}

func (uc *RoadmapInfoUsecase) GetAll(ctx context.Context) (*dto.GetAllRoadmapsInfoResponseDTO, error) {
	roadmaps, err := uc.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get all roadmaps: %v", err)
		return nil, fmt.Errorf("failed to get all roadmaps: %v", err)
	}

	var roadmapDTOs []dto.RoadmapInfoResponseDTO

	for _, roadmap := range roadmaps {
		roadmapDTO := dto.RoadmapInfoResponseDTO{
			ID:              roadmap.ID.String(),
			OwnerID:         roadmap.OwnerID.String(),
			CategoryID:      roadmap.CategoryID.String(),
			Name:            roadmap.Name,
			Description:     roadmap.Description,
			IsPublic:        roadmap.IsPublic,
			SubscriberCount: roadmap.SubscriberCount,
			CreatedAt:       roadmap.CreatedAt,
			UpdatedAt:       roadmap.UpdatedAt,
		}

		if roadmap.ReferencedRoadmapInfoID != nil {
			roadmapDTO.ReferencedRoadmapInfoID = roadmap.ReferencedRoadmapInfoID.String()
		}

		roadmapDTOs = append(roadmapDTOs, roadmapDTO)
	}

	return &dto.GetAllRoadmapsInfoResponseDTO{
		RoadmapsInfo: roadmapDTOs,
	}, nil
}

func (uc *RoadmapInfoUsecase) GetByID(ctx context.Context, roadmapID uuid.UUID) (*dto.GetByIDRoadmapInfoResponseDTO, error) {
	roadmap, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		log.Printf("Failed to get roadmap by ID %s: %v", roadmapID, err)
		return nil, fmt.Errorf("failed to get roadmap: %v", err)
	}

	if roadmap == nil {
		return nil, fmt.Errorf("roadmap not found")
	}

	roadmapDTO := dto.RoadmapInfoResponseDTO{
		ID:              roadmap.ID.String(),
		OwnerID:         roadmap.OwnerID.String(),
		CategoryID:      roadmap.CategoryID.String(),
		Name:            roadmap.Name,
		Description:     roadmap.Description,
		IsPublic:        roadmap.IsPublic,
		SubscriberCount: roadmap.SubscriberCount,
		CreatedAt:       roadmap.CreatedAt,
		UpdatedAt:       roadmap.UpdatedAt,
	}

	if roadmap.ReferencedRoadmapInfoID != nil {
		roadmapDTO.ReferencedRoadmapInfoID = roadmap.ReferencedRoadmapInfoID.String()
	}

	return &dto.GetByIDRoadmapInfoResponseDTO{RoadmapInfo: roadmapDTO}, nil
}

func (uc *RoadmapInfoUsecase) Create(ctx context.Context, request *dto.CreateRoadmapInfoRequestDTO) error {
	ownerID, err := uuid.Parse(request.OwnerID)
	if err != nil {
		return fmt.Errorf("invalid owner_id format: %v", err)
	}

	categoryID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category_id format: %v", err)
	}

	var referencedRoadmapInfoID *uuid.UUID

	if request.ReferencedRoadmapInfoID != nil {
		refID, err2 := uuid.Parse(*request.ReferencedRoadmapInfoID)
		if err2 != nil {
			return fmt.Errorf("invalid referenced_roadmap_id format: %v", err2)
		}
		referencedRoadmapInfoID = &refID
	}

	newRoadmapInfo := &entities.RoadmapInfo{
		OwnerID:                 ownerID,
		CategoryID:              categoryID,
		Name:                    request.Name,
		Description:             request.Description,
		IsPublic:                request.IsPublic,
		ReferencedRoadmapInfoID: referencedRoadmapInfoID,
	}

	err = uc.repo.Create(ctx, newRoadmapInfo)
	if err != nil {
		log.Printf("Failed to create roadmap: %v", err)
		return fmt.Errorf("failed to create roadmap: %v", err)
	}

	return nil
}

func (uc *RoadmapInfoUsecase) Update(ctx context.Context, roadmapID uuid.UUID, request *dto.UpdateRoadmapInfoRequestDTO) error {
	existing, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		log.Printf("Failed to get roadmap by ID %s: %v", roadmapID, err)
		return fmt.Errorf("failed to get roadmap: %v", err)
	}

	if existing == nil {
		return fmt.Errorf("roadmap not found")
	}

	updated := *existing

	if request.CategoryID != nil {
		categoryID, err2 := uuid.Parse(*request.CategoryID)
		if err2 != nil {
			return fmt.Errorf("invalid category_id format: %v", err2)
		}
		updated.CategoryID = categoryID
	}

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
		refID, err2 := uuid.Parse(*request.ReferencedRoadmapInfoID)

		if err2 != nil {
			return fmt.Errorf("invalid referenced_roadmap_id format: %v", err2)
		}

		updated.ReferencedRoadmapInfoID = &refID
	} else {
		updated.ReferencedRoadmapInfoID = nil
	}

	err = uc.repo.Update(ctx, &updated)
	if err != nil {
		log.Printf("Failed to update roadmap with ID %s: %v", roadmapID, err)
		return fmt.Errorf("failed to update roadmap: %v", err)
	}

	return nil
}

func (uc *RoadmapInfoUsecase) Delete(ctx context.Context, roadmapID uuid.UUID) error {
	existing, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		log.Printf("Failed to get roadmap by ID %s: %v", roadmapID, err)
		return fmt.Errorf("failed to get roadmap: %v", err)
	}

	if existing == nil {
		return fmt.Errorf("roadmap not found")
	}

	err = uc.repo.Delete(ctx, roadmapID)
	if err != nil {
		log.Printf("Failed to delete roadmap with ID %s: %v", roadmapID, err)
		return fmt.Errorf("failed to delete roadmap: %v", err)
	}

	return nil
}
