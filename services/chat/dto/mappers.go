package dto

import (
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

func GroupChatToDTO(chat *entities.GroupChat) GroupChatResponseDTO {
	return GroupChatResponseDTO{
		ID:            chat.ID,
		Title:         chat.Title,
		AvatarURL:     chat.AvatarURL,
		RoadmapNodeID: chat.RoadmapNodeID,
		CreatedAt:     chat.CreatedAt,
		UpdatedAt:     chat.UpdatedAt,
	}
}

func GroupChatListToDTO(chats []*entities.GroupChat) GroupChatListResponseDTO {
	var chatDTOs []GroupChatResponseDTO

	for _, chat := range chats {
		chatDTOs = append(chatDTOs, GroupChatToDTO(chat))
	}

	return GroupChatListResponseDTO{
		GroupChats: chatDTOs,
	}
}

func DirectChatToDTO(chat *entities.DirectChat) DirectChatResponseDTO {
	return DirectChatResponseDTO{
		ID:        chat.ID,
		User1ID:   chat.User1ID,
		User2ID:   chat.User2ID,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
}

func DirectChatListToDTO(chats []*entities.DirectChat) DirectChatListResponseDTO {
	var chatDTOs []DirectChatResponseDTO

	for _, chat := range chats {
		chatDTOs = append(chatDTOs, DirectChatToDTO(chat))
	}

	return DirectChatListResponseDTO{
		DirectChats: chatDTOs,
	}
}

func CreateGroupChatRequestToEntity(request *CreateGroupChatRequestDTO) *entities.GroupChat {
	return &entities.GroupChat{
		Title:         request.Title,
		AvatarURL:     request.AvatarURL,
		RoadmapNodeID: request.RoadmapNodeID,
	}
}

func CreateDirectChatRequestToEntity(request *CreateDirectChatRequestDTO, currentUserID uuid.UUID) *entities.DirectChat {
	return &entities.DirectChat{
		User1ID: currentUserID,
		User2ID: request.OtherUserID,
	}
}

func CreateGroupChatResponseFromEntity(chat *entities.GroupChat) CreateGroupChatResponseDTO {
	return CreateGroupChatResponseDTO{
		GroupChat: GroupChatToDTO(chat),
	}
}

func CreateDirectChatResponseFromEntity(chat *entities.DirectChat) CreateDirectChatResponseDTO {
	return CreateDirectChatResponseDTO{
		DirectChat: DirectChatToDTO(chat),
	}
}

func MemberToDTO(userID uuid.UUID, username, avatarURL string) MemberResponseDTO {
	return MemberResponseDTO{
		UserID:    userID,
		Username:  username,
		AvatarURL: avatarURL,
	}
}

func GroupChatMemberToDTO(member *entities.GroupChatMember, username, avatarURL string) MemberResponseDTO {
	return MemberResponseDTO{
		UserID:    member.UserID,
		Username:  username,
		AvatarURL: avatarURL,
	}
}

func GroupChatMemberListToDTO(members []*entities.GroupChatMember, userData map[uuid.UUID]MemberResponseDTO) ChatMemberListResponseDTO {
	var memberDTOs []MemberResponseDTO

	for _, member := range members {
		if userData, exists := userData[member.UserID]; exists {
			memberDTOs = append(memberDTOs, userData)
		} else {
			memberDTOs = append(memberDTOs, MemberResponseDTO{
				UserID: member.UserID,
			})
		}
	}

	return ChatMemberListResponseDTO{
		Members: memberDTOs,
	}
}

func DirectChatMembersToDTO(user1ID, user2ID uuid.UUID, user1Data, user2Data MemberResponseDTO) ChatMemberListResponseDTO {
	memberDTOs := make([]MemberResponseDTO, 2)

	memberDTOs[0] = user1Data
	memberDTOs[1] = user2Data

	return ChatMemberListResponseDTO{
		Members: memberDTOs,
	}
}

func MessageToDTO(message *entities.Message, username, avatarURL string) ChatMessageResponseDTO {
	return ChatMessageResponseDTO{
		ID:     message.ID,
		ChatID: message.ChatID,
		User: MemberResponseDTO{
			UserID:    message.UserID,
			Username:  username,
			AvatarURL: avatarURL,
		},
		Content:   message.Content,
		Metadata:  message.Metadata,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func MessageListToDTO(messages []*entities.Message, userData map[uuid.UUID]MemberResponseDTO) []ChatMessageResponseDTO {
	var messageDTOs []ChatMessageResponseDTO

	for _, message := range messages {
		if userInfo, exists := userData[message.UserID]; exists {
			messageDTOs = append(messageDTOs, ChatMessageResponseDTO{
				ID:        message.ID,
				ChatID:    message.ChatID,
				User:      userInfo,
				Content:   message.Content,
				Metadata:  message.Metadata,
				CreatedAt: message.CreatedAt,
				UpdatedAt: message.UpdatedAt,
			})
		} else {
			messageDTOs = append(messageDTOs, ChatMessageResponseDTO{
				ID:     message.ID,
				ChatID: message.ChatID,
				User: MemberResponseDTO{
					UserID: message.UserID,
				},
				Content:   message.Content,
				Metadata:  message.Metadata,
				CreatedAt: message.CreatedAt,
				UpdatedAt: message.UpdatedAt,
			})
		}
	}

	return messageDTOs
}

func GetChatMessagesResponseToDTO(messages []*entities.Message, userData map[uuid.UUID]MemberResponseDTO) GetChatMessagesResponseDTO {
	messageDTOs := MessageListToDTO(messages, userData)

	return GetChatMessagesResponseDTO{
		ChatMessages: messageDTOs,
	}
}

func SendMessageRequestToEntity(request *SendMessageRequestDTO) *entities.Message {
	return &entities.Message{
		ChatID:   request.ChatID,
		UserID:   request.UserID,
		Content:  request.Content,
		Metadata: request.Metadata,
	}
}
