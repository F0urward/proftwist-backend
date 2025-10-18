package dto

import (
	"github.com/F0urward/proftwist-backend/internal/entities"
)

func UserEntityToDTO(user *entities.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
	}
}

func UserTokenToDTO(user *entities.User, token string) *UserTokenDTO {
	return &UserTokenDTO{
		User:  UserEntityToDTO(user),
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
