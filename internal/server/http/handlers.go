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
	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.GetByID)).Methods("GET")
	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.Update)).Methods("PUT")
	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.Delete)).Methods("DELETE")
}
