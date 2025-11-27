package roadmap

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/google/uuid"
)

type MongoRepository interface {
	GetByID(context.Context, primitive.ObjectID) (*entities.Roadmap, error)
	Create(context.Context, *entities.Roadmap) error
	Update(context.Context, *entities.Roadmap) error
	Delete(context.Context, primitive.ObjectID) error
	CreateMaterial(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, material *entities.Material) (*entities.Material, error)
	DeleteMaterial(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, materialID uuid.UUID) error
	GetMaterialByID(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, materialID uuid.UUID) (*entities.Material, error)
	GetMaterialsByNode(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID) ([]*entities.Material, error)
}

type GigachatWebapi interface {
	GenerateRoadmapContent(ctx context.Context, req *dto.GenerateRoadmapDTO) (*entities.Roadmap, error)
}
