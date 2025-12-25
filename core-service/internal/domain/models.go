package domain

import "time"

type Student struct {
    ID        uint64    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"` 
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type StudentLog struct {
    ID                  uint64    `json:"id"`
    StudentID           uint64    `json:"student_id"`
    ActionType          string    `json:"action_type"` 
    MaterialID          string    `json:"material_id"`
    Correct             bool      `json:"correct"`
    TimeSpentSec        int       `json:"time_spent_sec"`
    Difficulty          int       `json:"difficulty"` 
    TimeSpentOnMat      int       `json:"time_spent_on_mat"`
    TimeSpentOnQuestion int       `json:"time_spent_on_question"`
    Attempts            int       `json:"attempts"`
    SelectedDistractor  string    `json:"selected_distractor"`
    Timestamp           time.Time `json:"timestamp"`
}

type StudentAnalytics struct {
    ID                uint64             `json:"id"`
    StudentID         uint64             `json:"student_id"`
    ClusterGroup      string             `json:"cluster_group"` 
    EngagementScore   int                `json:"engagement_score"` 
    AvgTimePerTask    float64            `json:"avg_time_per_task"`
    SuccessRate       float64            `json:"success_rate"` 
    TopicEfficiency   map[string]float64 `json:"topic_efficiency"`
    Recommendations   []string           `json:"recommendations"`
    AnalyzedAt        time.Time          `json:"analyzed_at"`
}

type MaterialAnalytics struct {
    MaterialID        string             `json:"material_id"`
    SuccessRate       float64            `json:"success_rate"`
    AvgAttempts       float64            `json:"avg_attempts"`
    DifficultyIndex   float64            `json:"difficulty_index"`
    DistractorStats   map[string]int     `json:"distractor_stats"`
    StudentsCompleted int                `json:"students_completed"`
    AvgTimeSpent      float64            `json:"avg_time_spent"`
}