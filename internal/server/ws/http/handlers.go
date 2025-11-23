package http

import (
	"net/http"

	wsServer "github.com/F0urward/proftwist-backend/internal/server/ws"
	"github.com/F0urward/proftwist-backend/internal/utils"
)

type WebSocketHandler struct {
	wsServer *wsServer.WsServer
}

func NewWebSocketHandler(wsServer *wsServer.WsServer) *WebSocketHandler {
	return &WebSocketHandler{
		wsServer: wsServer,
	}
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if err := h.wsServer.HandleWebSocket(w, r, userID); err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
}
