package main

import (
	"boards/internal/config"
	repository "boards/internal/repository/board"
	service "boards/internal/service/board"
	"boards/internal/transport/client"
	"boards/internal/transport/server"
	"boards/pkg/logger"
	"boards/pkg/postgres"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexMickh/scrum-protos/pkg/api/boards"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	ctx, err = logger.NewLogger(ctx)
	if err != nil {
		log.Fatal("Failed to init logger", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed to load config", zap.Error(err))
	}

	pool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to listen", zap.Error(err))
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Starting gRPC server on port :%d", cfg.GRPCPort))
	repo := repository.InitRepo(pool)
	serv := service.InitServ(repo, cfg)

	newClient, err := client.InitClient(cfg.AuthServiceAddr)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed to init client", zap.Error(err))
	}

	srv := server.InitServer(serv, newClient)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(logger.Interceptor(ctx), grpc.UnaryServerInterceptor(srv.AuthInterceptor)))
	boards.RegisterBoardsServer(grpcServer, srv)
	reflection.Register(grpcServer)

	go func() {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC server started", zap.Int("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Server stopping...")

	grpcServer.GracefulStop()
	newClient.Close()
	pool.Close()

	log.Println("Server gracefully stopped")
}
