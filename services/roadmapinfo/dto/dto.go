package dto

import "time"

type CreateRoadmapInfoRequestDTO struct {
	OwnerID                 string  `json:"owner_id" validate:"required,uuid4"`
	CategoryID              string  `json:"category_id" validate:"required,uuid4"`
	Name                    string  `json:"name" validate:"required,max=255"`
	Description             string  `json:"description"`
	IsPublic                bool    `json:"is_public"`
	ReferencedRoadmapInfoID *string `json:"referenced_roadmap_info_id,omitempty" validate:"omitempty,uuid4"`
}

type UpdateRoadmapInfoRequestDTO struct {
	CategoryID              *string `json:"category_id,omitempty" validate:"omitempty,uuid4"`
	Name                    *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description             *string `json:"description,omitempty"`
	IsPublic                *bool   `json:"is_public,omitempty"`
	ReferencedRoadmapInfoID *string `json:"referenced_roadmap_info_id,omitempty" validate:"omitempty,uuid4"`
}

type RoadmapInfoResponseDTO struct {
	ID                      string    `json:"id"`
	OwnerID                 string    `json:"owner_id"`
	CategoryID              string    `json:"category_id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	IsPublic                bool      `json:"is_public"`
	ReferencedRoadmapInfoID string    `json:"referenced_roadmap_info_id,omitempty"`
	SubscriberCount         int       `json:"subscriber_count"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type GetAllRoadmapsInfoResponseDTO struct {
	RoadmapsInfo []RoadmapInfoResponseDTO `json:"roadmaps_info"`
}

type GetByIDRoadmapInfoResponseDTO struct {
	RoadmapInfo RoadmapInfoResponseDTO `json:"roadmap_info"`
}
