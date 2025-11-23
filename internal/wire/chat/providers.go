package chat

import (
	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker/kafka"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	wsServer "github.com/F0urward/proftwist-backend/internal/server/ws"
	wsServerHTTPHandlers "github.com/F0urward/proftwist-backend/internal/server/ws/http"
	chat "github.com/F0urward/proftwist-backend/services/chat"
	chatHTTPHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
)

func AllHttpRegistrars(
	chatHandlers chat.Handlers,
	wsHandlers *wsServerHTTPHandlers.WebSocketHandler,
) []httpServer.HttpRegistrar {
	chatRegistrar := chatHTTPHandlers.NewChatHttpRegistrar(chatHandlers)
	wsRegistrar := wsServerHTTPHandlers.NewWebSocketHttpRegistrar(wsHandlers)

	return []httpServer.HttpRegistrar{
		chatRegistrar,
		wsRegistrar,
	}
}

func AllGrpcRegistrars(
	chatGrpcRegistrar grpcServer.GrpcRegistrar,
) []grpcServer.GrpcRegistrar {
	return []grpcServer.GrpcRegistrar{
		chatGrpcRegistrar,
	}
}

func AllWsRegistrars(
	chatWsRegistrar wsServer.WsRegistrar,
) []wsServer.WsRegistrar {
	return []wsServer.WsRegistrar{
		chatWsRegistrar,
	}
}

func ProvideNotificationConsumerConfig(cfg *config.Config) kafka.ConsumerConfig {
	return kafka.ConsumerConfig{
		Broker:  cfg.Kafka.Broker,
		Topic:   cfg.Kafka.Consumers.Notification.Topic,
		GroupID: cfg.Kafka.Consumers.Notification.GroupID,
	}
}

func ProvideNotificationProducerConfig(cfg *config.Config) kafka.ProducerConfig {
	return kafka.ProducerConfig{
		Broker: cfg.Kafka.Broker,
		Topic:  cfg.Kafka.Producers.Notification.Topic,
	}
}
