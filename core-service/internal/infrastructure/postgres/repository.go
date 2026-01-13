package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (interfaces.Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db unreachable: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

// --- СТУДЕНТЫ ---

func (r *PostgresRepository) SaveStudent(ctx context.Context, s *domain.Student) error {
	query := `INSERT INTO students (name, email) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, s.Name, s.Email).Scan(&s.ID)
}

func (r *PostgresRepository) GetStudentByID(ctx context.Context, id uint64) (*domain.Student, error) {
	s := &domain.Student{}
	query := `SELECT id, name, email FROM students WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *PostgresRepository) GetStudents(ctx context.Context) ([]uint64, error) {
    query := `SELECT DISTINCT student_id FROM student_logs`
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var studentIDs []uint64
    for rows.Next() {
        var id uint64
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        studentIDs = append(studentIDs, id)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return studentIDs, nil
}


func (r *PostgresRepository) SaveLog(ctx context.Context, log *domain.StudentLog) error {
	query := `
		INSERT INTO student_logs (student_id, action_type, correct, time_spent_sec, timestamp)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, 
		log.StudentID, log.ActionType, log.Correct, log.TimeSpentSec, log.Timestamp)
	return err
}

func (r *PostgresRepository) GetLogsByStudentID(ctx context.Context, id uint64, f, t time.Time) ([]*domain.StudentLog, error) {
	query := `
		SELECT student_id, action_type, correct, time_spent_sec, timestamp 
		FROM student_logs 
		WHERE student_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC`
	rows, err := r.db.QueryContext(ctx, query, id, f, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.StudentLog
	for rows.Next() {
		l := &domain.StudentLog{}
		if err := rows.Scan(&l.StudentID, &l.ActionType, &l.Correct, &l.TimeSpentSec, &l.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *PostgresRepository) GetLogsByMaterialID(ctx context.Context, m string) ([]*domain.StudentLog, error) {
	return nil, nil //Заглушка
}


func (r *PostgresRepository) SaveAnalytics(ctx context.Context, a *domain.StudentAnalytics) error {
	query := `
		INSERT INTO student_analytics (student_id, cluster_group, engagement_score, avg_time_per_task, success_rate, analyzed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, 
		a.StudentID, a.ClusterGroup, a.EngagementScore, a.AvgTimePerTask, a.SuccessRate, time.Now())
	return err
}

func (r *PostgresRepository) GetAnalyticsByStudentID(ctx context.Context, id uint64) (*domain.StudentAnalytics, error) {
	a := &domain.StudentAnalytics{}
	query := `
		SELECT student_id, cluster_group, engagement_score, avg_time_per_task, success_rate, analyzed_at
		FROM student_analytics WHERE student_id = $1
		ORDER BY analyzed_at DESC LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.StudentID, &a.ClusterGroup, &a.EngagementScore, &a.AvgTimePerTask, &a.SuccessRate, &a.AnalyzedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (r *PostgresRepository) UpdateAnalytics(ctx context.Context, a *domain.StudentAnalytics) error {
	query := `
		UPDATE student_analytics 
		SET cluster_group = $2, engagement_score = $3, avg_time_per_task = $4, success_rate = $5, analyzed_at = $6
		WHERE student_id = $1`
	_, err := r.db.ExecContext(ctx, query, 
		a.StudentID, a.ClusterGroup, a.EngagementScore, a.AvgTimePerTask, a.SuccessRate, time.Now())
	return err
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}