package domain

// Train represents the assembled train consisting of multiple wagons.
type Train struct {
	Locomotives []Locomotive    `json:"locomotives"` // List of locomotives (usually 1, but can be more for double-heading)
	Wagons      []SelectedWagon `json:"wagons"`
	TotalWeight float64         `json:"total_weight"` // Sum of all wagons' weight
	TotalBrake  float64         `json:"total_brake"`  // Sum of all wagons' brake weight
	TotalLength float64         `json:"total_length"` // Sum of all wagons' length
	AxleCount   int             `json:"axle_count"`   // Total number of axles
}

// CalculationResult holds the final output of the brake calculation.
type CalculationResult struct {
	BrakePercentage int    `json:"brake_percentage"` // Calculated brake percentage
	MaxSpeed        int    `json:"max_speed"`        // Max allowed speed based on rules
	IsSafe          bool   `json:"is_safe"`          // True if the train is allowed to depart
	Message         string `json:"message"`          // Error or success message
}

type DangerRule struct {
	CodeA  string `json:"code_a"`
	CodeB  string `json:"code_b"`
	Status string `json:"status"` // "*", "+", "-", "1", "2"
}
