package sets

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	vkClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/vkclient"
	awsClient "github.com/F0urward/proftwist-backend/internal/infrastructure/db/aws"
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
	awsClient.NewClient,
	vkClient.NewVKClient,
	gigachatclient.NewGigaChatClient,
)
