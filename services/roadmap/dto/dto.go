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

type RoadmapWithMaterialsDTO struct {
	ID                 primitive.ObjectID     `json:"_id,omitempty"`
	NodesWithMaterials []NodeWithMaterialsDTO `json:"nodes,omitempty"`
	Edges              []EdgeDTO              `json:"edges,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
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

type NodeWithMaterialsDTO struct {
	ID          uuid.UUID  `json:"id"`
	Type        string     `json:"type"`
	Description string     `json:"description,omitempty"`
	Position    Position   `json:"position"`
	Data        NodeData   `json:"data"`
	Measured    Measured   `json:"measured"`
	Selected    bool       `json:"selected"`
	Dragging    bool       `json:"dragging"`
	Materials   []Material `json:"materials,omitempty"`
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

type Material struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	AuthorID  uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAllRoadmapsResponseDTO struct {
	Roadmaps []RoadmapDTO `json:"roadmaps"`
}

type GetByIDRoadmapResponseDTO struct {
	Roadmap RoadmapDTO `json:"roadmap"`
}

type GetByIDRoadmapWithMaterialsResponseDTO struct {
	RoadmapWithMaterials RoadmapWithMaterialsDTO `json:"roadmap"`
}

type CreateRoadmapRequestDTO struct {
	AuthorID uuid.UUID               `json:"author_id"`
	IsPublic bool                    `json:"is_public"`
	Roadmap  RoadmapWithMaterialsDTO `json:"roadmap"`
}

type CreateRoadmapResponseDTO struct {
	RoadmapWithMaterials RoadmapWithMaterialsDTO `json:"roadmap"`
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

type CreateMaterialRequestDTO struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EnrichedMaterialResponseDTO struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	URL       string            `json:"url"`
	Author    MaterialAuthorDTO `json:"author"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type MaterialAuthorDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url,omitempty"`
}

type MaterialListResponseDTO struct {
	Materials []EnrichedMaterialResponseDTO `json:"materials"`
	Total     int                           `json:"total"`
}

type DeleteMaterialResponseDTO struct {
	Message string `json:"message"`
}
