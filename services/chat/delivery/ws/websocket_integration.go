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
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeSendMessage, wi.handleSendMessage)
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeTyping, wi.handleTyping)
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeJoin, wi.handleJoinChat)
	wi.wsServer.RegisterMessageHandler(dto.WebSocketMessageTypeLeave, wi.handleLeaveChat)
}

func (wi *WebSocketIntegration) handleSendMessage(client *websocket.Client, msg dto.WebSocketMessage) error {
	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("Handling send message")

	var sendData dto.SendMessageData
	if err := sendData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal send message: %w", err)
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

	message, err := wi.chatUC.SendMessage(context.Background(), sendReq)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Сообщение broadcast всем участникам чата (ВКЛЮЧАЯ отправителя)
	if err := wi.broadcastMessageToChat(chatID, message, client.UserID); err != nil {
		return fmt.Errorf("failed to broadcast message: %w", err)
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"message_id": message.ID,
		"user_id":    userID,
	}).Info("Message sent and broadcasted")

	return nil
}

func (wi *WebSocketIntegration) handleTyping(client *websocket.Client, msg dto.WebSocketMessage) error {
	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("Handling typing notification")

	var typingData dto.TypingData
	if err := typingData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal typing data: %w", err)
	}

	chatID, err := uuid.Parse(typingData.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	// Уведомление о печати broadcast всем участникам чата (КРОМЕ отправителя)
	if err := wi.broadcastTypingNotification(chatID, client.UserID, typingData.Typing); err != nil {
		return fmt.Errorf("failed to broadcast typing notification: %w", err)
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": client.UserID,
		"typing":  typingData.Typing,
	}).Info("Typing notification broadcasted")

	return nil
}

func (wi *WebSocketIntegration) handleJoinChat(client *websocket.Client, msg dto.WebSocketMessage) error {
	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("Handling join chat")

	var joinData dto.JoinChatData
	if err := joinData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal join chat data: %w", err)
	}

	chatID, err := uuid.Parse(joinData.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := wi.chatUC.JoinGroupChat(context.Background(), chatID, userID); err != nil {
		return fmt.Errorf("failed to join group chat: %w", err)
	}

	// Уведомление о присоединении broadcast всем участникам чата (ВКЛЮЧАЯ присоединившегося)
	if err := wi.broadcastUserJoined(chatID, client.UserID, ""); err != nil {
		wi.wsServer.Logger().WithError(err).Warn("Failed to broadcast user joined notification")
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("User joined chat successfully")

	return nil
}

func (wi *WebSocketIntegration) handleLeaveChat(client *websocket.Client, msg dto.WebSocketMessage) error {
	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("Handling leave chat")

	var leaveData dto.LeaveChatData
	if err := leaveData.UnmarshalJSON(msg.Data); err != nil {
		return fmt.Errorf("failed to unmarshal leave chat data: %w", err)
	}

	chatID, err := uuid.Parse(leaveData.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := wi.chatUC.LeaveGroupChat(context.Background(), chatID, userID); err != nil {
		return fmt.Errorf("failed to leave group chat: %w", err)
	}

	// Уведомление о выходе broadcast всем участникам чата (КРОМЕ вышедшего)
	if err := wi.broadcastUserLeft(chatID, client.UserID, ""); err != nil {
		wi.wsServer.Logger().WithError(err).Warn("Failed to broadcast user left notification")
	}

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("User left chat successfully")

	return nil
}

func (wi *WebSocketIntegration) broadcastMessageToChat(chatID uuid.UUID, message *chatdto.ChatMessageResponseDTO, senderUserID string) error {
	members, err := wi.chatUC.GetChatMembers(context.Background(), chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeMessageSent,
		Data:      wi.marshalMessageData(chatID, message, senderUserID),
		Timestamp: message.CreatedAt,
	}

	// Сообщения отправляются ВСЕМ участникам чата (включая отправителя)
	userIDs := wi.getMemberUserIDs(members)

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":        chatID,
		"sender":         senderUserID,
		"recipients":     len(userIDs),
		"message_id":     message.ID,
		"content_length": len(message.Content),
	}).Info("Broadcasting message to all chat members (including sender)")

	return wi.wsServer.SendToUsers(userIDs, wsMessage)
}

