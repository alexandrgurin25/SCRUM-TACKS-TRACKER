package board

import (
	"context"
	"fmt"
	"gateway/internal/domain/models"

	client "gateway/internal/grpc"

	"github.com/AlexMickh/scrum-protos/pkg/api/boards"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BoardService struct {
	conn   *grpc.ClientConn
	client boards.BoardsClient
}

func New(addr string) (*BoardService, error) {
	const op = "grpc.board.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := boards.NewBoardsClient(conn)

	return &BoardService{
		conn:   conn,
		client: client,
	}, nil
}

func (b *BoardService) CreateBoard(ctx context.Context, projectId, title string) (models.Board, error) {
	const op = "grpc.board.CreateBoard"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return models.Board{}, err
	}

	res, err := b.client.CreateBoard(ctx, &boards.CreateBoardRequest{Title: title})
	if err != nil {
		return models.Board{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.Board{
		ID:    res.GetId(),
		Title: title,
		Tasks: []*models.Task{},
	}, nil
}

func (b *BoardService) GetBoardByID(ctx context.Context, id string) (models.Board, error) {
	const op = "grpc.board.GetBoardByID"

	ctx, err := client.AddToken(ctx, op)
	if err != nil {
		return models.Board{}, err
	}

	res, err := b.client.GetBoard(ctx, &boards.GetBoardRequest{Id: id})
	if err != nil {
		return models.Board{}, fmt.Errorf("%s: %w", op, err)
	}

	var tasks []*models.Task
	for _, task := range res.GetTasks() {
		taskL := &models.Task{
			ID:          task.Id,
			Title:       task.GetTitle(),
			Description: task.GetDescription(),
			AuthorID:    task.GetAuthorId(),
			BoardID:     task.GetBoardId(),
			CreatedAt:   task.CreatedAt.AsTime(),
			UpdatedAt:   task.UpdatedAt.AsTime(),
			Deadline:    task.Deadline.AsTime(),
		}
		tasks = append(tasks, taskL)
	}

	return models.Board{
		ID:    res.GetId(),
		Title: res.GetTitle(),
		Tasks: tasks,
	}, nil
}

func (b *BoardService) Stop() {
	b.conn.Close()
}
