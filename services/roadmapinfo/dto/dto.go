package dto

import "time"

type CreateRoadmapInfoRequestDTO struct {
	AuthorID                string  `json:"-"`
	CategoryID              string  `json:"category_id" validate:"required,uuid4"`
	Name                    string  `json:"name" validate:"required"`
	Description             string  `json:"description"`
	IsPublic                bool    `json:"is_public"`
	ReferencedRoadmapInfoID *string `json:"referenced_roadmap_info_id,omitempty" validate:"omitempty,uuid4"`
}

type UpdateRoadmapInfoRequestDTO struct {
	CategoryID              *string `json:"category_id,omitempty" validate:"omitempty,uuid4"`
	Name                    *string `json:"name,omitempty" validate:"omitempty"`
	Description             *string `json:"description,omitempty"`
	IsPublic                *bool   `json:"is_public,omitempty"`
	ReferencedRoadmapInfoID *string `json:"referenced_roadmap_info_id,omitempty" validate:"omitempty,uuid4"`
}

type RoadmapInfoDTO struct {
	ID                      string    `json:"id"`
	RoadmapID               string    `json:"roadmap_id"`
	AuthorID                string    `json:"author_id"`
	CategoryID              string    `json:"category_id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	IsPublic                bool      `json:"is_public"`
	ReferencedRoadmapInfoID string    `json:"referenced_roadmap_info_id,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type CreateRoadmapInfoResponseDTO struct {
	RoadmapInfo RoadmapInfoDTO `json:"roadmap_info"`
}

type GetAllRoadmapsInfoResponseDTO struct {
	RoadmapsInfo []RoadmapInfoDTO `json:"roadmaps_info"`
}

type GetAllByCategoryIDRoadmapInfoResponseDTO struct {
	RoadmapsInfo []RoadmapInfoDTO `json:"roadmaps_info"`
}

type GetByIDRoadmapInfoResponseDTO struct {
	RoadmapInfo RoadmapInfoDTO `json:"roadmap_info"`
}
