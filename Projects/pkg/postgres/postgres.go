package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"projects/internal/config"
	"strings"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.MaxConns,
		cfg.MinConns,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {

		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		strings.Split(connString, "&")[0],
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create migrations: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}

	return pool, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
