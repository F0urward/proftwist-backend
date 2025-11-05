package dto

import (
	"encoding/json"
	"time"
)

type WebSocketMessageType string

const (
	WebSocketMessageTypeSendMessage WebSocketMessageType = "send_message"
	WebSocketMessageTypeTyping      WebSocketMessageType = "typing"
	WebSocketMessageTypeJoin        WebSocketMessageType = "join"
	WebSocketMessageTypeLeave       WebSocketMessageType = "leave"

	WebSocketMessageTypeMessageSent        WebSocketMessageType = "message_sent"
	WebSocketMessageTypeTypingNotification WebSocketMessageType = "typing_notification"
	WebSocketMessageTypeUserJoined         WebSocketMessageType = "user_joined"
	WebSocketMessageTypeUserLeft           WebSocketMessageType = "user_left"
	WebSocketMessageTypeError              WebSocketMessageType = "error"
)

type WebSocketMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Data      json.RawMessage      `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}

type SendMessageData struct {
	MessageID string                 `json:"message_id,omitempty"`
	ChatID    string                 `json:"chat_id"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type TypingData struct {
	ChatID string `json:"chat_id"`
	Typing bool   `json:"typing"`
}

type ReadReceiptData struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
}

type JoinChatData struct {
	ChatID string `json:"chat_id"`
}

type LeaveChatData struct {
	ChatID string `json:"chat_id"`
}

type MessageSentData struct {
	MessageID string                 `json:"message_id"`
	ChatID    string                 `json:"chat_id"`
	UserID    string                 `json:"user_id"`
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
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
	Typing bool   `json:"typing"`
}

type UserJoinedNotificationData struct {
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username,omitempty"`
	JoinedAt time.Time `json:"joined_at"`
}

type UserLeftNotificationData struct {
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username,omitempty"`
	LeftAt   time.Time `json:"left_at"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ChatID  string `json:"chat_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}
