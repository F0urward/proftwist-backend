package chat

import (
	"context"
	"time"

	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/F0urward/proftwist-backend/internal/server/websocket/dto"
	"github.com/F0urward/proftwist-backend/services/chat"
)

type WSNotifier struct {
	wsServer *websocket.Server
}

func NewWSNotifier(wsServer *websocket.Server) chat.Notifier {
	return &WSNotifier{wsServer: wsServer}
}

func (w *WSNotifier) NotifyMessageSent(ctx context.Context, userIDs []string, chatID, messageID, senderID, content, username, avatarURL string) error {
	messageData := dto.MessageSentData{
		MessageID: messageID,
		ChatID:    chatID,
		UserID:    senderID,
		Username:  username,
		AvatarURL: avatarURL,
		Content:   content,
		SentAt:    time.Now(),
	}

	data, err := messageData.MarshalJSON()
	if err != nil {
		return err
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeMessageSent,
		Data:      data,
		Timestamp: time.Now(),
	}

	return w.wsServer.SendToUsers(userIDs, wsMessage)
}

func (w *WSNotifier) NotifyTyping(ctx context.Context, userIDs []string, chatID, userID, username string, typing bool) error {
	typingData := dto.TypingNotificationData{
		ChatID:   chatID,
		UserID:   userID,
		Username: username,
		Typing:   typing,
	}

	data, err := typingData.MarshalJSON()
	if err != nil {
		return err
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeTypingNotification,
		Data:      data,
		Timestamp: time.Now(),
	}

	return w.wsServer.SendToUsers(userIDs, wsMessage)
}

func (w *WSNotifier) NotifyUserJoined(ctx context.Context, userIDs []string, chatID, userID, username string) error {
	joinData := dto.UserJoinedNotificationData{
		ChatID:   chatID,
		UserID:   userID,
		Username: username,
		JoinedAt: time.Now(),
	}

	data, err := joinData.MarshalJSON()
	if err != nil {
		return err
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeUserJoined,
		Data:      data,
		Timestamp: time.Now(),
	}

	return w.wsServer.SendToUsers(userIDs, wsMessage)
}

func (w *WSNotifier) NotifyUserLeft(ctx context.Context, userIDs []string, chatID, userID, username string) error {
	leftData := dto.UserLeftNotificationData{
		ChatID:   chatID,
		UserID:   userID,
		Username: username,
		LeftAt:   time.Now(),
	}

	data, err := leftData.MarshalJSON()
	if err != nil {
		return err
	}

	wsMessage := dto.WebSocketMessage{
		Type:      dto.WebSocketMessageTypeUserLeft,
		Data:      data,
		Timestamp: time.Now(),
	}

	return w.wsServer.SendToUsers(userIDs, wsMessage)
}
