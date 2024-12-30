package main

import (
	"context"
	"github.com/1Storm3/flibox-api/internal/app"
	"github.com/1Storm3/flibox-api/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// @title Swagger Flibox API
// @version 1.0
// @description Flibox API
// @host localhost:8080
// @BasePath /api
func main() {
	logger.Init("development")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		cancel()
	}()
	
	a := &app.App{}
	if err := a.Run(ctx); err != nil {
		logger.Error("Error running app", zap.Error(err))
	}

	<-ctx.Done()
	logger.Info("Shutting down gracefully...")
}
