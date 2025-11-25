package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/material"
)

type MaterialPostgresRepository struct {
	db *sql.DB
}

func NewMaterialPostgresRepository(db *sql.DB) material.Repository {
	return &MaterialPostgresRepository{
		db: db,
	}
}

func (r *MaterialPostgresRepository) CreateMaterial(ctx context.Context, material *entities.Material) (*entities.Material, error) {
	const op = "MaterialPostgresRepository.CreateMaterial"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":              op,
		"name":            material.Name,
		"roadmap_node_id": material.RoadmapNodeID,
		"author_id":       material.AuthorID,
	})

	err := r.db.QueryRowContext(ctx, queryCreateMaterial,
		material.Name,
		material.URL,
		material.RoadmapNodeID,
		material.AuthorID,
	).Scan(
		&material.ID,
		&material.CreatedAt,
		&material.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("failed to create material")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("material_id", material.ID.String()).Info("successfully created material")
	return material, nil
}

func (r *MaterialPostgresRepository) GetMaterialByID(ctx context.Context, materialID uuid.UUID) (*entities.Material, error) {
	const op = "MaterialPostgresRepository.GetMaterialByID"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("material_id", materialID)

	var material entities.Material

	err := r.db.QueryRowContext(ctx, queryGetMaterialByID, materialID).Scan(
		&material.ID,
		&material.Name,
		&material.URL,
		&material.RoadmapNodeID,
		&material.AuthorID,
		&material.CreatedAt,
		&material.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Debug("no material found")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).Error("failed to get material by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("material_id", material.ID).Debug("material retrieved")
	return &material, nil
}

func (r *MaterialPostgresRepository) GetMaterialsByNode(ctx context.Context, nodeID string) ([]*entities.Material, error) {
	const op = "MaterialPostgresRepository.GetMaterialsByNode"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("node_id", nodeID)

	rows, err := r.db.QueryContext(ctx, queryGetMaterialsByNode, nodeID)
	if err != nil {
		logger.WithError(err).Error("failed to query materials by roadmap node")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var materials []*entities.Material
	for rows.Next() {
		var material entities.Material

		err := rows.Scan(
			&material.ID,
			&material.Name,
			&material.URL,
			&material.RoadmapNodeID,
			&material.AuthorID,
			&material.CreatedAt,
			&material.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan material row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		materials = append(materials, &material)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("materials_count", len(materials)).Debug("materials by roadmap node retrieved")
	return materials, nil
}

func (r *MaterialPostgresRepository) GetMaterialsByAuthor(ctx context.Context, authorID uuid.UUID) ([]*entities.Material, error) {
	const op = "MaterialPostgresRepository.GetMaterialsByAuthor"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("author_id", authorID)

	rows, err := r.db.QueryContext(ctx, queryGetMaterialsByAuthor, authorID)
	if err != nil {
		logger.WithError(err).Error("failed to query materials by author")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithError(err).Warn("failed to close rows")
		}
	}()

	var materials []*entities.Material
	for rows.Next() {
		var material entities.Material

		err := rows.Scan(
			&material.ID,
			&material.Name,
			&material.URL,
			&material.RoadmapNodeID,
			&material.AuthorID,
			&material.CreatedAt,
			&material.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan material row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		materials = append(materials, &material)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("error iterating rows")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("materials_count", len(materials)).Debug("materials by author retrieved")
	return materials, nil
}

func (r *MaterialPostgresRepository) DeleteMaterial(ctx context.Context, materialID uuid.UUID) error {
	const op = "MaterialPostgresRepository.DeleteMaterial"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("material_id", materialID)

	result, err := r.db.ExecContext(ctx, queryDeleteMaterial, materialID)
	if err != nil {
		logger.WithError(err).Error("failed to delete material")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("failed to get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger.Info("successfully deleted material")
	return nil
}
