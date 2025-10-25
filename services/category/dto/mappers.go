package dto

import (
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
)

func CategoryToDTO(category *entities.Category) CategoryResponse {
	return CategoryResponse{
		CategoryID:  category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func CategoryListToDTO(categories []*entities.Category) GetAllCategoriesResponse {
	var categoryDTOs []CategoryResponse

	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, CategoryToDTO(category))
	}

	return GetAllCategoriesResponse{
		Categories: categoryDTOs,
	}
}

func CreateRequestToEntity(request *CreateCategoryRequestDTO) *entities.Category {
	return &entities.Category{
		Name:        request.Name,
		Description: request.Description,
	}
}

func UpdateRequestToEntity(existing *entities.Category, request *UpdateCategoryRequestDTO) *entities.Category {
	updated := *existing

	if request.Name != nil {
		updated.Name = *request.Name
	}

	if request.Description != nil {
		updated.Description = *request.Description
	}

	updated.UpdatedAt = time.Now()

	return &updated
}
