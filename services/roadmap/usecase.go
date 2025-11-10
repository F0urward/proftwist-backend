package roadmap

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type Usecase interface {
	GetByID(ctx context.Context, roadmapID primitive.ObjectID) (*dto.GetByIDRoadmapResponseDTO, error)
	Create(ctx context.Context, req *dto.CreateRoamapRequest) (*dto.RoadmapDTO, error)
	Update(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.UpdateRoadmapRequestDTO) error
	Delete(context.Context, primitive.ObjectID) error
	Generate(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.GenerateRoadmapRequestDTO) (*dto.GenerateRoadmapResponseDTO, error)
	RegenerateNodeIDs(roadmapDTO *dto.RoadmapDTO) *dto.RoadmapDTO
}
