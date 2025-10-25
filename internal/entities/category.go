package entities

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CategoryID  uuid.UUID `bson:"category_id" json:"category_id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
