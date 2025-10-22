package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
)

func NewDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	return client.Database(cfg.Mongo.DBName)
}

func NewClient(cfg *config.Config) *mongo.Client {
	const op = "mongo.NewClient"
	logger := logctx.GetLogger(context.Background()).WithField("op", op)

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
		logger.WithError(err).Error("failed to connect to mongo")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		_ = client.Disconnect(ctx)
		logger.WithError(err).Error("cannot ping mongo instance")
	}

	return client
}
