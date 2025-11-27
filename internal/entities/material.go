package entities

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID        uuid.UUID
	Name      string
	URL       string
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
