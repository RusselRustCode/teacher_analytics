import asyncio
import os
from app.infrastructure.database import PostgresRepository
from app.infrastructure.grpc_server import GRPCServer
from app.application.analyze_module import AnalyticsService
from app.api.grpc_handler import AnalyticsGRPCHandler
from app.infrastructure.kafka_consumer import AnalyticsConsumer
from app.infrastructure.redis_client import RedisCache

async def main():
    DSN = os.getenv("DATABASE_URL", "postgresql://admin:admin_password@student-analytics-postgres:5432/student_analytics")
    PORT = os.getenv("GRPC_PORT", "50052")
    KAFKA_BROKERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9094")
    REDIS_HOST = os.getenv("REDIS_HOST", "student-analytics-redis")
    REDIS_PORT = os.getenv("REDIS_PORT", "6379")

    repo = PostgresRepository(DSN)
    await repo.connect()


    cache = RedisCache(host=REDIS_HOST, port=REDIS_PORT)
    service = AnalyticsService(repo, cache)
    handler = AnalyticsGRPCHandler(service)

    server = GRPCServer(port=PORT, handler=handler)
    
    consumer = AnalyticsConsumer(brokers=KAFKA_BROKERS, topic="student-logs", service=service)

    print(f"--- Analytics Service запущен на порту {PORT} ---")
    
    try:
        await asyncio.gather(
            server.start(),        
            consumer.start(),      
        )
    except Exception as e:
        print(f"Ошибка: {e}")
    finally:
        await repo.close()

if __name__ == "__main__":
    asyncio.run(main())