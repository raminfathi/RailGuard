package ui

import (
	"errors"
	"fmt"
	"image/color"
	"railguard/internal/adapter/report"
	"railguard/internal/adapter/storage/sqlite" // Import needed for HistoryItem
	"railguard/internal/core/domain"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Decoder maps ported from Wagons.py logic
var wagonTypeMap = map[string]string{
	"1": "Covered", "2": "Open-Short", "3": "Open-High", "4": "Flat", "5": "Tank",
	"6": "Rail Carrier", "7": "Fridge", "8": "Ballast", "9": "Bulk",
}
var axleMap = map[string]string{
	"0": "3 Axles", "1": "2 Axles", "2": "4 Axles", "3": "4 Axles", "4": "4 Axles",
	"8": "6 Axles", "9": "6 Axles",
}

// Analyzes wagon number digits in real-time
func decodeWagonInfo(numStr string) string {
	if len(numStr) == 0 {
		return "..."
	}
	info := ""
	if len(numStr) >= 1 {
		if val, ok := wagonTypeMap[string(numStr[0])]; ok {
			info += "Type: " + val + " | "
		}
	}
	if len(numStr) >= 3 {
		if val, ok := axleMap[string(numStr[2])]; ok {
			info += "Axles: " + val
		}
	}
	return info
}

// Creates the visual guide for the train strip
func (a *App) makeLegend() fyne.CanvasObject {
	makeItem := func(c color.Color, text string) fyne.CanvasObject {
		rect := canvas.NewRectangle(c)
		rect.SetMinSize(fyne.NewSize(20, 20))
		return container.NewHBox(rect, widget.NewLabel(text))
	}

	silverColor := color.RGBA{R: 200, G: 200, B: 200, A: 255}
	blueColor := color.RGBA{R: 0, G: 120, B: 215, A: 255}

	return container.NewHBox(
		widget.NewLabelWithStyle("Legend:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		makeItem(theme.SuccessColor(), "Loco"),
		makeItem(theme.ErrorColor(), "Danger"),
		makeItem(theme.WarningColor(), "Brake Defect"),
		makeItem(silverColor, "Empty Wagon"),
		makeItem(blueColor, "Loaded Wagon"),
		layout.NewSpacer(),
	)
}

func (a *App) makeDashboard() fyne.CanvasObject {
	title := widget.NewLabel("Train Configuration System")
	title.TextStyle = fyne.TextStyle{Bold: true}
	slopeEntry := widget.NewEntry()
	slopeEntry.SetText("10")

	trainObjectsBox := container.NewHBox()
	visualScroll := container.NewHScroll(trainObjectsBox)
	visualScroll.SetMinSize(fyne.NewSize(0, 100))

	// Define refreshVisuals first so we can use it in callbacks
	var refreshVisuals func()
	refreshVisuals = func() {
		trainObjectsBox.Objects = nil

		// 1. Render Locomotives
		for i, l := range a.CurrentLocos {
			currentLoco := l
			idx := i
			btn := widget.NewButton(fmt.Sprintf("🚂 %d", l.Number), func() { a.showLocoActions(currentLoco, idx, refreshVisuals) })
			btn.Importance = widget.SuccessImportance // Green for Loco
			trainObjectsBox.Add(btn)
			trainObjectsBox.Add(container.NewCenter(&canvas.Line{StrokeColor: color.White, StrokeWidth: 2, Position2: fyne.NewPos(10, 0)}))
		}

		// 2. Render Wagons
		for i, w := range a.CurrentTrain {
			currentWagon := w
			idx := i
			label := fmt.Sprintf("🚃 %d", w.WagonSpec.Number)
			if w.HasDangerousGoods {
				label = "⚠️ " + strconv.Itoa(w.WagonSpec.Number)
			}

			btn := widget.NewButton(label, func() { a.openWagonEditForm(currentWagon, idx, refreshVisuals) })

			// Visual status logic
			if w.HasDangerousGoods {
				btn.Importance = widget.DangerImportance // Red
			} else if !w.IsMainBrakeHealthy {
				btn.Importance = widget.WarningImportance // Orange
			} else if !w.IsLoaded {
				btn.Importance = widget.MediumImportance // Silver (Empty)
			} else {
				btn.Importance = widget.HighImportance // Blue (Loaded)
			}
			trainObjectsBox.Add(btn)
			if i < len(a.CurrentTrain)-1 {
				trainObjectsBox.Add(container.NewCenter(&canvas.Line{StrokeColor: color.White, StrokeWidth: 2, Position2: fyne.NewPos(10, 0)}))
			}
		}
		trainObjectsBox.Refresh()
	}

	// --- WAGON INPUT SECTION ---
	wagonNumEntry := widget.NewEntry()
	wagonNumEntry.SetPlaceHolder("6-digit Wagon Number")
	searchFeedback := widget.NewLabel("...")
	addWagonBtn := widget.NewButtonWithIcon("Add to Composition", theme.ContentAddIcon(), nil)
	addWagonBtn.Disable()

	infoBtn := widget.NewButtonWithIcon("Technical Specs", theme.InfoIcon(), func() {
		num, _ := strconv.Atoi(wagonNumEntry.Text)
		w, err := a.WagonRepo.GetWagonByNumber(num)
		if err != nil {
			a.ShowError(err)
			return
		}

		details := fmt.Sprintf(`[ FULL SPECIFICATIONS ]
Wagon No: %d | Class: %s
Manufacturer: %s (%s) | Axles: %d
Length: %.2f m | RIV: %s
--------------------------------------
[ WEIGHTS ]
Empty: %.1f t | Max Load: %.1f t
Full: %.1f t | Volume: %.1f m3
--------------------------------------
[ BRAKING ]
Triple Valve: %s | Auto Brake: %s
Empty Brake: %.1f t | Loaded Brake: %.1f t
Hand Brake: %s (%.1f t)
--------------------------------------
[ MECHANICAL ]
Bogie: %s | Springs: %s
Bearing: %s | Coupling: %s
Wheel: %.0f mm | Pivot Dist: %.2f m`,
			w.Number, w.Type, w.Manufacturer, w.Year, w.Axles, w.Length, w.RIVCode,
			w.WeightEmpty, w.MaxCapacity, w.WeightLoaded, w.LoadVolume,
			w.ControlValveType, w.BrakeCylinderType, w.BrakeWeightEmpty, w.BrakeWeightLoaded,
			w.HandBrakeType, w.HandBrakeWeight,
			w.BogieType, w.SpringType, w.BearingType, w.CouplingType,
			w.WheelDiameter, w.BogiePivotDistance)

		scroll := container.NewVScroll(widget.NewLabel(details))
		scroll.SetMinSize(fyne.NewSize(450, 450))
		dialog.ShowCustom("Wagon Technical Sheet", "Close", scroll, a.MainWindow)
	})

	wagonNumEntry.OnChanged = func(s string) {
		searchFeedback.SetText(decodeWagonInfo(s))
		if len(s) == 6 {
			num, _ := strconv.Atoi(s)
			if _, err := a.WagonRepo.GetWagonByNumber(num); err == nil {
				addWagonBtn.Enable()
				infoBtn.Enable()
			} else {
				addWagonBtn.Disable()
				infoBtn.Disable()
			}
		} else {
			addWagonBtn.Disable()
			infoBtn.Disable()
		}
	}

	addWagonBtn.OnTapped = func() {
		num, _ := strconv.Atoi(wagonNumEntry.Text)
		w, _ := a.WagonRepo.GetWagonByNumber(num)
		a.openWagonEditForm(domain.SelectedWagon{WagonSpec: *w}, -1, refreshVisuals)
		wagonNumEntry.SetText("")
		refreshVisuals()
	}

	// --- LOCOMOTIVE TAB LOGIC (FIXED) ---
	addLocoBtn := widget.NewButtonWithIcon("Add Locomotive", theme.ContentAddIcon(), func() {
		idEntry := widget.NewEntry()
		idEntry.SetPlaceHolder("e.g. GM-12")
		numEntry := widget.NewEntry()
		numEntry.SetPlaceHolder("e.g. 204")
		weightEntry := widget.NewEntry()
		weightEntry.SetText("120")
		hotCheck := widget.NewCheck("Active (Hot)", nil)
		hotCheck.Checked = true

		dialog.ShowForm("Add Locomotive", "Add", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Model:", idEntry),
			widget.NewFormItem("Number:", numEntry),
			widget.NewFormItem("Weight (t):", weightEntry),
			widget.NewFormItem("Status:", hotCheck),
		}, func(ok bool) {
			if ok {
				w, _ := strconv.ParseFloat(weightEntry.Text, 64)
				n, _ := strconv.Atoi(numEntry.Text)
				// Simple logic: if hot, brake weight is approx 80% (simplified)
				br := 0.0
				if hotCheck.Checked {
					br = w * 0.8
				}

				newLoco := domain.Locomotive{
					ID: idEntry.Text, Number: n, Weight: w, BrakeWeight: br, IsHot: hotCheck.Checked,
				}
				a.CurrentLocos = append(a.CurrentLocos, newLoco)
				refreshVisuals()
			}
		}, a.MainWindow)
	})

	// Wrap in VBox to fix sizing issue
	locoTabContent := container.NewVBox(
		widget.NewLabel("Locomotive Management"),
		addLocoBtn,
	)

	// --- CONTROL BUTTONS ---

	// 1. Calculate
	calcBtn := widget.NewButtonWithIcon("CALCULATE", theme.ConfirmIcon(), func() {
		s, _ := strconv.Atoi(slopeEntry.Text)
		a.CurrentSlope = s
		isSafe, msg := a.Validator.ValidateComposition(a.CurrentTrain)
		if !isSafe {
			dialog.ShowError(errors.New(msg), a.MainWindow)
			return
		}

		res, train, _ := a.Calculator.CalculateTrainParameters(a.CurrentLocos, a.CurrentTrain, a.CurrentSlope)
		statusText := "✅ SAFETY PASSED"
		if !res.IsSafe {
			statusText = "❌ SAFETY FAILED"
		}

		dialog.ShowInformation("Result", fmt.Sprintf("%s\nMax Speed: %d km/h\nWeight: %.1f t", statusText, res.MaxSpeed, train.TotalWeight), a.MainWindow)
	})
	calcBtn.Importance = widget.HighImportance

	// 2. Save
	saveBtn := widget.NewButtonWithIcon("SAVE", theme.DocumentSaveIcon(), func() {
		if len(a.CurrentTrain) == 0 {
			return
		}
		s, _ := strconv.Atoi(slopeEntry.Text)
		a.CurrentSlope = s
		res, train, _ := a.Calculator.CalculateTrainParameters(a.CurrentLocos, a.CurrentTrain, a.CurrentSlope)
		a.showSaveDialog(train.TotalWeight, res.MaxSpeed)
	})

	// 3. History
	historyBtn := widget.NewButtonWithIcon("HISTORY", theme.HistoryIcon(), func() {
		a.showHistoryDialog(refreshVisuals)
	})

	// 4. PDF (RESTORED)
	pdfBtn := widget.NewButtonWithIcon("PDF LICENSE", theme.FileIcon(), func() {
		if len(a.CurrentTrain) == 0 {
			return
		}
		tn := widget.NewEntry()
		dr := widget.NewEntry()
		bs := widget.NewEntry()
		org := widget.NewEntry()
		dst := widget.NewEntry()

		dialog.ShowForm("Generate Brake License", "Generate", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Train No:", tn), widget.NewFormItem("Driver:", dr),
			widget.NewFormItem("Train Boss:", bs), widget.NewFormItem("Origin:", org), widget.NewFormItem("Dest:", dst),
		}, func(ok bool) {
			if ok {
				info := domain.TripInfo{TrainNumber: tn.Text, DriverName: dr.Text, TrainBossName: bs.Text, Origin: org.Text, Destination: dst.Text}
				s, _ := strconv.Atoi(slopeEntry.Text)
				a.CurrentSlope = s
				res, train, _ := a.Calculator.CalculateTrainParameters(a.CurrentLocos, a.CurrentTrain, a.CurrentSlope)

				err := report.NewPDFGenerator().GenerateBrakeLicense(train, res, info)
				if err == nil {
					a.ShowInfo("Success", "PDF License Generated Successfully!")
				} else {
					a.ShowError(err)
				}
			}
		}, a.MainWindow)
	})

	// Layout Assembly
	wagonBox := container.NewVBox(
		widget.NewLabel("Wagon Search:"),
		wagonNumEntry,
		container.NewHBox(infoBtn, addWagonBtn),
		searchFeedback,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Locomotives", locoTabContent),
		container.NewTabItem("Wagons", wagonBox),
	)

	// Top Section
	topSection := container.NewVBox(
		title,
		widget.NewForm(widget.NewFormItem("Track Slope (permil):", slopeEntry)),
	)

	// Bottom Section (Buttons)
	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		a.makeLegend(),
		visualScroll,
		widget.NewSeparator(),
		// Arrange buttons nicely
		container.NewHBox(
			historyBtn,
			layout.NewSpacer(),
			calcBtn,
			layout.NewSpacer(),
			saveBtn, pdfBtn, // Group Save and PDF together
		),
	)

	return container.NewBorder(topSection, bottomSection, nil, nil, tabs)
}

