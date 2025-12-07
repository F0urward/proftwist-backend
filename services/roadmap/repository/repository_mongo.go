package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmap"
)

const (
	collectionName = "roadmaps"
)

type RoadmapMongoRepository struct {
	collection *mongo.Collection
}

func NewRoadmapMongoRepository(db *mongo.Database) roadmap.MongoRepository {
	return &RoadmapMongoRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *RoadmapMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Roadmap, error) {
	const op = "RoadmapRepository.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": id.Hex(),
	})

	var roadmap entities.Roadmap

	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&roadmap)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.WithError(err).Error("failed to get roadmap by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &roadmap, nil
}

func (r *RoadmapMongoRepository) Create(ctx context.Context, roadmap *entities.Roadmap) error {
	const op = "RoadmapRepository.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.Hex(),
	})

	if roadmap.ID.IsZero() {
		roadmap.ID = primitive.NewObjectID()
	}

	if roadmap.CreatedAt.IsZero() {
		roadmap.CreatedAt = time.Now()
	}

	if roadmap.UpdatedAt.IsZero() {
		roadmap.UpdatedAt = time.Now()
	}

	for i := range roadmap.Nodes {
		if roadmap.Nodes[i].ID == uuid.Nil {
			roadmap.Nodes[i].ID = uuid.New()
		}
	}

	_, err := r.collection.InsertOne(ctx, roadmap)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RoadmapMongoRepository) Update(ctx context.Context, roadmap *entities.Roadmap) error {
	const op = "RoadmapRepository.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.Hex(),
	})

	roadmap.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"id": roadmap.ID},
		roadmap,
	)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		logger.Warn("roadmap not found for update")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("roadmap not found"))
	}

	return nil
}

func (r *RoadmapMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	const op = "RoadmapRepository.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": id.Hex(),
	})

	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		logger.Warn("roadmap not found for deletion")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("roadmap not found"))
	}

	return nil
}

func (r *RoadmapMongoRepository) CreateMaterial(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, material *entities.Material) (*entities.Material, error) {
	const op = "RoadmapRepository.CreateMaterial"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"roadmap_id":  roadmapID.Hex(),
		"node_id":     nodeID,
		"material_id": material.ID,
	})

	if material.ID == uuid.Nil {
		material.ID = uuid.New()
	}

	if material.CreatedAt.IsZero() {
		material.CreatedAt = time.Now()
	}

	if material.UpdatedAt.IsZero() {
		material.UpdatedAt = time.Now()
	}

	filter := bson.M{
		"id":       roadmapID,
		"nodes.id": nodeID,
	}

	update := bson.M{
		"$push": bson.M{
			"nodes.$.materials": material,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.WithError(err).Error("failed to create material")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		logger.Warn("roadmap or node not found for material creation")
		return nil, fmt.Errorf("roadmap or node not found")
	}

	logger.Info("successfully created material")
	return material, nil
}

func (r *RoadmapMongoRepository) DeleteMaterial(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, materialID uuid.UUID) error {
	const op = "RoadmapRepository.DeleteMaterial"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"roadmap_id":  roadmapID.Hex(),
		"node_id":     nodeID,
		"material_id": materialID,
	})

	filter := bson.M{"id": roadmapID}
	update := bson.M{
		"$pull": bson.M{
			"nodes.$[node].materials": bson.M{
				"id": materialID,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"node.id": nodeID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.WithError(err).Error("failed to delete material")
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		logger.Warn("roadmap not found for material deletion")
		return fmt.Errorf("roadmap not found")
	}

	if result.ModifiedCount == 0 {
		logger.Warn("material not found for deletion")
		return fmt.Errorf("material not found")
	}

	logger.Info("successfully deleted material")
	return nil
}

func (r *RoadmapMongoRepository) GetMaterialByID(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID, materialID uuid.UUID) (*entities.Material, error) {
	const op = "RoadmapRepository.GetMaterialByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"roadmap_id":  roadmapID.Hex(),
		"node_id":     nodeID,
		"material_id": materialID,
	})

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"id": roadmapID}}},
		{{Key: "$unwind", Value: "$nodes"}},
		{{Key: "$match", Value: bson.M{"nodes.id": nodeID}}},
		{{Key: "$unwind", Value: "$nodes.materials"}},
		{{Key: "$match", Value: bson.M{"nodes.materials.id": materialID}}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$nodes.materials"}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.WithError(err).Error("failed to aggregate material by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logger.WithError(err).Warn("failed to close cursor")
		}
	}()

	var material entities.Material
	if cursor.Next(ctx) {
		if err := cursor.Decode(&material); err != nil {
			logger.WithError(err).Error("failed to decode material")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return &material, nil
	}

	logger.Info("material not found")
	return nil, nil
}

func (r *RoadmapMongoRepository) GetMaterialsByNode(ctx context.Context, roadmapID primitive.ObjectID, nodeID uuid.UUID) ([]*entities.Material, error) {
	const op = "RoadmapRepository.GetMaterialsByNode"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.Hex(),
		"node_id":    nodeID,
	})

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"id": roadmapID}}},
		{{Key: "$unwind", Value: "$nodes"}},
		{{Key: "$match", Value: bson.M{"nodes.id": nodeID}}},
		{{Key: "$unwind", Value: "$nodes.materials"}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$nodes.materials"}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.WithError(err).Error("failed to aggregate materials by node")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logger.WithError(err).Warn("failed to close cursor")
		}
	}()

	var materials []*entities.Material
	if err := cursor.All(ctx, &materials); err != nil {
		logger.WithError(err).Error("failed to decode materials")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("count", len(materials)).Info("retrieved materials by node")
	return materials, nil
}
