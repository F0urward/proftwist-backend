package http

import (
	"net/http"
)

func (s *HttpServer) MapHandlers() {
	s.MUX.Handle("/roadmapsinfo", http.HandlerFunc(s.RoadmapInfoH.GetAll)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapInfoH.Create))).Methods("POST")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", http.HandlerFunc(s.RoadmapInfoH.GetByID)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo/roadmap/{roadmap_id}", http.HandlerFunc(s.RoadmapInfoH.GetByRoadmapID)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapInfoH.Update))).Methods("PUT")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapInfoH.Delete))).Methods("DELETE")

	s.MUX.Handle("/roadmaps", http.HandlerFunc(s.RoadmapH.GetAll)).Methods("GET")
	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.GetByID)).Methods("GET")
	s.MUX.Handle("/roadmaps/{roadmap_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapH.Update))).Methods("PUT")

	s.MUX.Handle("/auth/register", http.HandlerFunc(s.AuthH.Register)).Methods("POST")
	s.MUX.Handle("/auth/login", http.HandlerFunc(s.AuthH.Login)).Methods("POST")
	s.MUX.Handle("/auth/logout", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.Logout))).Methods("POST")
	s.MUX.Handle("/auth/me", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.GetMe))).Methods("GET")
	s.MUX.Handle("/auth/{user_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.GetByID))).Methods("GET")
	s.MUX.Handle("/auth", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.Update))).Methods("PUT")
	s.MUX.Handle("/auth/avatar", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.UploadAvatar))).Methods("POST")
	s.MUX.Handle("/auth/vk/link", http.HandlerFunc(s.AuthH.VKOauthLink)).Methods("GET")
	s.MUX.Handle("/auth/vk/callback", http.HandlerFunc(s.AuthH.VKOAuthCallback)).Methods("GET")

	s.MUX.Handle("/api/v1/chats", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.CreateChat))).Methods("POST")
	s.MUX.Handle("/api/v1/chats", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.GetChatsByUser))).Methods("GET")
	s.MUX.Handle("/api/v1/chats/{chat_id}/messages", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.GetChatMessages))).Methods("GET")
	s.MUX.Handle("/api/v1/chats/{chat_id}/members", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.AddMember))).Methods("POST")
	s.MUX.Handle("/api/v1/chats/{chat_id}/members/{user_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.RemoveMember))).Methods("DELETE")

	s.MUX.Handle("/ws", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.WebSocketH.HandleConnection))).Methods("GET")
}
