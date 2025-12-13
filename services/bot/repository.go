package bot

import (
	"context"
)

type GigachatWebapi interface {
	GetBotResponse(ctx context.Context, query, chatTitle string) (string, error)
}
