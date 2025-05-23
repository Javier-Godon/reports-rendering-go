package render_full_pdf

import (
	"fmt"
	"log"
	"strings"
	"encoding/json"
	"context"

	config "github.com/Javier-Godon/reports-rendering-go/framework"
	pb_system "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_system_usage"
	pb_user "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_user_usage"
	// render_xlsx "github.com/Javier-Godon/reports-rendering-go/render/xlsx"
	
	proto "github.com/Javier-Godon/reports-rendering-go/proto"
)

type RenderFullPdfHandler struct {
}

func NewRenderFullPdfHandler() *RenderFullPdfHandler {
	return &RenderFullPdfHandler{
	}
}

func (handler RenderFullPdfHandler) Handle(query RenderFullPdfQuery) (RenderFullPdfResult, error) {	
	address := config.AppConfig.DataProvider.ADDRESS 
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

	//  Map the gRPC responses to JSON strings.
	systemUsageJSON, err := mapGetCpuSystemUsageResultToJSON(systemUsage)
	if err != nil {
		log.Printf("Error converting system usage to JSON: %v", err)
	}

	userUsageJSON, err := mapGetCpuUserUsageResultToJSON(userUsage)
	if err != nil {
		log.Printf("Error converting user usage to JSON: %v", err)		
	}

	//  Create a list of JSON strings.
	dataListJSON := []string{systemUsageJSON, userUsageJSON}


	// render_xlsx.CpuSystemUsageReport()

	//  Do something with dataListJSON (e.g., return it, log it, etc.).
	fmt.Printf("Data: %v\n", dataListJSON)

	result := RenderFullPdfResult{
		Payload: strings.Join(dataListJSON, ","), // Join the JSON strings into a single string
	}

	return result, err
}





/// mapGetCpuSystemUsageResultToJSON converts the GetCpuSystemUsageResponse to a JSON string.
func mapGetCpuSystemUsageResultToJSON(resp *pb_system.GetCpuSystemUsageResponse) (string, error) {	
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal GetCpuSystemUsageResponse to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// mapGetCpuUserUsageResultToJSON converts the GetCpuUserUsageResponse to a JSON string.
func mapGetCpuUserUsageResultToJSON(resp *pb_user.GetCpuUserUsageResponse) (string, error) {	
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal GetCpuUserUsageResponse to JSON: %w", err)
	}
	return string(jsonBytes), nil
}