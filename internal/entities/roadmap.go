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
	ID          primitive.ObjectID
	Title       string
	Description string
	IsPublic    bool
	SubCount    int
	CategoryID  primitive.ObjectID
	AuthorID    uuid.UUID
	Nodes       []RoadmapNode
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoadmapNode struct {
	ID          uuid.UUID
	Title       string
	Description string
	Position    Position
	AuthorID    uuid.UUID
	Level       int
	Links       []Link
	Children    []RoadmapNode
}

type Link struct {
	URL      string
	Title    string
	Type     LinkType
	AuthorID uuid.UUID
}

type Position struct {
	X float64
	Y float64
}
