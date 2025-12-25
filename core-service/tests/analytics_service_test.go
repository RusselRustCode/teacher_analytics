package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RusselRustCode/teacher_analytics/core-service/internal/application"
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AnalyticsServiceTestSuite struct {
	suite.Suite
	ctx          context.Context
	repoMock     *mocks.Repository
	cacheMock    *mocks.Cache
	producerMock *mocks.MessageProducer
	clientMock   *mocks.AnalyticsClient
	service      interfaces.AnalyticsService
}

func (s *AnalyticsServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.repoMock = new(mocks.Repository)
	s.cacheMock = new(mocks.Cache)
	s.producerMock = new(mocks.MessageProducer)
	s.clientMock = new(mocks.AnalyticsClient)

	s.service = application.NewAnalyticsService(
		s.repoMock,
		s.cacheMock,
		s.producerMock,
		s.clientMock,
	)
}

func (s *AnalyticsServiceTestSuite) TestSendLog_Success() {
	log := &domain.StudentLog{
		StudentID:  1,
		ActionType:     "view_lesson",
		MaterialID: "math_101",
		Timestamp:  time.Now(),
	}

	s.repoMock.On("SaveLog", s.ctx, log).Return(nil)
	s.producerMock.On("SendLog", log).Return(nil)

	err := s.service.SendLog(s.ctx, log)

	assert.NoError(s.T(), err)
	s.repoMock.AssertExpectations(s.T())
	s.producerMock.AssertExpectations(s.T())
}

func (s *AnalyticsServiceTestSuite) TestGetAnalytics_CacheHit() {
	studentID := uint64(123)
	cachedData := `{"student_id":123, "status":"Active"}`

	s.cacheMock.On("Get", s.ctx, mock.Anything).Return(cachedData, nil)

	res, err := s.service.GetAnalytics(s.ctx, studentID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), studentID, res.StudentID)
	
	s.repoMock.AssertNotCalled(s.T(), "GetAnalyticsByStudentID", mock.Anything, mock.Anything)
}

func TestAnalyticsService(t *testing.T) {
	suite.Run(t, new(AnalyticsServiceTestSuite))
}