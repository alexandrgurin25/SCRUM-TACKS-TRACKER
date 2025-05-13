package server

import (
	"context"
	"fmt"
	"github.com/AlexMickh/scrum-protos/pkg/api/projects"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"projects/internal/service"
	"projects/pkg/logger"
)

type AuthClient interface {
	VerifyToken(ctx context.Context, token string) (string, error)
}

type Server struct {
	projects.UnimplementedProjectsServer
	service    service.Projects
	authClient AuthClient
}
type contextKey string

const (
	idKey contextKey = "authorID"
)

func InitServer(service service.Projects, auth AuthClient) *Server {
	return &Server{service: service, authClient: auth}
}

func (s *Server) CreateProject(ctx context.Context, req *projects.CreateProjectRequest) (*projects.CreateProjectResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if req.GetDescription() == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	authorID := ctx.Value(idKey).(string)

	id, err := s.service.CreateProject(ctx, authorID, req.GetDescription(), req.GetTitle())

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to create project", zap.Error(err))
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &projects.CreateProjectResponse{Id: id}, nil
}

func (s *Server) GetProject(ctx context.Context, req *projects.GetProjectRequest) (*projects.GetProjectResponse, error) {
	if req.GetId() == "" {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "id is required")
		return nil, status.Error(codes.InvalidArgument, "failed to get project")
	}

	project, participants, err := s.service.GetProject(ctx, req.GetId())

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get project", zap.Error(err))
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	var projectInit *projects.GetProjectResponse
	projectInit = &projects.GetProjectResponse{
		Id:             project.ID,
		AuthorId:       project.AuthorID,
		ParticipantsId: participants.UsersID,
		Title:          project.Title,
		Description:    project.Description,
		CreatedAt:      timestamppb.New(project.CreatedAt),
		UpdatedAt:      timestamppb.New(project.UpdatedAt),
	}

	return projectInit, nil
}

func (s *Server) AddParticipant(ctx context.Context, req *projects.AddParticipantRequest) (*empty.Empty, error) {
	if req.GetProjectId() == "" {
		return nil, status.Error(codes.InvalidArgument, "project id is required")
	}

	if req.GetParticipantEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "participant email is required")
	}

	err := s.service.AddParticipant(ctx, req.GetProjectId(), req.GetParticipantEmail())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to add participant", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to add participant")
	}

	return &empty.Empty{}, nil
}

func (s *Server) DeleteProject(ctx context.Context, req *projects.DeleteProjectRequest) (*empty.Empty, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}

	authorID, ok := ctx.Value(idKey).(string)
	if !ok {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "authorID is required")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	err := s.service.DeleteProject(ctx, req.GetId(), authorID)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to delete project", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete project")
	}

	return &empty.Empty{}, nil
}

func (s *Server) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
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

	authorId, err := s.authClient.VerifyToken(ctx, token[0])
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to verify token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if info.FullMethod == "/projects.Projects/CreateProject" || info.FullMethod == "/projects.Projects/DeleteProject" {
		ctx = context.WithValue(ctx, idKey, authorId)
	}

	return handler(ctx, req)
}
