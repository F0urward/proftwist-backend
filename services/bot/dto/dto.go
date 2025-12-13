package dto

import (
	"time"
)

type EventType string

const (
	MessageForBotType EventType = "chat.message_for_bot"
)

type BaseEvent struct {
	Type EventType `json:"type"`
}

type MessageForBotEvent struct {
	Type       EventType `json:"type"`
	ChatID     string    `json:"chat_id"`
	ChatTitle  string    `json:"chat_title"`
	Content    string    `json:"content"`
	ReceivedAt time.Time `json:"timestamp"`
}
