package graphql

import (
	"context"
	"errors"
	"fmt"
	"gateway/internal/graphql/graph"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/AlexMickh/logger/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

type Server struct {
	port           int
	authService    graph.AuthService
	taskService    graph.TaskService
	boardService   graph.BoardService
	projectService graph.ProjectService
	commentService graph.CommentService
}

func New(
	port int,
	authService graph.AuthService,
	taskService graph.TaskService,
	boardService graph.BoardService,
	projectService graph.ProjectService,
	commentService graph.CommentService,
) *Server {
	return &Server{
		port:           port,
		authService:    authService,
		taskService:    taskService,
		boardService:   boardService,
		projectService: projectService,
		commentService: commentService,
	}
}

type key string

const authKey = key("Authorization")

var ErrNoToken = errors.New("token is empty")

func middleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), logger.Key, log)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		ctx := context.WithValue(r.Context(), authKey, auth)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GetAuthToken(ctx context.Context) (string, error) {
	const op = "graphql.server.GetAuthToken"

	token := ctx.Value(authKey).(string)
	if token == "" {
		return "", fmt.Errorf("%s: %w", op, ErrNoToken)
	}

	return token, nil
}

func (s *Server) Run(ctx context.Context) {
	resolver := graph.NewResolver(s.authService, s.taskService, s.boardService, s.projectService, s.commentService)

	router := chi.NewRouter()
	router.Use(logger.ChiMiddleware(ctx))
	router.Use(middleware(logger.GetFromCtx(ctx)))
	router.Use(authMiddleware)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Websocket{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	logger.GetFromCtx(ctx).Info(ctx, fmt.Sprintf("connect to http://localhost:%d/ for GraphQL playground", s.port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to listen", zap.Error(err))
	}
}
