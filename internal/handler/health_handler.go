package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	pool *pgxpool.Pool
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

// Check verifies the application status and database readiness
// GET /health
func (h *HealthHandler) Check(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "connected"
	if h.pool != nil {
		if err := h.pool.Ping(ctx); err != nil {
			dbStatus = "disconnected"
			return c.Status(fiber.StatusServiceUnavailable).JSON(ErrorResponse("service unavailable", fiber.Map{
				"database": dbStatus,
			}))
		}
	} else {
		dbStatus = "uninitialized"
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse("application is healthy", fiber.Map{
		"status":   "ok",
		"database": dbStatus,
	}))
}
