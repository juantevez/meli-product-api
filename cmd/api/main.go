package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"meli-product-api/internal/application/service"
	"meli-product-api/internal/infrastructure/adapter/http/handler"
	jsonRepo "meli-product-api/internal/infrastructure/adapter/repository/json"
	"meli-product-api/internal/infrastructure/config"
	"meli-product-api/internal/infrastructure/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Setup logger
	logger := setupLogger(cfg.Logger)
	logger.Info("Starting MELI Product API",
		"version", "1.0.0",
		"port", cfg.Server.Port,
	)

	// Initialize repositories
	logger.Info("Initializing repositories...")

	productRepo, err := jsonRepo.NewProductRepository(cfg.Database.ProductsFile)
	if err != nil {
		logger.Error("Failed to initialize product repository", "error", err)
		log.Fatalf("Failed to initialize product repository: %v", err)
	}

	sellerRepo, err := jsonRepo.NewSellerRepository(cfg.Database.SellersFile)
	if err != nil {
		logger.Error("Failed to initialize seller repository", "error", err)
		log.Fatalf("Failed to initialize seller repository: %v", err)
	}

	reviewRepo, err := jsonRepo.NewReviewRepository(cfg.Database.ReviewsFile)
	if err != nil {
		logger.Error("Failed to initialize review repository", "error", err)
		log.Fatalf("Failed to initialize review repository: %v", err)
	}

	questionRepo, err := jsonRepo.NewQuestionRepository(cfg.Database.QuestionsFile)
	if err != nil {
		logger.Error("Failed to initialize question repository", "error", err)
		log.Fatalf("Failed to initialize question repository: %v", err)
	}

	logger.Info("✓ Repositories initialized successfully")

	// Initialize services
	logger.Info("Initializing services...")

	aggregatorService := service.NewProductAggregatorService(
		productRepo,
		sellerRepo,
		reviewRepo,
		questionRepo,
		logger,
	)

	searchService := service.NewProductSearchService(
		productRepo,
		logger,
	)

	logger.Info("✓ Services initialized successfully")

	// Initialize handlers
	productHandler := handler.NewProductHandler(
		aggregatorService,
		searchService,
		logger,
	)

	// Setup router
	r := router.NewRouter(productHandler, logger)

	// HTTP Server configuration
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("╔═══════════════════════════════════════════════════════════════")
		logger.Info("║ 🚀 MELI Product API Server Started")
		logger.Info("║ 📍 Listening on: http://" + addr)
		logger.Info("║ 🏥 Health Check: http://" + addr + "/health")
		logger.Info("║ 📚 API Docs: http://" + addr + "/api/v1/products")
		logger.Info("╚═══════════════════════════════════════════════════════════════")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		log.Fatal(err)
	}

	logger.Info("Server exited gracefully")
}

func setupLogger(cfg config.LoggerConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
