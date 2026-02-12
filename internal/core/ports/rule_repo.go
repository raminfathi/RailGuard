package ports

import "railguard/internal/core/domain"

// RuleRepository defines the interface for fetching brake and safety rules.
type RuleRepository interface {
	GetMaxSpeed(slope int, brakePercentage int) (int, error)
	// New method to fetch all danger rules
	GetAllDangerRules() ([]domain.DangerRule, error)
}
