package http

import (
	"net/http"
)

func (s *HttpServer) MapHandlers() {
	s.MUX.Handle("/roadmapsinfo", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapInfoH.GetAllByUserID))).Methods("GET")
	s.MUX.Handle("/roadmapsinfo", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapInfoH.Create))).Methods("POST")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", http.HandlerFunc(s.RoadmapInfoH.GetByID)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo/roadmap/{roadmap_id}", http.HandlerFunc(s.RoadmapInfoH.GetByRoadmapID)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo/public/category/{category_id}", http.HandlerFunc(s.RoadmapInfoH.GetAllPublicByCategoryID)).Methods("GET")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapInfoH.Update))).Methods("PUT")
	s.MUX.Handle("/roadmapsinfo/{roadmap_info_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapInfoH.Delete))).Methods("DELETE")

	s.MUX.Handle("/roadmaps/{roadmap_id}", http.HandlerFunc(s.RoadmapH.GetByID)).Methods("GET")
	s.MUX.Handle("/roadmaps/{roadmap_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapH.Update))).Methods("PUT")
	s.MUX.Handle("/roadmaps/{roadmap_id}/generate", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.RoadmapH.Generate))).Methods("PUT")

	s.MUX.Handle("/categories", http.HandlerFunc(s.CategoryH.GetAll)).Methods("GET")
	s.MUX.Handle("/categories", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.CategoryH.Create))).Methods("POST")
	s.MUX.Handle("/categories/{category_id}", http.HandlerFunc(s.CategoryH.GetByID)).Methods("GET")
	s.MUX.Handle("/categories/{category_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.CategoryH.Update))).Methods("PUT")
	s.MUX.Handle("/categories/{category_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.CategoryH.Delete))).Methods("DELETE")

	s.MUX.Handle("/auth/register", http.HandlerFunc(s.AuthH.Register)).Methods("POST")
	s.MUX.Handle("/auth/login", http.HandlerFunc(s.AuthH.Login)).Methods("POST")
	s.MUX.Handle("/auth/logout", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.Logout))).Methods("POST")
	s.MUX.Handle("/auth/me", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.GetMe))).Methods("GET")
	s.MUX.Handle("/auth/{user_id}", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.GetByID))).Methods("GET")
	s.MUX.Handle("/auth", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.Update))).Methods("PUT")
	s.MUX.Handle("/auth/avatar", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.AuthH.UploadAvatar))).Methods("POST")
	s.MUX.Handle("/auth/vk/link", http.HandlerFunc(s.AuthH.VKOauthLink)).Methods("GET")
	s.MUX.Handle("/auth/vk/callback", http.HandlerFunc(s.AuthH.VKOAuthCallback)).Methods("GET")

	s.MUX.Handle("/api/v1/group-chats/node/{node_id}", http.HandlerFunc(s.ChatH.GetGroupChatByNode)).Methods("GET")
	s.MUX.Handle("/api/v1/group-chats", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.GetGroupChatsByUser))).Methods("GET")
	s.MUX.Handle("/api/v1/group-chats/{chat_id}/members", http.HandlerFunc(s.ChatH.GetGroupChatMembers)).Methods("GET")
	s.MUX.Handle("/api/v1/group-chats/{chat_id}/messages", http.HandlerFunc(s.ChatH.GetGroupChatMessages)).Methods("GET")
	s.MUX.Handle("/api/v1/group-chats/{chat_id}/join", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.JoinGroupChat))).Methods("POST")
	s.MUX.Handle("/api/v1/group-chats/{chat_id}/leave", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.LeaveGroupChat))).Methods("POST")

	s.MUX.Handle("/api/v1/direct-chats", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.GetDirectChatsByUser))).Methods("GET")
	s.MUX.Handle("/api/v1/direct-chats/{chat_id}/member", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.GetDirectChatMembers))).Methods("GET")
	s.MUX.Handle("/api/v1/direct-chats/{chat_id}/messages", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.ChatH.GetDirectChatMessages))).Methods("GET")

	s.MUX.Handle("/ws", s.AuthMiddleware.AuthMiddleware(http.HandlerFunc(s.WebSocketH.HandleConnection))).Methods("GET")
}
