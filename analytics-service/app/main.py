import asyncio
import grpc
from proto import analytics_pb2_grpc
from app.api.grpc_handler import AnalyticsGRPCHandler
from app.application.analytics_service import AnalyticsService
from app.infrastructure.database import PostgresRepository # Твой будущий репо

async def serve():
    # Инициализируем слои (Dependency Injection)
    repo = PostgresRepository()
    service = AnalyticsService(repo)
    handler = AnalyticsGRPCHandler(service)

    server = grpc.aio.server()
    analytics_pb2_grpc.add_AnalyticsServiceServicer_to_server(handler, server)
    
    server.add_insecure_port('[::]:50052')
    print("Python Analytics Service started on port 50052")
    
    await server.start()
    await server.wait_for_termination()

if __name__ == "__main__":
    asyncio.run(serve())