package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMaterialRequestDTO struct {
	Name          string `json:"name" validate:"required"`
	URL           string `json:"url" validate:"required"`
	RoadmapNodeID string `json:"roadmap_node_id" validate:"required"`
}

type MaterialResponseDTO struct {
	ID            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`
	URL           string            `json:"url"`
	RoadmapNodeID string            `json:"roadmap_node_id"`
	Author        MaterialAuthorDTO `json:"author"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type MaterialAuthorDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url,omitempty"`
}

type MaterialListResponseDTO struct {
	Materials []MaterialResponseDTO `json:"materials"`
	Total     int                   `json:"total"`
}

type DeleteMaterialResponseDTO struct {
	Message string `json:"message"`
}
