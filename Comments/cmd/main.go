package main

import (
	"comments/internal/config"
	"comments/internal/grpc/client"
	"comments/internal/grpc/server"
	"comments/internal/service"
	"comments/internal/storage/postgres"
	"comments/internal/storage/redis"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexMickh/logger/pkg/logger"
	"github.com/AlexMickh/scrum-protos/pkg/api/comments"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()
	ctx, err := logger.New(ctx, cfg.Env)
	if err != nil {
		panic("failed to init logger")
	}

	logger.GetFromCtx(ctx).Info(ctx, "logger is working", zap.String("env", cfg.Env))

	logger.GetFromCtx(ctx).Info(ctx, "starting grpc server", zap.Int("port", cfg.Port))

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to listen", zap.Error(err))
	}

	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init db", zap.Error(err))
	}
	defer db.Close()

	cash, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init cash", zap.Error(err))
	}
	defer cash.Close()

	service := service.New(db, cash)

	client, err := client.New(cfg.AuthServiceAddr)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to init auth client", zap.Error(err))
	}
	defer client.Close()

	srv := server.New(service, client)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.Interceptor(ctx),
			// srv.AuthUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			logger.StreamInterceptor(ctx),
			// srv.AuthStreamInterceptor,
		),
	)
	comments.RegisterCommentsServer(server, srv)
	defer server.GracefulStop()

	go func() {
		if err := server.Serve(lis); err != nil {
			logger.GetFromCtx(ctx).Fatal(ctx, "failed to serve", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	logger.GetFromCtx(ctx).Info(ctx, "server stopped")
}
