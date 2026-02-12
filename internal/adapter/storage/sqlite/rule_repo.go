package sqlite

import (
	"database/sql"
	"fmt"
	"railguard/internal/core/domain"
)

type SQLiteRuleRepo struct {
	db *sql.DB
}

func NewSQLiteRuleRepo(db *sql.DB) *SQLiteRuleRepo {
	repo := &SQLiteRuleRepo{db: db}
	repo.initTable()
	repo.seedRules()
	return repo
}

func (r *SQLiteRuleRepo) initTable() {
	query := `CREATE TABLE IF NOT EXISTS danger_rules (
		code_a TEXT,
		code_b TEXT,
		status TEXT,
		PRIMARY KEY (code_a, code_b)
	);`
	r.db.Exec(query)
}

// GetMaxSpeed (Already implemented in main, but moved here for structure)
func (r *SQLiteRuleRepo) GetMaxSpeed(slope, brakePercent int) (int, error) {
	// Simple Logic (Can be moved to DB later)
	if brakePercent < 40 {
		return 0, nil
	}
	if slope > 20 {
		return 40, nil
	}
	return 60, nil
}

// GetAllDangerRules fetches the matrix from DB
func (r *SQLiteRuleRepo) GetAllDangerRules() ([]domain.DangerRule, error) {
	rows, err := r.db.Query("SELECT code_a, code_b, status FROM danger_rules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []domain.DangerRule
	for rows.Next() {
		var r domain.DangerRule
		if err := rows.Scan(&r.CodeA, &r.CodeB, &r.Status); err != nil {
			continue
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// seedRules inserts standard railway compatibility rules (Simplified UIC)
func (r *SQLiteRuleRepo) seedRules() {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM danger_rules").Scan(&count)
	if count > 0 {
		return
	} // Already seeded

	fmt.Println("Seeding Dangerous Goods Matrix...")

	// (-) = Forbidden / ممنوع
	// (+) = Allowed / مجاز
	// (1) = Separation 1 wagon / یک واگن فاصله

	rules := []domain.DangerRule{
		// Class 1 (Explosives) vs Others
		{"1", "1", "+"}, {"1", "2", "-"}, {"1", "3", "-"}, {"1", "4", "-"}, {"1", "5", "-"}, {"1", "6", "-"}, {"1", "7", "-"}, {"1", "8", "-"},

		// Class 2 (Gases)
		{"2", "1", "-"}, {"2", "2", "+"}, {"2", "3", "+"}, {"2", "4", "+"}, {"2", "5", "1"}, {"2", "6", "+"}, {"2", "7", "+"}, {"2", "8", "+"},

		// Class 3 (Flammable Liquids)
		{"3", "1", "-"}, {"3", "2", "+"}, {"3", "3", "+"}, {"3", "4", "+"}, {"3", "5", "1"}, {"3", "6", "+"}, {"3", "7", "+"}, {"3", "8", "+"},

		// Class 4 (Flammable Solids)
		{"4", "1", "-"}, {"4", "2", "+"}, {"4", "3", "+"}, {"4", "4", "+"}, {"4", "5", "1"}, {"4", "6", "+"}, {"4", "7", "+"}, {"4", "8", "+"},

		// Class 5 (Oxidizing)
		{"5", "1", "-"}, {"5", "2", "1"}, {"5", "3", "1"}, {"5", "4", "1"}, {"5", "5", "+"}, {"5", "6", "1"}, {"5", "7", "+"}, {"5", "8", "1"},

		// Class 6 (Toxic)
		{"6", "1", "-"}, {"6", "2", "+"}, {"6", "3", "+"}, {"6", "4", "+"}, {"6", "5", "1"}, {"6", "6", "+"}, {"6", "7", "+"}, {"6", "8", "+"},

		// Class 8 (Corrosive)
		{"8", "1", "-"}, {"8", "2", "+"}, {"8", "3", "+"}, {"8", "4", "+"}, {"8", "5", "1"}, {"8", "6", "+"}, {"8", "7", "+"}, {"8", "8", "+"},
	}

	stmt, _ := r.db.Prepare("INSERT OR REPLACE INTO danger_rules VALUES (?, ?, ?)")
	defer stmt.Close()

	for _, rule := range rules {
		stmt.Exec(rule.CodeA, rule.CodeB, rule.Status)
		// Insert reverse rule too (Symmetric)
		if rule.CodeA != rule.CodeB {
			stmt.Exec(rule.CodeB, rule.CodeA, rule.Status)
		}
	}
}
