package dto

import (
	"time"
)

type EventType string

const (
	MessagePublishedType EventType = "chat.message_published"
	UserTypingType       EventType = "chat.user_typing"
	UserJoinedType       EventType = "chat.user_joined"
	UserLeftType         EventType = "chat.user_left"
)

type BaseEvent struct {
	Type EventType `json:"type"`
}

type MessageSentEvent struct {
	Type      EventType `json:"type"`
	UserIDs   []string  `json:"user_ids"`
	ChatID    string    `json:"chat_id"`
	MessageID string    `json:"message_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	SentAt    time.Time `json:"timestamp"`
}

type TypingEvent struct {
	Type     EventType `json:"type"`
	UserIDs  []string  `json:"user_ids"`
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Typing   bool      `json:"typing"`
}

type UserJoinedEvent struct {
	Type     EventType `json:"type"`
	UserIDs  []string  `json:"user_ids"`
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	JoinedAt time.Time `json:"timestamp"`
}

type UserLeftEvent struct {
	Type     EventType `json:"type"`
	UserIDs  []string  `json:"user_ids"`
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	LeftAt   time.Time `json:"timestamp"`
}
