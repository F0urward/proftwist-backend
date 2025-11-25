package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/material"
)

type MaterialHttpRegistrar struct {
	handlers material.Handlers
}

func NewMaterialHttpRegistrar(handlers material.Handlers) httpServer.HttpRegistrar {
	return &MaterialHttpRegistrar{
		handlers: handlers,
	}
}

func (r *MaterialHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/materials/node/{node_id}", http.HandlerFunc(r.handlers.GetMaterialsByNode)).Methods("GET")

	s.MUX.Handle("/api/v1/materials", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.CreateMaterial))).Methods("POST")
	s.MUX.Handle("/api/v1/materials", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetUserMaterials))).Methods("GET")
	s.MUX.Handle("/api/v1/materials/{material_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.DeleteMaterial))).Methods("DELETE")
}
