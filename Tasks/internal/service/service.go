package service

import (
	"context"
	"fmt"
	"task/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	SaveTask(
		ctx context.Context,
		id string,
		authorID string,
		boardID string,
		title string,
		description string,
		deadline time.Time,
	) (models.Task, error)
	GetTasksByBoardID(ctx context.Context, boardID string) ([]models.Task, error)
	UpdateTask(
		ctx context.Context,
		id string,
		title string,
		description string,
		deadline time.Time,
	) (models.Task, error)
	ChangeBoard(ctx context.Context, id, boardID string, c chan error)
	DeleteTask(ctx context.Context, id string, c chan error)
}

type Cash interface {
	SaveTask(ctx context.Context, task models.Task) error
	SaveTasks(ctx context.Context, tasks []models.Task) error
	GetTasksByBoardID(ctx context.Context, boardID string) ([]models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	ChangeBoard(ctx context.Context, id, boardID string) error
	DeleteTask(ctx context.Context, id string) error
}

type Service struct {
	storage Storage
	cash    Cash
}

func New(storage Storage, cash Cash) *Service {
	return &Service{storage: storage, cash: cash}
}

func (s *Service) CreateTask(
	ctx context.Context,
	authorID string,
	boardID string,
	title string,
	description string,
	deadline time.Time,
) (string, error) {
	const op = "service.CreateTask"

	id := uuid.New().String()

	task, err := s.storage.SaveTask(
		ctx,
		id,
		authorID,
		boardID,
		title,
		description,
		deadline,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.cash.SaveTask(ctx, task)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Service) GetTasksByBoardID(ctx context.Context, boardID string) ([]models.Task, error) {
	const op = "service.GetTasksByBoardID"

	tasks, err := s.cash.GetTasksByBoardID(ctx, boardID)
	if err == nil && tasks != nil {
		return tasks, nil
	}

	tasks, err = s.storage.GetTasksByBoardID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.cash.SaveTasks(ctx, tasks)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *Service) UpdateTask(
	ctx context.Context,
	id string,
	title string,
	description string,
	deadline time.Time,
) (models.Task, error) {
	const op = "service.UpdateTask"

	task, err := s.storage.UpdateTask(ctx, id, title, description, deadline)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.cash.UpdateTask(ctx, task)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (s *Service) ChangeBoard(ctx context.Context, id, boardID string) error {
	const op = "service.ChangeBoard"

	c := make(chan error)
	go s.storage.ChangeBoard(ctx, id, boardID, c)

	err := s.cash.ChangeBoard(ctx, id, boardID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = <-c
	close(c)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	const op = "service.ChangeBoard"

	c := make(chan error)
	go s.storage.DeleteTask(ctx, id, c)

	err := s.cash.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = <-c
	close(c)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
