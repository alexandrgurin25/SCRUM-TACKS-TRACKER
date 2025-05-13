package board

import (
	"boards/internal/config"
	"boards/internal/models"
	"boards/internal/repository"
	"boards/internal/service"
	"boards/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Serv struct {
	repo repository.Boards
	cfg  *config.Config
}

func InitServ(repo repository.Boards, cfg *config.Config) service.Boards {
	return &Serv{repo: repo, cfg: cfg}
}

func (s *Serv) CreateBoard(ctx context.Context,
	authorID string, projectID string, title string) (string, error) {

	id := uuid.New().String()

	_, err := s.repo.SaveBoard(ctx, id, authorID, projectID, title)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to create board", zap.Error(err))
		return "", fmt.Errorf("failed to create board: %w", err)
	}

	return id, nil
}

func (s *Serv) GetBoard(ctx context.Context, id string) (*models.Board, *models.TasksList, error) {
	board, tasks, err := s.repo.GetBoard(ctx, id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get board", zap.Error(err))
		return nil, nil, fmt.Errorf("failed to get board: %w", err)
	}
	return board, tasks, nil
}

func (s *Serv) GetAllBoards(ctx context.Context, projectID string) (*models.BoardsList, error) {
	boards, err := s.repo.GetAllBoards(ctx, projectID)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get boards by project_id", zap.Error(err))
		return nil, fmt.Errorf("failed to get boards by project_id: %w", err)
	}
	return boards, nil
}
