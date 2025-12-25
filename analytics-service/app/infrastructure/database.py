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
                artifact_id as material_id,
                artifact_type,
                action_time as timestamp,
                (metadata->>'time_spent')::float as time_spent_sec,
                (metadata->>'correctness')::float as correctness,
                (metadata->>'attempts')::int as attempts,
                (metadata->>'selected_distractor') as selected_distractor
            FROM student_logs
            WHERE student_id = $1
            ORDER BY action_time ASC
        """
        async with self.pool.acquire() as conn:
            rows = await conn.fetch(query, student_id)
            return [dict(row) for row in rows]

    async def close(self):
        if self.pool:
            await self.pool.close()