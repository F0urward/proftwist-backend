package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/roadmap"
)

type RoadmapHttpRegistrar struct {
	handlers roadmap.Handlers
}

func NewRoadmapHttpRegistrar(handlers roadmap.Handlers) httpServer.HttpRegistrar {
	return &RoadmapHttpRegistrar{
		handlers: handlers,
	}
}

func (r *RoadmapHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/roadmaps/{roadmap_id}", http.HandlerFunc(r.handlers.GetByID)).Methods("GET")
	s.MUX.Handle("/api/v1/roadmaps/{roadmap_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Update))).Methods("PUT")
	s.MUX.Handle("/api/v1/roadmaps/{roadmap_id}/generate", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Generate))).Methods("PUT")
}
