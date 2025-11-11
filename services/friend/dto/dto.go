package dto

import (
	"time"

	"github.com/google/uuid"
)

type FriendResponseDTO struct {
	UserID         uuid.UUID  `json:"user_id"`
	Username       string     `json:"username"`
	AvatarURL      *string    `json:"avatar_url,omitempty"`
	SharedRoadmaps int        `json:"shared_roadmaps"`
	ChatID         *uuid.UUID `json:"chat_id,omitempty"`
}

type GetFriendsResponseDTO struct {
	Friends []FriendResponseDTO `json:"friends"`
}

type FriendRequestResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	FromUser  UserDTO   `json:"from_user"`
	ToUser    UserDTO   `json:"to_user"`
	Status    string    `json:"status"`
	Message   *string   `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetFriendRequestsResponseDTO struct {
	Received []FriendRequestResponseDTO `json:"received"`
	Sent     []FriendRequestResponseDTO `json:"sent"`
}

type CreateFriendRequestDTO struct {
	TargetUserID uuid.UUID `json:"target_user_id" validate:"required"`
	Message      *string   `json:"message,omitempty"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}
