package domain

// TripInfo holds metadata about the journey for the report.
type TripInfo struct {
	TrainNumber   string
	Origin        string
	Destination   string
	DriverName    string
	TrainBossName string // Name of the Train Boss
	Date          string
	Time          string
}
