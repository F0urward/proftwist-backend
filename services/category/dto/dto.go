package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetAllCategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

type CreateCategoryRequestDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryResponse struct {
	Category CategoryResponse `json:"category"`
}

type UpdateCategoryRequestDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
