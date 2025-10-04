package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoadmapInfo struct {
	ID                      uuid.UUID
	OwnerID                 uuid.UUID
	CategoryID              uuid.UUID
	Name                    string
	Description             string
	IsPublic                bool
	Color                   string
	ReferencedRoadmapInfoID *uuid.UUID
	SubscriberCount         int
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
