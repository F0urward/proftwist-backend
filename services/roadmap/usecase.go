package roadmap

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	GetAll(ctx context.Context) (*dto.GetAllRoadmapsResponseDTO, error)
	GetByID(ctx context.Context, roadmapID primitive.ObjectID) (*dto.GetByIDRoadmapResponseDTO, error)
	Create(ctx context.Context, req *dto.RoadmapDTO) (*dto.RoadmapDTO, error)
	Update(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.UpdateRoadmapRequestDTO) error
	Delete(context.Context, primitive.ObjectID) error
	Generate(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.GenerateRoadmapRequestDTO) (*dto.GenerateRoadmapResponseDTO, error)
}
