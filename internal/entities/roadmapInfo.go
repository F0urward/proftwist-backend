package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoadmapInfo struct {
	ID                      uuid.UUID
	RoadmapID               string
	AuthorID                uuid.UUID
	CategoryID              uuid.UUID
	Name                    string
	Description             string
	IsPublic                bool
	ReferencedRoadmapInfoID *uuid.UUID
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
