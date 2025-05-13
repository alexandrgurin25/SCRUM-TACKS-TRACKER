package postgres

import (
	"comments/internal/config"
	"comments/internal/domain/models"
	"comments/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func New(ctx context.Context, cfg config.DBConfig) (*Postgres, error) {
	const op = "storage.postgres.New"

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(cfg.MaxPools)

	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		connString,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{db: db}, nil

}

func (p *Postgres) SaveComment(
	ctx context.Context,
	id string,
	authorId string,
	taskId string,
	title string,
	description string,
) (*models.Comment, error) {
	const op = "storage.postgres.SaveComment"

	comment := &models.Comment{}
	err := p.db.GetContext(
		ctx,
		comment,
		`INSERT INTO comments
		(id, author_id, task_id, title, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, author_id, task_id, title, description, created_at, updated_at`,
		id, authorId, taskId, title, description,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (p *Postgres) GetAllByTaskId(ctx context.Context, taskId string, lastTime time.Time) ([]*models.Comment, error) {
	const op = "storage.postgres.GetAllByTaskId"

	var comments []*models.Comment
	target := time.Time{}
	if lastTime == target {
		err := p.db.SelectContext(ctx, &comments, "SELECT * FROM comments WHERE task_id = $1", taskId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		err := p.db.SelectContext(ctx, &comments,
			"SELECT * FROM comments WHERE task_id = $1 AND created_at > $2",
			taskId, lastTime)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return comments, nil
}

func (p *Postgres) Close() {
	p.db.Close()
}
