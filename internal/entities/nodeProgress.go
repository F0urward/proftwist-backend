package entities

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeProgressStatus string

const (
	NodeProgressPending    NodeProgressStatus = "ожидает"
	NodeProgressInProgress NodeProgressStatus = "в процессе"
	NodeProgressDone       NodeProgressStatus = "завершено"
	NodeProgressSkipped    NodeProgressStatus = "пропущено"
)

type NodeProgress struct {
	Status NodeProgressStatus `json:"status" bson:"status"`
}

type UserProgress struct {
	ID        primitive.ObjectID         `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    uuid.UUID                  `json:"user_id" bson:"user_id"`
	RoadmapID primitive.ObjectID         `json:"roadmap_id" bson:"roadmap_id"`
	Progress  map[uuid.UUID]NodeProgress `json:"progress" bson:"progress"`
	CreatedAt time.Time                  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at" bson:"updated_at"`
}
