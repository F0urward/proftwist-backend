package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "roadmaps"
)

type RoadmapRepository struct {
	collection *mongo.Collection
}

func NewRoadmapRepository(db *mongo.Database) roadmap.Repository {
	return &RoadmapRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *RoadmapRepository) GetAll(ctx context.Context) ([]*entities.Roadmap, error) {
	const op = "RoadmapRepository.GetAll"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.WithError(err).Error("failed to find roadmaps")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			logger.WithError(closeErr).Warn("failed to close cursor")
		}
	}()

	var roadmaps []*entities.Roadmap
	if err = cursor.All(ctx, &roadmaps); err != nil {
		logger.WithError(err).Error("failed to decode roadmaps")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return roadmaps, nil
}

func (r *RoadmapRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Roadmap, error) {
	const op = "RoadmapRepository.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": id.Hex(),
	})

	var roadmap entities.Roadmap

	// Поиск по полю id вместо _id
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

func (r *RoadmapRepository) Create(ctx context.Context, roadmap *entities.Roadmap) error {
	const op = "RoadmapRepository.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.Hex(),
	})

	// Генерация ID если не установлен
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

func (r *RoadmapRepository) Update(ctx context.Context, roadmap *entities.Roadmap) error {
	const op = "RoadmapRepository.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.Hex(),
	})

	roadmap.UpdatedAt = time.Now()

	// Обновление по полю id вместо _id
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

func (r *RoadmapRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	const op = "RoadmapRepository.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": id.Hex(),
	})

	// Удаление по полю id вместо _id
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
