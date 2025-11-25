package entities

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID            uuid.UUID
	Name          string
	URL           string
	RoadmapNodeID string
	AuthorID      uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
