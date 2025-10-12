package sets

import (
	vkClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/vkclient"
	mongoClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/mongo"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
	redisClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/redis"
	"github.com/google/wire"
)

var CommonSet = wire.NewSet(
	db.NewDatabase,
	mongoClient.NewClient,
	mongoClient.NewDatabase,
	redisClient.NewClient,
	vkClient.NewVKClient,
)
