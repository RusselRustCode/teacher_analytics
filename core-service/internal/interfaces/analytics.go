package interfaces

import (
    "context"
    "time"
    
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
)

type AnalyticsService interface {
    SendLog(ctx context.Context, log *domain.StudentLog) error
    GetStudentLogs(ctx context.Context, studentID uint64, from, to time.Time) ([]*domain.StudentLog, error)
    
    GetAnalytics(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error)
    TriggerAnalysis(ctx context.Context, studentID uint64) error
    
    GetStudents(ctx context.Context) ([]uint64, error)
    GetStudentByID(ctx context.Context, id uint64) (*domain.Student, error)
}

type Repository interface {
    SaveStudent(ctx context.Context, student *domain.Student) error
    GetStudentByID(ctx context.Context, id uint64) (*domain.Student, error)
    GetStudents(ctx context.Context) ([]uint64, error)
    
    SaveLog(ctx context.Context, log *domain.StudentLog) error
    GetLogsByStudentID(ctx context.Context, studentID uint64, from, to time.Time) ([]*domain.StudentLog, error)
    GetLogsByMaterialID(ctx context.Context, materialID string) ([]*domain.StudentLog, error)
    
    SaveAnalytics(ctx context.Context, analytics *domain.StudentAnalytics) error
    GetAnalyticsByStudentID(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error)
    UpdateAnalytics(ctx context.Context, analytics *domain.StudentAnalytics) error

    Close() error
}

type MessageProducer interface {
    Send(ctx context.Context, topic string, key []byte, value []byte) error
    SendJSON(ctx context.Context, topic string, data interface{}) error
    Close() error
}

type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Close() error
}

type AnalyticsClient interface {
    AnalyzeStudent(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error)
    HealthCheck(ctx context.Context) error
    Close() error
}