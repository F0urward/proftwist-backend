package dto

import "github.com/google/uuid"

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"image"`
}

type RegisterRequestDTO struct {
	Role     string `json:"role"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserTokenDTO struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}
