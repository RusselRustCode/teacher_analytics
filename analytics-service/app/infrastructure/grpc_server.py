import grpc
from proto import analytics_pb2_grpc
from app.api.grpc_handler import AnalyticsGRPCHandler

class GRPCServer:
    def __init__(self, port: str, handler: AnalyticsGRPCHandler):
        self.port = port
        self.handler = handler
        self.server = None

    async def start(self):
        self.server = grpc.aio.server()
        
        analytics_pb2_grpc.add_AnalyticsServiceServicer_to_server(
            self.handler, 
            self.server
        )
        
        listen_addr = f"[::]:{self.port}"
        self.server.add_insecure_port(listen_addr)
        
        print(f"Запустил gRPC сервер на {listen_addr}...")
        await self.server.start()

    async def wait_for_termination(self):
        if self.server:
            await self.server.wait_for_termination()

    async def stop(self):
        if self.server:
            await self.server.stop(5)
            print("gRPC сервер остановился")