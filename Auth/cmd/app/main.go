package main

import (
	"auth/internal/config"
	repository "auth/internal/repository/postgres"
	"auth/internal/service"
	transportServers "auth/internal/transport/servers"
	"auth/pkg/logger"
	"auth/pkg/postgres"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed to load config", zap.Error(err))
	}

	pg, err := postgres.New(ctx, cfg)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed to connect to PostgreSQL:", zap.Error(err))
	}
	defer pg.Close()

	userRepository := repository.NewUserRepository(pg)
	authService := service.NewAuthService(userRepository, cfg)

	grpcServer := transportServers.NewServer(*cfg, authService)

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	grpcServer.Stop()

	log.Println("Server exited properly")
}
