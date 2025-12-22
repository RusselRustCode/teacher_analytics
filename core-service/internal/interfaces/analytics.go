package interfaces

import (
    "context"
    "time"
    
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
)

// AnalyticsService - основной сервис аналитики
type AnalyticsService interface {
    // Логи
    SendLog(ctx context.Context, log *domain.StudentLog) error
    GetStudentLogs(ctx context.Context, studentID uint64, from, to time.Time) ([]*domain.StudentLog, error)
    
    // Аналитика
    GetAnalytics(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error)
    TriggerAnalysis(ctx context.Context, studentID uint64) error
    
    // Студенты
    GetStudents(ctx context.Context) ([]*domain.Student, error)
    GetStudentByID(ctx context.Context, id uint64) (*domain.Student, error)
}

// Repository - репозиторий для работы с данными
type Repository interface {
    // Студенты
    SaveStudent(ctx context.Context, student *domain.Student) error
    GetStudentByID(ctx context.Context, id uint64) (*domain.Student, error)
    GetStudents(ctx context.Context) ([]*domain.Student, error)
    
    // Логи
    SaveLog(ctx context.Context, log *domain.StudentLog) error
    GetLogsByStudentID(ctx context.Context, studentID uint64, from, to time.Time) ([]*domain.StudentLog, error)
    GetLogsByMaterialID(ctx context.Context, materialID string) ([]*domain.StudentLog, error)
    
    // Аналитика
    SaveAnalytics(ctx context.Context, analytics *domain.StudentAnalytics) error
    GetAnalyticsByStudentID(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error)
    UpdateAnalytics(ctx context.Context, analytics *domain.StudentAnalytics) error

    Close() error
}

// MessageProducer - продюсер сообщений (Kafka)
type MessageProducer interface {
    Send(ctx context.Context, topic string, key []byte, value []byte) error
    SendJSON(ctx context.Context, topic string, data interface{}) error
    Close() error
}

// Cache - кэш (Redis)
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}

// AnalyticsClient - gRPC клиент для сервиса аналитики
type AnalyticsClient interface {
    AnalyzeStudent(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error)
    HealthCheck(ctx context.Context) error
    Close() error
}