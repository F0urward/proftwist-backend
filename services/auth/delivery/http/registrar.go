package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/auth"
)

type AuthHttpRegistrar struct {
	handlers auth.Handlers
}

func NewAuthHttpRegistrar(handlers auth.Handlers) httpServer.HttpRegistrar {
	return &AuthHttpRegistrar{
		handlers: handlers,
	}
}

func (r *AuthHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/auth/register", http.HandlerFunc(r.handlers.Register)).Methods("POST")
	s.MUX.Handle("/api/v1/auth/login", http.HandlerFunc(r.handlers.Login)).Methods("POST")
	s.MUX.Handle("/api/v1/auth/logout", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Logout))).Methods("POST")
	s.MUX.Handle("/api/v1/auth/me", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetMe))).Methods("GET")
	s.MUX.Handle("/api/v1/auth/{user_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetByID))).Methods("GET")
	s.MUX.Handle("/api/v1/auth", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Update))).Methods("PUT")
	s.MUX.Handle("/api/v1/auth/avatar", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.UploadAvatar))).Methods("POST")
	s.MUX.Handle("/api/v1/auth/users/search", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.SearchUsers))).Methods("GET")
	s.MUX.Handle("/api/v1/auth/vk/link", http.HandlerFunc(r.handlers.VKOauthLink)).Methods("GET")
	s.MUX.Handle("/api/v1/auth/vk/callback", http.HandlerFunc(r.handlers.VKOAuthCallback)).Methods("GET")
}
