from typing import List
import pandas as pd
from app.domain.models import StudentAnalytics
from app.analyzers.engagement_analyzer import EngagementAnalyzer
from app.analyzers.cluster_analyzer import StudentClusterer # Обернем твой код в класс
from app.analyzers.material_analyzer import MaterialEffectivenessAnalyzer

class AnalyticsService:
    def __init__(self, repo):
        self.repo = repo

    async def get_student_analysis(self, student_id: int) -> StudentAnalytics:
        # 1. Загружаем логи студента из БД
        raw_logs = await self.repo.get_logs_by_student(student_id)
        if not raw_logs:
            raise ValueError("No logs found for student")

        df = pd.DataFrame(raw_logs)

        # 2. Считаем вовлеченность (твой код)
        eng_analyzer = EngagementAnalyzer(df)
        eng_analyzer.calculate_engagement_metrics()
        eng_summary = eng_analyzer.get_student_engagement_summary(student_id)

        # 3. Определяем кластер (твой код DBSCAN)
        clusterer = StudentClusterer()
        cluster_name = clusterer.predict_cluster(df)

        # 4. Собираем финальный объект
        return StudentAnalytics(
            student_id=student_id,
            cluster_group=cluster_name,
            engagement_score=int(eng_summary['efficiency']['эффективность обучения'] * 100),
            success_rate=df['correct'].mean(),
            topic_efficiency={}, # Сюда можно добавить логику из MaterialAnalyzer
            recommendations=self._generate_recs(cluster_name)
        )

    def _generate_recs(self, cluster: str) -> List[str]:
        recs = {
            "high_performer": ["Add more complex tasks", "Suggest peer-reviewing"],
            "needs_help": ["Review basic materials", "Schedule a consultation"],
        }
        return recs.get(cluster, ["Continue regular practice"])