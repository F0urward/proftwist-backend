package sets

import (
	"github.com/F0urward/proftwist-backend/services/chat/delivery/http"
	"github.com/F0urward/proftwist-backend/services/chat/repository"
	"github.com/F0urward/proftwist-backend/services/chat/usecase"
	"github.com/google/wire"
)

var ChatSet = wire.NewSet(
	repository.NewChatPostgresRepository,
	usecase.NewChatUseCase,
	http.NewChatHandler,
	http.NewWebSocketHandler,
	http.NewWebSocketIntegration,
)
