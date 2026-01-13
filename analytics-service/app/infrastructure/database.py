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
            
    async def save_analytics(self, student_id: int, result: dict):
        query = """
            INSERT INTO student_analytics 
            (student_id, cluster_group, engagement_score, avg_time_per_task, success_rate, analyzed_at) 
            VALUES ($1, $2, $3, $4, $5, NOW()) 
            ON CONFLICT (student_id) 
            DO UPDATE SET 
                cluster_group = EXCLUDED.cluster_group,
                engagement_score = EXCLUDED.engagement_score,
                avg_time_per_task = EXCLUDED.avg_time_per_task,
                success_rate = EXCLUDED.success_rate,
                analyzed_at = NOW();
            """
        await self.pool.execute(
            query, 
            student_id, 
            result.get('cluster', 'unknown'),
            float(result.get('engagement_score', 0)),
            float(result.get('avg_time', 0)),
            float(result.get('success_rate', 0)),
            # result.get('recommendations', [])
        )