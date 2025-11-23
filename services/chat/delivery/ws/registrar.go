package ws

import (
	wsServer "github.com/F0urward/proftwist-backend/internal/server/ws"
	"github.com/F0urward/proftwist-backend/internal/server/ws/dto"
	"github.com/F0urward/proftwist-backend/services/chat"
)

type ChatWsRegistrar struct {
	handlers chat.WSHandlers
}

func NewChatWsRegistrar(handlers chat.WSHandlers) wsServer.WsRegistrar {
	return &ChatWsRegistrar{
		handlers: handlers,
	}
}

func (r *ChatWsRegistrar) RegisterHandlers(s *wsServer.WsServer) {
	s.RegisterMessageHandler(dto.WebSocketMessageTypeSendMessage, r.handlers.HandleSendMessage)
	s.RegisterMessageHandler(dto.WebSocketMessageTypeTyping, r.handlers.HandleTyping)
}
