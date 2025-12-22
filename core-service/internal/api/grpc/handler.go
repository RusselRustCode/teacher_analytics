package grpc

import (
    "context"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "github.com/RusselRustCode/teacher_analytics/core-service/proto"
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
)

type GRPCHandler struct {
    pb.UnimplementedAnalyticsServiceServer
    service interfaces.AnalyticsService
}

func NewGRPCHandler(service interfaces.AnalyticsService) *GRPCHandler {
    return &GRPCHandler{
        service: service,
    }
}

func (h *GRPCHandler) AnalyzeStudent(ctx context.Context, req *pb.AnalyzeStudentRequest) (*pb.AnalyzeStudentResponse, error) {
    analytics, err := h.service.GetAnalytics(ctx, req.StudentId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get analytics: %v", err)
    }
    
    return &pb.AnalyzeStudentResponse{
        StudentId:        analytics.StudentID,
        Cluster:          analytics.ClusterGroup,
        EngagementScore:  int32(analytics.EngagementScore),
        SuccessRate:      analytics.SuccessRate,
        TopicEfficiency:  analytics.TopicEfficiency,
        Recommendations:  analytics.Recommendations,
        AnalyzedAt:       analytics.AnalyzedAt.Format("2006-01-02 15:04:05"),
    }, nil
}

func (h *GRPCHandler) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    // Проверяем доступность зависимостей
    _, err := h.service.GetStudents(ctx)
    if err != nil {
        return &pb.HealthCheckResponse{
            Healthy: false,
            Message: "Service dependencies not available",
        }, nil
    }
    
    return &pb.HealthCheckResponse{
        Healthy: true,
        Message: "Service is healthy",
    }, nil
}

func (h *GRPCHandler) BatchAnalyze(ctx context.Context, req *pb.BatchAnalyzeRequest) (*pb.BatchAnalyzeResponse, error) {
    var results []*pb.AnalyzeStudentResponse
    
    for _, studentID := range req.StudentIds {
        analytics, err := h.service.GetAnalytics(ctx, studentID)
        if err != nil {
            continue // Пропускаем ошибки в batch режиме
        }
        
        results = append(results, &pb.AnalyzeStudentResponse{
            StudentId:        analytics.StudentID,
            Cluster:          analytics.ClusterGroup,
            EngagementScore:  int32(analytics.EngagementScore),
            SuccessRate:      analytics.SuccessRate,
            TopicEfficiency:  analytics.TopicEfficiency,
            Recommendations:  analytics.Recommendations,
            AnalyzedAt:       analytics.AnalyzedAt.Format("2006-01-02 15:04:05"),
        })
    }
    
    return &pb.BatchAnalyzeResponse{
        Results: results,
    }, nil
}