package main

import (
	"context"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/db/mongo"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

func main() {
	const op = "seed.main"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	cfg := config.New()

	pgDB := postgres.NewDatabase(cfg)

	mongoClient := mongo.NewClient(cfg)

	mongoDB := mongo.NewDatabase(mongoClient, cfg)

	if err := SeedData(context.Background(), pgDB, mongoDB, cfg); err != nil {
		logger.WithError(err).Error("failed to seed data")
	} else {
		logger.Info("roadmaps successfully generated")
	}
}
