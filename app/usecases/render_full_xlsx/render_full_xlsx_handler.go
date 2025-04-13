package render_full_xlsx

import (
	"fmt"
	"log"
	"strings"
	"encoding/json"
	"context"

	pb_system "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_system_usage"
	pb_user "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_user_usage"
	
	proto "github.com/Javier-Godon/reports-rendering-go/proto"
)

type RenderFullXlsxHandler struct {
}

func NewRenderFullXlsxHandler() *RenderFullXlsxHandler {
	return &RenderFullXlsxHandler{
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

	//  Do something with dataListJSON (e.g., return it, log it, etc.).
	fmt.Printf("Data: %v\n", dataListJSON)

	result := RenderFullXlsxResult{
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