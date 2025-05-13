package main

import (
	"context"
	"gateway/internal/config"
	"gateway/internal/graphql"
	"gateway/internal/grpc/auth"
	"gateway/internal/grpc/board"
	"gateway/internal/grpc/comment"
	"gateway/internal/grpc/project"
	"gateway/internal/grpc/task"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexMickh/logger/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	path := ".env"
	cfg := config.MustLoad(path)

	ctx := context.Background()
	ctx, err := logger.New(ctx, cfg.Env)
	if err != nil {
		panic(err)
	}

	auth, err := auth.New(cfg.AuthServiceAddr)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init auth service client", zap.Error(err))
	}
	defer auth.Stop()

	ts, err := task.New(cfg.TaskServiceAddr)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init task service client", zap.Error(err))
	}
	defer ts.Stop()

	bs, err := board.New(cfg.BoardServiceAddr)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init board service client", zap.Error(err))
	}
	defer bs.Stop()

	ps, err := project.New(cfg.ProjectServiceAddr)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init project service client", zap.Error(err))
	}
	defer ps.Stop()

	cs, err := comment.New(cfg.CommentServiceAddr)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init comment service client", zap.Error(err))
	}
	defer cs.Close()

	server := graphql.New(cfg.Port, auth, ts, bs, ps, cs)
	go server.Run(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	logger.GetFromCtx(ctx).Info(ctx, "server stopped")
}
