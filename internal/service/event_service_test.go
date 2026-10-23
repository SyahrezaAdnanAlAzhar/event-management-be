package service_test

import (
	"context"
	"testing"
	"time"

	"event-management-be/internal/domain"
	"event-management-be/internal/service"
)

// mockEventRepository implements domain.EventRepository for unit testing
type mockEventRepository struct {
	createFunc func(ctx context.Context, event *domain.Event) (*domain.Event, error)
}

func (m *mockEventRepository) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, event)
	}
	return &domain.Event{
		ID:          1,
		Title:       event.Title,
		Description: event.Description,
		EventDate:   event.EventDate,
		Category:    event.Category,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func TestEventService_CreateEvent_Success(t *testing.T) {
	mockRepo := &mockEventRepository{}
	svc := service.NewEventService(mockRepo)

	eventDate, _ := time.Parse(time.RFC3339, "2026-11-15T09:00:00+07:00")
	dto := domain.CreateEventDTO{
		Title:       "Go Backend Workshop",
		Description: "An introductory workshop about backend development.",
		EventDate:   eventDate,
		Category:    domain.CategoryWorkshop,
	}

	result, err := svc.CreateEvent(context.Background(), dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.Title != dto.Title {
		t.Errorf("expected title %q, got %q", dto.Title, result.Title)
	}
	if result.Category != domain.CategoryWorkshop {
		t.Errorf("expected category %q, got %q", domain.CategoryWorkshop, result.Category)
	}
}

func TestEventService_CreateEvent_MissingTitle(t *testing.T) {
	mockRepo := &mockEventRepository{}
	svc := service.NewEventService(mockRepo)

	eventDate, _ := time.Parse(time.RFC3339, "2026-11-15T09:00:00+07:00")
	dto := domain.CreateEventDTO{
		Title:       "   ",
		Description: "Valid description",
		EventDate:   eventDate,
		Category:    domain.CategoryWorkshop,
	}

	_, err := svc.CreateEvent(context.Background(), dto)
	if err != domain.ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestEventService_CreateEvent_MissingEventDate(t *testing.T) {
	mockRepo := &mockEventRepository{}
	svc := service.NewEventService(mockRepo)

	dto := domain.CreateEventDTO{
		Title:       "Go Backend Workshop",
		Description: "Valid description",
		EventDate:   time.Time{},
		Category:    domain.CategoryWorkshop,
	}

	_, err := svc.CreateEvent(context.Background(), dto)
	if err != domain.ErrEventDateRequired {
		t.Fatalf("expected ErrEventDateRequired, got %v", err)
	}
}

func TestEventService_CreateEvent_InvalidCategory(t *testing.T) {
	mockRepo := &mockEventRepository{}
	svc := service.NewEventService(mockRepo)

	eventDate, _ := time.Parse(time.RFC3339, "2026-11-15T09:00:00+07:00")
	dto := domain.CreateEventDTO{
		Title:       "Go Backend Workshop",
		Description: "Valid description",
		EventDate:   eventDate,
		Category:    domain.EventCategory("unknown_category"),
	}

	_, err := svc.CreateEvent(context.Background(), dto)
	if err != domain.ErrInvalidCategory {
		t.Fatalf("expected ErrInvalidCategory, got %v", err)
	}
}
