package projects

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"projects/internal/models"
)

type Repo struct {
	pool *pgxpool.Pool
}

func InitRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) SaveProject(ctx context.Context, id string, authorID string, description string, title string) (*models.Project, error) {
	var project models.Project
	err := r.pool.QueryRow(ctx,
		`INSERT INTO projects
		(id, author_id, description, title)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, author_id, description, title, created_at, updated_at`,
		id, authorID, description, title).Scan(
		&project.ID,
		&project.AuthorID,
		&project.Description,
		&project.Title,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create project: %v", err)
	}

	return &project, nil
}

func (r *Repo) GetProject(ctx context.Context, id string) (*models.Project, *models.ParticipantsList, error) {
	var project models.Project
	err := r.pool.QueryRow(ctx,
		`SELECT id, author_id, description, title, created_at, updated_at 
		FROM projects WHERE id = $1`, id).Scan(
		&project.ID,
		&project.AuthorID,
		&project.Description,
		&project.Title,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get project: %v", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT user_id FROM users_projects WHERE project_id = $1`, id)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get project participants: %v", err)
	}
	defer rows.Close()

	participants := &models.ParticipantsList{
		ProjectID: id,
		UsersID:   make([]string, 0),
	}

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to scan participant: %v", err)
		}
		participants.UsersID = append(participants.UsersID, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error after iterating participants: %v", err)
	}

	return &project, participants, nil
}

func (r *Repo) AddParticipant(ctx context.Context, projectId string, participantEmail string) error {
	var projectExists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)`, projectId).Scan(&projectExists)
	if err != nil {
		return fmt.Errorf("failed to check project existence: %v", err)
	}
	if !projectExists {
		return fmt.Errorf("project with id %s does not exist", projectId)
	}

	var userID string
	err = r.pool.QueryRow(ctx, `SELECT uuid FROM users WHERE email = $1`, participantEmail).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %v", err)
	}

	var alreadyParticipant bool
	err = r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users_projects WHERE user_id = $1 AND project_id = $2)`,
		userID, projectId).Scan(&alreadyParticipant)
	if err != nil {
		return fmt.Errorf("failed to check participant existence: %v", err)
	}
	if alreadyParticipant {
		return fmt.Errorf("user %s is already a participant of project %s", participantEmail, projectId)
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO users_projects (user_id, project_id) VALUES ($1, $2)`,
		userID, projectId)
	if err != nil {
		return fmt.Errorf("failed to add participant: %v", err)
	}

	return nil
}

func (r *Repo) DeleteProject(ctx context.Context, id string, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users_projects WHERE project_id = $1`, id)
	if err != nil {
		return fmt.Errorf("unable to delete project participants: %v", err)
	}

	_, err = r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("unable to delete project: %v", err)
	}

	return nil
}
