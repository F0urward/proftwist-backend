package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/category"
)

type CategoryHttpRegistrar struct {
	handlers category.Handlers
}

func NewCategoryHttpRegistrar(handlers category.Handlers) httpServer.HttpRegistrar {
	return &CategoryHttpRegistrar{
		handlers: handlers,
	}
}

func (r *CategoryHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/categories", http.HandlerFunc(r.handlers.GetAll)).Methods("GET")
	s.MUX.Handle("/api/v1/categories", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Create))).Methods("POST")
	s.MUX.Handle("/api/v1/categories/{category_id}", http.HandlerFunc(r.handlers.GetByID)).Methods("GET")
	s.MUX.Handle("/api/v1/categories/{category_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Update))).Methods("PUT")
	s.MUX.Handle("/api/v1/categories/{category_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Delete))).Methods("DELETE")
}
