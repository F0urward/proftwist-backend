package dto

import (
	"io"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"image"`
}

type UserPublicDTO struct {
	ID               uuid.UUID            `json:"id"`
	Username         string               `json:"username"`
	AvatarUrl        string               `json:"image"`
	FriendshipStatus *FriendshipStatusDTO `json:"friendship_status,omitempty"`
}

type FriendshipStatusDTO struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id,omitempty"`
	IsSender  bool   `json:"is_sender"`
}

type RegisterRequestDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetMeResponseDTO struct {
	User UserDTO `json:"user"`
}

type UserTokenDTO struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}

type VKOauthLinkResponse struct {
	VKOauthURL string `json:"vk_oauth_url"`
}

type VKCallbackRequestDTO struct {
	Code     string `json:"code"`
	State    string `json:"state"`
	DeviceID string `json:"device_id,omitempty"`
}

type UpdateUserRequestDTO struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type GetUserByIDResponseDTO struct {
	User UserDTO `json:"user"`
}

type GetUsersByIDsResponseDTO struct {
	Users []UserDTO `json:"users"`
}

type SearchUsersResponseDTO struct {
	Users []UserPublicDTO `json:"users"`
}

type UploadAvatarRequestDTO struct {
	UserID      uuid.UUID `json:"-"`
	File        io.Reader `json:"-"`
	Name        string    `json:"-"`
	Size        int64     `json:"-"`
	ContentType string    `json:"-"`
	BucketName  string    `json:"-"`
}

type UploadAvatarResponseDTO struct {
	AvatarUrl string `json:"avatar_url"`
}
