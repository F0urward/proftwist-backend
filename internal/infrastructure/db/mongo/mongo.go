package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/F0urward/proftwist-backend/config"
)

func NewDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	return client.Database(cfg.Mongo.DBName)
}

func NewClient(cfg *config.Config) *mongo.Client {
	dsn := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin",
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
		cfg.Mongo.DBName,
	)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		_ = client.Disconnect(ctx)
		log.Fatalf("cannot ping mongo instance: %v", err)
	}

	return client
}
