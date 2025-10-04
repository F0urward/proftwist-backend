package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
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

func NewRoadmapRepository(db *mongo.Database) *RoadmapRepository {
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

	var roadmapDTOs []*dto.RoadmapDTO
	if err = cursor.All(ctx, &roadmapDTOs); err != nil {
		logger.WithError(err).Error("failed to decode roadmaps")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	roadmaps := make([]*entities.Roadmap, len(roadmapDTOs))
	for i, roadmapDTO := range roadmapDTOs {
		roadmaps[i] = dto.DTOToEntity(roadmapDTO)
	}

	return roadmaps, nil
}

func (r *RoadmapRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Roadmap, error) {
	const op = "RoadmapRepository.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": id.Hex(),
	})

	var roadmapDTO dto.RoadmapDTO

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&roadmapDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.WithError(err).Error("failed to get roadmap by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dto.DTOToEntity(&roadmapDTO), nil
}

func (r *RoadmapRepository) GetByAuthorID(ctx context.Context, authorID uuid.UUID) ([]*entities.Roadmap, error) {
	const op = "RoadmapRepository.GetByAuthorID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"author_id": authorID.String(),
	})

	cursor, err := r.collection.Find(ctx, bson.M{"author_id": authorID})
	if err != nil {
		logger.WithError(err).Error("failed to find roadmaps by author ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			logger.WithError(err).Warn("failed to close cursor")
		}
	}()

	roadmaps, err := r.decodeRoadmapsFromCursor(ctx, cursor)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return roadmaps, nil
}

func (r *RoadmapRepository) Create(ctx context.Context, roadmap *entities.Roadmap) error {
	const op = "RoadmapRepository.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.Hex(),
		"author_id":  roadmap.AuthorID.String(),
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

	if roadmap.Nodes != nil {
		for i := range roadmap.Nodes {
			if roadmap.Nodes[i].ID == uuid.Nil {
				roadmap.Nodes[i].ID = uuid.New()
			}

			if roadmap.Nodes[i].Children != nil {
				for j := range roadmap.Nodes[i].Children {
					if roadmap.Nodes[i].Children[j].ID == uuid.Nil {
						roadmap.Nodes[i].Children[j].ID = uuid.New()
					}
				}
			}
		}
	}

	roadmapDTO := dto.EntityToDTO(roadmap)

	_, err := r.collection.InsertOne(ctx, roadmapDTO)
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

	roadmapDTO := dto.EntityToDTO(roadmap)

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": roadmap.ID},
		roadmapDTO,
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

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
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

func (r *RoadmapRepository) SearchByTitle(ctx context.Context, title string) ([]*entities.Roadmap, error) {
	const op = "RoadmapRepository.SearchByTitle"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":    op,
		"title": title,
	})

	filter := bson.M{
		"title": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		logger.WithError(err).Error("failed to search roadmaps by title")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			logger.WithError(err).Warn("failed to close cursor")
		}
	}()

	roadmaps, err := r.decodeRoadmapsFromCursor(ctx, cursor)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return roadmaps, nil
}

func (r *RoadmapRepository) decodeRoadmapsFromCursor(ctx context.Context, cursor *mongo.Cursor) ([]*entities.Roadmap, error) {
	var roadmapDTOs []*dto.RoadmapDTO
	if err := cursor.All(ctx, &roadmapDTOs); err != nil {
		log.Printf("Failed to decode roadmaps: %v", err)
		return nil, fmt.Errorf("failed to decode roadmaps: %v", err)
	}

	roadmaps := make([]*entities.Roadmap, len(roadmapDTOs))
	for i, roadmapDTO := range roadmapDTOs {
		roadmaps[i] = dto.DTOToEntity(roadmapDTO)
	}

	return roadmaps, nil
}

func (r *RoadmapRepository) UpdatePrivacy(ctx context.Context, id primitive.ObjectID, isPublic bool) error {
	const op = "RoadmapRepository.UpdatePrivacy"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": id.Hex(),
		"is_public":  isPublic,
	})

	update := bson.M{
		"$set": bson.M{
			"is_public":  isPublic,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap privacy")
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		logger.Warn("roadmap not found for privacy update")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("roadmap not found"))
	}

	return nil
}
