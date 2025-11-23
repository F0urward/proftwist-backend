package broker

import "context"

type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

type Producer interface {
	Publish(ctx context.Context, key string, value []byte) error
	Close() error
}

type Consumer interface {
	ReadMessage(ctx context.Context) (Message, error)
	Close() error
}

type Handler interface {
	Handle(ctx context.Context, msg Message) error
}
