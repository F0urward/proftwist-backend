package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/chat"
)

type ChatHttpRegistrar struct {
	handlers chat.Handlers
}

func NewChatHttpRegistrar(handlers chat.Handlers) httpServer.HttpRegistrar {
	return &ChatHttpRegistrar{
		handlers: handlers,
	}
}

func (r *ChatHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/chats/group/node/{node_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetGroupChatByNode))).Methods("GET")
	s.MUX.Handle("/api/v1/chats/group", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetGroupChatsByUser))).Methods("GET")
	s.MUX.Handle("/api/v1/chats/group/{chat_id}/members", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetGroupChatMembers))).Methods("GET")
	s.MUX.Handle("/api/v1/chats/group/{chat_id}/messages", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetGroupChatMessages))).Methods("GET")
	s.MUX.Handle("/api/v1/chats/group/{chat_id}/join", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.JoinGroupChat))).Methods("POST")
	s.MUX.Handle("/api/v1/chats/group/{chat_id}/leave", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.LeaveGroupChat))).Methods("POST")

	s.MUX.Handle("/api/v1/chats/direct", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetDirectChatsByUser))).Methods("GET")
	s.MUX.Handle("/api/v1/chats/direct/{chat_id}/messages", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetDirectChatMessages))).Methods("GET")
}
