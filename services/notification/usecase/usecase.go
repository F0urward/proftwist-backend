package usecase

import (
	"context"
	"time"

	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/server/ws"
	wsDTO "github.com/F0urward/proftwist-backend/internal/server/ws/dto"
	"github.com/F0urward/proftwist-backend/services/notification"
	"github.com/F0urward/proftwist-backend/services/notification/dto"
)

type NotificationUsecase struct {
	wsServer *ws.WsServer
}

func NewNotificationUsecase(wsServer *ws.WsServer) notification.Usecase {
	return &NotificationUsecase{wsServer: wsServer}
}

func (uc *NotificationUsecase) HandleMessageSent(ctx context.Context, event dto.MessageSentEvent) error {
	const op = "NotificationUsecase.HandleMessageSent"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", event.ChatID).WithField("message_id", event.MessageID)

	logger.Infof("sending MessageSentEvent to users: %v", event.UserIDs)

	messageData := wsDTO.MessageSentData{
		MessageID: event.MessageID,
		ChatID:    event.ChatID,
		UserID:    event.SenderID,
		Username:  event.Username,
		AvatarURL: event.AvatarURL,
		Content:   event.Content,
		SentAt:    event.SentAt,
	}

	data, err := messageData.MarshalJSON()
	if err != nil {
		logger.WithError(err).Error("failed to marshal MessageSentData")
		return err
	}

	wsMessage := wsDTO.WebSocketMessage{
		Type:      wsDTO.WebSocketMessageTypeMessageSent,
		Data:      data,
		Timestamp: time.Now(),
	}

	if err := uc.wsServer.SendToUsers(event.UserIDs, wsMessage); err != nil {
		logger.WithError(err).Error("failed to send MessageSentEvent to users")
		return err
	}

	logger.Info("successfully sent MessageSentEvent")
	return nil
}

func (uc *NotificationUsecase) HandleTyping(ctx context.Context, event dto.TypingEvent) error {
	const op = "NotificationUsecase.HandleTyping"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", event.ChatID).WithField("user_id", event.UserID)

	logger.Infof("sending TypingEvent: %v", event.Typing)

	typingData := wsDTO.TypingNotificationData{
		ChatID:   event.ChatID,
		UserID:   event.UserID,
		Username: event.Username,
		Typing:   event.Typing,
	}

	data, err := typingData.MarshalJSON()
	if err != nil {
		logger.WithError(err).Error("failed to marshal TypingNotificationData")
		return err
	}

	wsMessage := wsDTO.WebSocketMessage{
		Type:      wsDTO.WebSocketMessageTypeTypingNotification,
		Data:      data,
		Timestamp: time.Now(),
	}

	if err := uc.wsServer.SendToUsers(event.UserIDs, wsMessage); err != nil {
		logger.WithError(err).Error("failed to send TypingEvent to users")
		return err
	}

	logger.Info("successfully sent TypingEvent")
	return nil
}

func (uc *NotificationUsecase) HandleUserJoined(ctx context.Context, event dto.UserJoinedEvent) error {
	const op = "NotificationUsecase.HandleUserJoined"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", event.ChatID).WithField("user_id", event.UserID)

	logger.Info("sending UserJoinedEvent")

	joinData := wsDTO.UserJoinedNotificationData{
		ChatID:   event.ChatID,
		UserID:   event.UserID,
		Username: event.Username,
		JoinedAt: event.JoinedAt,
	}

	data, err := joinData.MarshalJSON()
	if err != nil {
		logger.WithError(err).Error("failed to marshal UserJoinedNotificationData")
		return err
	}

	wsMessage := wsDTO.WebSocketMessage{
		Type:      wsDTO.WebSocketMessageTypeUserJoined,
		Data:      data,
		Timestamp: time.Now(),
	}

	if err := uc.wsServer.SendToUsers(event.UserIDs, wsMessage); err != nil {
		logger.WithError(err).Error("failed to send UserJoinedEvent to users")
		return err
	}

	logger.Info("successfully sent UserJoinedEvent")
	return nil
}

func (uc *NotificationUsecase) HandleUserLeft(ctx context.Context, event dto.UserLeftEvent) error {
	const op = "NotificationUsecase.HandleUserLeft"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("chat_id", event.ChatID).WithField("user_id", event.UserID)

	logger.Info("sending UserLeftEvent")

	leftData := wsDTO.UserLeftNotificationData{
		ChatID:   event.ChatID,
		UserID:   event.UserID,
		Username: event.Username,
		LeftAt:   event.LeftAt,
	}

	data, err := leftData.MarshalJSON()
	if err != nil {
		logger.WithError(err).Error("failed to marshal UserLeftNotificationData")
		return err
	}

	wsMessage := wsDTO.WebSocketMessage{
		Type:      wsDTO.WebSocketMessageTypeUserLeft,
		Data:      data,
		Timestamp: time.Now(),
	}

	if err := uc.wsServer.SendToUsers(event.UserIDs, wsMessage); err != nil {
		logger.WithError(err).Error("failed to send UserLeftEvent to users")
		return err
	}

	logger.Info("successfully sent UserLeftEvent")
	return nil
}