func (wi *WebSocketIntegration) broadcastTypingNotification(chatID uuid.UUID, senderUserID string, typing bool) error {
	members, err := wi.chatUC.GetChatMembers(context.Background(), chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	typingData := dto.TypingNotificationData{
		ChatID: chatID.String(),
		UserID: senderUserID,
		Typing: typing,
	}

	data, err := typingData.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal typing notification data: %w", err)
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeTypingNotification,
		Data:      data,
		Timestamp: time.Now(),
	}

	// Уведомления о печати отправляются всем КРОМЕ отправителя
	userIDs := wi.getMemberUserIDs(members, senderUserID)

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"sender":     senderUserID,
		"recipients": len(userIDs),
		"typing":     typing,
	}).Info("Broadcasting typing notification to all except sender")

	return wi.wsServer.SendToUsers(userIDs, wsMessage)
}

func (wi *WebSocketIntegration) broadcastUserJoined(chatID uuid.UUID, userIDStr string, username string) error {
	members, err := wi.chatUC.GetChatMembers(context.Background(), chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	joinData := dto.UserJoinedNotificationData{
		ChatID:   chatID.String(),
		UserID:   userIDStr,
		Username: username,
		JoinedAt: time.Now(),
	}

	data, err := joinData.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal user joined notification: %w", err)
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeUserJoined,
		Data:      data,
		Timestamp: time.Now(),
	}

	// Уведомления о присоединении отправляются ВСЕМ участникам (включая присоединившегося)
	userIDs := wi.getMemberUserIDs(members)

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"user_id":    userIDStr,
		"recipients": len(userIDs),
	}).Info("Broadcasting user joined notification to all members (including joined user)")

	return wi.wsServer.SendToUsers(userIDs, wsMessage)
}

func (wi *WebSocketIntegration) broadcastUserLeft(chatID uuid.UUID, userIDStr string, username string) error {
	members, err := wi.chatUC.GetChatMembers(context.Background(), chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	leftData := dto.UserLeftNotificationData{
		ChatID:   chatID.String(),
		UserID:   userIDStr,
		Username: username,
		LeftAt:   time.Now(),
	}

	data, err := leftData.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal user left notification: %w", err)
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeUserLeft,
		Data:      data,
		Timestamp: time.Now(),
	}

	// Уведомления о выходе отправляются всем КРОМЕ вышедшего
	userIDs := wi.getMemberUserIDs(members, userIDStr)

	wi.wsServer.Logger().WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"user_id":    userIDStr,
		"recipients": len(userIDs),
	}).Info("Broadcasting user left notification to all except left user")

	return wi.wsServer.SendToUsers(userIDs, wsMessage)
}

func (wi *WebSocketIntegration) getMemberUserIDs(members []chatdto.ChatMemberResponseDTO, excludeUserID ...string) []string {
	var userIDs []string
	for _, member := range members {
		userIDStr := member.UserID.String()

		// Если передан excludeUserID, пропускаем этого пользователя
		if len(excludeUserID) > 0 && userIDStr == excludeUserID[0] {
			continue
		}
		userIDs = append(userIDs, userIDStr)
	}
	return userIDs
}

func (wi *WebSocketIntegration) marshalMessageData(chatID uuid.UUID, message *chatdto.ChatMessageResponseDTO, senderUserID string) []byte {
	messageData := dto.MessageSentData{
		MessageID: message.ID.String(),
		ChatID:    chatID.String(),
		UserID:    senderUserID,
		Content:   message.Content,
		Metadata:  message.Metadata,
		SentAt:    message.CreatedAt,
	}

	data, _ := messageData.MarshalJSON()
	return data
}
