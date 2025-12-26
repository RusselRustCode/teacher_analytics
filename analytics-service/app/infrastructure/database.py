import asyncpg
import json

class PostgresRepository:
    def __init__(self, dsn: str):
        self.dsn = dsn
        self.pool = None

    async def connect(self):
        if not self.pool:
            self.pool = await asyncpg.create_pool(dsn=self.dsn)

    async def get_logs_by_student(self, student_id: int):
        query = """
            SELECT 
                student_id,
                action_type as artifact_type,
                timestamp,
                time_spent_sec,
                CASE WHEN correct THEN 1.0 ELSE 0.0 END as correctness,
                1 as attempts, -- заглушка
                'none' as selected_distractor -- заглушка
            FROM student_logs
            WHERE student_id = $1
            ORDER BY timestamp ASC
        """
        async with self.pool.acquire() as conn:
            rows = await conn.fetch(query, student_id)
            return [dict(row) for row in rows]

    async def close(self):
        if self.pool:
            await self.pool.close()