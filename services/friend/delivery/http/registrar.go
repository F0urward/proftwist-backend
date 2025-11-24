package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/friend"
)

type FriendHttpRegistrar struct {
	handlers friend.Handlers
}

func NewFriendHttpRegistrar(handlers friend.Handlers) httpServer.HttpRegistrar {
	return &FriendHttpRegistrar{
		handlers: handlers,
	}
}

func (r *FriendHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/friends", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetFriends))).Methods("GET")
	s.MUX.Handle("/api/v1/friends/{friend_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.DeleteFriend))).Methods("DELETE")

	s.MUX.Handle("/api/v1/friends/requests", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetFriendRequests))).Methods("GET")
	s.MUX.Handle("/api/v1/friends/requests", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.CreateFriendRequest))).Methods("POST")
	s.MUX.Handle("/api/v1/friends/requests/{request_id}/accept", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.AcceptFriendRequest))).Methods("POST")
	s.MUX.Handle("/api/v1/friends/requests/{request_id}/reject", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.RejectFriendRequest))).Methods("POST")
	s.MUX.Handle("/api/v1/friends/requests/{request_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.DeleteFriendRequest))).Methods("DELETE")

	s.MUX.Handle("/api/v1/friends/status/{target_user_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetFriendshipStatus))).Methods("GET")
}
