package grpc

import (
	"context"
	"fmt"
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

func NewRoadmapServer(usecase roadmap.Usecase) roadmapclient.RoadmapServiceServer {
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

	protoRoadmap := s.convertRoadmapWithMaterialsToProto(&createdRoadmap.RoadmapWithMaterials)

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

func (s *RoadmapServer) GetByIDWithMaterials(ctx context.Context, req *roadmapclient.GetByIDWithMaterialsRequest) (*roadmapclient.GetByIDWithMaterialsResponse, error) {
	roadmapID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &roadmapclient.GetByIDWithMaterialsResponse{
			Error: "invalid roadmap id format",
		}, nil
	}

	roadmapWithMaterials, err := s.uc.GetByIDWithMaterials(ctx, roadmapID)
	if err != nil {
		return &roadmapclient.GetByIDWithMaterialsResponse{
			Error: err.Error(),
		}, nil
	}

	protoRoadmap := s.convertRoadmapWithMaterialsToProto(&roadmapWithMaterials.RoadmapWithMaterials)

	return &roadmapclient.GetByIDWithMaterialsResponse{
		Roadmap: protoRoadmap,
	}, nil
}

func (s *RoadmapServer) RegenerateNodeIDs(ctx context.Context, req *roadmapclient.RegenerateNodeIDsRequest) (*roadmapclient.RegenerateNodeIDsResponse, error) {
	roadmapDTO, err := s.convertProtoRoadmapWithMaterialsToDTO(req.Roadmap)
	if err != nil {
		return &roadmapclient.RegenerateNodeIDsResponse{
			Error: err.Error(),
		}, nil
	}

	regeneratedRoadmap := s.uc.RegenerateNodeIDs(roadmapDTO)

	protoRoadmap := s.convertRoadmapWithMaterialsToProto(regeneratedRoadmap)

	return &roadmapclient.RegenerateNodeIDsResponse{
		Roadmap: protoRoadmap,
	}, nil
}

func (s *RoadmapServer) convertProtoRoadmapWithMaterialsToDTO(protoRoadmap *roadmapclient.RoadmapWithMaterials) (*dto.RoadmapWithMaterialsDTO, error) {
	if protoRoadmap == nil {
		return nil, fmt.Errorf("roadmap is nil")
	}

	roadmapDTO := &dto.RoadmapWithMaterialsDTO{
		CreatedAt: protoRoadmap.CreatedAt.AsTime(),
		UpdatedAt: protoRoadmap.UpdatedAt.AsTime(),
	}

	if protoRoadmap.Id != "" {
		objectID, err := primitive.ObjectIDFromHex(protoRoadmap.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid roadmap id format: %v", err)
		}
		roadmapDTO.ID = objectID
	}

	for _, protoNode := range protoRoadmap.Nodes {
		node, err := s.convertProtoNodeWithMaterialsToDTO(protoNode)
		if err != nil {
			return nil, err
		}
		roadmapDTO.NodesWithMaterials = append(roadmapDTO.NodesWithMaterials, *node)
	}

	for _, protoEdge := range protoRoadmap.Edges {
		roadmapDTO.Edges = append(roadmapDTO.Edges, dto.EdgeDTO{
			Source: protoEdge.Source,
			Target: protoEdge.Target,
			ID:     protoEdge.Id,
		})
	}

	return roadmapDTO, nil
}

func (s *RoadmapServer) convertCreateRequestToDTO(req *roadmapclient.CreateRequest) (*dto.CreateRoadmapRequestDTO, error) {
	roadmapDTO := dto.RoadmapWithMaterialsDTO{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, protoNode := range req.Nodes {
		node, err := s.convertProtoNodeWithMaterialsToDTO(protoNode)
		if err != nil {
			return nil, err
		}
		roadmapDTO.NodesWithMaterials = append(roadmapDTO.NodesWithMaterials, *node)
	}

	for _, protoEdge := range req.Edges {
		roadmapDTO.Edges = append(roadmapDTO.Edges, dto.EdgeDTO{
			Source: protoEdge.Source,
			Target: protoEdge.Target,
			ID:     protoEdge.Id,
		})
	}

	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("invalid author id")
	}

	return &dto.CreateRoadmapRequestDTO{AuthorID: authorID, IsPublic: req.IsPublic, Roadmap: roadmapDTO}, nil
}

func (s *RoadmapServer) convertProtoNodeWithMaterialsToDTO(protoNode *roadmapclient.NodeWithMaterials) (*dto.NodeWithMaterialsDTO, error) {
	nodeID, err := uuid.Parse(protoNode.Id)
	if err != nil {
		return nil, err
	}

	nodeDTO := &dto.NodeWithMaterialsDTO{
		ID:          nodeID,
		Type:        protoNode.Type,
		Description: protoNode.Description,
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
	}

	for _, protoMaterial := range protoNode.Materials {
		materialID, err := uuid.Parse(protoMaterial.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid material id: %v", err)
		}

		authorID, err := uuid.Parse(protoMaterial.AuthorId)
		if err != nil {
			return nil, fmt.Errorf("invalid material author id: %v", err)
		}

		material := dto.Material{
			ID:        materialID,
			Name:      protoMaterial.Name,
			URL:       protoMaterial.Url,
			AuthorID:  authorID,
			CreatedAt: protoMaterial.CreatedAt.AsTime(),
			UpdatedAt: protoMaterial.UpdatedAt.AsTime(),
		}
		nodeDTO.Materials = append(nodeDTO.Materials, material)
	}

	return nodeDTO, nil
}

func (s *RoadmapServer) convertRoadmapWithMaterialsToProto(roadmap *dto.RoadmapWithMaterialsDTO) *roadmapclient.RoadmapWithMaterials {
	protoRoadmap := &roadmapclient.RoadmapWithMaterials{
		Id:        roadmap.ID.Hex(),
		CreatedAt: timestamppb.New(roadmap.CreatedAt),
		UpdatedAt: timestamppb.New(roadmap.UpdatedAt),
	}

	for _, node := range roadmap.NodesWithMaterials {
		protoRoadmap.Nodes = append(protoRoadmap.Nodes, s.convertNodeWithMaterialsToProto(node))
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

func (s *RoadmapServer) convertNodeWithMaterialsToProto(node dto.NodeWithMaterialsDTO) *roadmapclient.NodeWithMaterials {
	protoNode := &roadmapclient.NodeWithMaterials{
		Id:          node.ID.String(),
		Type:        node.Type,
		Description: node.Description,
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

	for _, material := range node.Materials {
		protoMaterial := &roadmapclient.Material{
			Id:        material.ID.String(),
			Name:      material.Name,
			Url:       material.URL,
			AuthorId:  material.AuthorID.String(),
			CreatedAt: timestamppb.New(material.CreatedAt),
			UpdatedAt: timestamppb.New(material.UpdatedAt),
		}
		protoNode.Materials = append(protoNode.Materials, protoMaterial)
	}

	return protoNode
}
