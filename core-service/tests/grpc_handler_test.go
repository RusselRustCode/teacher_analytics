package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RusselRustCode/teacher_analytics/core-service/internal/api/grpc"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/mocks"
	"github.com/RusselRustCode/teacher_analytics/core-service/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GRPCHandlerTestSuite struct {
	suite.Suite
	serviceMock *mocks.AnalyticsService
	handler     *grpc.GRPCHandler
}

func (s *GRPCHandlerTestSuite) SetupTest() {
	s.serviceMock = new(mocks.AnalyticsService)
	s.handler = grpc.NewGRPCHandler(s.serviceMock)
}

func (s *GRPCHandlerTestSuite) TestAnalyzeStudent_Success() {
	ctx := context.Background()
	var studentID uint64 = 1

	// Данные, которые якобы вернет наш сервис (UseCase)
	expectedAnalytics := &domain.StudentAnalytics{
		StudentID:    studentID,
		ClusterGroup: "high_performer", // Согласно вашему struct
		EngagementScore: 90,
		SuccessRate:     0.95,
		AnalyzedAt:      time.Now(),
	}

	// Настраиваем мок. В интерфейсе метод называется GetAnalytics
	s.serviceMock.On("GetAnalytics", ctx, studentID).Return(expectedAnalytics, nil)

	// Создаем правильный запрос из вашего proto
	req := &proto.AnalyzeStudentRequest{
		StudentId: studentID,
	}

	// Вызываем актуальный метод хендлера
	resp, err := s.handler.AnalyzeStudent(ctx, req)

	// Проверки
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), studentID, resp.StudentId)
	assert.Equal(s.T(), "high_performer", resp.Cluster)
	assert.Equal(s.T(), int32(90), resp.EngagementScore)
	
	s.serviceMock.AssertExpectations(s.T())
}

func TestGRPCHandlerSuite(t *testing.T) {
	suite.Run(t, new(GRPCHandlerTestSuite))
}