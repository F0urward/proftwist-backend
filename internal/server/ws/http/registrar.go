package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
)

type WebSocketHttpRegistrar struct {
	handler *WebSocketHandler
}

func NewWebSocketHttpRegistrar(handler *WebSocketHandler) httpServer.HttpRegistrar {
	return &WebSocketHttpRegistrar{
		handler: handler,
	}
}

func (r *WebSocketHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/chats/ws", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handler.HandleConnection))).Methods("GET")
}
