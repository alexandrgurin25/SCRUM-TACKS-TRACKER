package server

import (
	"comments/internal/domain/models"
	"comments/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlexMickh/logger/pkg/logger"
	"github.com/AlexMickh/scrum-protos/pkg/api/comments"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateComment(
		ctx context.Context,
		authorId string,
		taskId string,
		title string,
		description string,
	) (string, error)
	GetAllByTaskId(ctx context.Context, taskId string, lastTime time.Time) ([]*models.Comment, error)
}

type AuthClient interface {
	VerifyToken(ctx context.Context, token string) (string, error)
}

type Server struct {
	comments.UnimplementedCommentsServer
	service    Service
	authClient AuthClient
}

func New(service Service, authClient AuthClient) *Server {
	return &Server{
		service:    service,
		authClient: authClient,
	}
}

func (s *Server) CreateComment(
	ctx context.Context,
	req *comments.CreateCommentRequest,
) (*comments.CreateCommentResponse, error) {
	const op = "grpc.server.CreateComment"

	ctx = logger.GetFromCtx(ctx).With(ctx, zap.String("op", op))

	if req.GetTaskId() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "task_id is empty")
		return nil, status.Error(codes.InvalidArgument, "task_id is required")
	}
	if req.GetTitle() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "title is empty")
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetDescription() == "" {
		logger.GetFromCtx(ctx).Error(ctx, "description is empty")
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	// authorId := ctx.Value(idKey).(string)

	// TODO: add auth interceptor
	id, err := s.service.CreateComment(
		ctx,
		"1",
		req.GetTaskId(),
		req.GetTitle(),
		req.GetDescription(),
	)
	if err != nil {
		logger.GetFromCtx(ctx).Error(ctx, "can't create comment", zap.Error(err))
		return nil, status.Error(codes.Internal, "server error")
	}

	return &comments.CreateCommentResponse{
		Id: id,
	}, nil
}

// FIXME: don't work with goroutine
func (s *Server) GetComments(
	req *comments.GetCommentsRequest,
	srv grpc.ServerStreamingServer[comments.GetCommentsResponse],
) error {
	// const op = "grpc.server.GetComments"
	// var wg sync.WaitGroup
	// var mut sync.Mutex

	lastTime := time.Time{}
	ctx := srv.Context()

	// ctx = logger.GetFromCtx(ctx).With(ctx, zap.String("op", op))

	if req.GetTaskId() == "" {
		// logger.GetFromCtx(ctx).Error(ctx, "task_id is empty")
		return status.Error(codes.InvalidArgument, "description is required")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		comms, err := s.service.GetAllByTaskId(ctx, req.GetTaskId(), lastTime)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				fmt.Println("yes")
				continue
			}
			// logger.GetFromCtx(ctx).Error(ctx, "failed to get comments", zap.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		if len(comms) > 0 {
			l := len(comms)
			// wg.Add(l)
			fmt.Println(l)
			fmt.Printf("%v\n", comms)
			for _, comment := range comms {
				// fmt.Println(i)
				func() {
					// defer wg.Done()
					// mut.Lock()
					// defer mut.Unlock()

					lastTime = comment.CreatedAt
					srv.Send(&comments.GetCommentsResponse{
						Id:          comment.ID,
						AuthorId:    comment.AuthorID,
						TaskId:      comment.TaskID,
						Title:       comment.Title,
						Description: comment.Description,
					})
				}()
			}
			// wg.Wait()
		}

	}
}

type key string

const idKey = key("id")

func (s *Server) AuthUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	const op = "grpc.server.AuthUnaryInterceptor"

	ctx = logger.GetFromCtx(ctx).With(ctx, zap.String("op", op))

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

	if info.FullMethod == "/comments.Comments/CreateTask" {
		ctx = context.WithValue(ctx, idKey, id)
	}

	return handler(ctx, req)
}

func (s *Server) AuthStreamInterceptor(
	srv any,
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	const op = "grpc.server.AuthUnaryInterceptor"

	ctx := stream.Context()
	ctx = logger.GetFromCtx(ctx).With(ctx, zap.String("op", op))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.GetFromCtx(ctx).Error(ctx, "failed to get metadata")
		return status.Error(codes.Unauthenticated, "no metadata")
	}

	token, ok := md["authorization"]
	if !ok {
		logger.GetFromCtx(ctx).Error(ctx, "failed to get auth token")
		return status.Error(codes.Unauthenticated, "no token")
	}

	_, err := s.authClient.VerifyToken(ctx, token[0])
	if err != nil {
		logger.GetFromCtx(ctx).Error(ctx, "failed to verify token", zap.Error(err))
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	return handler(srv, stream)
}
