package roadmap

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	GetAll(context.Context) ([]*entities.Roadmap, error)
	GetByID(context.Context, primitive.ObjectID) (*entities.Roadmap, error)
	Create(context.Context, *entities.Roadmap) (*entities.Roadmap, error)
	Update(context.Context, *entities.Roadmap) (*entities.Roadmap, error)
	Delete(context.Context, primitive.ObjectID) error
	Generate(ctx context.Context, roadmapID primitive.ObjectID, req *dto.GenerateRoadmapRequest) (*dto.GenerateRoadmapResponse, error)
}
