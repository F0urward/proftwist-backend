package adapter

import (
	"context"
	"time"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	botDTO "github.com/F0urward/proftwist-backend/services/bot/dto"
	"github.com/F0urward/proftwist-backend/services/chat"
)

type BotPublisher struct {
	producer broker.Producer
}

func NewBotPublisher(producer broker.Producer) chat.BotPublisher {
	return &BotPublisher{producer: producer}
}

func (b *BotPublisher) PublishMessageForBot(ctx context.Context, chatID, chatTitle, content string) error {
	event := botDTO.MessageForBotEvent{
		Type:       botDTO.MessageForBotType,
		ChatID:     chatID,
		ChatTitle:  chatTitle,
		Content:    content,
		ReceivedAt: time.Now(),
	}

	data, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	return b.producer.Publish(ctx, chatID, data)
}