// --- NEW HELPER FUNCTIONS FOR SAVE & HISTORY ---

func (a *App) showSaveDialog(currentWeight float64, maxSpeed int) {
	tnEntry := widget.NewEntry()
	tnEntry.SetPlaceHolder("e.g. 4055")
	driverEntry := widget.NewEntry()

	dialog.ShowForm("Save Train Composition", "Save", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Train Number:", tnEntry),
		widget.NewFormItem("Driver Name:", driverEntry),
	}, func(confirm bool) {
		if confirm {
			item := sqlite.HistoryItem{
				TrainNumber: tnEntry.Text,
				DriverName:  driverEntry.Text,
				Slope:       a.CurrentSlope,
				TotalWeight: currentWeight,
				MaxSpeed:    maxSpeed,
				Locos:       a.CurrentLocos,
				Wagons:      a.CurrentTrain,
			}

			if repo, ok := a.WagonRepo.(*sqlite.WagonRepository); ok {
				err := repo.SaveTrainComposition(item)
				if err == nil {
					a.ShowInfo("Success", "Train Saved to History!")
				} else {
					a.ShowError(err)
				}
			}
		}
	}, a.MainWindow)
}

func (a *App) showHistoryDialog(loadCallback func()) {
	repo, ok := a.WagonRepo.(*sqlite.WagonRepository)
	if !ok {
		return
	}

	history, err := repo.GetAllHistory()
	if err != nil {
		a.ShowError(err)
		return
	}

	list := widget.NewList(
		func() int { return len(history) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabelWithStyle("Train No", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Details"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			h := history[i]
			box := o.(*fyne.Container)
			lblTitle := box.Objects[0].(*widget.Label)
			lblDetails := box.Objects[1].(*widget.Label)

			dateStr := h.CreatedAt.Format("2006-01-02 15:04")
			lblTitle.SetText(fmt.Sprintf("Train #%s | Driver: %s", h.TrainNumber, h.DriverName))
			lblDetails.SetText(fmt.Sprintf("%s | Weight: %.0f t | Speed: %d km/h", dateStr, h.TotalWeight, h.MaxSpeed))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selected := history[id]
		dialog.ShowConfirm("Load Train?", fmt.Sprintf("Load Train #%s?\nCurrent unsaved changes will be lost.", selected.TrainNumber), func(b bool) {
			if b {
				a.CurrentLocos = selected.Locos
				a.CurrentTrain = selected.Wagons
				a.CurrentSlope = selected.Slope
				loadCallback()
			}
			list.Unselect(id)
		}, a.MainWindow)
	}

	d := dialog.NewCustom("Saved Trains History", "Close", container.NewMax(list), a.MainWindow)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}

// --- WAGON EDIT FORM ---

func (a *App) openWagonEditForm(wagon domain.SelectedWagon, index int, refresh func()) {
	var d dialog.Dialog
	checkDangerous := widget.NewCheck("Dangerous Goods", nil)
	dangerCodes := []string{"1", "2b", "2a", "2at", "3a", "3bc", "4-1", "4-2", "4-3", "5-1", "5-2", "6-1", "6-1 HCN", "6-2", "7", "8", "9"}
	dangerSelect := widget.NewSelect(dangerCodes, nil)
	dangerSelect.Disable()

	loadRadio := widget.NewRadioGroup([]string{"Empty", "Loaded"}, func(s string) {
		if s == "Empty" {
			checkDangerous.SetChecked(false)
			checkDangerous.Disable()
			dangerSelect.Disable()
		} else {
			checkDangerous.Enable()
		}
	})
	checkDangerous.OnChanged = func(b bool) {
		if b {
			dangerSelect.Enable()
		} else {
			dangerSelect.Disable()
		}
	}

	checkMainBrake := widget.NewCheck("Air Brake Healthy", nil)
	checkHandBrake := widget.NewCheck("Hand Brake Healthy", nil)
	checkHandle := widget.NewCheck("Brake Handle Healthy", nil)

	// Populate existing state
	loadRadio.Selected = "Loaded"
	checkMainBrake.Checked = true
	checkHandBrake.Checked = true
	checkHandle.Checked = true
	if index >= 0 {
		if wagon.IsLoaded {
			loadRadio.Selected = "Loaded"
		} else {
			loadRadio.Selected = "Empty"
		}
		checkDangerous.Checked = wagon.HasDangerousGoods
		dangerSelect.Selected = wagon.DangerousGoodsCode
		checkMainBrake.Checked = wagon.IsMainBrakeHealthy
		checkHandBrake.Checked = wagon.IsHandBrakeHealthy
		checkHandle.Checked = wagon.IsBrakeHandleHealthy
	}

	closeDialog := func() {
		if d != nil {
			d.Hide()
		}
	}
	var extraButtons []fyne.CanvasObject
	if index >= 0 {
		btnDelete := widget.NewButtonWithIcon("Remove", theme.DeleteIcon(), func() {
			a.CurrentTrain = append(a.CurrentTrain[:index], a.CurrentTrain[index+1:]...)
			refresh()
			closeDialog()
		})
		btnDelete.Importance = widget.DangerImportance
		btnLeft := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			if index > 0 {
				a.CurrentTrain[index], a.CurrentTrain[index-1] = a.CurrentTrain[index-1], a.CurrentTrain[index]
				refresh()
				closeDialog()
			}
		})
		btnRight := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			if index < len(a.CurrentTrain)-1 {
				a.CurrentTrain[index], a.CurrentTrain[index+1] = a.CurrentTrain[index+1], a.CurrentTrain[index]
				refresh()
				closeDialog()
			}
		})
		extraButtons = append(extraButtons, btnLeft, btnRight, layout.NewSpacer(), btnDelete)
	}

	saveBtn := widget.NewButtonWithIcon("Apply Changes", theme.DocumentSaveIcon(), func() {
		updated := wagon
		updated.IsLoaded = loadRadio.Selected == "Loaded"
		updated.IsMainBrakeHealthy = checkMainBrake.Checked
		updated.IsHandBrakeHealthy = checkHandBrake.Checked
		updated.IsBrakeHandleHealthy = checkHandle.Checked
		updated.HasDangerousGoods = checkDangerous.Checked
		updated.DangerousGoodsCode = dangerSelect.Selected

		if updated.IsLoaded {
			updated.EffectiveWeight = wagon.WagonSpec.WeightLoaded
			updated.EffectiveBrakeWeight = wagon.WagonSpec.BrakeWeightLoaded
		} else {
			updated.EffectiveWeight = wagon.WagonSpec.WeightEmpty
			updated.EffectiveBrakeWeight = wagon.WagonSpec.BrakeWeightEmpty
		}
		if !updated.IsMainBrakeHealthy {
			updated.EffectiveBrakeWeight = 0
		}

		if index == -1 {
			a.CurrentTrain = append(a.CurrentTrain, updated)
		} else {
			a.CurrentTrain[index] = updated
		}
		refresh()
		closeDialog()
	})

	form := widget.NewForm(
		widget.NewFormItem("Load Status:", loadRadio),
		widget.NewFormItem("Braking Systems:", container.NewVBox(checkMainBrake, checkHandBrake, checkHandle)),
		widget.NewFormItem("Cargo Type:", container.NewVBox(checkDangerous, dangerSelect)),
	)

	content := container.NewVBox(form, widget.NewSeparator(), container.NewHBox(extraButtons...), widget.NewSeparator(), saveBtn)
	d = dialog.NewCustom("Wagon Config #"+strconv.Itoa(wagon.WagonSpec.Number), "Cancel", content, a.MainWindow)
	d.Resize(fyne.NewSize(400, 550))
	d.Show()
}

func (a *App) showLocoActions(loco domain.Locomotive, index int, refresh func()) {
	var d dialog.Dialog
	btnDelete := widget.NewButtonWithIcon("Remove Locomotive", theme.DeleteIcon(), func() {
		a.CurrentLocos = append(a.CurrentLocos[:index], a.CurrentLocos[index+1:]...)
		refresh()
		if d != nil {
			d.Hide()
		}
	})
	btnDelete.Importance = widget.DangerImportance
	info := widget.NewLabel(fmt.Sprintf("Locomotive #%d\nModel: %s\nWeight: %.1f t", loco.Number, loco.ID, loco.Weight))
	d = dialog.NewCustom("Locomotive Options", "Close", container.NewVBox(info, widget.NewSeparator(), btnDelete), a.MainWindow)
	d.Show()
}
