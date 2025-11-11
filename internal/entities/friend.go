package entities

import (
	"time"

	"github.com/google/uuid"
)

type FriendStatus string

const (
	FriendStatusPending  FriendStatus = "pending"
	FriendStatusAccepted FriendStatus = "accepted"
	FriendStatusRejected FriendStatus = "rejected"
)

type FriendRequest struct {
	ID         uuid.UUID
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	Status     FriendStatus
	Message    *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Friend struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FriendID  uuid.UUID
	ChatID    *uuid.UUID
	CreatedAt time.Time
}
