// seed/main.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/F0urward/proftwist-backend/services/chat/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "roadmaps"
)

type RoadmapJSON struct {
	Name string           `json:"name"`
	Data entities.Roadmap `json:"data"`
}

type RoadmapCollection struct {
	Roadmaps []RoadmapJSON `json:"roadmaps"`
}

type RoadmapInfo struct {
	ID        string
	RoadmapID string
}

func SeedData(ctx context.Context, pgDB *sql.DB, mongoDB *mongo.Database, cfg *config.Config) error {
	const op = "SeedData"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	chatRepo := repository.NewChatPostgresRepository(pgDB)

	roadmapData, err := loadRoadmapData()
	if err != nil {
		logger.WithError(err).Error("failed to load roadmap data")
		return fmt.Errorf("%s: %w", op, err)
	}

	collection := mongoDB.Collection(collectionName)

	for _, roadmapJSON := range roadmapData.Roadmaps {
		roadmapLogger := logger.WithField("roadmap_name", roadmapJSON.Name)

		roadmapInfo, err := getRoadmapInfoByName(ctx, pgDB, roadmapJSON.Name)
		if err != nil {
			roadmapLogger.WithError(err).Warn("failed to get roadmap_info")
			continue
		}

		if roadmapInfo == nil {
			roadmapLogger.Warn("no roadmap_info found, skipping")
			continue
		}

		mongoID, err := primitive.ObjectIDFromHex(roadmapInfo.RoadmapID)
		if err != nil {
			roadmapLogger.WithError(err).Warn("invalid roadmap_id format, generating new ID")
			mongoID = primitive.NewObjectID()
		}

		roadmapDoc := &entities.Roadmap{
			ID:        mongoID,
			Nodes:     roadmapJSON.Data.Nodes,
			Edges:     roadmapJSON.Data.Edges,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = collection.InsertOne(ctx, roadmapDoc)
		if err != nil {
			roadmapLogger.WithError(err).Error("failed to insert roadmap")
			return fmt.Errorf("%s: failed to insert roadmap %s: %w", op, roadmapJSON.Name, err)
		}

		createdChatsCount, err := createGroupChatsForNodes(ctx, chatRepo, roadmapDoc.Nodes)
		if err != nil {
			roadmapLogger.WithError(err).Warn("failed to create group chats for nodes")
		} else {
			roadmapLogger.WithField("chats_created", createdChatsCount).Info("group chats created for nodes")
		}

		roadmapLogger.WithFields(map[string]interface{}{
			"roadmap_id":  roadmapInfo.RoadmapID,
			"mongo_id":    roadmapDoc.ID.Hex(),
			"nodes_count": len(roadmapDoc.Nodes),
			"edges_count": len(roadmapDoc.Edges),
			"chats_count": createdChatsCount,
		}).Info("successfully seeded roadmap")
	}

	logger.WithField("total_roadmaps", len(roadmapData.Roadmaps)).Info("seeding completed")
	return nil
}

func createGroupChatsForNodes(ctx context.Context, chatRepo chat.Repository, nodes []entities.RoadmapNode) (int, error) {
	const op = "createGroupChatsForNodes"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	createdCount := 0

	for _, node := range nodes {
		nodeLogger := logger.WithFields(map[string]interface{}{
			"node_id":   node.ID,
			"node_name": node.Data.Label,
		})
		existingChat, err := chatRepo.GetGroupChatByNode(ctx, node.ID.String())
		if err != nil {
			nodeLogger.WithError(err).Warn("failed to check existing group chat")
			continue
		}

		if existingChat != nil {
			nodeLogger.WithField("chat_id", existingChat.ID).Info("group chat already exists")
			continue
		}

		title := node.Data.Label
		roadmapNodeID := node.ID.String()

		groupChat := &entities.GroupChat{
			Title:         &title,
			AvatarURL:     nil,
			RoadmapNodeID: &roadmapNodeID,
		}

		createdChat, err := chatRepo.CreateGroupChat(ctx, groupChat)
		if err != nil {
			nodeLogger.WithError(err).Warn("failed to create group chat")
			continue
		}

		nodeLogger.WithFields(map[string]interface{}{
			"chat_id": createdChat.ID,
		}).Info("created group chat")
		createdCount++
	}

	logger.WithField("created_chats", createdCount).Info("group chats creation completed")
	return createdCount, nil
}

func getRoadmapInfoByName(ctx context.Context, db *sql.DB, name string) (*RoadmapInfo, error) {
	const op = "getRoadmapInfoByName"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":           op,
		"roadmap_name": name,
	})

	query := `SELECT id, roadmap_id FROM roadmap_info WHERE name = $1`

	var roadmapInfo RoadmapInfo
	err := db.QueryRowContext(ctx, query, name).Scan(&roadmapInfo.ID, &roadmapInfo.RoadmapID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("roadmap_info not found")
			return nil, nil
		}
		logger.WithError(err).Error("database query failed")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("roadmap_id", roadmapInfo.RoadmapID).Info("roadmap_info found")
	return &roadmapInfo, nil
}

func loadRoadmapData() (*RoadmapCollection, error) {
	const op = "loadRoadmapData"

	jsonPath := filepath.Join("data", "roadmaps.json")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read roadmap data file: %w", op, err)
	}

	var roadmapCollection RoadmapCollection
	if err := json.Unmarshal(data, &roadmapCollection); err != nil {
		return nil, fmt.Errorf("%s: failed to unmarshal roadmap data: %w", op, err)
	}

	return &roadmapCollection, nil
}
