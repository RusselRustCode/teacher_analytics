import json
from aiokafka import AIOKafkaConsumer

class AnalyticsConsumer:
    def __init__(self, brokers: str, topic: str, service):
        self.consumer = AIOKafkaConsumer(
            topic,
            bootstrap_servers=brokers,
            group_id="analytics_group"
        )
        self.service = service # Ссылка на application.AnalyticsService

    async def start(self):
        await self.consumer.start()
        try:
            async for msg in self.consumer:
                data = json.loads(msg.value)
                # Например, при поступлении нового лога - пересчитываем кэш
                print(f"Received log for student {data['student_id']}, triggering update...")
                # await self.service.process_new_log(data) 
        finally:
            await self.consumer.stop()