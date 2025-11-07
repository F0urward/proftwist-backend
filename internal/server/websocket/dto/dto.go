package dto

import (
	"encoding/json"
	"time"
)

type ChatType string

const (
	ChatTypeGroup  ChatType = "group"
	ChatTypeDirect ChatType = "direct"
)

type WebSocketMessageType string

const (
	WebSocketMessageTypeSendMessage WebSocketMessageType = "send_message"
	WebSocketMessageTypeTyping      WebSocketMessageType = "typing"

	WebSocketMessageTypeMessageSent        WebSocketMessageType = "message_sent"
	WebSocketMessageTypeTypingNotification WebSocketMessageType = "typing_notification"
	WebSocketMessageTypeUserJoined         WebSocketMessageType = "user_joined"
	WebSocketMessageTypeUserLeft           WebSocketMessageType = "user_left"
)

type WebSocketMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Data      json.RawMessage      `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}

type SendMessageData struct {
	MessageID string                 `json:"message_id,omitempty"`
	ChatID    string                 `json:"chat_id"`
	ChatType  ChatType               `json:"chat_type"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type TypingData struct {
	ChatID   string   `json:"chat_id"`
	ChatType ChatType `json:"chat_type"`
	Typing   bool     `json:"typing"`
}

type MessageSentData struct {
	MessageID string                 `json:"message_id"`
	ChatID    string                 `json:"chat_id"`
	UserID    string                 `json:"user_id"`
	Username  string                 `json:"username"`
	AvatarURL string                 `json:"avatar_url"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	SentAt    time.Time              `json:"sent_at"`
}

type MessageDeliveredData struct {
	MessageID   string    `json:"message_id"`
	ChatID      string    `json:"chat_id"`
	UserID      string    `json:"user_id"`
	DeliveredAt time.Time `json:"delivered_at"`
}

type TypingNotificationData struct {
	ChatID   string `json:"chat_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Typing   bool   `json:"typing"`
}

type UserJoinedNotificationData struct {
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	JoinedAt time.Time `json:"joined_at"`
}

type UserLeftNotificationData struct {
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	LeftAt   time.Time `json:"left_at"`
}
