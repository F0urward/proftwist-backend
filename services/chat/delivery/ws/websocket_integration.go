package http

import (
	"context"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/F0urward/proftwist-backend/internal/server/websocket/dto"
	"github.com/F0urward/proftwist-backend/services/chat"
	chatdto "github.com/F0urward/proftwist-backend/services/chat/dto"
	"github.com/google/uuid"
)

type WebSocketIntegration struct {
	chatUC   chat.Usecase
	wsServer *websocket.Server
}

func NewWebSocketIntegration(chatUC chat.Usecase, wsServer *websocket.Server) *WebSocketIntegration {
	integration := &WebSocketIntegration{
		chatUC:   chatUC,
		wsServer: wsServer,
	}

	integration.registerHandlers()
	return integration
}

func (wi *WebSocketIntegration) registerHandlers() {
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeChat, wi.handleChatMessage)
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeRead, wi.handleReadMessage)
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeJoin, wi.handleJoinChat)
}

func (wi *WebSocketIntegration) handleChatMessage(client *websocket.Client, msg dto.WebSocketMessage) error {
	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("Handling chat message")

	var chatData dto.ChatMessageData
	if err := chatData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal chat message: %w", err)
	}

	chatID, err := uuid.Parse(chatData.ChatID)
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
		Content:  chatData.Content,
		Metadata: chatData.Metadata,
	}

	message, err := wi.chatUC.SendMessage(context.Background(), sendReq)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if err := wi.broadcastMessageToChat(chatID, message, client.UserID); err != nil {
		return fmt.Errorf("failed to broadcast message: %w", err)
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"message_id": message.ID,
		"user_id":    userID,
	}).Info("Chat message processed and broadcasted")

	return nil
}

func (wi *WebSocketIntegration) handleReadMessage(client *websocket.Client, msg dto.WebSocketMessage) error {
	var readData dto.ReadReceiptData
	if err := readData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal read message: %w", err)
	}

	chatID, err := uuid.Parse(readData.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"user_id":    userID,
		"message_id": readData.MessageID,
	}).Info("Read receipt processed")

	return nil
}

func (wi *WebSocketIntegration) handleJoinChat(client *websocket.Client, msg dto.WebSocketMessage) error {
	var joinData dto.JoinChatData
	if err := joinData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal join message: %w", err)
	}

	chatID, err := uuid.Parse(joinData.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	chats, err := wi.chatUC.GetUserChats(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user chats: %w", err)
	}

	isMember := false
	for _, chat := range chats {
		if chat.ID == chatID {
			isMember = true
			break
		}
	}

	if !isMember {
		return fmt.Errorf("user is not a member of this chat")
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("User joined chat via WebSocket")

	return nil
}

func (wi *WebSocketIntegration) broadcastMessageToChat(chatID uuid.UUID, message *chatdto.ChatMessageResponseDTO, senderUserID string) error {
	members, err := wi.chatUC.GetChatMembers(context.Background(), chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeChat,
		Data:      wi.marshalMessageData(message),
		Timestamp: message.CreatedAt,
		UserID:    senderUserID,
		ChatID:    chatID.String(),
	}

	userIDs := wi.getMemberUserIDs(members, senderUserID)

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":        chatID,
		"sender":         senderUserID,
		"recipients":     len(userIDs),
		"message_id":     message.ID,
		"content_length": len(message.Content),
	}).Info("Broadcasting message to chat members")

	return wi.wsServer.SendToUsers(userIDs, wsMessage)
}

func (wi *WebSocketIntegration) getMemberUserIDs(members []chatdto.ChatMemberResponseDTO, excludeUserID string) []string {
	var userIDs []string
	for _, member := range members {
		userIDStr := member.UserID.String()
		if userIDStr != excludeUserID {
			userIDs = append(userIDs, userIDStr)
		}
	}
	return userIDs
}

func (wi *WebSocketIntegration) marshalMessageData(message *chatdto.ChatMessageResponseDTO) []byte {
	messageData := dto.ChatMessageData{
		MessageID: message.ID.String(),
		ChatID:    message.ChatID.String(),
		Content:   message.Content,
		Metadata:  message.Metadata,
	}

	data, _ := messageData.MarshalJSON()
	return data
}

func (wi *WebSocketIntegration) BroadcastNewChat(chat *chatdto.ChatResponseDTO, memberIDs []uuid.UUID) error {
	chatData := dto.ChatMessageData{
		ChatID:  chat.ID.String(),
		Content: fmt.Sprintf("New chat created: %s", chat.Title),
		Metadata: map[string]interface{}{
			"action":    "chat_created",
			"chat_type": chat.Type,
			"title":     chat.Title,
		},
	}

	data, _ := chatData.MarshalJSON()
	message := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeChat,
		Data:      data,
		Timestamp: time.Now(),
		ChatID:    chat.ID.String(),
	}

	var userIDs []string
	for _, memberID := range memberIDs {
		userIDs = append(userIDs, memberID.String())
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chat.ID,
		"members":    len(userIDs),
		"chat_title": chat.Title,
	}).Info("Broadcasting new chat notification")

	return wi.wsServer.SendToUsers(userIDs, message)
}
