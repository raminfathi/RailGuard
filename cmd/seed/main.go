package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

const (
	dbPath     = "./railway.db"
	wagonsFile = "./assets/data/m.f.wagon-bari.xlsx"
	brakeFile  = "./assets/data/braek_per.xlsx"
	dangerFile = "./assets/data/dangers.xlsx"
)

func main() {
	// 1. Remove old DB if exists to start fresh
	os.Remove(dbPath)

	// 2. Connect (Create) DB
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	createSchema(db)
	
	// 3. Import Data
	importWagons(db)
	importBrakeRules(db)
	importDangerMatrix(db)

	fmt.Println("\n✅ Database seeded successfully!")
}

func createSchema(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS wagon_specs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			wagon_type TEXT,
			axis_count INTEGER,
			start_number INTEGER,
			end_number INTEGER,
			brake_weight_empty REAL,
			brake_weight_loaded REAL,
			wagon_weight_empty REAL,
			wagon_weight_loaded REAL,
			wagon_length REAL,
			max_capacity REAL
		);`,
		`CREATE TABLE IF NOT EXISTS brake_rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slope INTEGER,
			brake_percentage INTEGER,
			max_speed INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS dangerous_goods (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code_a TEXT,
			code_b TEXT,
			status TEXT
		);`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatalf("Schema creation failed: %v", err)
		}
	}
	fmt.Println("Created Database Schema.")
}

func importWagons(db *sql.DB) {
	fmt.Println("Importing Wagons from Excel...")
	f, err := excelize.OpenFile(wagonsFile)
	if err != nil {
		log.Fatalf("Cannot open wagons file: %v", err)
	}
	defer f.Close()

	// Get the first sheet name
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Cannot read rows: %v", err)
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO wagon_specs 
		(wagon_type, axis_count, start_number, end_number, brake_weight_empty, brake_weight_loaded, 
		wagon_weight_empty, wagon_weight_loaded, wagon_length, max_capacity) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	// Iterate rows (Skip header row 0)
	for i, row := range rows {
		if i == 0 || len(row) < 10 {
			continue 
		}

		// Helper to safely parse numbers
		getFloat := func(idx int) float64 {
			if idx >= len(row) { return 0 }
			val := strings.TrimSpace(row[idx])
			if val == "" || val == "-" { return 0 }
			f, _ := strconv.ParseFloat(val, 64)
			return f
		}
		getInt := func(idx int) int {
			if idx >= len(row) { return 0 }
			val := strings.TrimSpace(row[idx])
			i, _ := strconv.Atoi(val)
			return i
		}

		// Mapping based on your Excel columns:
		// 0: Type, 1: Axis, 2: Start, 3: End, ... 
		// Note: Adjust indices if Excel columns change
		wagonType := row[0]
		axisCount := getInt(1)
		startNum := getInt(2)
		endNum := getInt(3)
		// Assuming brake weights are at index 16 & 17 based on typical layout, 
		// verify against your specific excel if numbers are 0.
		// Based on csv snippet: Type(0), Axis(1), Start(2), End(3)...
		// Let's assume standard order, update indices if needed after first run.
		brakeEmpty := getFloat(16) 
		brakeLoaded := getFloat(17)
		wEmpty := getFloat(21)
		wLoaded := getFloat(22)
		length := getFloat(23)
		capacity := getFloat(20)

		_, err = stmt.Exec(wagonType, axisCount, startNum, endNum, brakeEmpty, brakeLoaded, wEmpty, wLoaded, length, capacity)
		if err != nil {
			log.Printf("Error inserting row %d: %v", i, err)
		}
	}
	tx.Commit()
}

func importBrakeRules(db *sql.DB) {
	fmt.Println("Importing Brake Rules...")
	f, err := excelize.OpenFile(brakeFile)
	if err != nil {
		log.Fatalf("Cannot open brake file: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO brake_rules (slope, brake_percentage, max_speed) VALUES (?, ?, ?)")
	defer stmt.Close()

	// Row 0 contains percentages (Headers)
	// Col 0 contains Slopes
	var percentages []int
	for c, cell := range rows[0] {
		if c == 0 { continue } // Skip first cell (corner)
		val, _ := strconv.Atoi(cell)
		percentages = append(percentages, val)
	}

	for r, row := range rows {
		if r == 0 { continue } // Skip header
		
		slope, _ := strconv.Atoi(row[0]) // First column is slope

		for c, cell := range row {
			if c == 0 { continue } // Skip slope column
			if c-1 >= len(percentages) { break }
			
			speed, _ := strconv.Atoi(cell)
			perc := percentages[c-1]

			_, err = stmt.Exec(slope, perc, speed)
			if err != nil {
				log.Printf("Error inserting rule: %v", err)
			}
		}
	}
	tx.Commit()
}

func importDangerMatrix(db *sql.DB) {
	fmt.Println("Importing Danger Matrix...")
	f, err := excelize.OpenFile(dangerFile)
	if err != nil {
		log.Fatalf("Cannot open danger file: %v", err)
	}
	defer f.Close()

	// Assuming Sheet 1 is the matrix
	rows, err := f.GetRows("Sheet1") 
	if err != nil {
		log.Fatal(err)
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO dangerous_goods (code_a, code_b, status) VALUES (?, ?, ?)")
	defer stmt.Close()

	// Row 0: Headers (Code B)
	var headers []string
	for _, h := range rows[0] {
		headers = append(headers, h)
	}

	for r, row := range rows {
		if r == 0 { continue }
		codeA := row[0]

		for c, cell := range row {
			if c == 0 { continue }
			if c >= len(headers) { break }
			
			codeB := headers[c]
			status := cell
			
			_, err = stmt.Exec(codeA, codeB, status)
			if err != nil {
				log.Printf("Error inserting danger rule: %v", err)
			}
		}
	}
	tx.Commit()
}