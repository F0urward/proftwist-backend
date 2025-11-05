package chat

import "context"

type Notifier interface {
	NotifyMessageSent(ctx context.Context, userIDs []string, chatID, messageID, senderID, content string) error
	NotifyTyping(ctx context.Context, userIDs []string, chatID, userID string, typing bool) error
	NotifyUserJoined(ctx context.Context, userIDs []string, chatID, userID string) error
	NotifyUserLeft(ctx context.Context, userIDs []string, chatID, userID string) error
}
