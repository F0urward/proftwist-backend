package worker

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	"github.com/F0urward/proftwist-backend/services/bot"
)

type BotWorker struct {
	consumer broker.Consumer
	h        bot.Handlers
}

func NewBotWorker(consumer broker.Consumer, h bot.Handlers) *BotWorker {
	return &BotWorker{consumer: consumer, h: h}
}

func (w *BotWorker) Start(ctx context.Context) {
	go func() {
		for {
			msg, err := w.consumer.ReadMessage(ctx)
			if err != nil {
				continue
			}
			_ = w.h.HandleMessage(ctx, msg)
		}
	}()
}
