package friend

import (
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
)

func AllHttpRegistrars(
	friendHttpRegistrar httpServer.HttpRegistrar,
) []httpServer.HttpRegistrar {
	return []httpServer.HttpRegistrar{
		friendHttpRegistrar,
	}
}
