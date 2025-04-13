package xlsx

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// CpuUsageData holds the CPU usage information for a single CPU.
type CpuUsageData struct {
	CPU      string
	AvgUsage float64
	MaxUsage float64
	MinUsage float64
}

// CpuSystemUsageReport holds the overall CPU system usage report data.
type CpuSystemUsageReport struct {
	DateFrom int64
	DateTo   int64
	Data     []CpuUsageData
}

// Render creates an XLSX report of CPU system usage using excelize.
func (r *CpuSystemUsageReport) Render() ([]byte, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Error closing file:", err) // Log the error, but don't stop execution.
		}
	}()
	// Create a new sheet.
	sheetName := "CPU Usage"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}
	file.SetActiveSheet(index)

	// --- Formatting ---
	boldFont := &excelize.Font{Bold: true}

	// Set column widths.  Excelize uses different units, so these are approximations.
	if err := file.SetColWidth(sheetName, "A", "A", 15); err != nil {
		return nil, fmt.Errorf("failed to set column width: %w", err)
	}
	if err := file.SetColWidth(sheetName, "B", "D", 18); err != nil {
		return nil, fmt.Errorf("failed to set column width: %w", err)
	}

	// --- Headers ---
	var styleID int // Declare styleID to be used later
	headers := []string{"CPU", "Average Usage (%)", "Max Usage (%)", "Min Usage (%)"}
	for colNum, header := range headers {
		cell, err := excelize.CoordinatesToCellName(colNum+1, 1) // Excelize is 1-based
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinates: %w", err)
		}
		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("failed to set cell value for header: %w", err)
		}
		styleID, err := file.NewStyle(&excelize.Style{Font: boldFont})
		if err != nil {
			return nil, fmt.Errorf("failed to create style for header: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cell, cell, styleID); err != nil {
			return nil, fmt.Errorf("failed to set cell style for header: %w", err)
		}
	}

	// --- Write Data ---
	for rowNum, usage := range r.Data {
		row := rowNum + 2 // Start from row 2 (1-based)
		cellCPU, err := excelize.CoordinatesToCellName(1, row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinates: %w", err)
		}
		if err := file.SetCellValue(sheetName, cellCPU, usage.CPU); err != nil {
			return nil, fmt.Errorf("failed to set cell value for CPU: %w", err)
		}

		cellAvg, err := excelize.CoordinatesToCellName(2, row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinates: %w", err)
		}
		if err := file.SetCellValue(sheetName, cellAvg, usage.AvgUsage); err != nil {
			return nil, fmt.Errorf("failed to set cell value for AvgUsage: %w", err)
		}
		styleID, err = file.NewStyle(&excelize.Style{NumFmt: 10}) // Use integer format code (e.g., 10 for two decimal places)
		if err != nil {
			return nil, fmt.Errorf("failed to create style for AvgUsage: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cellAvg, cellAvg, styleID); err != nil {
			return nil, fmt.Errorf("failed to set cell style for AvgUsage: %w", err)
		}

		cellMax, err := excelize.CoordinatesToCellName(3, row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinates: %w", err)
		}
		if err := file.SetCellValue(sheetName, cellMax, usage.MaxUsage); err != nil {
			return nil, fmt.Errorf("failed to set cell value for MaxUsage: %w", err)
		}
		styleID, err = file.NewStyle(&excelize.Style{NumFmt: 10}) // Use integer format code (e.g., 10 for two decimal places)
		if err != nil {
			return nil, fmt.Errorf("failed to create style for MaxUsage: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cellMax, cellMax, styleID); err != nil {
			return nil, fmt.Errorf("failed to set cell style for MaxUsage: %w", err)
		}

		cellMin, err := excelize.CoordinatesToCellName(4, row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinates: %w", err)
		}
		if err := file.SetCellValue(sheetName, cellMin, usage.MinUsage); err != nil {
			return nil, fmt.Errorf("failed to set cell value for MinUsage: %w", err)
		}
		styleID, err = file.NewStyle(&excelize.Style{NumFmt: 10}) // Use integer format code for two decimal places
		if err != nil {
			return nil, fmt.Errorf("failed to create style for MinUsage: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cellMin, cellMin, styleID); err != nil {
			log.Printf("failed to set cell style for MinUsage: %v", err)
		}
		log.Printf("failed to set cell style for MinUsage: %v", err)
		continue
	}

	// --- Optional Bar Chart ---
	if len(r.Data) > 0 {
		chart := excelize.Chart{
			Type: excelize.Col3DClustered,
			Title: []excelize.RichTextRun{
				{Text: "CPU Average Usage"},
			},
			// Removed invalid PlotArea field
			// Removed unsupported Legend field
			XAxis: excelize.ChartAxis{
				Title: []excelize.RichTextRun{
					{Text: "CPU"},
				},
			},
			YAxis: excelize.ChartAxis{
				Title: []excelize.RichTextRun{
					{Text: "Average Usage (%)"},
				},
			},
			Series: []excelize.ChartSeries{
				{
					Name:   "Average Usage (%)",
					Values: fmt.Sprintf("'%s'!B2:B%d", sheetName, len(r.Data)+1),
				},
			},
		}

		if err := file.AddChart(sheetName, fmt.Sprintf("E%d", len(r.Data)+3), &chart); err != nil { // Position the chart
			return nil, fmt.Errorf("failed to add chart: %w", err)
		}
	}

	// Get the byte data of the XLSX file.
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}
	return buffer.Bytes(), nil
}

// func main() {
// 	// Example usage for testing
// 	testData := []CpuUsageData{
// 		{"cpu0", 11.71, 45.34, 2.83},
// 		{"cpu1", 12.31, 37.98, 1.00},
// 	}
// 	report := CpuSystemUsageReport{
// 		DateFrom: time.Now().Add(-time.Hour).Unix(),
// 		DateTo:   time.Now().Unix(),
// 		Data:     testData,
// 	}

// 	excelBytes, err := report.Render()
// 	if err != nil {
// 		log.Fatalf("Error: %v", err) // Use log.Fatalf for a non-recoverable error.
// 	}

// 	err = os.WriteFile("cpu_usage_report_formatted.xlsx", excelBytes, 0644)
// 	if err != nil {
// 		log.Fatalf("Error writing file: %v", err)
// 	}

// 	fmt.Println("Test report generated: cpu_usage_report_formatted.xlsx")
// }
