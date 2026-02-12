package ui

import (
	"railguard/internal/core/domain"
	"railguard/internal/core/ports"
	"railguard/internal/core/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

type App struct {
	FyneApp    fyne.App
	MainWindow fyne.Window

	WagonRepo  ports.WagonRepository
	Calculator *services.BrakeCalculatorService
	Validator  *services.SafetyValidatorService

	CurrentTrain []domain.SelectedWagon
	CurrentLocos []domain.Locomotive
	CurrentSlope int
}

func NewApp(wRepo ports.WagonRepository, calc *services.BrakeCalculatorService, val *services.SafetyValidatorService) *App {
	// os.Setenv("FYNE_FONT", "./assets/Vazir.ttf")

	myApp := app.New()

	// --- Apply Custom Theme ---
	myApp.Settings().SetTheme(&MyTheme{})

	myWindow := myApp.NewWindow("RailGuard Pro - Train Safety System")

	application := &App{
		FyneApp:      myApp,
		MainWindow:   myWindow,
		WagonRepo:    wRepo,
		Calculator:   calc,
		Validator:    val,
		CurrentSlope: 10,
	}

	dashboard := application.makeDashboard()

	// --- Background Image Setup ---
	// Ensure 'assets/background.jpg' exists in your project directory.
	// bgImage := canvas.NewImageFromResource(resourceBackgroundJpg)
	// bgImage.FillMode = canvas.ImageFillStretch
	// // Set translucency to fade the image so text remains readable (0.70 = 70% visible)
	// bgImage.Translucency = 0.60

	// // Layering: Background at the bottom, Dashboard on top
	// finalLayout := container.NewMax(bgImage, dashboard)

	myWindow.SetContent(dashboard)
	myWindow.Resize(fyne.NewSize(1000, 750))

	return application
}

func (a *App) Run() {
	a.MainWindow.ShowAndRun()
}

func (a *App) ShowError(err error) {
	dialog.ShowError(err, a.MainWindow)
}

func (a *App) ShowInfo(title, message string) {
	dialog.ShowInformation(title, message, a.MainWindow)
}
