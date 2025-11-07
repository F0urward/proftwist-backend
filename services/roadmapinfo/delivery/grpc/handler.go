package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo/dto"
)

type RoadmapInfoServer struct {
	uc roadmapinfo.Usecase
	roadmapinfoclient.UnimplementedRoadmapInfoServiceServer
}

func NewRoadmapInfoServer(usecase roadmapinfo.Usecase) *RoadmapInfoServer {
	return &RoadmapInfoServer{uc: usecase}
}

func (s *RoadmapInfoServer) GetByRoadmapID(ctx context.Context, req *roadmapinfoclient.GetByRoadmapIDRequest) (*roadmapinfoclient.GetByRoadmapIDResponse, error) {
	roadmapInfo, err := s.uc.GetByRoadmapID(ctx, req.RoadmapId)
	if err != nil {
		return &roadmapinfoclient.GetByRoadmapIDResponse{
			Error: err.Error(),
		}, nil
	}

	protoRoadmapInfo := convertRoadmapInfoToProto(&roadmapInfo.RoadmapInfo)

	return &roadmapinfoclient.GetByRoadmapIDResponse{
		RoadmapInfo: protoRoadmapInfo,
	}, nil
}

func convertRoadmapInfoToProto(dto *dto.RoadmapInfoDTO) *roadmapinfoclient.RoadmapInfo {
	if dto == nil {
		return nil
	}

	return &roadmapinfoclient.RoadmapInfo{
		Id:                      dto.ID,
		RoadmapId:               dto.RoadmapID,
		AuthorId:                dto.AuthorID,
		CategoryId:              dto.CategoryID,
		Name:                    dto.Name,
		Description:             dto.Description,
		IsPublic:                dto.IsPublic,
		ReferencedRoadmapInfoId: dto.ReferencedRoadmapInfoID,
		SubscriberCount:         int32(dto.SubscriberCount),
		CreatedAt:               timestamppb.New(dto.CreatedAt),
		UpdatedAt:               timestamppb.New(dto.UpdatedAt),
	}
}
