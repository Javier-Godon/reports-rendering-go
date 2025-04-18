package xlsx

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type CpuSystemUsageData struct {
	CPU      string
	AvgUsage float64
	MaxUsage float64
	MinUsage float64
}

type CpuSystemUsageReport struct {
	DateFrom int64
	DateTo   int64
	Data     []CpuSystemUsageData
}

func (r *CpuSystemUsageReport) Render(file *excelize.File) error {
	sheetName := "CPU System Usage"
	_, err := file.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w", err)
	}

	boldFont := &excelize.Font{Bold: true}

	if err := file.SetColWidth(sheetName, "A", "A", 15); err != nil {
		return fmt.Errorf("failed to set column width: %w", err)
	}
	if err := file.SetColWidth(sheetName, "B", "D", 18); err != nil {
		return fmt.Errorf("failed to set column width: %w", err)
	}

	headers := []string{"CPU", "Average Usage (%)", "Max Usage (%)", "Min Usage (%)"}
	for colNum, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colNum+1, 1)
		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			return fmt.Errorf("failed to set cell value for header: %w", err)
		}
		styleID, _ := file.NewStyle(&excelize.Style{Font: boldFont})
		if err := file.SetCellStyle(sheetName, cell, cell, styleID); err != nil {
			return fmt.Errorf("failed to set cell style for header: %w", err)
		}
	}

	dataCellStyleID, err := file.NewStyle(&excelize.Style{
		NumFmt: 10,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create data cell style: %w", err)
	}

	textCellStyleID, err := file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create text cell style: %w", err)
	}

	for rowNum, usage := range r.Data {
		row := rowNum + 2

		cellCPU, _ := excelize.CoordinatesToCellName(1, row)
		if err := file.SetCellValue(sheetName, cellCPU, usage.CPU); err != nil {
			return fmt.Errorf("failed to set CPU: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cellCPU, cellCPU, textCellStyleID); err != nil {
			return fmt.Errorf("failed to set style for CPU: %w", err)
		}

		cellAvg, _ := excelize.CoordinatesToCellName(2, row)
		if err := file.SetCellValue(sheetName, cellAvg, usage.AvgUsage); err != nil {
			return fmt.Errorf("failed to set AvgUsage: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cellAvg, cellAvg, dataCellStyleID); err != nil {
			return fmt.Errorf("failed to set style for AvgUsage: %w", err)
		}

		cellMax, _ := excelize.CoordinatesToCellName(3, row)
		if err := file.SetCellValue(sheetName, cellMax, usage.MaxUsage); err != nil {
			return fmt.Errorf("failed to set MaxUsage: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cellMax, cellMax, dataCellStyleID); err != nil {
			return fmt.Errorf("failed to set style for MaxUsage: %w", err)
		}

		cellMin, _ := excelize.CoordinatesToCellName(4, row)
		if err := file.SetCellValue(sheetName, cellMin, usage.MinUsage); err != nil {
			return fmt.Errorf("failed to set MinUsage: %w", err)
		}
		if err := file.SetCellStyle(sheetName, cellMin, cellMin, dataCellStyleID); err != nil {
			log.Printf("failed to set style for MinUsage: %v", err)
		}
	}

	chartRowOffset := len(r.Data) + 4
	chartCell := fmt.Sprintf("A%d", chartRowOffset)

	if err := file.AddChart(sheetName, chartCell, &excelize.Chart{
		Type: excelize.Col3DClustered,
		Title: []excelize.RichTextRun{
			{Text: "CPU Average Usage"},
		},
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
				Name:       "Average Usage (%)",
				Values:     fmt.Sprintf("'%s'!B2:B%d", sheetName, len(r.Data)+1),
				Categories: fmt.Sprintf("'%s'!A2:A%d", sheetName, len(r.Data)+1),
				Line:       excelize.ChartLine{Width: 2},
			},
		},
		PlotArea: excelize.ChartPlotArea{
			ShowVal:         true,
			ShowCatName:     false,
			ShowSerName:     false,
			ShowLeaderLines: false,
		},
		Legend: excelize.ChartLegend{
			Position: "top",
		},
		Dimension: excelize.ChartDimension{
			Width:  960,
			Height: 560,
		},
	}); err != nil {
		return fmt.Errorf("failed to add chart: %w", err)
	}

	return nil
}
