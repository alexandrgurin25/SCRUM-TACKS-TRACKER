package service

import (
	"comments/internal/domain/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	SaveComment(
		ctx context.Context,
		id string,
		authorId string,
		taskId string,
		title string,
		description string,
	) (*models.Comment, error)
	GetAllByTaskId(ctx context.Context, taskId string, lastTime time.Time) ([]*models.Comment, error)
}

type Cash interface {
	SaveComment(ctx context.Context, comment *models.Comment) error
	SaveComments(ctx context.Context, comments []*models.Comment) error
	GetAllByTaskId(ctx context.Context, taskId string, lastTime time.Time) ([]*models.Comment, error)
}

type Service struct {
	storage Storage
	cash    Cash
}

func New(storage Storage, cash Cash) *Service {
	return &Service{
		storage: storage,
		cash:    cash,
	}
}

func (s *Service) CreateComment(
	ctx context.Context,
	authorId string,
	taskId string,
	title string,
	description string,
) (string, error) {
	const op = "service.CreateComment"

	id := uuid.New().String()

	comment, err := s.storage.SaveComment(
		ctx,
		id,
		authorId,
		taskId,
		title,
		description,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.cash.SaveComment(ctx, comment)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return comment.ID, nil
}

func (s *Service) GetAllByTaskId(ctx context.Context, taskId string, lastTime time.Time) ([]*models.Comment, error) {
	const op = "service.GetAllByTaskId"

	// ! can't load from cash
	comments, err := s.cash.GetAllByTaskId(ctx, taskId, lastTime)
	if err == nil && comments != nil {
		fmt.Println("from cash")
		return comments, nil
	}
	// fmt.Println(err, comments)

	comments, err = s.storage.GetAllByTaskId(ctx, taskId, lastTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(comments) > 0 {
		fmt.Println("from db")
	}

	err = s.cash.SaveComments(ctx, comments)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}
