package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/F0urward/proftwist-backend/config"
)

func New(cfg *config.Config) (*mongo.Client, error) {
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
		return nil, fmt.Errorf("error creating and connecting MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}

	return client, nil
}
