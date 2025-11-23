package chat

import (
	"context"
	"time"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/broker"
	"github.com/F0urward/proftwist-backend/services/chat"
	notificationDTO "github.com/F0urward/proftwist-backend/services/notification/dto"
)

type BrokerNotifier struct {
	producer broker.Producer
}

func NewBrokerNotifier(producer broker.Producer) chat.Notifier {
	return &BrokerNotifier{producer: producer}
}

func (k *BrokerNotifier) NotifyMessageSent(ctx context.Context, userIDs []string, chatID, messageID, senderID, content, username, avatarURL string) error {
	event := notificationDTO.MessageSentEvent{
		Type:      notificationDTO.MessagePublishedType,
		UserIDs:   userIDs,
		ChatID:    chatID,
		MessageID: messageID,
		SenderID:  senderID,
		Content:   content,
		Username:  username,
		AvatarURL: avatarURL,
		SentAt:    time.Now(),
	}

	data, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	return k.producer.Publish(ctx, chatID, data)
}

func (k *BrokerNotifier) NotifyTyping(ctx context.Context, userIDs []string, chatID, userID, username string, typing bool) error {
	event := notificationDTO.TypingEvent{
		Type:     notificationDTO.UserTypingType,
		UserIDs:  userIDs,
		ChatID:   chatID,
		UserID:   userID,
		Username: username,
		Typing:   typing,
	}

	data, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	return k.producer.Publish(ctx, chatID, data)
}

func (k *BrokerNotifier) NotifyUserJoined(ctx context.Context, userIDs []string, chatID, userID, username string) error {
	event := notificationDTO.UserJoinedEvent{
		Type:     notificationDTO.UserJoinedType,
		UserIDs:  userIDs,
		ChatID:   chatID,
		UserID:   userID,
		Username: username,
		JoinedAt: time.Now(),
	}

	data, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	return k.producer.Publish(ctx, chatID, data)
}

func (k *BrokerNotifier) NotifyUserLeft(ctx context.Context, userIDs []string, chatID, userID, username string) error {
	event := notificationDTO.UserLeftEvent{
		Type:     notificationDTO.UserLeftType,
		UserIDs:  userIDs,
		ChatID:   chatID,
		UserID:   userID,
		Username: username,
		LeftAt:   time.Now(),
	}

	data, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	return k.producer.Publish(ctx, chatID, data)
}
