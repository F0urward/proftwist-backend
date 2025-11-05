package chat

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/F0urward/proftwist-backend/internal/server/websocket/dto"
)

type Handlers interface {
	CreateChat(w http.ResponseWriter, r *http.Request)
	AddMember(w http.ResponseWriter, r *http.Request)
	RemoveMember(w http.ResponseWriter, r *http.Request)
	GetChatsByUser(w http.ResponseWriter, r *http.Request)
	GetChatMessages(w http.ResponseWriter, r *http.Request)
	JoinGroupChat(w http.ResponseWriter, r *http.Request)
	LeaveGroupChat(w http.ResponseWriter, r *http.Request)
}

type WSHandlers interface {
	HandleSendMessage(client *websocket.Client, msg dto.WebSocketMessage) error
	HandleTyping(client *websocket.Client, msg dto.WebSocketMessage) error
}
