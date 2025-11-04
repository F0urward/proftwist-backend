package chat

import (
	"net/http"
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
