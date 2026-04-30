package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/ai"
)

type AIHttpRegistrar struct {
	handlers ai.Handlers
}

func NewAIHttpRegistrar(handlers ai.Handlers) *AIHttpRegistrar {
	return &AIHttpRegistrar{
		handlers: handlers,
	}
}

func (r *AIHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/ai/roadmap-node-description", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GenerateRoadmapNodeDescription))).Methods("POST")
}
