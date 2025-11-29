package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/F0urward/proftwist-backend/services/chat/dto"
)

type ChatServer struct {
	uc chat.Usecase
	chatclient.UnimplementedChatServiceServer
}

func NewChatServer(usecase chat.Usecase) chatclient.ChatServiceServer {
	return &ChatServer{uc: usecase}
}

func (s *ChatServer) SendGroupChatMessage(ctx context.Context, req *chatclient.SendGroupChatMessageRequest) (*chatclient.SendGroupChatMessageResponse, error) {
	chatID, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid chat ID: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	var metadata map[string]interface{}
	if req.Metadata != nil {
		metadata = req.Metadata.AsMap()
	}

	sendMessageReq := dto.SendMessageRequestDTO{
		ChatID:   chatID,
		UserID:   userID,
		Content:  req.Content,
		Metadata: metadata,
	}

	message, err := s.uc.SendGroupChatMessage(ctx, &sendMessageReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	pbMessage, err := s.convertChatMessageToProto(*message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert message: %v", err)
	}

	return &chatclient.SendGroupChatMessageResponse{
		Message: pbMessage,
	}, nil
}

func (s *ChatServer) convertChatMessageToProto(msg dto.ChatMessageResponseDTO) (*chatclient.ChatMessage, error) {
	var pbMetadata *structpb.Struct
	if msg.Metadata != nil {
		var err error
		pbMetadata, err = structpb.NewStruct(msg.Metadata)
		if err != nil {
			return nil, err
		}
	}

	return &chatclient.ChatMessage{
		Id:          msg.ID.String(),
		GroupChatId: msg.ChatID.String(),
		UserId:      msg.User.UserID.String(),
		Content:     msg.Content,
		Metadata:    pbMetadata,
		CreatedAt:   timestamppb.New(msg.CreatedAt),
		UpdatedAt:   timestamppb.New(msg.UpdatedAt),
	}, nil
}

func (s *ChatServer) CreateGroupChat(ctx context.Context, req *chatclient.CreateGroupChatRequest) (*chatclient.CreateGroupChatResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &chatclient.CreateGroupChatResponse{
			Error: "invalid user id format",
		}, nil
	}

	createReq, err := s.convertCreateGroupChatRequestToDTO(req)
	if err != nil {
		return &chatclient.CreateGroupChatResponse{
			Error: err.Error(),
		}, nil
	}

	createdChat, err := s.uc.CreateGroupChat(ctx, userID, createReq)
	if err != nil {
		return &chatclient.CreateGroupChatResponse{
			Error: err.Error(),
		}, nil
	}

	protoGroupChat := s.convertGroupChatToProto(&createdChat.GroupChat)

	return &chatclient.CreateGroupChatResponse{
		GroupChat: protoGroupChat,
	}, nil
}

func (s *ChatServer) DeleteGroupChat(ctx context.Context, req *chatclient.DeleteGroupChatRequest) (*chatclient.DeleteGroupChatResponse, error) {
	chatID, err := uuid.Parse(req.ChatId)
	if err != nil {
		return &chatclient.DeleteGroupChatResponse{
			Success: false,
			Error:   "invalid chat id format",
		}, nil
	}

	err = s.uc.DeleteGroupChat(ctx, chatID)
	if err != nil {
		return &chatclient.DeleteGroupChatResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &chatclient.DeleteGroupChatResponse{
		Success: true,
	}, nil
}

func (s *ChatServer) CreateDirectChat(ctx context.Context, req *chatclient.CreateDirectChatRequest) (*chatclient.CreateDirectChatResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &chatclient.CreateDirectChatResponse{
			Error: "invalid user id format",
		}, nil
	}

	createReq, err := s.convertCreateDirectChatRequestToDTO(req)
	if err != nil {
		return &chatclient.CreateDirectChatResponse{
			Error: err.Error(),
		}, nil
	}

	createdChat, err := s.uc.CreateDirectChat(ctx, userID, createReq)
	if err != nil {
		return &chatclient.CreateDirectChatResponse{
			Error: err.Error(),
		}, nil
	}

	protoDirectChat := s.convertDirectChatToProto(&createdChat.DirectChat)

	return &chatclient.CreateDirectChatResponse{
		DirectChat: protoDirectChat,
	}, nil
}

