package roadmap

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	GetAll(context.Context) ([]*entities.Roadmap, error)
	GetByID(context.Context, primitive.ObjectID) (*entities.Roadmap, error)
	Create(ctx context.Context, req *dto.RoadmapDTO) (*dto.RoadmapDTO, error)
	Update(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.UpdateRoadmapRequest) error
	Delete(context.Context, primitive.ObjectID) error
	Generate(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.GenerateRoadmapRequest) (*dto.GenerateRoadmapResponse, error)
}
