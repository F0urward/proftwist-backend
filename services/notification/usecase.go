package notification

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/notification/dto"
)

type Usecase interface {
	HandleMessageSent(ctx context.Context, event dto.MessageSentEvent) error
	HandleTyping(ctx context.Context, event dto.TypingEvent) error
	HandleUserJoined(ctx context.Context, event dto.UserJoinedEvent) error
	HandleUserLeft(ctx context.Context, event dto.UserLeftEvent) error
}
