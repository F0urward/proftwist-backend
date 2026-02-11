package chat

import (
	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker/kafka"
	"github.com/F0urward/proftwist-backend/internal/metrics"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	wsServer "github.com/F0urward/proftwist-backend/internal/server/ws"
	wsServerHTTPHandlers "github.com/F0urward/proftwist-backend/internal/server/ws/http"
	chat "github.com/F0urward/proftwist-backend/services/chat"
	chatAdapters "github.com/F0urward/proftwist-backend/services/chat/adapter"
	chatHTTPHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() metrics.Metrics {
	reg := prometheus.NewRegistry()

	wrapped := prometheus.WrapRegistererWith(
		prometheus.Labels{
			"service": "proftwist-chat-service",
		},
		reg,
	)

	return metrics.NewMetrics(reg, wrapped)
}

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

func ProvideNotificationPublisher(cfg *config.Config) chat.NotificationPublisher {
	producerConfig := kafka.ProducerConfig{
		Broker: cfg.Kafka.Broker,
		Topic:  cfg.Kafka.Producers.Notification.Topic,
	}
	producer := kafka.NewProducer(producerConfig)
	return chatAdapters.NewNotificationPublisher(producer)
}

func ProvideBotPublisher(cfg *config.Config) chat.BotPublisher {
	producerConfig := kafka.ProducerConfig{
		Broker: cfg.Kafka.Broker,
		Topic:  cfg.Kafka.Producers.Bot.Topic,
	}
	producer := kafka.NewProducer(producerConfig)
	return chatAdapters.NewBotPublisher(producer)
}

func ProvideNotificationConsumerConfig(cfg *config.Config) kafka.ConsumerConfig {
	return kafka.ConsumerConfig{
		Broker:  cfg.Kafka.Broker,
		Topic:   cfg.Kafka.Consumers.Notification.Topic,
		GroupID: cfg.Kafka.Consumers.Notification.GroupID,
	}
}

func ProvideBotConsumerConfig(cfg *config.Config) kafka.ConsumerConfig {
	return kafka.ConsumerConfig{
		Broker:  cfg.Kafka.Broker,
		Topic:   cfg.Kafka.Consumers.Bot.Topic,
		GroupID: cfg.Kafka.Consumers.Bot.GroupID,
	}
}
