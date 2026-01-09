package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"

	_ "github.com/lib/pq"
)

func NewDatabase(cfg *config.Config) *sql.DB {
	const op = "postgres.NewDatabase"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)

	db, err := sql.Open(cfg.Postgres.Driver, dsn)
	if err != nil {
		logger.WithError(err).Error("failed to connect to postgres")
	}

	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Postgres.ConnMaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(cfg.Postgres.ConnMaxIdleTime) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		logger.WithError(err).Error("cannot ping postgres instance")
	}

	return db
}
