package kafka

import (
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
)
//Простая заглушка
type MessageProducer interface {
	SendLog(log *domain.StudentLog) error
	Close() error
}