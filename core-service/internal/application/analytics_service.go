package application

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
)

type AnalyticsServiceImpl struct {
    repo    interfaces.Repository
    cache   interfaces.Cache
    producer interfaces.MessageProducer
    client  interfaces.AnalyticsClient
}

func NewAnalyticsService(
    repo interfaces.Repository,
    cache interfaces.Cache,
    producer interfaces.MessageProducer,
    client interfaces.AnalyticsClient,
) interfaces.AnalyticsService {
    return &AnalyticsServiceImpl{
        repo:     repo,
        cache:    cache,
        producer: producer,
        client:   client,
    }
}

func (s *AnalyticsServiceImpl) SendLog(ctx context.Context, log *domain.StudentLog) error {
    if log.StudentID == 0 || log.ActionType == "" {
        return fmt.Errorf("Неверная запись: требуются поля student_id и action_type.")
    }
    
    if log.Timestamp.IsZero() {
        log.Timestamp = time.Now()
    }
    
    if err := s.repo.SaveLog(ctx, log); err != nil {
        return fmt.Errorf("Не получилось созранить лог: %w", err)
    }
    
    kafkaData := map[string]interface{}{
        "type":       "student_log",
        "student_id": log.StudentID,
        "log_data":   log,
        "timestamp":  time.Now().Unix(),
    }
    
    if err := s.producer.SendJSON(ctx, "student-logs", kafkaData); err != nil {
        return fmt.Errorf("не получилось отправить в кафку: %w", err)
    }
    
    cacheKey := fmt.Sprintf("аналитика:%d", log.StudentID)
    s.cache.Delete(ctx, cacheKey)
    
    return nil
}

func (s *AnalyticsServiceImpl) GetAnalytics(ctx context.Context, studentID uint64) (*domain.StudentAnalytics, error) {
    cacheKey := fmt.Sprintf("аналитика:%d", studentID)
    
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil && cached != "" {
        var analytics domain.StudentAnalytics
        if err := json.Unmarshal([]byte(cached), &analytics); err == nil {
            return &analytics, nil
        }
    }
    
    analytics, err := s.repo.GetAnalyticsByStudentID(ctx, studentID)
    if err == nil && analytics != nil {
        // Кэшируем
        analyticsJSON, _ := json.Marshal(analytics)
        s.cache.Set(ctx, cacheKey, analyticsJSON, 5*time.Minute)
        return analytics, nil
    }
    
    // Если нет аналитики, запускаем анализ
    if err := s.TriggerAnalysis(ctx, studentID); err != nil {
        return nil, fmt.Errorf("не удалось запустить анализ: %w", err)
    }
    
    // Возвращаем заглушку
    return &domain.StudentAnalytics{
        StudentID:       studentID,
        ClusterGroup:    "processing",
        EngagementScore: 0,
        SuccessRate:     0,
        Recommendations: []string{"Анализ запущен, пожалуйста, подождите..."},
        AnalyzedAt:      time.Now(),
    }, nil
}

func (s *AnalyticsServiceImpl) TriggerAnalysis(ctx context.Context, studentID uint64) error {
    analysisCmd := map[string]interface{}{
        "type":       "analysis_command",
        "command":    "analyze_student",
        "student_id": studentID,
        "timestamp":  time.Now().Unix(),
    }
    
    if err := s.producer.SendJSON(ctx, "analysis-commands", analysisCmd); err != nil {
        return fmt.Errorf("failed to send analysis command: %w", err)
    }
    
    return nil
}

func (s *AnalyticsServiceImpl) GetStudentLogs(ctx context.Context, studentID uint64, from, to time.Time) ([]*domain.StudentLog, error) {
    return s.repo.GetLogsByStudentID(ctx, studentID, from, to)
}

func (s *AnalyticsServiceImpl) GetStudents(ctx context.Context) ([]*domain.Student, error) {
    return s.repo.GetStudents(ctx)
}

func (s *AnalyticsServiceImpl) GetStudentByID(ctx context.Context, id uint64) (*domain.Student, error) {
    return s.repo.GetStudentByID(ctx, id)
}