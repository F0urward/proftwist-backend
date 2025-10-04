package sets

import (
	"log"

	"github.com/F0urward/proftwist-backend/config"
	mongoClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/mongo"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProvideMongoClient(cfg *config.Config) *mongo.Client {
	client, err := mongoClient.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client
}

var CommonSet = wire.NewSet(
	db.New,
	ProvideMongoClient,
	mongoClient.NewDatabase,
)
