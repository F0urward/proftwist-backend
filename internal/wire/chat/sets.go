package chat

import (
	"github.com/google/wire"

	chatAdapter "github.com/F0urward/proftwist-backend/services/chat/adapter"
	chatGRPCHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/grpc"
	chatHTTPHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
	chatWSHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/ws"
	chatRepository "github.com/F0urward/proftwist-backend/services/chat/repository"
	chatUsecase "github.com/F0urward/proftwist-backend/services/chat/usecase"

	notificationHandlers "github.com/F0urward/proftwist-backend/services/notification/delivery/broker"
	notificationUsecase "github.com/F0urward/proftwist-backend/services/notification/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	friendClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/friendclient"
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
	chatAdapter.NewBrokerNotifier,
)

var WsSet = wire.NewSet(
	wsServerHTTPHandlers.NewWebSocketHandler,
)

var NotificationSet = wire.NewSet(
	notificationHandlers.NewNotificationHandlers,
	notificationUsecase.NewNotificationUsecase,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	authClient.NewAuthClient,
	friendClient.NewFriendClient,
)

var BrokerSet = wire.NewSet(
	kafka.NewConsumer,
	kafka.NewProducer,
)
