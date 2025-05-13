package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"task/internal/config"
	"task/internal/grpc/client"
	"task/internal/grpc/server"
	"task/internal/service"
	"task/internal/storage/postgres"
	"task/internal/storage/redis"

	"github.com/AlexMickh/scrum-protos/pkg/api/tasks"

	"github.com/AlexMickh/logger/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	path := ".env"
	cfg := config.MustLoad(path)

	ctx := context.Background()
	ctx, err := logger.New(ctx, cfg.Env)
	if err != nil {
		panic(err)
	}

	logger.GetFromCtx(ctx).Info(ctx, "logger is working", zap.String("env", cfg.Env))

	logger.GetFromCtx(ctx).Info(ctx, "starting grpc server", zap.Int("port", cfg.Port))

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to listen", zap.Error(err))
	}

	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to connect to the postgres", zap.Error(err))
	}

	cash, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to connect to the redis")
	}

	service := service.New(db, cash)

	client, err := client.New(cfg.AuthServiceAddr)
	if err != nil {
		logger.GetFromCtx(ctx).Fatal(ctx, "failed to connect to the auth service")
	}

	srv := server.New(service, client)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logger.Interceptor(ctx),
		grpc.UnaryServerInterceptor(srv.AuthInterceptor),
	))
	tasks.RegisterTasksServer(server, srv)

	go func() {
		if err := server.Serve(lis); err != nil {
			logger.GetFromCtx(ctx).Fatal(ctx, "failed to serve", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	server.GracefulStop()
	client.Close()
	db.Close()
	cash.Close()

	logger.GetFromCtx(ctx).Info(ctx, "server stopped")

}
