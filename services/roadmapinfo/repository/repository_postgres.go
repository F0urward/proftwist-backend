package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
)

type RoadmapInfoPostgresRepository struct {
	db *sql.DB
}

func NewRoadmapInfoPostgresRepository(db *sql.DB) roadmapinfo.Repository {
	return &RoadmapInfoPostgresRepository{db: db}
}

func (r *RoadmapInfoPostgresRepository) GetAllPublic(ctx context.Context) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetAllPublic"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

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

func (r *RoadmapInfoPostgresRepository) GetAllPublicByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetAllPublicByCategoryID"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *RoadmapInfoPostgresRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetAllByUserID"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *RoadmapInfoPostgresRepository) GetByID(ctx context.Context, roadmapID uuid.UUID) (*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetByID"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *RoadmapInfoPostgresRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetByIDs"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":        op,
		"ids_count": len(ids),
	})

	if len(ids) == 0 {
		logger.Info("no IDs provided")
		return []*entities.RoadmapInfo{}, nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	rows, err := r.db.QueryContext(ctx, queryGetByIDs, pq.Array(idStrings))
	if err != nil {
		logger.WithError(err).Error("failed to query roadmaps by IDs")
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

	logger.WithField("roadmaps_count", len(roadmaps)).Info("successfully retrieved roadmaps by IDs")
	return roadmaps, nil
}

func (r *RoadmapInfoPostgresRepository) GetByRoadmapID(ctx context.Context, roadmapID string) (*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.GetByRoadmapID"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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
		logger.Info("roadmap not found by roadmap ID")
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

	logger.Info("successfully retrieved roadmap by roadmap ID")
	return roadmap, nil
}

func (r *RoadmapInfoPostgresRepository) Create(ctx context.Context, roadmap *entities.RoadmapInfo) (*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.Create"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *RoadmapInfoPostgresRepository) Update(ctx context.Context, roadmap *entities.RoadmapInfo) error {
	const op = "RoadmapInfoRepository.Update"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *RoadmapInfoPostgresRepository) Delete(ctx context.Context, roadmapID uuid.UUID) error {
	const op = "RoadmapInfoRepository.Delete"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
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

func (r *RoadmapInfoPostgresRepository) CreateSubscription(ctx context.Context, userID, roadmapInfoID uuid.UUID) error {
	const op = "RoadmapInfoRepository.CreateSubscription"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"user_id":         userID.String(),
		"roadmap_info_id": roadmapInfoID.String(),
	})

	result, err := r.db.ExecContext(ctx, queryCreateSubscription, userID, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to create subscription")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("subscription already exists or failed to create")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("subscription already exists"))
	}

	logger.Info("successfully created subscription")
	return nil
}

func (r *RoadmapInfoPostgresRepository) DeleteSubscription(ctx context.Context, userID, roadmapInfoID uuid.UUID) error {
	const op = "RoadmapInfoRepository.DeleteSubscription"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"user_id":         userID.String(),
		"roadmap_info_id": roadmapInfoID.String(),
	})

	result, err := r.db.ExecContext(ctx, queryDeleteSubscription, userID, roadmapInfoID)
	if err != nil {
		logger.WithError(err).Error("failed to delete subscription")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		logger.Warn("subscription not found for deletion")
		return fmt.Errorf("%s: %w", op, fmt.Errorf("subscription not found"))
	}

	logger.Info("successfully deleted subscription")
	return nil
}

func (r *RoadmapInfoPostgresRepository) SubscriptionExists(ctx context.Context, userID, roadmapInfoID uuid.UUID) (bool, error) {
	const op = "RoadmapInfoRepository.SubscriptionExists"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"user_id":         userID.String(),
		"roadmap_info_id": roadmapInfoID.String(),
	})

	var exists bool
	err := r.db.QueryRowContext(ctx, querySubscriptionExists, userID, roadmapInfoID).Scan(&exists)
	if err != nil {
		logger.WithError(err).Error("failed to check subscription existence")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("exists", exists).Info("checked subscription existence")
	return exists, nil
}

func (r *RoadmapInfoPostgresRepository) GetSubscribedRoadmapIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	const op = "RoadmapInfoRepository.GetSubscribedRoadmapIDs"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	rows, err := r.db.QueryContext(ctx, queryGetSubscribedRoadmapIDs, userID)
	if err != nil {
		logger.WithError(err).Error("failed to query subscribed roadmap IDs")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.WithError(closeErr).Warn("failed to close rows")
		}
	}()

	var roadmapIDs []uuid.UUID

	for rows.Next() {
		var roadmapID uuid.UUID
		if err = rows.Scan(&roadmapID); err != nil {
			logger.WithError(err).Error("failed to scan roadmap ID")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		roadmapIDs = append(roadmapIDs, roadmapID)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("roadmap_ids_count", len(roadmapIDs)).Info("successfully retrieved subscribed roadmap IDs")
	return roadmapIDs, nil
}

func (r *RoadmapInfoPostgresRepository) SearchPublic(ctx context.Context, query string, categoryID *uuid.UUID) ([]*entities.RoadmapInfo, error) {
	const op = "RoadmapInfoRepository.SearchPublic"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":          op,
		"query":       query,
		"category_id": categoryID,
	})

	var rows *sql.Rows
	var err error

	if categoryID != nil {
		rows, err = r.db.QueryContext(ctx, querySearchPublicWithCategory, query, categoryID)
	} else {
		rows, err = r.db.QueryContext(ctx, querySearchPublic, query)
	}

	if err != nil {
		logger.WithError(err).Error("failed to query public roadmaps for search")
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

	logger.WithField("roadmaps_count", len(roadmaps)).Info("successfully searched public roadmaps")
	return roadmaps, nil
}
