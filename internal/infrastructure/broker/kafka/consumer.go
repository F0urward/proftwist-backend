package kafka

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	Broker  string
	Topic   string
	GroupID string
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg ConsumerConfig) broker.Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Broker},
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) ReadMessage(ctx context.Context) (broker.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return broker.Message{}, err
	}

	return broker.Message{
		Topic: msg.Topic,
		Key:   msg.Key,
		Value: msg.Value,
	}, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
