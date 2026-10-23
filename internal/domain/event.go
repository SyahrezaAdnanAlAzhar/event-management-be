package domain

import (
	"context"
	"time"
)

// EventCategory represents valid event categories
type EventCategory string

const (
	CategoryWorkshop    EventCategory = "workshop"
	CategoryCompetition EventCategory = "competition"
	CategorySeminar     EventCategory = "seminar"
	CategoryCommunity   EventCategory = "community"
)

// IsValid validates whether the category is one of the supported categories
func (c EventCategory) IsValid() bool {
	switch c {
	case CategoryWorkshop, CategoryCompetition, CategorySeminar, CategoryCommunity:
		return true
	default:
		return false
	}
}

// Event represents the core domain model for an event
type Event struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	EventDate   time.Time     `json:"event_date"`
	Category    EventCategory `json:"category"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// CreateEventDTO holds user input for creating an event
type CreateEventDTO struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	EventDate   time.Time     `json:"event_date"`
	Category    EventCategory `json:"category"`
}

// EventRepository defines the data access contract for events
type EventRepository interface {
	Create(ctx context.Context, event *Event) (*Event, error)
}
