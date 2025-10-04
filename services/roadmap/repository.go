package roadmap

import (
	"context"
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	GetAll(context.Context) ([]*entities.Roadmap, error)
	GetByID(context.Context, primitive.ObjectID) (*entities.Roadmap, error)
	GetByAuthorID(context.Context, uuid.UUID) ([]*entities.Roadmap, error)
	Create(context.Context, *entities.Roadmap) error
	Update(context.Context, *entities.Roadmap) error
	Delete(context.Context, primitive.ObjectID) error
	SearchByTitle(context.Context, string) ([]*entities.Roadmap, error)
	UpdatePrivacy(context.Context, primitive.ObjectID, bool) error
}
