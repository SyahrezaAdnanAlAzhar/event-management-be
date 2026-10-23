package postgres

import (
	"context"
	"fmt"
	"time"

	"event-management-be/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EventRepository implements domain.EventRepository using PostgreSQL pgxpool
type EventRepository struct {
	pool *pgxpool.Pool
}

// NewEventRepository creates a new instance of PostgreSQL EventRepository
func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{pool: pool}
}

// Create inserts a new event into the PostgreSQL database using a parameterized query
func (r *EventRepository) Create(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	query := `
		INSERT INTO events (title, description, event_date, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, event_date, category, created_at, updated_at
	`

	now := time.Now().UTC()
	var (
		id          int64
		title       string
		description *string
		eventDate   time.Time
		category    string
		createdAt   time.Time
		updatedAt   time.Time
	)

	var descParam *string
	if e.Description != "" {
		descParam = &e.Description
	}

	err := r.pool.QueryRow(ctx, query,
		e.Title,
		descParam,
		e.EventDate,
		string(e.Category),
		now,
		now,
	).Scan(
		&id,
		&title,
		&description,
		&eventDate,
		&category,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert event: %w", err)
	}

	result := &domain.Event{
		ID:        id,
		Title:     title,
		EventDate: eventDate,
		Category:  domain.EventCategory(category),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if description != nil {
		result.Description = *description
	}

	return result, nil
}
