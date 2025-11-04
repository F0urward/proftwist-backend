package dto

import (
	"encoding/json"
	"time"
)

type WebSocketMessageType string

const (
	// Action types (handled by backend)
	WebSocketMessageTypeSendMessage WebSocketMessageType = "send_message" // Отправка сообщения
	WebSocketMessageTypeTyping      WebSocketMessageType = "typing"       // Печать
	WebSocketMessageTypeRead        WebSocketMessageType = "read"         // Прочтение сообщения
	WebSocketMessageTypeJoin        WebSocketMessageType = "join"         // Присоединение к чату
	WebSocketMessageTypeLeave       WebSocketMessageType = "leave"        // Выход из чата

	// Notification types (broadcast only)
	WebSocketMessageTypeMessageSent        WebSocketMessageType = "message_sent"        // Сообщение отправлено
	WebSocketMessageTypeMessageDelivered   WebSocketMessageType = "message_delivered"   // Сообщение доставлено
	WebSocketMessageTypeTypingNotification WebSocketMessageType = "typing_notification" // Уведомление о печати
	WebSocketMessageTypeUserJoined         WebSocketMessageType = "user_joined"         // Уведомление о присоединении пользователя
	WebSocketMessageTypeUserLeft           WebSocketMessageType = "user_left"           // Уведомление о выходе пользователя
	WebSocketMessageTypeError              WebSocketMessageType = "error"               // Ошибка
)

type WebSocketMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Data      json.RawMessage      `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}

// Action Data Structures (for client-initiated actions)

type SendMessageData struct {
	MessageID string                 `json:"message_id,omitempty"`
	ChatID    string                 `json:"chat_id"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type TypingData struct {
	ChatID string `json:"chat_id"`
	Typing bool   `json:"typing"` // true - начал печатать, false - закончил
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

// Notification Data Structures (for server-initiated broadcasts)

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
	Typing bool   `json:"typing"` // true - пользователь печатает, false - перестал печатать
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
