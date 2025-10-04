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
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Nodes     []RoadmapNode      `bson:"nodes,omitempty"`
	Edges     []RoadmapEdge      `bson:"edges,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type RoadmapNode struct {
	ID       uuid.UUID `bson:"id"`
	Type     string    `bson:"type"`
	Position Position  `bson:"position"`
	Data     NodeData  `bson:"data"`
	Measured Measured  `bson:"measured"`
	Selected bool      `bson:"selected"`
	Dragging bool      `bson:"dragging"`
}

type NodeData struct {
	Label string `bson:"label"`
	Type  string `bson:"type"`
}

type Measured struct {
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
}

type RoadmapEdge struct {
	Source string `bson:"source"`
	Target string `bson:"target"`
	ID     string `bson:"id"`
}

type Position struct {
	X float64 `bson:"x"`
	Y float64 `bson:"y"`
}
