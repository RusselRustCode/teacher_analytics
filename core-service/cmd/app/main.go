package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"fmt"

	"github.com/gin-gonic/gin"
	google_grpc "google.golang.org/grpc" // Псевдоним для стандартной библиотеки

	internal_grpc "github.com/RusselRustCode/teacher_analytics/core-service/internal/api/grpc"
	internal_http "github.com/RusselRustCode/teacher_analytics/core-service/internal/api/http"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/application"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/config"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/infrastructure/kafka"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/infrastructure/postgres"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/infrastructure/redis"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/infrastructure/grpc"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
	"github.com/RusselRustCode/teacher_analytics/core-service/proto"
)

func main() {
	cfg := config.LoadConfig()

	repo, err := postgres.NewPostgresRepository(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Не получилось подключсится к бд: %v", err)
	}
	defer repo.Close()

	kafkaBrokers := []string{os.Getenv("KAFKA_BOOTSTRAP_SERVERS")}
    redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
    analyticsAddr := fmt.Sprintf("%s:%s", os.Getenv("ANALYTICS_GRPC_HOST"), os.Getenv("ANALYTICS_GRPC_PORT"))

	var redisCache interfaces.Cache = nil    
	var kafkaProducer interfaces.MessageProducer = nil 
	var analyticsClient interfaces.AnalyticsClient = nil

	redisCache = redis.NewRedisCache(redisAddr, os.Getenv("REDIS_PASSWORD"), 0)
    defer redisCache.Close()

	kafkaProducer = kafka.NewKafkaProducer(kafkaBrokers)
	defer kafkaProducer.Close()

	analyticsClient, err = grpc.NewGRPCAnalyticsClient(analyticsAddr)
    if err != nil {
        log.Fatalf("Не получилось подключиться к analytics client: %v", err)
    }
    defer analyticsClient.Close()

	analyticsService := application.NewAnalyticsService(
		repo,
		redisCache,
		kafkaProducer,
		analyticsClient,
	)

	go startGRPCServer(cfg.GRPCPort, analyticsService)
	go startHTTPServer(cfg.HTTPPort, analyticsService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Выключение серверов...")
}

func startGRPCServer(port string, service interfaces.AnalyticsService) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Не смог прослушать: %v", err)
	}

	s := google_grpc.NewServer()
	
	grpcHandler := internal_grpc.NewGRPCHandler(service)
	proto.RegisterAnalyticsServiceServer(s, grpcHandler)

	log.Printf("gRPC server прослушивает :%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Не удалось предоставить gRPC: %v", err)
	}
}

func startHTTPServer(port string, service interfaces.AnalyticsService) {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	handler := internal_http.NewHTTPHandler(service)

	api := router.Group("/api")
	{
		api.POST("/log", handler.SendLog)
		api.GET("/analytics/:student_id", handler.GetAnalytics)
		api.GET("/students/:student_id/logs", handler.GetStudentLogs)
		api.GET("/students", handler.GetStudents)
	}

	log.Printf("HTTP-сервер прослушивает :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Не удалось запустить HTTP-сервер.: %v", err)
	}
}