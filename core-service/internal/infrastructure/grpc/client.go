package grpc

import (
    "context"
    "fmt"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    pb "github.com/RusselRustCode/teacher_analytics/core-service/proto"
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
)

type GRPCAnalyticsClient struct {
    conn   *grpc.ClientConn
    client pb.AnalyticsServiceClient
}

func NewGRPCAnalyticsClient(addr string) (interfaces.AnalyticsClient, error) {
    conn, err := grpc.NewClient(
        addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithTimeout(10*time.Second),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC client: %w", err)
    }
    
    return &GRPCAnalyticsClient{
        conn:   conn,
        client: pb.NewAnalyticsServiceClient(conn),
    }, nil
}

func (c *GRPCAnalyticsClient) AnalyzeStudent(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error) {
    req := &pb.AnalyzeStudentRequest{
        StudentId: studentID,
    }
    
    resp, err := c.client.AnalyzeStudent(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("gRPC call failed: %w", err)
    }
    
    // Конвертируем protobuf в доменную модель
    analytics := &domain.StudentAnalytics{
        StudentID:       resp.StudentId,
        ClusterGroup:    resp.Cluster,
        EngagementScore: int(resp.EngagementScore),
        SuccessRate:     resp.SuccessRate,
        TopicEfficiency: resp.TopicEfficiency,
        Recommendations: resp.Recommendations,
        AnalyzedAt:      time.Now(),
    }
    
    return analytics, nil
}

func (c *GRPCAnalyticsClient) HealthCheck(ctx context.Context) error {
    req := &pb.HealthCheckRequest{}
    
    _, err := c.client.HealthCheck(ctx, req)
    return err
}

func (c *GRPCAnalyticsClient) Close() error {
    return c.conn.Close()
}