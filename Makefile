.PHONY: up down build logs clean proto test

# Запуск проекта
up:
	docker-compose up -d

# Остановка проекта
down:
	docker-compose down

# Перезапуск
restart:
	docker-compose restart

# Просмотр логов
logs:
	docker-compose logs -f

# Сборка всех сервисов
build:
	docker-compose build --no-cache

# Очистка
clean:
	docker-compose down -v
	docker system prune -f

# Генерация protobuf
proto:
	cd core-service && ./scripts/generate.sh
	cd analytics-service && python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. app/proto/analytics.proto

# Тесты
test:
	cd core-service && go test ./...
	cd analytics-service && python -m pytest

# Миграции БД
migrate:
	docker exec -it student-analytics-postgres psql -U admin -d student_analytics -f /docker-entrypoint-initdb.d/init.sql

# Импорт тестовых данных
import-test-data:
	docker exec -it student-analytics-postgres psql -U admin -d student_analytics -c "SELECT COUNT(*) FROM students;"