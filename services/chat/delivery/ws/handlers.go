package http

import (
	"context"
	"fmt"

	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/F0urward/proftwist-backend/internal/server/websocket/dto"
	"github.com/F0urward/proftwist-backend/services/chat"
	chatdto "github.com/F0urward/proftwist-backend/services/chat/dto"
	"github.com/google/uuid"
)

type ChatWSHandlers struct {
	chatUC   chat.Usecase
	wsServer *websocket.Server
}

func NewChatWSHandlers(chatUC chat.Usecase, wsServer *websocket.Server) chat.WSHandlers {
	integration := &ChatWSHandlers{
		chatUC:   chatUC,
		wsServer: wsServer,
	}

	integration.registerHandlers()
	return integration
}

func (wi *ChatWSHandlers) registerHandlers() {
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeSendMessage, wi.HandleSendMessage)
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeTyping, wi.HandleTyping)
}

func (wi *ChatWSHandlers) HandleSendMessage(client *websocket.Client, msg dto.WebSocketMessage) error {
	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("Handling send message")

	var sendData dto.SendMessageData
	if err := sendData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal send message: %w", err)
	}

	if err := wi.validateChatType(sendData.ChatType); err != nil {
		return err
	}

	chatID, err := uuid.Parse(sendData.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	sendReq := chatdto.SendMessageRequestDTO{
		ChatID:   chatID,
		UserID:   userID,
		Content:  sendData.Content,
		Metadata: sendData.Metadata,
	}

	var message *chatdto.ChatMessageResponseDTO

	switch sendData.ChatType {
	case dto.ChatTypeGroup:
		message, err = wi.chatUC.SendGroupMessage(context.Background(), &sendReq)
	case dto.ChatTypeDirect:
		message, err = wi.chatUC.SendDirectMessage(context.Background(), &sendReq)
	default:
		return fmt.Errorf("unsupported chat type: %s", sendData.ChatType)
	}

	if err != nil {
		return fmt.Errorf("failed to send %s message: %w", sendData.ChatType, err)
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"chat_type":  sendData.ChatType,
		"message_id": message.ID,
		"user_id":    userID,
	}).Info("Message sent successfully")

	return nil
}

func (wi *ChatWSHandlers) HandleTyping(client *websocket.Client, msg dto.WebSocketMessage) error {
	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("Handling typing notification")

	var typingData dto.TypingData
	if err := typingData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal typing data: %w", err)
	}

	if err := wi.validateChatType(typingData.ChatType); err != nil {
		return err
	}

	chatID, err := uuid.Parse(typingData.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	isGroup := typingData.ChatType == dto.ChatTypeGroup
	if err := wi.chatUC.BroadcastTyping(context.Background(), chatID, userID, typingData.Typing, isGroup); err != nil {
		wi.wsServer.Logger().WithError(err).Warn("Failed to broadcast typing notification")
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": client.UserID,
		"typing":  typingData.Typing,
	}).Info("Typing notification handled")

	return nil
}

func (wi *ChatWSHandlers) validateChatType(chatType dto.ChatType) error {
	switch chatType {
	case dto.ChatTypeGroup, dto.ChatTypeDirect:
		return nil
	default:
		return fmt.Errorf("invalid chat type: %s", chatType)
	}
}
