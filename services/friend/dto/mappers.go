package dto

import (
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

func FriendToDTO(userID uuid.UUID, userData *UserDTO, sharedRoadmaps int, chatID *uuid.UUID) FriendResponseDTO {
	return FriendResponseDTO{
		UserID:         userID,
		Username:       userData.Username,
		AvatarURL:      userData.AvatarURL,
		SharedRoadmaps: sharedRoadmaps,
		ChatID:         chatID,
	}
}

func FriendsToDTO(friendIDs []uuid.UUID, userData map[uuid.UUID]*UserDTO, sharedRoadmaps map[uuid.UUID]int, chatIDs map[uuid.UUID]*uuid.UUID) GetFriendsResponseDTO {
	var friendDTOs []FriendResponseDTO

	for _, friendID := range friendIDs {
		if userInfo, exists := userData[friendID]; exists {
			sharedCount := sharedRoadmaps[friendID]
			chatID := chatIDs[friendID]
			friendDTOs = append(friendDTOs, FriendToDTO(friendID, userInfo, sharedCount, chatID))
		}
	}

	return GetFriendsResponseDTO{
		Friends: friendDTOs,
	}
}

func FriendRequestToDTO(request *entities.FriendRequest, fromUserData, toUserData *UserDTO) FriendRequestResponseDTO {
	return FriendRequestResponseDTO{
		ID:        request.ID,
		FromUser:  UserToDTO(fromUserData),
		ToUser:    UserToDTO(toUserData),
		Status:    string(request.Status),
		Message:   request.Message,
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
	}
}

func FriendRequestsToDTO(requests []*entities.FriendRequest, userData map[uuid.UUID]*UserDTO) []FriendRequestResponseDTO {
	var requestDTOs []FriendRequestResponseDTO

	for _, request := range requests {
		fromUserData := userData[request.FromUserID]
		toUserData := userData[request.ToUserID]

		if fromUserData != nil && toUserData != nil {
			requestDTOs = append(requestDTOs, FriendRequestToDTO(request, fromUserData, toUserData))
		}
	}

	return requestDTOs
}

func GetFriendRequestsResponseToDTO(receivedRequests, sentRequests []*entities.FriendRequest, userData map[uuid.UUID]*UserDTO) GetFriendRequestsResponseDTO {
	return GetFriendRequestsResponseDTO{
		Received: FriendRequestsToDTO(receivedRequests, userData),
		Sent:     FriendRequestsToDTO(sentRequests, userData),
	}
}

func UserToDTO(userData *UserDTO) UserDTO {
	if userData == nil {
		return UserDTO{}
	}
	return UserDTO{
		ID:        userData.ID,
		Username:  userData.Username,
		AvatarURL: userData.AvatarURL,
	}
}

func CreateFriendRequestToEntity(userID uuid.UUID, request *CreateFriendRequestDTO) *entities.FriendRequest {
	return &entities.FriendRequest{
		FromUserID: userID,
		ToUserID:   request.TargetUserID,
		Message:    request.Message,
		Status:     entities.FriendStatusPending,
	}
}
