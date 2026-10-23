package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-management-be/internal/config"
	"event-management-be/internal/database"
	"event-management-be/internal/handler"
	"event-management-be/internal/repository/postgres"
	"event-management-be/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL successfully")

	// 3. Auto-apply migrations
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}
	if err := database.RunMigrations(context.Background(), pool, migrationsPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migrations checked/applied successfully")

	// 4. Clean Architecture Dependency Wiring
	eventRepo := postgres.NewEventRepository(pool)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService)
	healthHandler := handler.NewHealthHandler(pool)

	// 5. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Campus Event Management API",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	// 6. Register Routes
	app.Get("/health", healthHandler.Check)

	api := app.Group("/api/v1")
	api.Post("/events", eventHandler.CreateEvent)

	// 7. Graceful Shutdown setup
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	go func() {
		log.Printf("Starting HTTP server on %s", serverAddr)
		if err := app.Listen(serverAddr, fiber.ListenConfig{
			DisableStartupMessage: false,
		}); err != nil {
			log.Printf("Server shut down error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-shutdownChan
	log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
