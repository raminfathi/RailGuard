package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"railguard/internal/core/domain"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WagonRepository struct {
	db *sql.DB
}

func NewWagonRepository(dbPath string) (*WagonRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &WagonRepository{db: db}
	if err := repo.initDB(); err != nil {
		return nil, err
	}

	// Seed Data if empty
	repo.seedDefaultData()
	repo.initHistoryTable()
	return repo, nil
}

func (r *WagonRepository) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS wagons (
		number INTEGER PRIMARY KEY,
		type TEXT,
		axles INTEGER,
		weight_empty REAL,
		weight_loaded REAL,
		max_capacity REAL,
		load_volume REAL,
		brake_weight_empty REAL,
		brake_weight_loaded REAL,
		length REAL,
		load_length REAL,
		load_width REAL,
		floor_height REAL,
		internal_height REAL,
		bogie_pivot_dist REAL,
		wheel_diameter REAL,
		riv_code TEXT,
		manufacturer TEXT,
		year TEXT,
		bogie_type TEXT,
		bearing_type TEXT,
		spring_type TEXT,
		hand_brake_type TEXT,
		hand_brake_weight REAL,
		control_valve TEXT,
		brake_cylinder TEXT,
		coupling_type TEXT
	);`
	_, err := r.db.Exec(query)
	return err
}

// GetWagonByNumber retrieves full details
func (r *WagonRepository) GetWagonByNumber(number int) (*domain.Wagon, error) {
	query := `SELECT * FROM wagons WHERE number = ?`
	row := r.db.QueryRow(query, number)

	var w domain.Wagon
	// Note: The order must match the CREATE TABLE columns EXACTLY
	err := row.Scan(
		&w.Number, &w.Type, &w.Axles,
		&w.WeightEmpty, &w.WeightLoaded, &w.MaxCapacity, &w.LoadVolume,
		&w.BrakeWeightEmpty, &w.BrakeWeightLoaded,
		&w.Length, &w.LoadLength, &w.LoadWidth, &w.FloorHeight, &w.InternalHeight,
		&w.BogiePivotDistance, &w.WheelDiameter,
		&w.RIVCode, &w.Manufacturer, &w.Year, &w.BogieType, &w.BearingType, &w.SpringType,
		&w.HandBrakeType, &w.HandBrakeWeight, &w.ControlValveType, &w.BrakeCylinderType, &w.CouplingType,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// seedDefaultData adds the Excel data provided by user
func (r *WagonRepository) seedDefaultData() {
	// Check if data exists
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM wagons").Scan(&count)
	if count > 0 {
		return
	}

	fmt.Println("Seeding Database with Excel Data...")

	// Helper to insert a range of wagons
	insertRange := func(from, to int, w domain.Wagon) {
		stmt, _ := r.db.Prepare(`INSERT OR IGNORE INTO wagons VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
		defer stmt.Close()

		for i := from; i <= to; i++ {
			stmt.Exec(
				i, w.Type, w.Axles,
				w.WeightEmpty, w.WeightLoaded, w.MaxCapacity, w.LoadVolume,
				w.BrakeWeightEmpty, w.BrakeWeightLoaded,
				w.Length, w.LoadLength, w.LoadWidth, w.FloorHeight, w.InternalHeight,
				w.BogiePivotDistance, w.WheelDiameter,
				w.RIVCode, w.Manufacturer, w.Year, w.BogieType, w.BearingType, w.SpringType,
				w.HandBrakeType, w.HandBrakeWeight, w.ControlValveType, w.BrakeCylinderType, w.CouplingType,
			)
		}
	}

	// 1. Car Carrier (Romania) - 140001 to 140168
	w1 := domain.Wagon{
		Type: "ویژه حمل خودرو", Axles: 3, Length: 27.0, LoadLength: 25.7, LoadWidth: 2.68,
		WeightEmpty: 25.7, WeightLoaded: 55.8, MaxCapacity: 31.6, LoadVolume: 0,
		BrakeWeightEmpty: 36, BrakeWeightLoaded: 36,
		RIVCode: "Hccrrs", Manufacturer: "رومانی", Year: "1385",
		BearingType: "سیلندریکال(80*240*120)", BogieType: "تک محور", SpringType: "قوس منفی 120",
		HandBrakeType: "پیچی جانبی", HandBrakeWeight: 20.1,
		ControlValveType: "KE1CSL", BrakeCylinderType: "DRV2AT-600", CouplingType: "زنجیری",
		FloorHeight: 1.17, InternalHeight: 4.35, WheelDiameter: 840,
	}
	insertRange(140001, 140168, w1)

	// 2. Car Carrier (Serbia) - 140501 to 140550
	w2 := w1 // Copy basic
	w2.Manufacturer = "صربستان"
	w2.Year = "88-1386"
	w2.BogieType = "---"
	w2.SpringType = "قوس صفر"
	w2.HandBrakeWeight = 19.7
	w2.ControlValveType = "KE1CSL-KE2CSL"
	w2.BrakeCylinderType = "DRV2A-450"
	w2.BrakeWeightEmpty = 31
	w2.BrakeWeightLoaded = 31
	w2.WeightEmpty = 24.84
	w2.WeightLoaded = 46
	w2.MaxCapacity = 26
	w2.LoadVolume = 20
	w2.Length = 23.8
	w2.InternalHeight = 0
	w2.FloorHeight = 1.25
	insertRange(140501, 140550, w2)

	// 3. Covered Wagon (Gas - Germany) - 147001 to 147600
	w3 := domain.Wagon{
		Type: "مسقف (Gas)", Axles: 4, Length: 16.8, LoadLength: 15.3, LoadWidth: 2.64,
		WeightEmpty: 23, WeightLoaded: 90, MaxCapacity: 67, LoadVolume: 105,
		BrakeWeightEmpty: 26, BrakeWeightLoaded: 53,
		RIVCode: "Gas", Manufacturer: "آلمان - ایران", Year: "1364",
		BearingType: "سیلندریکال", BogieType: "H655", SpringType: "قوس منفی 120",
		HandBrakeType: "پیچی در ایوان", HandBrakeWeight: 24,
		ControlValveType: "KE1CSL", BrakeCylinderType: "DRV2A-600", CouplingType: "یونی کوپلر",
		BogiePivotDistance: 11.5, FloorHeight: 1.2, InternalHeight: 3.07, WheelDiameter: 920,
	}
	insertRange(147001, 147600, w3)

	// You can add more ranges here based on your excel...
}
func (r *WagonRepository) GetDB() *sql.DB {
	return r.db
}