func (s *ChatServer) DeleteDirectChat(ctx context.Context, req *chatclient.DeleteDirectChatRequest) (*chatclient.DeleteDirectChatResponse, error) {
	chatID, err := uuid.Parse(req.ChatId)
	if err != nil {
		return &chatclient.DeleteDirectChatResponse{
			Success: false,
			Error:   "invalid chat id format",
		}, nil
	}

	err = s.uc.DeleteDirectChat(ctx, chatID)
	if err != nil {
		return &chatclient.DeleteDirectChatResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &chatclient.DeleteDirectChatResponse{
		Success: true,
	}, nil
}

func (s *ChatServer) GetGroupChatByNode(ctx context.Context, req *chatclient.GetGroupChatByNodeRequest) (*chatclient.GetGroupChatByNodeResponse, error) {
	groupChat, err := s.uc.GetGroupChatByNode(ctx, req.NodeId)
	if err != nil {
		return &chatclient.GetGroupChatByNodeResponse{
			Error: err.Error(),
		}, nil
	}

	protoGroupChat := s.convertGroupChatToProto(groupChat)

	return &chatclient.GetGroupChatByNodeResponse{
		GroupChat: protoGroupChat,
	}, nil
}

func (s *ChatServer) convertCreateGroupChatRequestToDTO(req *chatclient.CreateGroupChatRequest) (*dto.CreateGroupChatRequestDTO, error) {
	memberIDs := make([]uuid.UUID, len(req.MemberIds))
	for i, memberIDStr := range req.MemberIds {
		memberID, err := uuid.Parse(memberIDStr)
		if err != nil {
			return nil, err
		}
		memberIDs[i] = memberID
	}

	var title *string
	if req.Title != "" {
		title = &req.Title
	}

	var avatarURL *string
	if req.AvatarUrl != "" {
		avatarURL = &req.AvatarUrl
	}

	var roadmapNodeID *string
	if req.RoadmapNodeId != "" {
		roadmapNodeID = &req.RoadmapNodeId
	}

	return &dto.CreateGroupChatRequestDTO{
		Title:         title,
		AvatarURL:     avatarURL,
		RoadmapNodeID: roadmapNodeID,
		MemberIDs:     memberIDs,
	}, nil
}

func (s *ChatServer) convertCreateDirectChatRequestToDTO(req *chatclient.CreateDirectChatRequest) (*dto.CreateDirectChatRequestDTO, error) {
	otherUserID, err := uuid.Parse(req.OtherUserId)
	if err != nil {
		return nil, err
	}

	return &dto.CreateDirectChatRequestDTO{
		OtherUserID: otherUserID,
	}, nil
}

func (s *ChatServer) convertGroupChatToProto(chat *dto.GroupChatResponseDTO) *chatclient.GroupChat {
	protoChat := &chatclient.GroupChat{
		Id:        chat.ID.String(),
		CreatedAt: timestamppb.New(chat.CreatedAt),
		UpdatedAt: timestamppb.New(chat.UpdatedAt),
	}

	if chat.Title != nil {
		protoChat.Title = *chat.Title
	}
	if chat.AvatarURL != nil {
		protoChat.AvatarUrl = *chat.AvatarURL
	}
	if chat.RoadmapNodeID != nil {
		protoChat.RoadmapNodeId = *chat.RoadmapNodeID
	}

	return protoChat
}

func (s *ChatServer) convertDirectChatToProto(chat *dto.DirectChatResponseDTO) *chatclient.DirectChat {
	return &chatclient.DirectChat{
		Id:        chat.ID.String(),
		User1Id:   chat.Members[0].UserID.String(),
		User2Id:   chat.Members[1].UserID.String(),
		CreatedAt: timestamppb.New(chat.CreatedAt),
		UpdatedAt: timestamppb.New(chat.UpdatedAt),
	}
}
