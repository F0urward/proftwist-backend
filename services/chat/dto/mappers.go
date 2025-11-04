package dto

import (
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
)

func ChatToDTO(chat *entities.Chat) ChatResponseDTO {
	return ChatResponseDTO{
		ID:          chat.ID,
		Type:        string(chat.Type),
		Title:       chat.Title,
		Description: chat.Description,
		AvatarURL:   chat.AvatarURL,
		CreatedBy:   chat.CreatedBy,
		CreatedAt:   chat.CreatedAt,
		UpdatedAt:   chat.UpdatedAt,
	}
}

func ChatListToDTO(chats []*entities.Chat) []ChatResponseDTO {
	var chatDTOs []ChatResponseDTO

	for _, chat := range chats {
		chatDTOs = append(chatDTOs, ChatToDTO(chat))
	}

	return chatDTOs
}

func ChatMemberToDTO(member *entities.ChatMember) ChatMemberResponseDTO {
	return ChatMemberResponseDTO{
		ID:       member.ID,
		UserID:   member.UserID,
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
		LastRead: member.LastRead,
	}
}

func ChatMemberListToDTO(members []*entities.ChatMember) []ChatMemberResponseDTO {
	var memberDTOs []ChatMemberResponseDTO

	for _, member := range members {
		memberDTOs = append(memberDTOs, ChatMemberToDTO(member))
	}

	return memberDTOs
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

func CreateChatRequestToEntity(request *CreateChatRequestDTO) (*entities.Chat, error) {
	chatType := entities.ChatType(request.Type)

	chat := &entities.Chat{
		ID:          uuid.New(),
		Type:        chatType,
		Title:       request.Title,
		Description: request.Description,
		AvatarURL:   request.AvatarURL,
		CreatedBy:   request.CreatedByID,
	}

	return chat, nil
}

func GetChatMessagesResponseToDTO(messages []*entities.Message) GetChatMessagesResponseDTO {
	messageDTOs := MessageListToDTO(messages)

	return GetChatMessagesResponseDTO{
		ChatMessages: messageDTOs,
	}
}
