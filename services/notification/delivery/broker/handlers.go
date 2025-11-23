package kafka

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/notification"
	"github.com/F0urward/proftwist-backend/services/notification/dto"
)

type NotificationHandlers struct {
	notificationUC notification.Usecase
}

func NewNotificationHandlers(notificationUC notification.Usecase) notification.Handlers {
	return &NotificationHandlers{notificationUC: notificationUC}
}

func (h *NotificationHandlers) HandleMessage(ctx context.Context, msg broker.Message) error {
	const op = "NotificationHandlers.HandleMessage"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("message_key", msg.Key)

	logger.Info("received new message")

	var baseEvent dto.BaseEvent
	if err := baseEvent.UnmarshalJSON(msg.Value); err != nil {
		logger.WithError(err).Error("failed to unmarshal base event")
		return err
	}

	logger = logger.WithField("event_type", baseEvent.Type)

	switch baseEvent.Type {
	case dto.MessagePublishedType:
		var event dto.MessageSentEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal MessageSentEvent")
			return err
		}
		logger.Info("handling MessageSentEvent")
		return h.notificationUC.HandleMessageSent(ctx, event)

	case dto.UserTypingType:
		var event dto.TypingEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal TypingEvent")
			return err
		}
		logger.Info("handling TypingEvent")
		return h.notificationUC.HandleTyping(ctx, event)

	case dto.UserJoinedType:
		var event dto.UserJoinedEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal UserJoinedEvent")
			return err
		}
		logger.Info("handling UserJoinedEvent")
		return h.notificationUC.HandleUserJoined(ctx, event)

	case dto.UserLeftType:
		var event dto.UserLeftEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal UserLeftEvent")
			return err
		}
		logger.Info("handling UserLeftEvent")
		return h.notificationUC.HandleUserLeft(ctx, event)
	}

	logger.Warn("received message with unknown event type, ignoring")
	return nil
}
