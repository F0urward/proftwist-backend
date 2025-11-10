package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
)

type RoadmapInfoRepository struct {
	db *sql.DB
}

func NewRoadmapInfoRepository(db *sql.DB) roadmapinfo.Repository {
	return &RoadmapInfoRepository{db: db}
}

func (r *RoadmapInfoRepository) GetAllPublic(ctx context.Context) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetAllPublic"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetAllPublic)
	if err != nil {
		logger.WithError(err).Error("failed to query public roadmaps")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.WithError(closeErr).Warn("failed to close rows")
		}
	}()

	roadmaps := []*entities.RoadmapInfo{}

	for rows.Next() {
		roadmap := &entities.RoadmapInfo{}
		var referencedRoadmapInfoID sql.NullString

		if err = rows.Scan(
			&roadmap.ID,
			&roadmap.RoadmapID,
			&roadmap.AuthorID,
			&roadmap.CategoryID,
			&roadmap.Name,
			&roadmap.Description,
			&roadmap.IsPublic,
			&referencedRoadmapInfoID,
			&roadmap.CreatedAt,
			&roadmap.UpdatedAt,
		); err != nil {
			logger.WithError(err).Error("failed to scan roadmap row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if referencedRoadmapInfoID.Valid {
			parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
			if err != nil {
				logger.WithError(err).WithField("referenced_roadmap_id", referencedRoadmapInfoID.String).Error("invalid referenced roadmap ID in database")
				return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("invalid referenced_roadmap_id in database: %w", err))
			}
			roadmap.ReferencedRoadmapInfoID = &parsedUUID
		} else {
			roadmap.ReferencedRoadmapInfoID = nil
		}

		roadmaps = append(roadmaps, roadmap)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("roadmaps_count", len(roadmaps)).Info("successfully retrieved public roadmaps")
	return roadmaps, nil
}

func (r *RoadmapInfoRepository) GetAllPublicByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetAllPublicByCategoryID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"category_id": categoryID.String(),
	})

	rows, err := r.db.QueryContext(ctx, queryGetAllPublicByCategoryID, categoryID)
	if err != nil {
		logger.WithError(err).Error("failed to query roadmaps by category ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.WithError(closeErr).Warn("failed to close rows")
		}
	}()

	roadmaps := []*entities.RoadmapInfo{}

	for rows.Next() {
		roadmap := &entities.RoadmapInfo{}
		var referencedRoadmapInfoID sql.NullString

		if err = rows.Scan(
			&roadmap.ID,
			&roadmap.RoadmapID,
			&roadmap.AuthorID,
			&roadmap.CategoryID,
			&roadmap.Name,
			&roadmap.Description,
			&roadmap.IsPublic,
			&referencedRoadmapInfoID,
			&roadmap.CreatedAt,
			&roadmap.UpdatedAt,
		); err != nil {
			logger.WithError(err).Error("failed to scan roadmap row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if referencedRoadmapInfoID.Valid {
			parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
			if err != nil {
				logger.WithError(err).WithField("referenced_roadmap_id", referencedRoadmapInfoID.String).Error("invalid referenced roadmap ID in database")
				return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("invalid referenced_roadmap_id in database: %w", err))
			}
			roadmap.ReferencedRoadmapInfoID = &parsedUUID
		} else {
			roadmap.ReferencedRoadmapInfoID = nil
		}

		roadmaps = append(roadmaps, roadmap)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("roadmaps_count", len(roadmaps)).Info("successfully retrieved roadmaps by category")
	return roadmaps, nil
}

func (r *RoadmapInfoRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetAllByUserID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	rows, err := r.db.QueryContext(ctx, queryGetAllByUserID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to query roadmaps by user ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.WithError(closeErr).Warn("failed to close rows")
		}
	}()

	roadmaps := []*entities.RoadmapInfo{}

	for rows.Next() {
		roadmap := &entities.RoadmapInfo{}
		var referencedRoadmapInfoID sql.NullString

		if err = rows.Scan(
			&roadmap.ID,
			&roadmap.RoadmapID,
			&roadmap.AuthorID,
			&roadmap.CategoryID,
			&roadmap.Name,
			&roadmap.Description,
			&roadmap.IsPublic,
			&referencedRoadmapInfoID,
			&roadmap.CreatedAt,
			&roadmap.UpdatedAt,
		); err != nil {
			logger.WithError(err).Error("failed to scan roadmap row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if referencedRoadmapInfoID.Valid {
			parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
			if err != nil {
				logger.WithError(err).WithField("referenced_roadmap_id", referencedRoadmapInfoID.String).Error("invalid referenced roadmap ID in database")
				return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("invalid referenced_roadmap_id in database: %w", err))
			}
			roadmap.ReferencedRoadmapInfoID = &parsedUUID
		} else {
			roadmap.ReferencedRoadmapInfoID = nil
		}

		roadmaps = append(roadmaps, roadmap)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("roadmaps_count", len(roadmaps)).Info("successfully retrieved roadmaps by user ID")
	return roadmaps, nil
}

