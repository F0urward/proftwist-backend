package dto

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoadmapDTO struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Nodes     []NodeDTO          `json:"nodes,omitempty" bson:"nodes,omitempty"`
	Edges     []EdgeDTO          `json:"edges,omitempty" bson:"edges,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type NodeDTO struct {
	ID          uuid.UUID `json:"id" bson:"id"`
	Type        string    `json:"type" bson:"type"`
	Position    Position  `json:"position" bson:"position"`
	Data        NodeData  `json:"data" bson:"data"`
	Measured    Measured  `json:"measured" bson:"measured"`
	Selected    bool      `json:"selected" bson:"selected"`
	Dragging    bool      `json:"dragging" bson:"dragging"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
}

type NodeData struct {
	Label string `json:"label" bson:"label"`
	Type  string `json:"type" bson:"type"`
}

type Measured struct {
	Width  float64 `json:"width" bson:"width"`
	Height float64 `json:"height" bson:"height"`
}

type EdgeDTO struct {
	Source string `json:"source" bson:"source"`
	Target string `json:"target" bson:"target"`
	ID     string `json:"id" bson:"id"`
}

type Position struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type GetAllRoadmapsResponseDTO struct {
	Roadmaps []RoadmapDTO `json:"roadmaps"`
}

type GetByIDRoadmapResponseDTO struct {
	Roadmap RoadmapDTO `json:"roadmap"`
}

type CreateRoamapRequest struct {
	AuthorID uuid.UUID  `json:"author_id"`
	IsPublic bool       `json:"is_public"`
	Roadmap  RoadmapDTO `json:"roadmap"`
}

type UpdateRoadmapRequestDTO struct {
	Nodes []NodeDTO `json:"nodes,omitempty"`
	Edges []EdgeDTO `json:"edges,omitempty"`
}

type UpdateRoadmapResponseDTO struct {
	Roadmap RoadmapDTO `json:"roadmap"`
}

type GenerateRoadmapRequestDTO struct {
	Content    string `json:"description"`
	Complexity string `json:"complexity"`
}

type GenerateRoadmapResponseDTO struct {
	RoadmapID primitive.ObjectID `json:"roadmapId"`
}

type GenerateRoadmapDTO struct {
	Topic       string
	Description string
	Content     string
	Complexity  string
}
