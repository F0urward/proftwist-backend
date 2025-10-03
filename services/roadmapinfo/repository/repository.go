package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
	"github.com/google/uuid"
)

type RoadmapInfoRepository struct {
	db *sql.DB
}

func NewRoadmapInfoRepository(db *sql.DB) roadmapinfo.Repository {
	return &RoadmapInfoRepository{db: db}
}

func (r *RoadmapInfoRepository) GetAll(ctx context.Context) ([]*entities.RoadmapInfo, error) {
	roadmaps := []*entities.RoadmapInfo{}

	rows, err := r.db.QueryContext(ctx, queryGetAll)
	if err != nil {
		log.Printf("Failed to get rows: %v", err)
		return nil, fmt.Errorf("failed to get rows: %v", err)
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Printf("Failed to close rows: %v", closeErr)
		}
	}()

	for rows.Next() {
		roadmap := &entities.RoadmapInfo{}
		var referencedRoadmapInfoID sql.NullString

		if err := rows.Scan(
			&roadmap.ID,
			&roadmap.OwnerID,
			&roadmap.CategoryID,
			&roadmap.Name,
			&roadmap.Description,
			&roadmap.IsPublic,
			&roadmap.Color,
			&referencedRoadmapInfoID,
			&roadmap.SubscriberCount,
			&roadmap.CreatedAt,
			&roadmap.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		if referencedRoadmapInfoID.Valid {
			parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
			if err != nil {
				return nil, fmt.Errorf("invalid referenced_roadmap_id in database: %v", err)
			}
			roadmap.ReferencedRoadmapInfoID = &parsedUUID
		} else {
			roadmap.ReferencedRoadmapInfoID = nil
		}

		roadmaps = append(roadmaps, roadmap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return roadmaps, nil
}

func (r *RoadmapInfoRepository) GetByID(ctx context.Context, roadmapID uuid.UUID) (*entities.RoadmapInfo, error) {
	roadmap := &entities.RoadmapInfo{}
	var referencedRoadmapInfoID sql.NullString

	err := r.db.QueryRowContext(ctx, queryGetByID, roadmapID).Scan(
		&roadmap.ID,
		&roadmap.OwnerID,
		&roadmap.CategoryID,
		&roadmap.Name,
		&roadmap.Description,
		&roadmap.IsPublic,
		&roadmap.Color,
		&referencedRoadmapInfoID,
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

	if referencedRoadmapInfoID.Valid {
		parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid referenced_roadmap_id in database: %v", err)
		}
		roadmap.ReferencedRoadmapInfoID = &parsedUUID
	} else {
		roadmap.ReferencedRoadmapInfoID = nil
	}

	return roadmap, nil
}

func (r *RoadmapInfoRepository) Create(ctx context.Context, roadmap *entities.RoadmapInfo) error {
	var refRoadmapInfoID interface{}
	if roadmap.ReferencedRoadmapInfoID != nil {
		refRoadmapInfoID = *roadmap.ReferencedRoadmapInfoID
	} else {
		refRoadmapInfoID = nil
	}

	_, err := r.db.ExecContext(ctx, queryCreate,
		roadmap.OwnerID,
		roadmap.CategoryID,
		roadmap.Name,
		roadmap.Description,
		roadmap.IsPublic,
		roadmap.Color,
		refRoadmapInfoID,
		roadmap.SubscriberCount,
	)
	if err != nil {
		log.Printf("Failed to create roadmap: %v", err)
		return fmt.Errorf("failed to create roadmap: %v", err)
	}
	return nil
}

func (r *RoadmapInfoRepository) Update(ctx context.Context, roadmap *entities.RoadmapInfo) error {
	roadmap.UpdatedAt = time.Now()

	var refRoadmapInfoID interface{}
	if roadmap.ReferencedRoadmapInfoID != nil {
		refRoadmapInfoID = *roadmap.ReferencedRoadmapInfoID
	} else {
		refRoadmapInfoID = nil
	}

	result, err := r.db.ExecContext(ctx, queryUpdate,
		roadmap.ID,
		roadmap.CategoryID,
		roadmap.Name,
		roadmap.Description,
		roadmap.IsPublic,
		roadmap.Color,
		refRoadmapInfoID,
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

func (r *RoadmapInfoRepository) Delete(ctx context.Context, roadmapID uuid.UUID) error {
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
