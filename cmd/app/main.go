package main

import (
	"log"
	"os"
	"path/filepath"
	"railguard/internal/adapter/storage/sqlite"
	"railguard/internal/core/services"
	"railguard/internal/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	// 1. FIX: Initialize Fyne App with a Unique ID
	// This ID is required for Preferences and Storage API on Android.
	myApp := app.NewWithID("com.ramin.railguard")

	// 2. Database Path Logic (Android/Desktop)
	var dbPath string
	storageRoot := myApp.Storage().RootURI()

	if storageRoot != nil && storageRoot.Scheme() == "file" {
		// Android / Mobile path
		dbPath = filepath.Join(storageRoot.Path(), "railguard.db")
	} else {
		// Desktop path
		dbPath = "./railguard.db"
	}

	// 3. Initialize Repositories
	// The Repository will now create the file if it doesn't exist.
	wagonRepo, err := sqlite.NewWagonRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize Wagon Repository: %v", err)
	}

	ruleRepo := sqlite.NewRuleRepository(dbPath)

	// 4. Initialize Services
	brakeCalculator := services.NewBrakeCalculatorService(ruleRepo)
	safetyValidator, err := services.NewSafetyValidatorService(ruleRepo)
	if err != nil {
		log.Fatalf("Failed to initialize Safety Validator: %v", err)
	}

	// 5. Initialize UI
	application := ui.NewApp(wagonRepo, brakeCalculator, safetyValidator)

	// Inject the app instance with the correct ID
	application.FyneApp = myApp

	// 6. Run
	if application.MainWindow == nil {
		log.Fatal("Main Window is nil")
	}
	// os.Setenv("FYNE_SCALE", "0.9")
	application.MainWindow.ShowAndRun()

	os.Exit(0)
}
