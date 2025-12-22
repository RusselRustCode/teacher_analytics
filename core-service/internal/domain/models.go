package domain

import "time"

// Student - модель студента
type Student struct {
    ID        uint64    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"` // student, teacher, admin
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// StudentLog - лог действий студента
type StudentLog struct {
    ID                  uint64    `json:"id"`
    StudentID           uint64    `json:"student_id"`
    ActionType          string    `json:"action_type"` // view_material, test_answer, watch_video
    MaterialID          string    `json:"material_id"`
    Correct             bool      `json:"correct"`
    TimeSpentSec        int       `json:"time_spent_sec"`
    Difficulty          int       `json:"difficulty"` // 1-5
    TimeSpentOnMat      int       `json:"time_spent_on_mat"`
    TimeSpentOnQuestion int       `json:"time_spent_on_question"`
    Attempts            int       `json:"attempts"`
    SelectedDistractor  string    `json:"selected_distractor"`
    Timestamp           time.Time `json:"timestamp"`
}

// StudentAnalytics - результаты анализа студента
type StudentAnalytics struct {
    ID                uint64             `json:"id"`
    StudentID         uint64             `json:"student_id"`
    ClusterGroup      string             `json:"cluster_group"` // high_performer, average, needs_help
    EngagementScore   int                `json:"engagement_score"` // 0-100
    AvgTimePerTask    float64            `json:"avg_time_per_task"`
    SuccessRate       float64            `json:"success_rate"` // 0.0-1.0
    TopicEfficiency   map[string]float64 `json:"topic_efficiency"`
    Recommendations   []string           `json:"recommendations"`
    AnalyzedAt        time.Time          `json:"analyzed_at"`
}

// MaterialAnalytics - аналитика материала
type MaterialAnalytics struct {
    MaterialID        string             `json:"material_id"`
    SuccessRate       float64            `json:"success_rate"`
    AvgAttempts       float64            `json:"avg_attempts"`
    DifficultyIndex   float64            `json:"difficulty_index"`
    DistractorStats   map[string]int     `json:"distractor_stats"`
    StudentsCompleted int                `json:"students_completed"`
    AvgTimeSpent      float64            `json:"avg_time_spent"`
}