package server

import (
	"boards/internal/service"
	"boards/pkg/logger"
	"context"
	"fmt"
	"github.com/AlexMickh/scrum-protos/pkg/api/boards"
	"github.com/AlexMickh/scrum-protos/pkg/api/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthClient interface {
	VerifyToken(ctx context.Context, token string) (string, error)
}

type Server struct {
	boards.UnimplementedBoardsServer
	service    service.Boards
	authClient AuthClient
}
type key string

const idKey = key("id")

func InitServer(service service.Boards, auth AuthClient) *Server {
	return &Server{service: service, authClient: auth}
}

func (s *Server) CreateBoard(ctx context.Context, req *boards.CreateBoardRequest) (*boards.CreateBoardResponse, error) {
	if req.GetProjectId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Project ID cannot be empty")
	}

	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	authorID := ctx.Value(idKey).(string)

	id, err := s.service.CreateBoard(ctx, authorID, req.GetProjectId(), req.GetTitle())

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to create board", zap.Error(err))
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	return &boards.CreateBoardResponse{Id: id}, nil
}

func (s *Server) GetBoard(ctx context.Context, req *boards.GetBoardRequest) (*boards.GetBoardResponse, error) {
	if req.GetId() == "" {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "id is required")
		return nil, status.Error(codes.InvalidArgument, "failed to get board")
	}

	board, tasks, err := s.service.GetBoard(ctx, req.GetId())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get board", zap.Error(err))
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	var tasksInit []*types.Task
	for _, task := range tasks.Tasks {
		taskInit := &types.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			AuthorId:    task.AuthorID,
			BoardId:     task.BoardID,
			CreatedAt:   timestamppb.New(task.CreatedAt),
			UpdatedAt:   timestamppb.New(task.UpdatedAt),
			Deadline:    timestamppb.New(task.Deadline),
		}
		tasksInit = append(tasksInit, taskInit)
	}

	var boardInit *boards.GetBoardResponse
	boardInit = &boards.GetBoardResponse{
		Id:    board.ID,
		Title: board.Title,
		Tasks: tasksInit,
	}

	return &boards.GetBoardResponse{Id: boardInit.Id}, nil
}

func (s *Server) GetAllBoards(ctx context.Context, req *boards.GetAllBoardsRequest) (*boards.GetAllBoardsResponse, error) {
	if req.GetId() == "" {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "id is required")
		return nil, status.Error(codes.InvalidArgument, "failed to get all boards")
	}

	boardList, err := s.service.GetAllBoards(ctx, req.GetId())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get all boards", zap.Error(err))
		return nil, fmt.Errorf("failed to get all boards: %w", err)
	}

	response := &boards.GetAllBoardsResponse{
		Boards: make([]*boards.GetBoardResponse, 0, len(boardList.Boards)),
	}

	for _, board := range boardList.Boards {
		response.Boards = append(response.Boards, &boards.GetBoardResponse{
			Id:    board.ID,
			Title: board.Title,
		})
	}

	return response, nil
}

func (s *Server) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if info.FullMethod == "/boards.Boards/GetAllBoards" {
		return handler(ctx, req)
	}
	fmt.Println("AuthInterceptor called")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get metadata")
		return nil, status.Error(codes.Unauthenticated, "no metadata")
	}

	token, ok := md["authorization"]
	if !ok {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get auth token")
		return nil, status.Error(codes.Unauthenticated, "no token")
	}

	id, err := s.authClient.VerifyToken(ctx, token[0])
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to verify token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if info.FullMethod == "/boards.Boards/CreateBoard" {
		ctx = context.WithValue(ctx, idKey, id)
	}

	return handler(ctx, req)
}
