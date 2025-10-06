package http

import (
	"net/http"
)

func (s *HttpServer) MapHandlers() {
	s.MUX.Handle("/roadmapsinfo", http.HandlerFunc(s.RoadmapInfoH.GetAll)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo", http.HandlerFunc(s.RoadmapInfoH.Create)).Methods("POST")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", http.HandlerFunc(s.RoadmapInfoH.GetByID)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", http.HandlerFunc(s.RoadmapInfoH.Update)).Methods("PUT")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", http.HandlerFunc(s.RoadmapInfoH.Delete)).Methods("DELETE")

	s.MUX.Handle("/roadmaps", http.HandlerFunc(s.RoadmapH.GetAll)).Methods("GET")
	s.MUX.Handle("/roadmaps", http.HandlerFunc(s.RoadmapH.Create)).Methods("POST")
	s.MUX.Handle("/roadmaps/search", http.HandlerFunc(s.RoadmapH.SearchByTitle)).Methods("GET")
	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.GetByID)).Methods("GET")
	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.Update)).Methods("PUT")
	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.Delete)).Methods("DELETE")
	s.MUX.Handle("/roadmaps/author/{author_id}", http.HandlerFunc(s.RoadmapH.GetByAuthorID)).Methods("GET")
	s.MUX.Handle("/roadmaps/{roadmap_id}/privacy", http.HandlerFunc(s.RoadmapH.UpdatePrivacy)).Methods("PATCH")

	s.MUX.Handle("/auth/register", http.HandlerFunc(s.AuthH.Register)).Methods("POST")
	s.MUX.Handle("/auth/login", http.HandlerFunc(s.AuthH.Login)).Methods("POST")
	s.MUX.Handle("/auth/logout", http.HandlerFunc(s.AuthH.Logout)).Methods("POST")
}
