package service

import (
	"boards/internal/models"
	"context"
)

type Boards interface {
	CreateBoard(ctx context.Context, authorID string, projectID string, title string) (string, error)
	GetBoard(ctx context.Context, id string) (*models.Board, *models.TasksList, error)
	GetAllBoards(ctx context.Context, projectID string) (*models.BoardsList, error)
}
