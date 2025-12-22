package kafka

import (
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
)

// MessageProducer описывает методы для отправки сообщений в очередь
type MessageProducer interface {
	SendLog(log *domain.StudentLog) error
	Close() error
}