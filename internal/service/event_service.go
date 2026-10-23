package service

import (
	"context"
	"strings"

	"event-management-be/internal/domain"
)

// EventService defines business logic operations for events
type EventService interface {
	CreateEvent(ctx context.Context, dto domain.CreateEventDTO) (*domain.Event, error)
}

type eventService struct {
	repo domain.EventRepository
}

// NewEventService creates a new EventService instance
func NewEventService(repo domain.EventRepository) EventService {
	return &eventService{repo: repo}
}

// CreateEvent validates business rules and creates an event
func (s *eventService) CreateEvent(ctx context.Context, dto domain.CreateEventDTO) (*domain.Event, error) {
	// Validate required title
	trimmedTitle := strings.TrimSpace(dto.Title)
	if trimmedTitle == "" {
		return nil, domain.ErrTitleRequired
	}

	// Validate required event date
	if dto.EventDate.IsZero() {
		return nil, domain.ErrEventDateRequired
	}

	// Validate category
	if !dto.Category.IsValid() {
		return nil, domain.ErrInvalidCategory
	}

	event := &domain.Event{
		Title:       trimmedTitle,
		Description: strings.TrimSpace(dto.Description),
		EventDate:   dto.EventDate,
		Category:    dto.Category,
	}

	created, err := s.repo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	return created, nil
}
