package category

import (
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
)

func AllHttpRegistrars(
	categoryHttpRegistrar httpServer.HttpRegistrar,
) []httpServer.HttpRegistrar {
	return []httpServer.HttpRegistrar{
		categoryHttpRegistrar,
	}
}
