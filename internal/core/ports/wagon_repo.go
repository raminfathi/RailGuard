package ports

import "railguard/internal/core/domain"

// WagonRepository defines the interface for interacting with wagon data.
// This allows us to swap the database implementation without changing the core logic.
type WagonRepository interface {
	// GetWagonByNumber finds a wagon specification based on its 6-digit number.
	// Since data is stored as ranges, it checks if the number falls within a range.
	GetWagonByNumber(number int) (*domain.Wagon, error)
}
