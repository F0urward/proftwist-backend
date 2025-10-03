package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoadmapInfo struct {
	ID                      uuid.UUID  `json:"id"`
	OwnerID                 uuid.UUID  `json:"owner_id"`
	CategoryID              uuid.UUID  `json:"category_id"`
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	IsPublic                bool       `json:"is_public"`
	Color                   string     `json:"color"`
	ReferencedRoadmapInfoID *uuid.UUID `json:"referenced_roadmap_info_id,omitempty"`
	SubscriberCount         int        `json:"subscriber_count"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}
