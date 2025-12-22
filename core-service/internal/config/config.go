package config

import "os"

type Config struct {
	HTTPPort    string
	GRPCPort    string
	DBDSN       string
	RedisAddr   string
	KafkaAddr   string
	AnalyticsAddr string
}

func LoadConfig() *Config {
	return &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		GRPCPort:      getEnv("GRPC_PORT", "50051"),
		DBDSN:         "postgres://admin:admin@postgres:5432/student_analytics?sslmode=disable",
		RedisAddr:     getEnv("REDIS_HOST", "redis") + ":6379",
		KafkaAddr:     getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9094"),
		AnalyticsAddr: getEnv("ANALYTICS_GRPC_HOST", "analytics-service") + ":50052",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}