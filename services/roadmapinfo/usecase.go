package roadmapinfo

import (
	"context"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
)

type Usecase interface {
	GetAllPublic(ctx context.Context) (*dto.GetAllRoadmapsInfoResponseDTO, error)
	GetAllPublicByCategoryID(ctx context.Context, categoryID uuid.UUID) (*dto.GetAllRoadmapsInfoResponseDTO, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) (*dto.GetAllRoadmapsInfoResponseDTO, error)
	GetByID(context.Context, uuid.UUID) (*dto.GetByIDRoadmapInfoResponseDTO, error)
	GetByRoadmapID(ctx context.Context, roadmapID string) (*dto.GetByIDRoadmapInfoResponseDTO, error)
	CreatePrivate(ctx context.Context, request *dto.CreatePrivateRoadmapInfoRequestDTO) (*dto.CreatePrivateRoadmapInfoResponseDTO, error)
	UpdatePrivate(context.Context, uuid.UUID, uuid.UUID, *dto.UpdatePrivateRoadmapInfoRequestDTO) error
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	Fork(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID) (*dto.CreatePrivateRoadmapInfoResponseDTO, error)
	Publish(ctx context.Context, roadmapInfoID uuid.UUID, userID uuid.UUID) (*dto.CreatePrivateRoadmapInfoResponseDTO, error)
}
