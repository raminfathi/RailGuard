package domain

// Locomotive represents the engine of the train.
type Locomotive struct {
	ID          string  `json:"id"`           // e.g., "GM-1", "Alstom"
	Number      int     `json:"number"`       // e.g., 2065
	Weight      float64 `json:"weight"`       // Weight in tons (e.g., 120 tons)
	BrakeWeight float64 `json:"brake_weight"` // Brake power in tons
	IsHot       bool    `json:"is_hot"`       // True = Active (Pulling), False = Dead (Towed)
}
