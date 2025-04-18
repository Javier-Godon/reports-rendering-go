package xlsx

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// CpuUsageData holds the CPU usage information for a single CPU.
type CpuUserUsageData struct {
	CPU      string
	AvgUsage float64
	MaxUsage float64
	MinUsage float64
}

// CpuUserUsageReport holds the overall CPU system usage report data.
type CpuUserUsageReport struct {
	DateFrom int64
	DateTo   int64
	Data     []CpuUserUsageData
}

// Render creates an XLSX report of CPU system usage using excelize.
func (r *CpuUserUsageReport) Render(file *excelize.File) error {
	sheetName := "CPU User Usage"
	_, err := file.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w", err)
	}

	boldFont := &excelize.Font{Bold: true}

	// Set column widths
	_ = file.SetColWidth(sheetName, "A", "A", 15)
	_ = file.SetColWidth(sheetName, "B", "D", 18)

	// Headers
	headers := []string{"CPU", "Average Usage (%)", "Max Usage (%)", "Min Usage (%)"}
	headerStyle, err := file.NewStyle(&excelize.Style{Font: boldFont})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w", err)
	}
	for col, title := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		file.SetCellValue(sheetName, cell, title)
		file.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Cell Styles
	numStyle, _ := file.NewStyle(&excelize.Style{
		NumFmt: 10,
		Border: borders(),
	})
	textStyle, _ := file.NewStyle(&excelize.Style{
		Border: borders(),
	})

	// Data
	for i, usage := range r.Data {
		row := i + 2
		_ = setStyledCell(file, sheetName, 1, row, usage.CPU, textStyle)
		_ = setStyledCell(file, sheetName, 2, row, usage.AvgUsage, numStyle)
		_ = setStyledCell(file, sheetName, 3, row, usage.MaxUsage, numStyle)
		_ = setStyledCell(file, sheetName, 4, row, usage.MinUsage, numStyle)
	}

	// Chart
	chartRow := len(r.Data) + 4
	chartCell := fmt.Sprintf("A%d", chartRow)

	err = file.AddChart(sheetName, chartCell, &excelize.Chart{
		Type:  excelize.Col3DClustered,
		Title: []excelize.RichTextRun{{Text: "CPU Average Usage"}},
		XAxis: excelize.ChartAxis{
			Title: []excelize.RichTextRun{{Text: "CPU"}},
		},
		YAxis: excelize.ChartAxis{
			Title: []excelize.RichTextRun{{Text: "Average Usage (%)"}},
		},
		Series: []excelize.ChartSeries{{
			Name:       "Average Usage (%)",
			Values:     fmt.Sprintf("'%s'!B2:B%d", sheetName, len(r.Data)+1),
			Categories: fmt.Sprintf("'%s'!A2:A%d", sheetName, len(r.Data)+1),
			Line:       excelize.ChartLine{Width: 2},
		}},
		PlotArea: excelize.ChartPlotArea{
			ShowVal: true,
		},
		Legend: excelize.ChartLegend{
			Position: "top",
		},
		Dimension: excelize.ChartDimension{
			Width:  960,
			Height: 560,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add chart: %w", err)
	}

	return nil
}

func setStyledCell(file *excelize.File, sheet string, col, row int, value interface{}, style int) error {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	if err := file.SetCellValue(sheet, cell, value); err != nil {
		return err
	}
	return file.SetCellStyle(sheet, cell, cell, style)
}

func borders() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	}
}
