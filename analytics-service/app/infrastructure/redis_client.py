import redis.asyncio as redis
import json

class RedisCache:
    def __init__(self, host: str, port: int):
        self.client = redis.Redis(host=host, port=port, decode_responses=True)

    async def set_analytics(self, student_id: int, data: dict):
        await self.client.set(f"analytics:{student_id}", json.dumps(data), ex=3600)
        
    async def get_analytics(self, student_id: int):
        data = await self.client.get(f"analytics:{student_id}")
        return json.loads(data) if data else None