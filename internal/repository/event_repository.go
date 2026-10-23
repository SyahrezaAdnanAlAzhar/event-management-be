package repository

import (
	"context"

	"event-management-be/internal/domain"
)

// EventRepository defines the data access contract for events
type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) (*domain.Event, error)
}
