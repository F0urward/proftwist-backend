package roadmap

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	GetAll(context.Context) ([]*entities.Roadmap, error)
	GetByID(context.Context, primitive.ObjectID) (*entities.Roadmap, error)
	Create(context.Context, *entities.Roadmap) (*entities.Roadmap, error)
	Update(context.Context, *entities.Roadmap) (*entities.Roadmap, error)
	Delete(context.Context, primitive.ObjectID) error
}
