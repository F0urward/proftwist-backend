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
	ID         uuid.UUID    `json:"id" db:"id"`
	FromUserID uuid.UUID    `json:"from_user_id" db:"from_user_id"`
	ToUserID   uuid.UUID    `json:"to_user_id" db:"to_user_id"`
	Status     FriendStatus `json:"status" db:"status"`
	Message    *string      `json:"message" db:"message"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" db:"updated_at"`
}

type Friend struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	FriendID  uuid.UUID `json:"friend_id" db:"friend_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
