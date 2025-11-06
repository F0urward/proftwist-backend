package chat

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/F0urward/proftwist-backend/internal/server/websocket/dto"
)

type Handlers interface {
	GetGroupChatByNode(w http.ResponseWriter, r *http.Request)
	GetGroupChatsByUser(w http.ResponseWriter, r *http.Request)
	GetGroupChatMembers(w http.ResponseWriter, r *http.Request)
	GetGroupChatMessages(w http.ResponseWriter, r *http.Request)
	JoinGroupChat(w http.ResponseWriter, r *http.Request)
	LeaveGroupChat(w http.ResponseWriter, r *http.Request)

	GetDirectChatsByUser(w http.ResponseWriter, r *http.Request)
	GetDirectChatMembers(w http.ResponseWriter, r *http.Request)
	GetDirectChatMessages(w http.ResponseWriter, r *http.Request)
}

type WSHandlers interface {
	HandleSendMessage(client *websocket.Client, msg dto.WebSocketMessage) error
	HandleTyping(client *websocket.Client, msg dto.WebSocketMessage) error
}
