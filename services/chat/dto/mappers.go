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

func GroupChatMemberToDTO(member *entities.GroupChatMember) ChatMemberResponseDTO {
	return ChatMemberResponseDTO{
		UserID: member.UserID,
	}
}

func GroupChatMemberListToDTO(members []*entities.GroupChatMember) ChatMemberListResponseDTO {
	var memberDTOs []ChatMemberResponseDTO

	for _, member := range members {
		memberDTOs = append(memberDTOs, GroupChatMemberToDTO(member))
	}

	return ChatMemberListResponseDTO{
		Members: memberDTOs,
	}
}

func DirectChatMembersToDTO(user1ID uuid.UUID, user2ID uuid.UUID) ChatMemberListResponseDTO {
	memberDTOs := make([]ChatMemberResponseDTO, 2)

	memberDTOs[0] = ChatMemberResponseDTO{UserID: user1ID}
	memberDTOs[1] = ChatMemberResponseDTO{UserID: user2ID}

	return ChatMemberListResponseDTO{
		Members: memberDTOs,
	}
}

func MessageToDTO(message *entities.Message) ChatMessageResponseDTO {
	return ChatMessageResponseDTO{
		ID:        message.ID,
		ChatID:    message.ChatID,
		UserID:    message.UserID,
		Content:   message.Content,
		Metadata:  message.Metadata,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func MessageListToDTO(messages []*entities.Message) []ChatMessageResponseDTO {
	var messageDTOs []ChatMessageResponseDTO

	for _, message := range messages {
		messageDTOs = append(messageDTOs, MessageToDTO(message))
	}

	return messageDTOs
}

func GetChatMessagesResponseToDTO(messages []*entities.Message) GetChatMessagesResponseDTO {
	messageDTOs := MessageListToDTO(messages)

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
