package http

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/F0urward/proftwist-backend/internal/utils"
)

type WebSocketHandler struct {
	wsServer *websocket.Server
}

func NewWebSocketHandler(wsServer *websocket.Server) *WebSocketHandler {
	return &WebSocketHandler{
		wsServer: wsServer,
	}
}

func (h *WebSocketHandler) WebSocketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
		if !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if err := h.wsServer.HandleWebSocket(w, r, userID); err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
	})
}
