package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"event-management-be/internal/domain"
	"event-management-be/internal/handler"

	"github.com/gofiber/fiber/v3"
)

type mockService struct {
	createEventFunc func(ctx context.Context, dto domain.CreateEventDTO) (*domain.Event, error)
}

func (m *mockService) CreateEvent(ctx context.Context, dto domain.CreateEventDTO) (*domain.Event, error) {
	if m.createEventFunc != nil {
		return m.createEventFunc(ctx, dto)
	}
	return &domain.Event{
		ID:          1,
		Title:       dto.Title,
		Description: dto.Description,
		EventDate:   dto.EventDate,
		Category:    dto.Category,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func TestHealthHandler(t *testing.T) {
	app := fiber.New()
	healthHandler := handler.NewHealthHandler(nil)
	app.Get("/health", healthHandler.Check)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var jsonResp handler.JSONResponse
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !jsonResp.Success {
		t.Errorf("expected success to be true")
	}
}

func TestEventHandler_CreateEvent_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockService{}
	eventHandler := handler.NewEventHandler(mockSvc)
	app.Post("/api/v1/events", eventHandler.CreateEvent)

	payload := map[string]interface{}{
		"title":       "Go Backend Workshop",
		"description": "An introductory workshop about backend development.",
		"event_date":  "2026-11-15T09:00:00+07:00",
		"category":    "workshop",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var jsonResp handler.JSONResponse
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !jsonResp.Success {
		t.Errorf("expected success to be true")
	}
}

func TestEventHandler_CreateEvent_MissingTitle(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockService{}
	eventHandler := handler.NewEventHandler(mockSvc)
	app.Post("/api/v1/events", eventHandler.CreateEvent)

	payload := map[string]interface{}{
		"title":       "",
		"description": "An introductory workshop about backend development.",
		"event_date":  "2026-11-15T09:00:00+07:00",
		"category":    "workshop",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestEventHandler_CreateEvent_InvalidCategory(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockService{}
	eventHandler := handler.NewEventHandler(mockSvc)
	app.Post("/api/v1/events", eventHandler.CreateEvent)

	payload := map[string]interface{}{
		"title":       "Go Backend Workshop",
		"description": "An introductory workshop about backend development.",
		"event_date":  "2026-11-15T09:00:00+07:00",
		"category":    "unsupported_cat",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
