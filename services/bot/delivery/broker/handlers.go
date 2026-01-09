package kafka

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/bot"
	"github.com/F0urward/proftwist-backend/services/bot/dto"
)

type BotHandlers struct {
	botUC bot.Usecase
}

func NewBotHandlers(botUC bot.Usecase) bot.Handlers {
	return &BotHandlers{
		botUC: botUC,
	}
}

func (h *BotHandlers) HandleMessage(ctx context.Context, msg broker.Message) error {
	const op = "BotHandlers.HandleMessage"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithField("message_key", msg.Key)

	logger.Info("received new message")

	var baseEvent dto.BaseEvent
	if err := baseEvent.UnmarshalJSON(msg.Value); err != nil {
		logger.WithError(err).Error("failed to unmarshal base event")
		return err
	}

	logger = logger.WithField("event_type", baseEvent.Type)

	switch baseEvent.Type {
	case dto.MessageForBotType:
		var event dto.MessageForBotEvent
		if err := event.UnmarshalJSON(msg.Value); err != nil {
			logger.WithError(err).Error("failed to unmarshal MessageForBotEvent")
			return err
		}

		logger.WithFields(map[string]interface{}{
			"chat_id":        event.ChatID,
			"content_length": len(event.Content),
		}).Info("handling MessageForBotEvent")

		if err := h.botUC.HandleBotTrigger(ctx, event); err != nil {
			logger.WithError(err).Error("failed to handle MessageForBotEvent")
			return err
		}

		logger.Info("successfully handled MessageForBotEvent")
		return nil

	default:
		logger.WithField("event_type", baseEvent.Type).Warn("received message with unknown event type, ignoring")
		return nil
	}
}
