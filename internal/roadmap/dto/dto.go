package dto

import "time"

type CreateRoadmapRequestDTO struct {
	OwnerID             string  `json:"owner_id" validate:"required,uuid4"`
	CategoryID          string  `json:"category_id" validate:"required,uuid4"`
	Name                string  `json:"name" validate:"required,max=255"`
	Description         string  `json:"description"`
	IsPublic            bool    `json:"is_public"`
	Color               string  `json:"color,omitempty"`
	ReferencedRoadmapID *string `json:"referenced_roadmap_id,omitempty" validate:"omitempty,uuid4"`
}

type UpdateRoadmapRequestDTO struct {
	ID                  string  `json:"id" validate:"required,uuid4"`
	OwnerID             string  `json:"owner_id" validate:"required,uuid4"`
	CategoryID          string  `json:"category_id,omitempty" validate:"omitempty,uuid4"`
	Name                string  `json:"name,omitempty" validate:"omitempty,max=255"`
	Description         string  `json:"description,omitempty"`
	IsPublic            bool    `json:"is_public,omitempty"`
	Color               string  `json:"color,omitempty"`
	ReferencedRoadmapID *string `json:"referenced_roadmap_id,omitempty" validate:"omitempty,uuid4"`
}

type RoadmapResponseDTO struct {
	ID                  string    `json:"id"`
	OwnerID             string    `json:"owner_id"`
	CategoryID          string    `json:"category_id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	IsPublic            bool      `json:"is_public"`
	Color               string    `json:"color"`
	ReferencedRoadmapID string    `json:"referenced_roadmap_id,omitempty"`
	SubscriberCount     int       `json:"subscriber_count"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type GetAllRoadmapsResponseDTO struct {
	Roadmaps []RoadmapResponseDTO `json:"roadmaps"`
}

type GetByIDRoadmapResponseDTO struct {
	Roadmap RoadmapResponseDTO `json:"roadmap"`
}
