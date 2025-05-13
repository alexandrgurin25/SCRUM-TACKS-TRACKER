package board

import (
	"boards/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func InitRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) SaveBoard(ctx context.Context, id string, authorID string, projectID string, title string) (*models.Board, error) {
	var board models.Board
	err := r.pool.QueryRow(ctx,
		`INSERT INTO boards
		(id, author_id, project_id, title)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, author_id, project_id, title, created_at, updated_at`,
		id, authorID, projectID, title).Scan(
		&board.ID,
		&board.AuthorID,
		&board.ProjectID,
		&board.Title,
		&board.CreatedAt,
		&board.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create board: %v", err)
	}

	return &board, nil
}

func (r *Repo) GetBoard(ctx context.Context, id string) (*models.Board, *models.TasksList, error) {
	var board models.Board
	err := r.pool.QueryRow(ctx,
		`SELECT id, author_id, project_id, title, created_at, updated_at 
         FROM boards WHERE id = $1`, id).Scan(
		&board.ID,
		&board.AuthorID,
		&board.ProjectID,
		&board.Title,
		&board.CreatedAt,
		&board.UpdatedAt,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to get board: %v", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, title, description, board_id, author_id, status, created_at, updated_at, deadline
         FROM tasks WHERE board_id = $1`, id)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get tasks: %v", err)
	}
	defer rows.Close()

	tasksList := &models.TasksList{
		BoardID: board.ID,
		Tasks:   make([]models.Task, 0),
	}

	for rows.Next() {
		var task models.Task
		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.BoardID,
			&task.AuthorID,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Deadline,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to scan task: %v", err)
		}
		tasksList.Tasks = append(tasksList.Tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error after iterating tasks: %v", err)
	}

	return &board, tasksList, nil
}

func (r *Repo) GetAllBoards(ctx context.Context, projectID string) (*models.BoardsList, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, author_id, project_id, title, created_at, updated_at 
         FROM boards WHERE project_id = $1`, projectID)

	if err != nil {
		return nil, fmt.Errorf("unable to get boards: %v", err)
	}
	defer rows.Close()

	var boardList models.BoardsList
	for rows.Next() {
		var board models.Board
		err = rows.Scan(
			&board.ID,
			&board.AuthorID,
			&board.ProjectID,
			&board.Title,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan board: %v", err)
		}
		boardList.Boards = append(boardList.Boards, board)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating boards: %v", err)
	}

	boardList.ProjectID = projectID

	return &boardList, nil
}
