package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/google/uuid"
)

type roadmapRepository struct {
	db *sql.DB
}

func NewRoadmapRepository(db *sql.DB) roadmap.Repository {
	return &roadmapRepository{db: db}
}

func (r *roadmapRepository) GetAll(ctx context.Context) ([]*entities.Roadmap, error) {
	roadmaps := []*entities.Roadmap{}

	rows, err := r.db.QueryContext(ctx, queryGetAllRoadmap)
	if err != nil {
		log.Printf("Failed to get rows: %v", err)
		return nil, fmt.Errorf("failed to get rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		roadmap := &entities.Roadmap{}
		var referencedRoadmapID sql.NullString

		if err := rows.Scan(
			&roadmap.ID,
			&roadmap.OwnerID,
			&roadmap.CategoryID,
			&roadmap.Name,
			&roadmap.Description,
			&roadmap.IsPublic,
			&roadmap.Color,
			&referencedRoadmapID,
			&roadmap.SubscriberCount,
			&roadmap.CreatedAt,
			&roadmap.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		if referencedRoadmapID.Valid {
			parsedUUID, err := uuid.Parse(referencedRoadmapID.String)
			if err != nil {
				return nil, fmt.Errorf("invalid referenced_roadmap_id in database: %v", err)
			}
			roadmap.ReferencedRoadmapID = &parsedUUID
		} else {
			roadmap.ReferencedRoadmapID = nil
		}

		roadmaps = append(roadmaps, roadmap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return roadmaps, nil
}

func (r *roadmapRepository) GetByID(ctx context.Context, roadmapID uuid.UUID) (*entities.Roadmap, error) {
	roadmap := &entities.Roadmap{}
	var referencedRoadmapID sql.NullString

	err := r.db.QueryRowContext(ctx, queryGetByID, roadmapID).Scan(
		&roadmap.ID,
		&roadmap.OwnerID,
		&roadmap.CategoryID,
		&roadmap.Name,
		&roadmap.Description,
		&roadmap.IsPublic,
		&roadmap.Color,
		&referencedRoadmapID,
		&roadmap.SubscriberCount,
		&roadmap.CreatedAt,
		&roadmap.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Failed to get roadmap by ID %s: %v", roadmapID, err)
		return nil, fmt.Errorf("failed to get roadmap: %v", err)
	}

	if referencedRoadmapID.Valid {
		parsedUUID, err := uuid.Parse(referencedRoadmapID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid referenced_roadmap_id in database: %v", err)
		}
		roadmap.ReferencedRoadmapID = &parsedUUID
	} else {
		roadmap.ReferencedRoadmapID = nil
	}

	return roadmap, nil
}

func (r *roadmapRepository) Create(ctx context.Context, roadmap *entities.Roadmap) error {
	var refRoadmapID interface{}
	if roadmap.ReferencedRoadmapID != nil {
		refRoadmapID = *roadmap.ReferencedRoadmapID
	} else {
		refRoadmapID = nil
	}

	_, err := r.db.ExecContext(ctx, queryCreate,
		roadmap.ID,
		roadmap.OwnerID,
		roadmap.CategoryID,
		roadmap.Name,
		roadmap.Description,
		roadmap.IsPublic,
		roadmap.Color,
		refRoadmapID,
		roadmap.SubscriberCount,
		roadmap.CreatedAt,
		roadmap.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to create roadmap: %v", err)
		return fmt.Errorf("failed to create roadmap: %v", err)
	}
	return nil
}

func (r *roadmapRepository) Update(ctx context.Context, roadmap *entities.Roadmap) error {
	roadmap.UpdatedAt = time.Now()

	var refRoadmapID interface{}
	if roadmap.ReferencedRoadmapID != nil {
		refRoadmapID = *roadmap.ReferencedRoadmapID
	} else {
		refRoadmapID = nil
	}

	result, err := r.db.ExecContext(ctx, queryUpdate,
		roadmap.ID,
		roadmap.CategoryID,
		roadmap.Name,
		roadmap.Description,
		roadmap.IsPublic,
		roadmap.Color,
		refRoadmapID,
		roadmap.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to update roadmap with ID %s: %v", roadmap.ID, err)
		return fmt.Errorf("failed to update roadmap: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("roadmap not found")
	}

	return nil
}

func (r *roadmapRepository) Delete(ctx context.Context, roadmapID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, queryDelete, roadmapID)
	if err != nil {
		log.Printf("Failed to delete roadmap with ID %s: %v", roadmapID, err)
		return fmt.Errorf("failed to delete roadmap: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("roadmap not found")
	}
	return nil
}
