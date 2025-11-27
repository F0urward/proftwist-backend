package entities

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkType string

const (
	LinkTypeChat   LinkType = "chat"
	LinkTypeSource LinkType = "source"
)

type Roadmap struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	Nodes     []RoadmapNode      `json:"nodes,omitempty"`
	Edges     []RoadmapEdge      `json:"edges,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type RoadmapNode struct {
	ID          uuid.UUID  `json:"id"`
	Type        string     `json:"type"`
	Position    Position   `json:"position"`
	Data        NodeData   `json:"data"`
	Measured    Measured   `json:"measured"`
	Selected    bool       `json:"selected"`
	Dragging    bool       `json:"dragging"`
	Description string     `json:"description"`
	Materials   []Material `bson:"materials"`
}

type NodeData struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

type Measured struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type RoadmapEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	ID     string `json:"id"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
