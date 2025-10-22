package roadmap

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoRepository interface {
	GetAll(context.Context) ([]*entities.Roadmap, error)
	GetByID(context.Context, primitive.ObjectID) (*entities.Roadmap, error)
	Create(context.Context, *entities.Roadmap) error
	Update(context.Context, *entities.Roadmap) error
	Delete(context.Context, primitive.ObjectID) error
}

type GigachatWebapi interface {
	GenerateRoadmapContent(ctx context.Context, req *dto.GenerateRoadmapDTO) (*entities.Roadmap, error)
}
