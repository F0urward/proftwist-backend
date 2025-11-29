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

		logger.WithFields(map[string]interface{}{
			"chat_id":     event.ChatID,
			"message_id":  event.MessageID,
			"users_count": len(event.UserIDs),
		}).Info("handling MessageSentEvent")

		if err := h.notificationUC.HandleMessageSent(ctx, event); err != nil {
			logger.WithError(err).Error("failed to handle MessageSentEvent")
			return err
		}

		logger.Info("successfully handled MessageSentEvent")
		return nil

	case dto.UserTypingType:
		var event dto.TypingEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal TypingEvent")
			return err
		}

		logger.WithFields(map[string]interface{}{
			"chat_id":     event.ChatID,
			"user_id":     event.UserID,
			"typing":      event.Typing,
			"users_count": len(event.UserIDs),
		}).Info("handling TypingEvent")

		if err := h.notificationUC.HandleTyping(ctx, event); err != nil {
			logger.WithError(err).Error("failed to handle TypingEvent")
			return err
		}

		logger.Info("successfully handled TypingEvent")
		return nil

	case dto.UserJoinedType:
		var event dto.UserJoinedEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal UserJoinedEvent")
			return err
		}

		logger.WithFields(map[string]interface{}{
			"chat_id":     event.ChatID,
			"user_id":     event.UserID,
			"users_count": len(event.UserIDs),
		}).Info("handling UserJoinedEvent")

		if err := h.notificationUC.HandleUserJoined(ctx, event); err != nil {
			logger.WithError(err).Error("failed to handle UserJoinedEvent")
			return err
		}

		logger.Info("successfully handled UserJoinedEvent")
		return nil

	case dto.UserLeftType:
		var event dto.UserLeftEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal UserLeftEvent")
			return err
		}

		logger.WithFields(map[string]interface{}{
			"chat_id":     event.ChatID,
			"user_id":     event.UserID,
			"users_count": len(event.UserIDs),
		}).Info("handling UserLeftEvent")

		if err := h.notificationUC.HandleUserLeft(ctx, event); err != nil {
			logger.WithError(err).Error("failed to handle UserLeftEvent")
			return err
		}

		logger.Info("successfully handled UserLeftEvent")
		return nil

	default:
		logger.Warn("received message with unknown event type, ignoring")
		return nil
	}
}
