package server

import (
	"context"
	"errors"
	"task/internal/domain/models"
	"task/internal/storage"
	"time"

	"github.com/AlexMickh/scrum-protos/pkg/api/tasks"
	"github.com/AlexMickh/scrum-protos/pkg/api/types"

	"github.com/AlexMickh/logger/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	CreateTask(
		ctx context.Context,
		authorID string,
		boardID string,
		title string,
		description string,
		deadline time.Time,
	) (string, error)
	GetTasksByBoardID(ctx context.Context, boardID string) ([]models.Task, error)
	UpdateTask(
		ctx context.Context,
		id string,
		title string,
		description string,
		deadline time.Time,
	) (models.Task, error)
	ChangeBoard(ctx context.Context, id, boardID string) error
	DeleteTask(ctx context.Context, id string) error
}

type AuthClient interface {
	VerifyToken(ctx context.Context, token string) (string, error)
}

type Server struct {
	tasks.UnimplementedTasksServer
	service    Service
	authClient AuthClient
}

type key string

const idKey = key("id")

func New(service Service, authClient AuthClient) *Server {
	return &Server{service: service, authClient: authClient}
}

func (s *Server) CreateTask(ctx context.Context, req *tasks.CreateTaskRequest) (*tasks.CreateTaskResponse, error) {
	const op = "grpc.server.CreateTask"

	ctx = logger.GetFromCtx(ctx).With(ctx, zap.String("op", op))

	if req.GetBoardId() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "board_id is empty")
		return nil, status.Error(codes.InvalidArgument, "board_id is required")
	}
	if req.GetTitle() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "title is empty")
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetDescription() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "description is empty")
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}
	if req.GetDeadline().String() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "deadline is empty")
		return nil, status.Error(codes.InvalidArgument, "deadline is required")
	}

	authorID := ctx.Value(idKey).(string)

	id, err := s.service.CreateTask(
		ctx,
		authorID,
		req.GetBoardId(),
		req.GetTitle(),
		req.GetDescription(),
		req.Deadline.AsTime(),
	)
	if err != nil {
		logger.GetFromCtx(ctx).Error(ctx, "failed to create task", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	return &tasks.CreateTaskResponse{
		Id:          id,
		BoardId:     req.GetBoardId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Deadline:    req.Deadline,
	}, nil
}

func (s *Server) GetTasksByBoardID(
	ctx context.Context,
	req *tasks.GetTasksByBoardIDRequest,
) (*tasks.GetAllTasksResponse, error) {
	const op = "grpc.server.GetTasksByBoardID"

	ctx = logger.GetFromCtx(ctx).With(ctx,
		zap.String("op", op),
		zap.String("board_id", req.GetBoardId()),
	)

	if req.GetBoardId() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "board_id is empty")
		return nil, status.Error(codes.InvalidArgument, "board_id is required")
	}

	tasksList, err := s.service.GetTasksByBoardID(ctx, req.GetBoardId())
	if err != nil {
		logger.GetFromCtx(ctx).Error(ctx, "failed to get tasks by board id", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get tasks by board id")
	}

	var tasksType []*types.Task
	for _, task := range tasksList {
		taskType := &types.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			AuthorId:    task.AuthorID,
			BoardId:     task.BoardID,
			CreatedAt:   timestamppb.New(task.CreatedAt),
			UpdatedAt:   timestamppb.New(task.UpdatedAt),
			Deadline:    timestamppb.New(task.Deadline),
		}
		tasksType = append(tasksType, taskType)
	}

	return &tasks.GetAllTasksResponse{Tasks: tasksType}, nil
}

func (s *Server) UpdateTask(ctx context.Context, req *tasks.UpdateTaskRequest) (*tasks.UpdateTaskResponse, error) {
	const op = "grpc.server.UpdateTask"

	ctx = logger.GetFromCtx(ctx).With(ctx,
		zap.String("op", op),
		zap.String("id", req.GetId()),
	)

	if req.GetId() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "id is empty")
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	task, err := s.service.UpdateTask(
		ctx,
		req.GetId(),
		req.GetTitle(),
		req.GetDescription(),
		req.GetDeadline().AsTime(),
	)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			logger.GetFromCtx(ctx).Error(ctx, "user not found", zap.String("id", req.GetId()))
			return nil, status.Error(codes.Internal, "user not found")
		}
		logger.GetFromCtx(ctx).Error(ctx, "failed to update task", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	taskType := &types.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		AuthorId:    task.AuthorID,
		BoardId:     task.BoardID,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
		Deadline:    timestamppb.New(task.Deadline),
	}

	return &tasks.UpdateTaskResponse{Task: taskType}, nil
}

func (s *Server) ChangeBoard(ctx context.Context, req *tasks.ChangeBoardRequest) (*emptypb.Empty, error) {
	const op = "grpc.server.ChangeBoard"

	ctx = logger.GetFromCtx(ctx).With(ctx,
		zap.String("op", op),
		zap.String("id", req.GetId()),
	)

	if req.GetId() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "id is empty")
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetBoardId() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "board_id is empty")
		return nil, status.Error(codes.InvalidArgument, "board_id is required")
	}

	err := s.service.ChangeBoard(ctx, req.GetId(), req.GetBoardId())
	if err != nil {
		logger.GetFromCtx(ctx).Error(ctx, "failed to change board", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to change board")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteTask(ctx context.Context, req *tasks.DeleteTaskRequest) (*emptypb.Empty, error) {
	const op = "grpc.server.DeleteTask"

	ctx = logger.GetFromCtx(ctx).With(ctx,
		zap.String("op", op),
		zap.String("id", req.GetId()),
	)

	if req.GetId() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "id is empty")
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.service.DeleteTask(ctx, req.GetId())
	if err != nil {
		logger.GetFromCtx(ctx).Error(ctx, "failed to delete task", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete task")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if info.FullMethod == "/tasks.Tasks/GetTasksByBoardID" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.GetFromCtx(ctx).Error(ctx, "failed to get metadata")
		return nil, status.Error(codes.Unauthenticated, "no metadata")
	}

	token, ok := md["authorization"]
	if !ok {
		logger.GetFromCtx(ctx).Error(ctx, "failed to get auth token")
		return nil, status.Error(codes.Unauthenticated, "no token")
	}

	id, err := s.authClient.VerifyToken(ctx, token[0])
	if err != nil {
		logger.GetFromCtx(ctx).Error(ctx, "failed to verify token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if info.FullMethod == "/tasks.Tasks/CreateTask" {
		ctx = context.WithValue(ctx, idKey, id)
	}

	return handler(ctx, req)
}
