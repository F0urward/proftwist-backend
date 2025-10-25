package chat

import (
	"net/http"
)

type ChatHandlerInterface interface {
	RegisterRoutes(mux *http.ServeMux)
	CreateChat(w http.ResponseWriter, r *http.Request)
	GetUserChats(w http.ResponseWriter, r *http.Request)
	GetChat(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	GetChatMessages(w http.ResponseWriter, r *http.Request)
	AddMember(w http.ResponseWriter, r *http.Request)
	RemoveMember(w http.ResponseWriter, r *http.Request)
	DeleteChat(w http.ResponseWriter, r *http.Request)
	JoinChannel(w http.ResponseWriter, r *http.Request)
}
