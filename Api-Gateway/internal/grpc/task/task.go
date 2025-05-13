package task

import (
	"context"
	"fmt"
	"gateway/internal/domain/models"
	"gateway/internal/graphql"
	client "gateway/internal/grpc"

	"time"

	"github.com/AlexMickh/scrum-protos/pkg/api/tasks"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskService struct {
	conn   *grpc.ClientConn
	client tasks.TasksClient
}

func New(addr string) (*TaskService, error) {
	const op = "grpc.task.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := tasks.NewTasksClient(conn)

	return &TaskService{
		conn:   conn,
		client: client,
	}, nil
}

func (t *TaskService) CreateTask(
	ctx context.Context,
	boardID string,
	title string,
	description string,
	deadline time.Time,
) (models.Task, error) {
	const op = "grpc.task.CreateTask"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return models.Task{}, err
	}

	res, err := t.client.CreateTask(ctx, &tasks.CreateTaskRequest{
		BoardId:     boardID,
		Title:       title,
		Description: description,
		Deadline:    timestamppb.New(deadline),
	})
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.Task{
		ID:          res.GetId(),
		Title:       title,
		Description: description,
		AuthorID:    "",
		BoardID:     boardID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Deadline:    deadline,
	}, nil
}

func (t *TaskService) UpdateTask(
	ctx context.Context,
	id string,
	title string,
	description string,
	deadlime time.Time,
) (*models.Task, error) {
	const op = "grpc.task.UpdateTask"

	token, err := graphql.GetAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	res, err := t.client.UpdateTask(ctx, &tasks.UpdateTaskRequest{
		Id:          id,
		Title:       title,
		Description: description,
		Deadline:    timestamppb.New(deadlime),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Task{
		ID:          res.GetTask().GetId(),
		Title:       res.GetTask().GetTitle(),
		Description: res.GetTask().GetDescription(),
		AuthorID:    res.GetTask().GetAuthorId(),
		BoardID:     res.GetTask().GetBoardId(),
		CreatedAt:   res.GetTask().GetCreatedAt().AsTime(),
		UpdatedAt:   res.GetTask().GetUpdatedAt().AsTime(),
		Deadline:    res.GetTask().GetDeadline().AsTime(),
	}, nil
}

func (t *TaskService) ChangeBoard(ctx context.Context, id, boardId string) error {
	const op = "grpc.task.ChangeBoard"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return err
	}

	_, err = t.client.ChangeBoard(ctx, &tasks.ChangeBoardRequest{
		Id:      id,
		BoardId: boardId,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (t *TaskService) DeleteTask(ctx context.Context, id string) error {
	const op = "grpc.task.DeleteTask"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return err
	}

	_, err = t.client.DeleteTask(ctx, &tasks.DeleteTaskRequest{Id: id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (t *TaskService) Stop() {
	t.conn.Close()
}
