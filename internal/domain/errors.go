package domain

import "errors"

var (
	// ErrEventNotFound is returned when an event is not found
	ErrEventNotFound = errors.New("event not found")

	// ErrInvalidCategory is returned when category is not supported
	ErrInvalidCategory = errors.New("invalid event category: must be workshop, competition, seminar, or community")

	// ErrTitleRequired is returned when title is empty
	ErrTitleRequired = errors.New("event title is required")

	// ErrEventDateRequired is returned when event date is zero
	ErrEventDateRequired = errors.New("event date is required")
)
