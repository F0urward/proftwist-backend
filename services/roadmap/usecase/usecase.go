package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type RoadmapUsecase struct {
	repo roadmap.Repository
}

func NewRoadmapUsecase(repo roadmap.Repository) roadmap.Usecase {
	return &RoadmapUsecase{
		repo: repo,
	}
}

func (uc *RoadmapUsecase) GetAll(ctx context.Context) (*dto.GetAllRoadmapsResponseDTO, error) {
	roadmaps, err := uc.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get all roadmaps: %v", err)
		return nil, fmt.Errorf("failed to get all roadmaps: %v", err)
	}

	var roadmapDTOs []dto.RoadmapResponseDTO
	for _, roadmap := range roadmaps {
		roadmapDTO := dto.RoadmapResponseDTO{
			ID:              roadmap.ID.String(),
			OwnerID:         roadmap.OwnerID.String(),
			CategoryID:      roadmap.CategoryID.String(),
			Name:            roadmap.Name,
			Description:     roadmap.Description,
			IsPublic:        roadmap.IsPublic,
			Color:           roadmap.Color,
			SubscriberCount: roadmap.SubscriberCount,
			CreatedAt:       roadmap.CreatedAt,
			UpdatedAt:       roadmap.UpdatedAt,
		}

		if roadmap.ReferencedRoadmapID != nil {
			roadmapDTO.ReferencedRoadmapID = roadmap.ReferencedRoadmapID.String()
		}

		roadmapDTOs = append(roadmapDTOs, roadmapDTO)
	}

	return &dto.GetAllRoadmapsResponseDTO{Roadmaps: roadmapDTOs}, nil
}

func (uc *RoadmapUsecase) GetByID(ctx context.Context, roadmapID uuid.UUID) (*dto.GetByIDRoadmapResponseDTO, error) {
	roadmap, err := uc.repo.GetByID(ctx, roadmapID)
	if err != nil {
		log.Printf("Failed to get roadmap by ID %s: %v", roadmapID, err)
		return nil, fmt.Errorf("failed to get roadmap: %v", err)
	}
	if roadmap == nil {
		return nil, fmt.Errorf("roadmap not found")
	}

	roadmapDTO := dto.RoadmapResponseDTO{
		ID:              roadmap.ID.String(),
		OwnerID:         roadmap.OwnerID.String(),
		CategoryID:      roadmap.CategoryID.String(),
		Name:            roadmap.Name,
		Description:     roadmap.Description,
		IsPublic:        roadmap.IsPublic,
		Color:           roadmap.Color,
		SubscriberCount: roadmap.SubscriberCount,
		CreatedAt:       roadmap.CreatedAt,
		UpdatedAt:       roadmap.UpdatedAt,
	}

	if roadmap.ReferencedRoadmapID != nil {
		roadmapDTO.ReferencedRoadmapID = roadmap.ReferencedRoadmapID.String()
	}

	return &dto.GetByIDRoadmapResponseDTO{Roadmap: roadmapDTO}, nil
}

func (uc *RoadmapUsecase) Create(ctx context.Context, request *dto.CreateRoadmapRequestDTO) error {
	ownerID, err := uuid.Parse(request.OwnerID)
	if err != nil {
		return fmt.Errorf("invalid owner_id format: %v", err)
	}
	categoryID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category_id format: %v", err)
	}

	var referencedRoadmapID *uuid.UUID
	if request.ReferencedRoadmapID != nil {
		refID, err := uuid.Parse(*request.ReferencedRoadmapID)
		if err != nil {
			return fmt.Errorf("invalid referenced_roadmap_id format: %v", err)
		}
		referencedRoadmapID = &refID
	}

	newRoadmap := &entities.Roadmap{
		OwnerID:             ownerID,
		CategoryID:          categoryID,
		Name:                request.Name,
		Description:         request.Description,
		IsPublic:            request.IsPublic,
		Color:               request.Color,
		ReferencedRoadmapID: referencedRoadmapID,
	}

	err = uc.repo.Create(ctx, newRoadmap)
	if err != nil {
		log.Printf("Failed to create roadmap: %v", err)
		return fmt.Errorf("failed to create roadmap: %v", err)
	}
	return nil
}

func (uc *RoadmapUsecase) Update(ctx context.Context, request *dto.UpdateRoadmapRequestDTO) error {
	roadmapID, err := uuid.Parse(request.ID)
	if err != nil {
		return fmt.Errorf("invalid roadmap id format: %v", err)
	}

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
		categoryID, err := uuid.Parse(*request.CategoryID)
		if err != nil {
			return fmt.Errorf("invalid category_id format: %v", err)
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
	if request.Color != nil {
		updated.Color = *request.Color
	}

	if request.ReferencedRoadmapID != nil {
		refID, err := uuid.Parse(*request.ReferencedRoadmapID)
		if err != nil {
			return fmt.Errorf("invalid referenced_roadmap_id format: %v", err)
		}
		updated.ReferencedRoadmapID = &refID
	} else {
		updated.ReferencedRoadmapID = nil
	}

	err = uc.repo.Update(ctx, &updated)
	if err != nil {
		log.Printf("Failed to update roadmap with ID %s: %v", roadmapID, err)
		return fmt.Errorf("failed to update roadmap: %v", err)
	}
	return nil
}

func (uc *RoadmapUsecase) Delete(ctx context.Context, roadmapID uuid.UUID) error {
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
