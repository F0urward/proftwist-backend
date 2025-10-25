// internal/infrastructure/client/websocketclient/dto/websocket_dto.go
package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WebSocketMessageType string

const (
	WebSocketMessageTypeChat   WebSocketMessageType = "chat"
	WebSocketMessageTypeTyping WebSocketMessageType = "typing"
	WebSocketMessageTypeRead   WebSocketMessageType = "read"
	WebSocketMessageTypeJoin   WebSocketMessageType = "join"
	WebSocketMessageTypeLeave  WebSocketMessageType = "leave"
	WebSocketMessageTypeError  WebSocketMessageType = "error"
)

type WebSocketMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Data      json.RawMessage      `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
	UserID    string               `json:"user_id,omitempty"`
	ChatID    string               `json:"chat_id,omitempty"`
}

type ChatMessageData struct {
	ChatID   string                 `json:"chat_id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type TypingMessageData struct {
	ChatID string `json:"chat_id"`
	Typing bool   `json:"typing"`
}

type ReadReceiptData struct {
	ChatID    string    `json:"chat_id"`
	MessageID uuid.UUID `json:"message_id"`
	ReadAt    time.Time `json:"read_at"`
}

type JoinChatData struct {
	ChatID string `json:"chat_id"`
}

type ErrorMessageData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
