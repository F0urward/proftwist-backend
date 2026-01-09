package ws

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	websocket "github.com/F0urward/proftwist-backend/internal/server/ws"
	"github.com/F0urward/proftwist-backend/internal/server/ws/dto"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/chat"
	chatDTO "github.com/F0urward/proftwist-backend/services/chat/dto"
)

type ChatWsHandlers struct {
	chatUC chat.Usecase
}

func NewChatWsHandlers(chatUC chat.Usecase) chat.WSHandlers {
	return &ChatWsHandlers{
		chatUC: chatUC,
	}
}

func (wsh *ChatWsHandlers) HandleSendMessage(client *websocket.WsClient, msg dto.WebSocketMessage) error {
	const op = "ChatWsHandlers.HandleSendMessage"
	ctx := context.Background()
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	logger.WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("handling send message")

	var sendData dto.SendMessageData
	if err := sendData.UnmarshalJSON(msg.Data); err != nil {
		logger.WithError(err).Error("failed to unmarshal send message data")
		return fmt.Errorf("failed to unmarshal send message: %w", err)
	}

	if err := wsh.validateChatType(sendData.ChatType); err != nil {
		logger.WithError(err).Error("invalid chat type")
		return err
	}

	chatID, err := uuid.Parse(sendData.ChatID)
	if err != nil {
		logger.WithError(err).Error("invalid chat ID")
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		logger.WithError(err).Error("invalid user ID")
		return fmt.Errorf("invalid user ID: %w", err)
	}

	sendReq := chatDTO.SendMessageRequestDTO{
		ChatID:   chatID,
		UserID:   userID,
		Content:  sendData.Content,
		Metadata: sendData.Metadata,
	}

	var message *chatDTO.ChatMessageResponseDTO

	switch sendData.ChatType {
	case dto.ChatTypeGroup:
		message, err = wsh.chatUC.SendGroupChatMessage(ctx, &sendReq)
	case dto.ChatTypeDirect:
		message, err = wsh.chatUC.SendDirectMessage(ctx, &sendReq)
	default:
		err = fmt.Errorf("unsupported chat type: %s", sendData.ChatType)
	}

	if err != nil {
		logger.WithError(err).Error("failed to send message")
		return fmt.Errorf("failed to send %s message: %w", sendData.ChatType, err)
	}

	logger.WithFields(map[string]interface{}{
		"chat_id":    chatID,
		"chat_type":  sendData.ChatType,
		"message_id": message.ID,
		"user_id":    userID,
	}).Info("message sent successfully")

	return nil
}

func (wsh *ChatWsHandlers) HandleTyping(client *websocket.WsClient, msg dto.WebSocketMessage) error {
	const op = "ChatWsHandlers.HandleTyping"
	ctx := context.Background()
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	logger.WithFields(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"type":      msg.Type,
	}).Info("handling typing notification")

	var typingData dto.TypingData
	if err := typingData.UnmarshalJSON(msg.Data); err != nil {
		logger.WithError(err).Error("failed to unmarshal typing data")
		return fmt.Errorf("failed to unmarshal typing data: %w", err)
	}

	if err := wsh.validateChatType(typingData.ChatType); err != nil {
		logger.WithError(err).Error("invalid chat type")
		return err
	}

	chatID, err := uuid.Parse(typingData.ChatID)
	if err != nil {
		logger.WithError(err).Error("invalid chat ID")
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(client.UserID)
	if err != nil {
		logger.WithError(err).Error("invalid user ID")
		return fmt.Errorf("invalid user ID: %w", err)
	}

	isGroup := typingData.ChatType == dto.ChatTypeGroup
	if err := wsh.chatUC.BroadcastTyping(ctx, chatID, userID, typingData.Typing, isGroup); err != nil {
		logger.WithError(err).Warn("failed to broadcast typing notification")
	}

	logger.WithFields(map[string]interface{}{
		"chat_id": chatID,
		"user_id": client.UserID,
		"typing":  typingData.Typing,
	}).Info("typing notification handled successfully")

	return nil
}

func (wsh *ChatWsHandlers) validateChatType(chatType dto.ChatType) error {
	switch chatType {
	case dto.ChatTypeGroup, dto.ChatTypeDirect:
		return nil
	default:
		return fmt.Errorf("invalid chat type: %s", chatType)
	}
}
