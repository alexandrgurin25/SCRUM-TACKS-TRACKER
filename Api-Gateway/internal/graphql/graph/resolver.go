package graph

import (
	"context"
	"gateway/internal/domain/models"
	"gateway/internal/graphql/graph/model"

	// "gateway/internal/grpc/comment"

	// "gateway/internal/grpc/comment"
	"time"

	"github.com/AlexMickh/scrum-protos/pkg/api/comments"
)

type TaskService interface {
	CreateTask(
		ctx context.Context,
		boardID string,
		title string,
		description string,
		deadline time.Time,
	) (models.Task, error)
	UpdateTask(
		ctx context.Context,
		id string,
		title string,
		description string,
		deadlime time.Time,
	) (*models.Task, error)
	ChangeBoard(ctx context.Context, id, boardId string) error
	DeleteTask(ctx context.Context, id string) error
}

type AuthService interface {
	Register(ctx context.Context, email, username, password string) (models.User, error)
	Login(ctx context.Context, username, password string) (string, string, error)
	UpdateTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type BoardService interface {
	CreateBoard(ctx context.Context, projectId, title string) (models.Board, error)
	GetBoardByID(ctx context.Context, id string) (models.Board, error)
}

type ProjectService interface {
	CreateProject(ctx context.Context, title, descrition string) (string, error)
	GetProject(ctx context.Context, id string) (*models.Project, error)
	AddParticipantToProject(ctx context.Context, id, email string) error
	DeleteProject(ctx context.Context, id string) error
}

type CommentService interface {
	CreateComment(ctx context.Context, taskId, title, description string) (string, error)
	GetComments(
		ctx context.Context,
		taskId string,
		ch chan *model.Comment,
		done chan struct{},
	)
	GetClient() comments.CommentsClient
}

func NewResolver(
	authService AuthService,
	taskService TaskService,
	boardService BoardService,
	projectService ProjectService,
	commentService CommentService,
) *Resolver {
	return &Resolver{
		authService:    authService,
		taskService:    taskService,
		boardService:   boardService,
		projectService: projectService,
		commentService: commentService,
	}
}

type Resolver struct {
	authService    AuthService
	taskService    TaskService
	boardService   BoardService
	projectService ProjectService
	commentService CommentService
}
