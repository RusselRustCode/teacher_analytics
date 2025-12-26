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
    user := getEnv("DB_USER", "admin")
    pass := getEnv("DB_PASSWORD", "admin_password")
    host := getEnv("DB_HOST", "student-analytics-postgres") 
    port := getEnv("DB_PORT", "5432")
    name := getEnv("DB_NAME", "student_analytics")

    dsn := "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + name + "?sslmode=disable"

    return &Config{
        HTTPPort:      getEnv("HTTP_PORT", "8080"),
        GRPCPort:      getEnv("GRPC_PORT", "50051"),
        DBDSN:         dsn,
        RedisAddr:     getEnv("REDIS_HOST", "redis") + ":" + getEnv("REDIS_PORT", "6379"),
        KafkaAddr:     getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9094"),
        AnalyticsAddr: getEnv("ANALYTICS_GRPC_HOST", "analytics-service") + ":" + getEnv("ANALYTICS_GRPC_PORT", "50052"),
    }
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}