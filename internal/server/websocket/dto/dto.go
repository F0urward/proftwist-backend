package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WebSocketMessageType string

const (
	WebSocketMessageTypeChat WebSocketMessageType = "chat"
	WebSocketMessageTypeRead WebSocketMessageType = "read"
	WebSocketMessageTypeJoin WebSocketMessageType = "join"
	// WebSocketMessageTypeLeave  WebSocketMessageType = "leave"
	WebSocketMessageTypeError WebSocketMessageType = "error"
)

type WebSocketMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Data      json.RawMessage      `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
	UserID    string               `json:"user_id,omitempty"`
	ChatID    string               `json:"chat_id,omitempty"`
}

type ChatMessageData struct {
	MessageID string                 `json:"message_id"`
	ChatID    string                 `json:"chat_id"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ReadReceiptData struct {
	ChatID    string    `json:"chat_id"`
	MessageID uuid.UUID `json:"message_id"`
	ReadAt    time.Time `json:"read_at"`
}

type JoinChatData struct {
	ChatID string `json:"chat_id"`
}
