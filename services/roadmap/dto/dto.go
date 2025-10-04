package dto

import (
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RoadmapDTO struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	IsPublic    bool               `json:"IsPublic" bson:"IsPublic"`
	SubCount    int                `json:"subCount" bson:"subCount"`
	CategoryID  primitive.ObjectID `json:"categoryId" bson:"categoryId,omitempty"`
	AuthorID    uuid.UUID          `json:"author_id" bson:"author_id"`
	Nodes       []RoadmapNodeDTO   `json:"nodes,omitempty" bson:"nodes,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type RoadmapNodeDTO struct {
	ID          uuid.UUID        `json:"id" bson:"id"`
	Title       string           `json:"title" bson:"title"`
	Description string           `json:"description,omitempty" bson:"description,omitempty"`
	Position    Position         `json:"position" bson:"position"`
	AuthorID    uuid.UUID        `json:"author_id" bson:"author_id"`
	Level       int              `json:"level" bson:"level"`
	Links       []Link           `json:"links,omitempty" bson:"links,omitempty"`
	Children    []RoadmapNodeDTO `json:"children,omitempty" bson:"children,omitempty"`
}

type Link struct {
	URL      string            `json:"url" bson:"url"`
	Title    string            `json:"title" bson:"title"`
	Type     entities.LinkType `json:"type" bson:"type"`
	AuthorID uuid.UUID         `json:"author_id" bson:"author_id"`
}

type Position struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type CreateRoadmapRequest struct {
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description"`
	IsPublic    bool               `json:"IsPublic" bson:"IsPublic"`
	SubCount    int                `json:"subCount" bson:"subCount"`
	CategoryID  primitive.ObjectID `json:"categoryId" bson:"categoryId"`
	Nodes       []RoadmapNodeDTO   `json:"nodes,omitempty"`
}

type UpdateRoadmapRequest struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	IsPublic    bool               `json:"IsPublic" bson:"IsPublic"`
	CategoryID  primitive.ObjectID `json:"categoryId" bson:"categoryId"`
	Nodes       []RoadmapNodeDTO   `json:"nodes,omitempty"`
}

type UpdatePrivacyRequest struct {
	IsPublic bool `json:"isPublic" validate:"required"`
}

type UpdatePrivacyResponse struct {
	Message   string             `json:"message"`
	IsPublic  bool               `json:"isPublic"`
	RoadmapID primitive.ObjectID `json:"roadmapId"`
}
