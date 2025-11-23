package worker

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	"github.com/F0urward/proftwist-backend/services/notification"
)

type NotificationWorker struct {
	consumer broker.Consumer
	h        notification.Handlers
}

func NewNotificationWorker(consumer broker.Consumer, h notification.Handlers) *NotificationWorker {
	return &NotificationWorker{consumer: consumer, h: h}
}

func (w *NotificationWorker) Start(ctx context.Context) {
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
