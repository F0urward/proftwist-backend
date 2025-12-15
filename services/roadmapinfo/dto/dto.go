package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePrivateRoadmapInfoRequestDTO struct {
	AuthorID                string  `json:"-"`
	CategoryID              string  `json:"category_id" validate:"required,uuid4"`
	Name                    string  `json:"name" validate:"required"`
	Description             string  `json:"description"`
	IsPublic                bool    `json:"-"`
	ReferencedRoadmapInfoID *string `json:"referenced_roadmap_info_id,omitempty" validate:"omitempty,uuid4"`
}

type UpdatePrivateRoadmapInfoRequestDTO struct {
	CategoryID              *string `json:"category_id,omitempty" validate:"omitempty,uuid4"`
	Name                    *string `json:"name,omitempty" validate:"omitempty"`
	Description             *string `json:"description,omitempty"`
	ReferencedRoadmapInfoID *string `json:"referenced_roadmap_info_id,omitempty" validate:"omitempty,uuid4"`
}

type AuthorDTO struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url,omitempty"`
}

type RoadmapInfoDTO struct {
	ID                      string    `json:"id"`
	RoadmapID               string    `json:"roadmap_id"`
	Author                  AuthorDTO `json:"author"`
	CategoryID              string    `json:"category_id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	IsPublic                bool      `json:"is_public"`
	ReferencedRoadmapInfoID string    `json:"referenced_roadmap_info_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type CreatePrivateRoadmapInfoResponseDTO struct {
	RoadmapInfo RoadmapInfoDTO `json:"roadmap_info"`
}

type GetAllRoadmapsInfoResponseDTO struct {
	RoadmapsInfo []RoadmapInfoDTO `json:"roadmaps_info"`
}

type GetByIDRoadmapInfoResponseDTO struct {
	RoadmapInfo RoadmapInfoDTO `json:"roadmap_info"`
}

type GetSubscribedRoadmapsInfoResponseDTO struct {
	RoadmapsInfo []RoadmapInfoDTO `json:"roadmaps_info"`
}
