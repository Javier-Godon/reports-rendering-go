package grpc_client

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb_system "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_system_usage"
	pb_user "github.com/Javier-Godon/reports-rendering-go/proto/get_cpu_user_usage"
)

// GRPCClient holds the gRPC connection and client stubs.  It's safe for concurrent use.
type GRPCClient struct {
	conn         *grpc.ClientConn
	systemClient pb_system.GetCpuSystemUsageServiceClient
	userClient   pb_user.GetCpuUserUsageServiceClient
}

// NewGRPCClient creates a new GRPCClient.
func NewGRPCClient(address string) (*GRPCClient, error) {
	// Use insecure.NewCredentials() for a non-secure connection.  For production, use appropriate credentials.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to gRPC server: %v", err)
		return nil, err // Important: Return the error!
	}

	client := &GRPCClient{
		conn:         conn,
		systemClient: pb_system.NewGetCpuSystemUsageServiceClient(conn),
		userClient:   pb_user.NewGetCpuUserUsageServiceClient(conn),
	}
	return client, nil
}

// GetCpuSystemUsage retrieves CPU system usage.
func (c *GRPCClient) GetCpuSystemUsage(ctx context.Context, dateFrom int64, dateTo int64) (*pb_system.GetCpuSystemUsageResponse, error) {
	
	req := &pb_system.GetCpuSystemUsageRequest{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	resp, err := c.systemClient.GetCpuSystemUsage(ctx, req)
	if err != nil {
		log.Printf("failed to get cpu system usage: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetCpuUserUsage retrieves CPU user usage.
func (c *GRPCClient) GetCpuUserUsage(ctx context.Context, dateFrom int64, dateTo int64) (*pb_user.GetCpuUserUsageResponse, error) {
	
	req := &pb_user.GetCpuUserUsageRequest{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	resp, err := c.userClient.GetCpuUserUsage(ctx, req)
	if err != nil {
		log.Printf("failed to get cpu user usage: %v", err)
		return nil, err
	}
	return resp, nil
}

// Close closes the underlying gRPC connection.  It's good practice to close connections when you're done with them.
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

// GRPCClientSingleton is a singleton to ensure only one gRPC client.
type GRPCClientSingleton struct {
	client *GRPCClient
	once   sync.Once // Ensures the client is initialized only once.
	err    error
}

// Instance returns the singleton instance of the GRPCClient.
func (s *GRPCClientSingleton) Instance(address string) (*GRPCClient, error) {
	s.once.Do(func() {
		s.client, s.err = NewGRPCClient(address) // Initialize the client
	})
	return s.client, s.err
}

var grpcClientSingleton GRPCClientSingleton // Global instance of the singleton
