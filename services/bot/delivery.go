package bot

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
)

type Handlers interface {
	HandleMessage(ctx context.Context, msg broker.Message) error
}
