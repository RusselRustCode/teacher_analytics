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

	expectedAnalytics := &domain.StudentAnalytics{
		StudentID:    studentID,
		ClusterGroup: "high_performer", 
		EngagementScore: 90,
		SuccessRate:     0.95,
		AnalyzedAt:      time.Now(),
	}

	s.serviceMock.On("GetAnalytics", ctx, studentID).Return(expectedAnalytics, nil)

	req := &proto.AnalyzeStudentRequest{
		StudentId: studentID,
	}

	resp, err := s.handler.AnalyzeStudent(ctx, req)

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