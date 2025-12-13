package chat

import "context"

type NotificationPublisher interface {
	NotifyMessageSent(ctx context.Context, userIDs []string, chatID, messageID, senderID, content, username, avatarURL string) error
	NotifyTyping(ctx context.Context, userIDs []string, chatID, userID, username string, typing bool) error
	NotifyUserJoined(ctx context.Context, userIDs []string, chatID, userID, username string) error
	NotifyUserLeft(ctx context.Context, userIDs []string, chatID, userID, username string) error
}

type BotPublisher interface {
	PublishMessageForBot(ctx context.Context, chatID, chatTitle, content string) error
}
