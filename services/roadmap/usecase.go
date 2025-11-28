package roadmap

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type Usecase interface {
	GetByID(ctx context.Context, roadmapID primitive.ObjectID) (*dto.GetByIDRoadmapResponseDTO, error)
	GetByIDWithMaterials(ctx context.Context, roadmapID primitive.ObjectID) (*dto.GetByIDRoadmapWithMaterialsResponseDTO, error)
	Create(ctx context.Context, req *dto.CreateRoadmapRequestDTO) (*dto.CreateRoadmapResponseDTO, error)
	Update(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.UpdateRoadmapRequestDTO) error
	Delete(context.Context, primitive.ObjectID) error
	Generate(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, req *dto.GenerateRoadmapRequestDTO) (*dto.GenerateRoadmapResponseDTO, error)
	RegenerateNodeIDs(roadmapDTO *dto.RoadmapWithMaterialsDTO) *dto.RoadmapWithMaterialsDTO
	CreateMaterial(ctx context.Context, userID uuid.UUID, roadmapID primitive.ObjectID, nodeID uuid.UUID, req dto.CreateMaterialRequestDTO) (*dto.EnrichedMaterialResponseDTO, error)
	DeleteMaterial(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, materialID uuid.UUID, userID uuid.UUID) error
	GetMaterialsByNode(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID) (*dto.MaterialListResponseDTO, error)
}
