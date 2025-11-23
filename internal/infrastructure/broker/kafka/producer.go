package kafka

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Broker string
	Topic  string
}

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(cfg ProducerConfig) broker.Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP([]string{cfg.Broker}...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Producer{writer: writer, topic: cfg.Topic}
}

func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Topic: p.topic,
		Key:   []byte(key),
		Value: value,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
