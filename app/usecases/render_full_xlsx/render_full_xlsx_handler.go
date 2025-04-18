package render_full_xlsx

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/xuri/excelize/v2"

	pb_system "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_system_usage"
	pb_user "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_user_usage"
	render_xlsx "github.com/Javier-Godon/reports-rendering-go/render/xlsx"

	proto "github.com/Javier-Godon/reports-rendering-go/proto"
)

type RenderFullXlsxHandler struct{}

func NewRenderFullXlsxHandler() *RenderFullXlsxHandler {
	return &RenderFullXlsxHandler{}
}

type RenderResult struct {
	File      *excelize.File
	SheetName string
	Err       error
}

func copySheet(src, dst *excelize.File, sheetName string) {
	rows, err := src.GetRows(sheetName)
	if err != nil {
		log.Printf("Error reading rows from %s: %v", sheetName, err)
		return
	}
	dst.NewSheet(sheetName)
	for rowIdx, row := range rows {
		for colIdx, cell := range row {
			cellRef, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			dst.SetCellValue(sheetName, cellRef, cell)
		}
	}
}

func (handler RenderFullXlsxHandler) Handle(query RenderFullXlsxQuery) (RenderFullXlsxResult, error) {
	address := "localhost:50051"
	clientSingleton := proto.GRPCClientSingleton{}
	client, err := clientSingleton.Instance(address)
	if err != nil {
		log.Fatalf("failed to get gRPC client instance: %v", err)
	}
	defer client.Close()

	ctx := context.Background() // Define a context
	systemUsage, err := client.GetCpuSystemUsage(ctx, int64(query.DateFrom), int64(query.DateTo))
	if err != nil {
		log.Printf("Error getting system usage: %v", err)
	}

	userUsage, err := client.GetCpuUserUsage(ctx, int64(query.DateFrom), int64(query.DateTo))
	if err != nil {
		log.Printf("Error getting user usage: %v", err)
	}

	// Map the gRPC responses to CpuUsageData structs.
	systemUsageData, err := mapGetCpuSystemUsageResponse(systemUsage)
	if err != nil {
		log.Printf("Error converting system usage: %v", err)
	}

	userUsageData, err := mapGetCpuUserUsageResponse(userUsage)
	if err != nil {
		log.Printf("Error converting user usage: %v", err)
	}

	reportSystemData := render_xlsx.CpuSystemUsageReport{
		DateFrom: int64(query.DateFrom),
		DateTo:   int64(query.DateTo),
		Data:     systemUsageData,
	}

	reportUserData := render_xlsx.CpuUserUsageReport{
		DateFrom: int64(query.DateFrom),
		DateTo:   int64(query.DateTo),
		Data:     userUsageData,
	}

	f := excelize.NewFile()
	var mu sync.Mutex
	var wg sync.WaitGroup
	var systemErr, userErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		systemErr = reportSystemData.Render(f)
	}()

	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		userErr = reportUserData.Render(f)
	}()

	wg.Wait()

	if systemErr != nil {
		return RenderFullXlsxResult{}, fmt.Errorf("error rendering system data: %w", systemErr)
	}
	if userErr != nil {
		return RenderFullXlsxResult{}, fmt.Errorf("error rendering user data: %w", userErr)
	}

	// Remove default sheet (if it's still there)
	if defaultSheet := f.GetSheetName(0); defaultSheet != "" {
		_ = f.DeleteSheet(defaultSheet)
	}

	if idx, err := f.GetSheetIndex("CPU System Usage"); err == nil {
		f.SetActiveSheet(idx)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return RenderFullXlsxResult{}, fmt.Errorf("failed to write file: %w", err)
	}

	xlsxBytes := buf.Bytes()

	if err := os.WriteFile("final_report.xlsx", xlsxBytes, 0o644); err != nil {
		log.Printf("Error writing XLSX file: %v", err)
	}

	fmt.Println("Excel report generated successfully")
	return RenderFullXlsxResult{Payload: xlsxBytes}, nil
	
}

// / mapGetCpuSystemUsageResponse maps the gRPC response to a slice of CpuUsageData.
func mapGetCpuSystemUsageResponse(usage *pb_system.GetCpuSystemUsageResponse) ([]render_xlsx.CpuSystemUsageData, error) {
	data := make([]render_xlsx.CpuSystemUsageData, len(usage.Usages))
	for i, u := range usage.Usages {
		data[i] = render_xlsx.CpuSystemUsageData{
			CPU:      u.Cpu,
			AvgUsage: u.AvgUsage,
			MaxUsage: u.MaxUsage,
			MinUsage: u.MinUsage,
		}
	}
	return data, nil
}

// mapGetCpuUserUsageResponse maps the gRPC response to a slice of CpuUsageData.
func mapGetCpuUserUsageResponse(usage *pb_user.GetCpuUserUsageResponse) ([]render_xlsx.CpuUserUsageData, error) {
	data := make([]render_xlsx.CpuUserUsageData, len(usage.Usages))
	for i, u := range usage.Usages {
		data[i] = render_xlsx.CpuUserUsageData{
			CPU:      u.Cpu,
			AvgUsage: u.AvgUsage,
			MaxUsage: u.MaxUsage,
			MinUsage: u.MinUsage,
		}
	}
	return data, nil
}
