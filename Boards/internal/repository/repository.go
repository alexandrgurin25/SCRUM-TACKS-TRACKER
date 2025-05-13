package repository

import (
	"boards/internal/models"
	"context"
)

type Boards interface {
	SaveBoard(ctx context.Context, id string, authorID string, projectID string, title string) (*models.Board, error)
	GetBoard(ctx context.Context, id string) (*models.Board, *models.TasksList, error)
	GetAllBoards(ctx context.Context, projectID string) (*models.BoardsList, error)
}
