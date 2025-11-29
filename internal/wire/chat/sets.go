package chat

import (
	"github.com/google/wire"

	chatGRPCHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/grpc"
	chatHTTPHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
	chatWSHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/ws"
	chatRepository "github.com/F0urward/proftwist-backend/services/chat/repository"
	chatUsecase "github.com/F0urward/proftwist-backend/services/chat/usecase"

	notificationHandlers "github.com/F0urward/proftwist-backend/services/notification/delivery/broker"
	notificationUsecase "github.com/F0urward/proftwist-backend/services/notification/usecase"

	botHandlers "github.com/F0urward/proftwist-backend/services/bot/delivery/broker"
	botRepository "github.com/F0urward/proftwist-backend/services/bot/repository"
	botUsecase "github.com/F0urward/proftwist-backend/services/bot/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	chatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	friendClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/friendclient"
	gigachatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"

	wsServerHTTPHandlers "github.com/F0urward/proftwist-backend/internal/server/ws/http"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker/kafka"
)

var ChatSet = wire.NewSet(
	chatRepository.NewChatPostgresRepository,
	chatUsecase.NewChatUsecase,
	chatHTTPHandlers.NewChatHandler,
	chatGRPCHandlers.NewChatServer,
	chatGRPCHandlers.NewChatGrpcRegistrar,
	chatWSHandlers.NewChatWsHandlers,
	chatWSHandlers.NewChatWsRegistrar,
)

var WsSet = wire.NewSet(
	wsServerHTTPHandlers.NewWebSocketHandler,
)

var NotificationSet = wire.NewSet(
	notificationHandlers.NewNotificationHandlers,
	notificationUsecase.NewNotificationUsecase,
)

var BotSet = wire.NewSet(
	botHandlers.NewBotHandlers,
	botUsecase.NewBotUsecase,
	botRepository.NewGigachatWebapi,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	authClient.NewAuthClient,
	friendClient.NewFriendClient,
	chatClient.NewChatClient,
	gigachatClient.NewGigaChatClient,
)

var BrokerSet = wire.NewSet(
	kafka.NewConsumer,
	kafka.NewProducer,
)
