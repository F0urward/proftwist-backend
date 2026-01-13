package dto

import (
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

func UserToDTO(user *entities.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
	}
}

func UserListToDTO(users []*entities.User) []UserDTO {
	var userDTOs []UserDTO

	for _, user := range users {
		userDTOs = append(userDTOs, UserToDTO(user))
	}

	return userDTOs
}

func UserListToPublicDTO(users []*entities.User, friendshipStatus map[uuid.UUID]*FriendshipStatusDTO) []UserPublicDTO {
	var userPublicDTOs []UserPublicDTO

	for _, user := range users {
		if user == nil {
			continue
		}

		dto := UserPublicDTO{
			ID:        user.ID,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl,
		}

		if status, exists := friendshipStatus[user.ID]; exists {
			dto.FriendshipStatus = status
		}

		userPublicDTOs = append(userPublicDTOs, dto)
	}

	return userPublicDTOs
}

func UserTokenToDTO(user *entities.User, token string) *UserTokenDTO {
	return &UserTokenDTO{
		User:  UserToDTO(user),
		Token: token,
	}
}

func RegisterRequestToEntity(request *RegisterRequestDTO, passwordHash string) *entities.User {
	return &entities.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: passwordHash,
	}
}

func LoginRequestToEntity(request *LoginRequestDTO, passwordHash string) *entities.User {
	return &entities.User{
		Email:        request.Email,
		PasswordHash: passwordHash,
	}
}

func UpdateUserRequestToEntity(existing *entities.User, request *UpdateUserRequestDTO) *entities.User {
	updated := *existing

	if request.Username != "" {
		updated.Username = request.Username
	}
	if request.Email != "" {
		updated.Email = request.Email
	}
	return &updated
}

func UploadAvatarRequestToUploadInputEntity(req *UploadAvatarRequestDTO) *entities.UploadInput {
	return &entities.UploadInput{
		File:        req.File,
		Name:        req.Name,
		Size:        req.Size,
		ContentType: req.ContentType,
		BucketName:  req.BucketName,
	}
}
