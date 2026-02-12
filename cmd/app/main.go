package main

import (
	"log"

	"railguard/internal/adapter/storage/sqlite"
	"railguard/internal/core/services"
	"railguard/internal/ui"
)

func main() {
	// 1. Initialize Database
	dbPath := "railway.db"

	// Repositories
	wagonRepo, err := sqlite.NewWagonRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Rule Repo (Using the same DB connection from wagonRepo)
	ruleRepo := sqlite.NewSQLiteRuleRepo(wagonRepo.GetDB())

	// 2. Initialize Services
	calculator := services.NewBrakeCalculatorService(ruleRepo)

	validator, err := services.NewSafetyValidatorService(ruleRepo)
	if err != nil {
		log.Fatalf("Failed to create validator: %v", err)
	}

	// 3. Initialize UI (With Custom Theme & Background)
	myApp := ui.NewApp(wagonRepo, calculator, validator)

	// 4. Run
	myApp.Run()
}
