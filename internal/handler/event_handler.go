package handler

import (
	"errors"
	"strings"
	"time"

	"event-management-be/internal/domain"
	"event-management-be/internal/service"

	"github.com/gofiber/fiber/v3"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
	service service.EventService
}

// NewEventHandler creates a new EventHandler
func NewEventHandler(service service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// createEventRequest defines incoming payload format
type createEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
	Category    string `json:"category"`
}

// CreateEvent handles POST /api/v1/events
func (h *EventHandler) CreateEvent(c fiber.Ctx) error {
	var req createEventRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse("invalid request payload", "failed to parse request body"))
	}

	// Request validation
	trimmedTitle := strings.TrimSpace(req.Title)
	if trimmedTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse("validation failed", "title is required"))
	}

	trimmedDate := strings.TrimSpace(req.EventDate)
	if trimmedDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse("validation failed", "event_date is required"))
	}

	parsedDate, err := time.Parse(time.RFC3339, trimmedDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse("validation failed", "event_date must be in RFC3339 format (e.g. 2026-11-15T09:00:00+07:00)"))
	}

	trimmedCategory := strings.TrimSpace(req.Category)
	if trimmedCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse("validation failed", "category is required"))
	}

	category := domain.EventCategory(trimmedCategory)
	if !category.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse("validation failed", "invalid category: must be workshop, competition, seminar, or community"))
	}

	dto := domain.CreateEventDTO{
		Title:       trimmedTitle,
		Description: strings.TrimSpace(req.Description),
		EventDate:   parsedDate,
		Category:    category,
	}

	event, err := h.service.CreateEvent(c.Context(), dto)
	if err != nil {
		if errors.Is(err, domain.ErrTitleRequired) ||
			errors.Is(err, domain.ErrEventDateRequired) ||
			errors.Is(err, domain.ErrInvalidCategory) {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse("validation failed", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse("failed to create event", "internal server error"))
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse("event created successfully", event))
}
