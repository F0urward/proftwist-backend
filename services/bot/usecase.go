package bot

import (
	"context"

	"github.com/F0urward/proftwist-backend/services/bot/dto"
)

type Usecase interface {
	HandleBotTrigger(ctx context.Context, event dto.MessageForBotEvent) error
}
