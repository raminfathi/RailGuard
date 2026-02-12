package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

var (
	// Palette
	white       = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	silver      = color.RGBA{R: 200, G: 200, B: 200, A: 255} // NEW: Visible Silver for Empty Wagons
	railBlue    = color.RGBA{R: 0, G: 120, B: 215, A: 255}   // Brighter Blue for Loaded Wagons
	darkOverlay = color.RGBA{R: 20, G: 25, B: 40, A: 230}    // Dark Background
	orange      = color.RGBA{R: 255, G: 140, B: 0, A: 255}
)

func (m MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return darkOverlay

	case theme.ColorNameForeground:
		return white

	// --- BUTTON COLORS ---
	case theme.ColorNameButton:
		// Standard Button -> Used for "Empty Wagon"
		// Now it returns Silver/Light Gray so it pops on dark background
		return silver

	case theme.ColorNamePrimary:
		// Primary Button -> Used for "Loaded Wagon"
		return railBlue

	case theme.ColorNameInputBackground:
		return color.RGBA{R: 50, G: 60, B: 80, A: 255}

	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 160, G: 160, B: 160, A: 255}

	case theme.ColorNameScrollBar:
		return orange

	default:
		return theme.DarkTheme().Color(name, variant)
	}
}

func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 13
	}
	return theme.DefaultTheme().Size(name)
}
