package services

import (
	"math"
	"railguard/internal/core/domain"
	"railguard/internal/core/ports"
)

// BrakeCalculatorService handles the core logic for train brake calculations.
type BrakeCalculatorService struct {
	ruleRepo ports.RuleRepository
}

// NewBrakeCalculatorService creates a new instance of the service.
func NewBrakeCalculatorService(repo ports.RuleRepository) *BrakeCalculatorService {
	return &BrakeCalculatorService{ruleRepo: repo}
}

// CalculateTrainParameters computes the total weight, brake weight, and validation.
// FIX: Input type changed to []domain.SelectedWagon
func (s *BrakeCalculatorService) CalculateTrainParameters(locos []domain.Locomotive, wagons []domain.SelectedWagon, slope int) (*domain.CalculationResult, *domain.Train, error) {

	train := &domain.Train{
		Locomotives: locos,
		Wagons:      wagons,
		TotalWeight: 0,
		TotalBrake:  0,
		TotalLength: 0,
		AxleCount:   0,
	}
	// 1. Calculate Locomotives
	for _, loco := range locos {
		train.TotalWeight += loco.Weight
		train.TotalBrake += loco.BrakeWeight
		train.TotalLength += 20 // Approx length if not specified, or add Length to struct
		train.AxleCount += 6    // Usually 6 axles for main line locos
	}

	// 2. Calculate Wagons (Existing logic)
	for _, w := range wagons {
		train.TotalWeight += w.EffectiveWeight
		train.TotalBrake += w.EffectiveBrakeWeight
		train.TotalLength += w.WagonSpec.Length
		train.AxleCount += w.WagonSpec.Axles
	}

	if train.TotalWeight == 0 {
		return &domain.CalculationResult{
			IsSafe:  false,
			Message: "Train weight is zero.",
		}, train, nil
	}

	// 2. Calculate Brake Percentage formula: (TotalBrake / TotalWeight) * 100
	rawPercentage := (train.TotalBrake / train.TotalWeight) * 100
	brakePercentage := int(math.Floor(rawPercentage)) // Round down to be safe

	// 3. Get Max Speed from Rules (Database)
	maxSpeed, err := s.ruleRepo.GetMaxSpeed(slope, brakePercentage)
	if err != nil {
		return nil, nil, err
	}

	// 4. Determine Safety
	result := &domain.CalculationResult{
		BrakePercentage: brakePercentage,
		MaxSpeed:        maxSpeed,
	}

	if maxSpeed > 0 {
		result.IsSafe = true
		result.Message = "Train is safe to depart."
	} else {
		result.IsSafe = false
		result.Message = "Brake percentage is insufficient for this slope."
	}

	return result, train, nil
}
