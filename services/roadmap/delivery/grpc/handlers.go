package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

type RoadmapServer struct {
	uc roadmap.Usecase
	roadmapclient.UnimplementedRoadmapServiceServer
}

func NewRoadmapServer(usecase roadmap.Usecase) *RoadmapServer {
	return &RoadmapServer{uc: usecase}
}

func (s *RoadmapServer) Create(ctx context.Context, req *roadmapclient.CreateRequest) (*roadmapclient.CreateResponse, error) {
	roadmapDTO, err := s.convertCreateRequestToDTO(req)
	if err != nil {
		return &roadmapclient.CreateResponse{
			Error: err.Error(),
		}, nil
	}

	createdRoadmap, err := s.uc.Create(ctx, roadmapDTO)
	if err != nil {
		return &roadmapclient.CreateResponse{
			Error: err.Error(),
		}, nil
	}

	protoRoadmap := s.convertRoadmapToProto(createdRoadmap)

	return &roadmapclient.CreateResponse{
		Roadmap: protoRoadmap,
	}, nil
}

func (s *RoadmapServer) Delete(ctx context.Context, req *roadmapclient.DeleteRequest) (*roadmapclient.DeleteResponse, error) {
	roadmapID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &roadmapclient.DeleteResponse{
			Success: false,
			Error:   "invalid roadmap id format",
		}, nil
	}

	err = s.uc.Delete(ctx, roadmapID)
	if err != nil {
		return &roadmapclient.DeleteResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &roadmapclient.DeleteResponse{
		Success: true,
	}, nil
}

func (s *RoadmapServer) GetByID(ctx context.Context, req *roadmapclient.GetByIDRequest) (*roadmapclient.GetByIDResponse, error) {
	roadmapID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &roadmapclient.GetByIDResponse{
			Error: "invalid roadmap id format",
		}, nil
	}

	roadmap, err := s.uc.GetByID(ctx, roadmapID)
	if err != nil {
		return &roadmapclient.GetByIDResponse{
			Error: err.Error(),
		}, nil
	}

	protoRoadmap := s.convertRoadmapToProto(&roadmap.Roadmap)

	return &roadmapclient.GetByIDResponse{
		Roadmap: protoRoadmap,
	}, nil
}

func (s *RoadmapServer) convertCreateRequestToDTO(req *roadmapclient.CreateRequest) (*dto.RoadmapDTO, error) {
	roadmapDTO := &dto.RoadmapDTO{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, protoNode := range req.Nodes {
		node, err := s.convertProtoNodeToDTO(protoNode)
		if err != nil {
			return nil, err
		}
		roadmapDTO.Nodes = append(roadmapDTO.Nodes, *node)
	}

	for _, protoEdge := range req.Edges {
		roadmapDTO.Edges = append(roadmapDTO.Edges, dto.EdgeDTO{
			Source: protoEdge.Source,
			Target: protoEdge.Target,
			ID:     protoEdge.Id,
		})
	}

	return roadmapDTO, nil
}

func (s *RoadmapServer) convertProtoNodeToDTO(protoNode *roadmapclient.Node) (*dto.NodeDTO, error) {
	nodeID, err := uuid.Parse(protoNode.Id)
	if err != nil {
		return nil, err
	}

	return &dto.NodeDTO{
		ID:   nodeID,
		Type: protoNode.Type,
		Position: dto.Position{
			X: protoNode.Position.X,
			Y: protoNode.Position.Y,
		},
		Data: dto.NodeData{
			Label: protoNode.Data.Label,
			Type:  protoNode.Data.Type,
		},
		Measured: dto.Measured{
			Width:  protoNode.Measured.Width,
			Height: protoNode.Measured.Height,
		},
		Selected: protoNode.Selected,
		Dragging: protoNode.Dragging,
	}, nil
}

func (s *RoadmapServer) convertRoadmapToProto(roadmap *dto.RoadmapDTO) *roadmapclient.Roadmap {
	protoRoadmap := &roadmapclient.Roadmap{
		Id:        roadmap.ID.Hex(),
		CreatedAt: timestamppb.New(roadmap.CreatedAt),
		UpdatedAt: timestamppb.New(roadmap.UpdatedAt),
	}

	for _, node := range roadmap.Nodes {
		protoRoadmap.Nodes = append(protoRoadmap.Nodes, s.convertNodeToProto(node))
	}

	for _, edge := range roadmap.Edges {
		protoRoadmap.Edges = append(protoRoadmap.Edges, &roadmapclient.Edge{
			Source: edge.Source,
			Target: edge.Target,
			Id:     edge.ID,
		})
	}

	return protoRoadmap
}

func (s *RoadmapServer) convertNodeToProto(node dto.NodeDTO) *roadmapclient.Node {
	return &roadmapclient.Node{
		Id:   node.ID.String(),
		Type: node.Type,
		Position: &roadmapclient.Position{
			X: node.Position.X,
			Y: node.Position.Y,
		},
		Data: &roadmapclient.NodeData{
			Label: node.Data.Label,
			Type:  node.Data.Type,
		},
		Measured: &roadmapclient.Measured{
			Width:  node.Measured.Width,
			Height: node.Measured.Height,
		},
		Selected: node.Selected,
		Dragging: node.Dragging,
	}
}
