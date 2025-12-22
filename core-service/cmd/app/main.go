package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	google_grpc "google.golang.org/grpc" // Псевдоним для стандартной библиотеки

	// Твои пакеты
	internal_grpc "github.com/RusselRustCode/teacher_analytics/core-service/internal/api/grpc"
	internal_http "github.com/RusselRustCode/teacher_analytics/core-service/internal/api/http"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/application"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/config"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/infrastructure/postgres"
	"github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
	"github.com/RusselRustCode/teacher_analytics/core-service/proto"
)

func main() {
	// 1. Загрузка конфигурации
	cfg := config.LoadConfig()

	// 2. Инициализация репозитория (PostgreSQL)
	repo, err := postgres.NewPostgresRepository(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	// ПРИМЕЧАНИЕ: Если у тебя нет реальных Redis/Kafka сейчас, 
	// передай nil или пустые структуры, чтобы запуститься для теста.
	var redisCache interfaces.Cache = nil      // Замени на свою инициализацию
	var kafkaProducer interfaces.MessageProducer = nil // Замени на свою инициализацию
	var analyticsClient interfaces.AnalyticsClient = nil

	// 3. Создание сервиса (Бизнес-логика)
	analyticsService := application.NewAnalyticsService(
		repo,
		redisCache,
		kafkaProducer,
		analyticsClient,
	)

	// 4. Запуск серверов
	go startGRPCServer(cfg.GRPCPort, analyticsService)
	go startHTTPServer(cfg.HTTPPort, analyticsService)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
}

func startGRPCServer(port string, service interfaces.AnalyticsService) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Используем стандартный grpc для сервера
	s := google_grpc.NewServer()
	
	// Используем твой пакет для хендлера
	grpcHandler := internal_grpc.NewGRPCHandler(service)
	proto.RegisterAnalyticsServiceServer(s, grpcHandler)

	log.Printf("gRPC server listening on :%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startHTTPServer(port string, service interfaces.AnalyticsService) {
	router := gin.Default()

	// Мидлвар для CORS
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

	log.Printf("HTTP server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}