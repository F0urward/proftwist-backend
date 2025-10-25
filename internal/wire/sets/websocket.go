package sets

import (
	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/google/wire"
)

var WebSocketSet = wire.NewSet(
	websocket.NewWebSocketServer,
)