type HistoryItem struct {
	ID          int
	TrainNumber string
	DriverName  string
	CreatedAt   time.Time
	Slope       int
	TotalWeight float64
	MaxSpeed    int
	Locos       []domain.Locomotive
	Wagons      []domain.SelectedWagon
}

// Initialize the history table
func (r *WagonRepository) initHistoryTable() {
	query := `
	CREATE TABLE IF NOT EXISTS train_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		train_number TEXT,
		driver_name TEXT,
		created_at DATETIME,
		slope INTEGER,
		total_weight REAL,
		max_speed INTEGER,
		locos_json TEXT,
		wagons_json TEXT
	);`
	_, err := r.db.Exec(query)
	if err != nil {
		fmt.Println("Error creating history table:", err)
	}
}

// SaveTrainComposition saves the current setup to DB
func (r *WagonRepository) SaveTrainComposition(h HistoryItem) error {
	locosBytes, _ := json.Marshal(h.Locos)
	wagonsBytes, _ := json.Marshal(h.Wagons)

	query := `INSERT INTO train_history (train_number, driver_name, created_at, slope, total_weight, max_speed, locos_json, wagons_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, h.TrainNumber, h.DriverName, time.Now(), h.Slope, h.TotalWeight, h.MaxSpeed, string(locosBytes), string(wagonsBytes))
	return err
}

// GetAllHistory retrieves the list of saved trains
func (r *WagonRepository) GetAllHistory() ([]HistoryItem, error) {
	rows, err := r.db.Query("SELECT id, train_number, driver_name, created_at, slope, total_weight, max_speed, locos_json, wagons_json FROM train_history ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []HistoryItem
	for rows.Next() {
		var h HistoryItem
		var locosJson, wagonsJson string

		err := rows.Scan(&h.ID, &h.TrainNumber, &h.DriverName, &h.CreatedAt, &h.Slope, &h.TotalWeight, &h.MaxSpeed, &locosJson, &wagonsJson)
		if err != nil {
			continue
		}

		// Unmarshal JSON back to Go structs
		json.Unmarshal([]byte(locosJson), &h.Locos)
		json.Unmarshal([]byte(wagonsJson), &h.Wagons)

		history = append(history, h)
	}
	return history, nil
}
