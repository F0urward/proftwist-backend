package chat

import (
	"net/http"

	websocket "github.com/F0urward/proftwist-backend/internal/server/ws"
	"github.com/F0urward/proftwist-backend/internal/server/ws/dto"
)

type Handlers interface {
	GetGroupChatByNode(w http.ResponseWriter, r *http.Request)
	GetGroupChatsByUser(w http.ResponseWriter, r *http.Request)
	GetGroupChatMembers(w http.ResponseWriter, r *http.Request)
	GetGroupChatMessages(w http.ResponseWriter, r *http.Request)
	JoinGroupChat(w http.ResponseWriter, r *http.Request)
	LeaveGroupChat(w http.ResponseWriter, r *http.Request)

	GetDirectChatsByUser(w http.ResponseWriter, r *http.Request)
	GetDirectChatMessages(w http.ResponseWriter, r *http.Request)
}

type WSHandlers interface {
	HandleSendMessage(client *websocket.WsClient, msg dto.WebSocketMessage) error
	HandleTyping(client *websocket.WsClient, msg dto.WebSocketMessage) error
}
