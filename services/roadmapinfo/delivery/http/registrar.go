package http

import (
	"net/http"

	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
)

type RoadmapInfoHttpRegistrar struct {
	handlers roadmapinfo.Handlers
}

func NewRoadmapInfoHttpRegistrar(handlers roadmapinfo.Handlers) httpServer.HttpRegistrar {
	return &RoadmapInfoHttpRegistrar{
		handlers: handlers,
	}
}

func (r *RoadmapInfoHttpRegistrar) RegisterRoutes(s *httpServer.HttpServer) {
	s.MUX.Handle("/api/v1/roadmapsinfo/public", http.HandlerFunc(r.handlers.GetAllPublic)).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo/public/search", http.HandlerFunc(r.handlers.SearchPublic)).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo/public/category/{category_id}", http.HandlerFunc(r.handlers.GetAllPublicByCategoryID)).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo/public/subscribed", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetSubscribedRoadmaps))).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo/public/{roadmap_info_id}/subscribe", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Subscribe))).Methods("POST")
	s.MUX.Handle("/api/v1/roadmapsinfo/public/{roadmap_info_id}/unsubscribe", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Unsubscribe))).Methods("DELETE")
	s.MUX.Handle("/api/v1/roadmapsinfo/public/{roadmap_info_id}/subscription", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.CheckSubscription))).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.GetAllByUserID))).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo/{roadmap_info_id}", http.HandlerFunc(r.handlers.GetByID)).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo/roadmap/{roadmap_id}", http.HandlerFunc(r.handlers.GetByRoadmapID)).Methods("GET")
	s.MUX.Handle("/api/v1/roadmapsinfo/private", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.CreatePrivate))).Methods("POST")
	s.MUX.Handle("/api/v1/roadmapsinfo/private/{roadmap_info_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.UpdatePrivate))).Methods("PUT")
	s.MUX.Handle("/api/v1/roadmapsinfo/private/{roadmap_info_id}/publish", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Publish))).Methods("POST")
	s.MUX.Handle("/api/v1/roadmapsinfo/public/{roadmap_info_id}/fork", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Fork))).Methods("POST")
	s.MUX.Handle("/api/v1/roadmapsinfo/{roadmap_info_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(r.handlers.Delete))).Methods("DELETE")
}
