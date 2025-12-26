import json
from aiokafka import AIOKafkaConsumer

class AnalyticsConsumer:
    def __init__(self, brokers: str, topic: str, service):
        self.consumer = AIOKafkaConsumer(
            topic,
            bootstrap_servers=brokers,
            group_id="analytics_group_v2", 
            auto_offset_reset="earliest",  # Читать всё с начала, если группа новая
            value_deserializer=lambda m: json.loads(m.decode('utf-8'))
        )
        self.service = service 

    async def start(self):
        print(f"--- Попытка подключения к Kafka... ---")
        try:
            await self.consumer.start()
            print(f"--- Успешное подключение! Слушаю топик... ---")
            
            async for msg in self.consumer:
                data = msg.value 
                s_id = data.get('student_id') or data.get('log_data', {}).get('student_id')
                print(f"--- [KAFKA] Получены данные для студента {s_id} ---")
                
                await self.service.process_new_log(data) 
        except Exception as e:
            print(f"--- [KAFKA ERROR] Ошибка: {e} ---")
        finally:
            await self.consumer.stop()