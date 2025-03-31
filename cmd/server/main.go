package main

import (
	"fmt"
	"os"
	"os/signal"
	"sweng-task/internal/app"
	"syscall"

	"go.uber.org/zap"
	"sweng-task/internal/config"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	log := logger.Sugar()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Infow("Configuration loaded",
		"environment", cfg.App.Environment,
		"log_level", cfg.App.LogLevel,
		"server_port", cfg.Server.Port,
	)

	fiberApp := app.SetupApp(cfg, log)

	// Start server
	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Infof("Starting server on %s", address)
		if err := fiberApp.Listen(address); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	if err := fiberApp.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Info("Server gracefully stopped")
}
