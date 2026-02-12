package report

import (
	"fmt"
	"railguard/internal/core/domain"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// GenerateBrakeLicense creates a PDF file with the train safety report.
func (g *PDFGenerator) GenerateBrakeLicense(train *domain.Train, res *domain.CalculationResult, info domain.TripInfo) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// --- 1. Header Section ---
	pdf.SetFont("Arial", "B", 16)
	// CellFormat args: width, height, text, border, ln, align, fill, link, linkStr
	pdf.CellFormat(190, 10, "RAILWAY BRAKE LICENSE & SAFETY CERTIFICATE", "0", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(190, 8, "Islamic Republic of Iran Railways", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// --- 2. Trip Information ---
	// For simple key-values, we use simple Cells
	pdf.SetFont("Arial", "B", 10)

	// Row 1
	pdf.Cell(30, 8, "Date:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 8, time.Now().Format("2006-01-02"))

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Time:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 8, time.Now().Format("15:04:05"))
	pdf.Ln(8)

	// Row 2
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Train No:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 8, info.TrainNumber)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Loco No:")
	pdf.SetFont("Arial", "", 10)

	locoNum := "-"
	if len(train.Locomotives) > 0 {
		locoNum = fmt.Sprintf("%d", train.Locomotives[0].Number)
	}
	pdf.Cell(60, 8, locoNum)
	pdf.Ln(8)

	// Row 3
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Origin:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 8, info.Origin)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Destination:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 8, info.Destination)
	pdf.Ln(8)

	// Row 4
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Driver:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 8, info.DriverName)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Train Boss:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 8, info.TrainBossName)
	pdf.Ln(15)

	// --- 3. Technical Data Table ---
	w := []float64{40, 35, 35, 40, 40} // Column widths

	// Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(220, 220, 220) // Gray background
	headers := []string{"Total Axles", "Total Wagons", "Total Weight (t)", "Total Length (m)", "Max Speed (km/h)"}

	for i, header := range headers {
		// CellFormat for Header (with Fill)
		pdf.CellFormat(w[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1) // Move to next line

	// Table Body
	pdf.SetFont("Arial", "", 12)

	// CellFormat for Body (No Fill)
	pdf.CellFormat(w[0], 12, fmt.Sprintf("%d", train.AxleCount), "1", 0, "C", false, 0, "")
	pdf.CellFormat(w[1], 12, fmt.Sprintf("%d", len(train.Wagons)), "1", 0, "C", false, 0, "")
	pdf.CellFormat(w[2], 12, fmt.Sprintf("%.2f", train.TotalWeight), "1", 0, "C", false, 0, "")
	pdf.CellFormat(w[3], 12, fmt.Sprintf("%.2f", train.TotalLength), "1", 0, "C", false, 0, "")

	// Highlight Speed (Bold)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(w[4], 12, fmt.Sprintf("%d", res.MaxSpeed), "1", 0, "C", false, 0, "")
	pdf.Ln(20)

	// --- 4. Brake Specifics ---
	pdf.SetFont("Arial", "B", 11)

	// FIX: Use CellFormat instead of Cell for alignment
	pdf.CellFormat(0, 10, "BRAKE PERFORMANCE DETAILS:", "0", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 8, fmt.Sprintf("Total Brake Weight:  %.2f tons", train.TotalBrake))
	pdf.Ln(6)
	pdf.Cell(50, 8, fmt.Sprintf("Brake Percentage:    %d %%", res.BrakePercentage))
	pdf.Ln(6)

	status := "REJECTED"
	if res.IsSafe {
		status = "ACCEPTED (SAFE TO DEPART)"
		pdf.SetTextColor(0, 128, 0) // Green
	} else {
		pdf.SetTextColor(255, 0, 0) // Red
	}
	pdf.Cell(50, 8, fmt.Sprintf("Final Status:        %s", status))
	pdf.SetTextColor(0, 0, 0) // Reset color
	pdf.Ln(20)

	// --- 5. Signatures ---
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 5, "I certify that the brake test has been performed correctly and the train is safe.", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Signature Lines
	y := pdf.GetY()
	pdf.Line(20, y, 80, y)   // Line 1
	pdf.Line(120, y, 180, y) // Line 2

	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(90, 5, "Train Examiner / Technical Officer", "0", 0, "C", false, 0, "")
	pdf.CellFormat(90, 5, "Train Boss / Station Master", "0", 1, "C", false, 0, "")

	// Output to file
	filename := fmt.Sprintf("BrakeLicense_%s.pdf", info.TrainNumber)
	return pdf.OutputFileAndClose(filename)
}
