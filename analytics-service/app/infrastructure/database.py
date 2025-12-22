import asyncpg # Быстрая асинхронная библиотека для Postgres
import pandas as pd

class PostgresRepository:
    def __init__(self, dsn: str):
        self.dsn = dsn
        self.pool = None

    async def connect(self):
        if not self.pool:
            self.pool = await asyncpg.create_pool(dsn=self.dsn)

    async def get_logs_by_student(self, student_id: int):
        # Запрос всех логов для конкретного студента
        query = "SELECT * FROM student_logs WHERE student_id = $1"
        async with self.pool.acquire() as conn:
            rows = await conn.fetch(query, student_id)
            # Конвертируем записи asyncpg в список словарей (для Pandas)
            return [dict(row) for row in rows]

    async def get_all_logs(self):
        # Для кластеризации нам могут понадобиться данные всех студентов
        query = "SELECT * FROM student_logs"
        async with self.pool.acquire() as conn:
            rows = await conn.fetch(query)
            return [dict(row) for row in rows]
            
    async def close(self):
        if self.pool:
            await self.pool.close()