func (r *RoadmapInfoRepository) GetByID(ctx context.Context, roadmapID uuid.UUID) (*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.String(),
	})

	roadmap := &entities.RoadmapInfo{}
	var referencedRoadmapInfoID sql.NullString

	err := r.db.QueryRowContext(ctx, queryGetByID, roadmapID).Scan(
		&roadmap.ID,
		&roadmap.RoadmapID,
		&roadmap.AuthorID,
		&roadmap.CategoryID,
		&roadmap.Name,
		&roadmap.Description,
		&roadmap.IsPublic,
		&referencedRoadmapInfoID,
		&roadmap.CreatedAt,
		&roadmap.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Info("roadmap not found")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if referencedRoadmapInfoID.Valid {
		parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
		if err != nil {
			logger.WithError(err).WithField("referenced_roadmap_id", referencedRoadmapInfoID.String).Error("invalid referenced roadmap ID in database")
			return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("invalid referenced_roadmap_id in database: %w", err))
		}
		roadmap.ReferencedRoadmapInfoID = &parsedUUID
	} else {
		roadmap.ReferencedRoadmapInfoID = nil
	}

	logger.Info("successfully retrieved roadmap")
	return roadmap, nil
}

func (r *RoadmapInfoRepository) GetByRoadmapID(ctx context.Context, roadmapID string) (*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetByRoadmapID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID,
	})

	roadmap := &entities.RoadmapInfo{}
	var referencedRoadmapInfoID sql.NullString

	err := r.db.QueryRowContext(ctx, queryGetByRoadmapID, roadmapID).Scan(
		&roadmap.ID,
		&roadmap.RoadmapID,
		&roadmap.AuthorID,
		&roadmap.CategoryID,
		&roadmap.Name,
		&roadmap.Description,
		&roadmap.IsPublic,
		&referencedRoadmapInfoID,
		&roadmap.CreatedAt,
		&roadmap.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Debug("roadmap not found by roadmap ID")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).Error("failed to get roadmap by roadmap ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if referencedRoadmapInfoID.Valid {
		parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
		if err != nil {
			logger.WithError(err).WithField("referenced_roadmap_id", referencedRoadmapInfoID.String).Error("invalid referenced roadmap ID in database")
			return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("invalid referenced_roadmap_id in database: %w", err))
		}
		roadmap.ReferencedRoadmapInfoID = &parsedUUID
	} else {
		roadmap.ReferencedRoadmapInfoID = nil
	}

	logger.Debug("successfully retrieved roadmap by roadmap ID")
	return roadmap, nil
}

func (r *RoadmapInfoRepository) Create(ctx context.Context, roadmap *entities.RoadmapInfo) (*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.Create"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"author_id": roadmap.AuthorID.String(),
	})

	var refRoadmapInfoID interface{}
	if roadmap.ReferencedRoadmapInfoID != nil {
		refRoadmapInfoID = *roadmap.ReferencedRoadmapInfoID
	} else {
		refRoadmapInfoID = nil
	}

	createdRoadmap := &entities.RoadmapInfo{}
	var referencedRoadmapInfoID sql.NullString

	err := r.db.QueryRowContext(ctx, queryCreate,
		roadmap.AuthorID,
		roadmap.CategoryID,
		roadmap.Name,
		roadmap.Description,
		roadmap.IsPublic,
		refRoadmapInfoID,
		roadmap.RoadmapID,
	).Scan(
		&createdRoadmap.ID,
		&createdRoadmap.RoadmapID,
		&createdRoadmap.AuthorID,
		&createdRoadmap.CategoryID,
		&createdRoadmap.Name,
		&createdRoadmap.Description,
		&createdRoadmap.IsPublic,
		&referencedRoadmapInfoID,
		&createdRoadmap.CreatedAt,
		&createdRoadmap.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to create roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if referencedRoadmapInfoID.Valid {
		parsedUUID, err := uuid.Parse(referencedRoadmapInfoID.String)
		if err != nil {
			logger.WithError(err).WithField("referenced_roadmap_id", referencedRoadmapInfoID.String).Error("invalid referenced roadmap ID in database")
			return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("invalid referenced_roadmap_id in database: %w", err))
		}
		createdRoadmap.ReferencedRoadmapInfoID = &parsedUUID
	} else {
		createdRoadmap.ReferencedRoadmapInfoID = nil
	}

	logger.Info("successfully created roadmap")
	return createdRoadmap, nil
}

func (r *RoadmapInfoRepository) Update(ctx context.Context, roadmap *entities.RoadmapInfo) error {
	const op = "RoadmapInfoRepository.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmap.ID.String(),
	})

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
		refRoadmapInfoID,
		roadmap.RoadmapID,
		roadmap.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to update roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("roadmap not found for update")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("roadmap not found"))
	}

	logger.Info("successfully updated roadmap")
	return nil
}

func (r *RoadmapInfoRepository) Delete(ctx context.Context, roadmapID uuid.UUID) error {
	const op = "RoadmapInfoRepository.Delete"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"roadmap_id": roadmapID.String(),
	})

	result, err := r.db.ExecContext(ctx, queryDelete, roadmapID)
	if err != nil {
		logger.WithError(err).Error("failed to delete roadmap")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("roadmap not found for deletion")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("roadmap not found"))
	}

	logger.Info("successfully deleted roadmap")
	return nil
}
