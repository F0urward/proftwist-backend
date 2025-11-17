package chat

import "context"

type Broadcaster interface {
	SendToUser(userID string, msg interface{}) error
	SendToUsers(userIDs []string, msg interface{}) error
	Broadcast(msg interface{}) error
}

type Notifier interface {
	NotifyMessageSent(ctx context.Context, userIDs []string, chatID, messageID, senderID, content, username, avatarURL string) error
	NotifyTyping(ctx context.Context, userIDs []string, chatID, userID, username string, typing bool) error
	NotifyUserJoined(ctx context.Context, userIDs []string, chatID, userID, username string) error
	NotifyUserLeft(ctx context.Context, userIDs []string, chatID, userID, username string) error
}
