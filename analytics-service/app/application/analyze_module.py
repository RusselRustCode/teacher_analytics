import pandas as pd
from app.domain.models import StudentAnalytics
from app.analyzers.engagement_analyzer import EngagementAnalyzer
from app.analyzers.cluster_analyzer import StudentClusterer
from app.analyzers.material_analyzer import MaterialEffectivenessAnalyzer

class AnalyticsService:
    def __init__(self, repo, cache):
        self.repo = repo
        self.clusterer = StudentClusterer()
        self.cache = cache

    async def get_student_analysis(self, student_id: int):
        
        cached_data = await self.cache.get_analytics(student_id)
        if cached_data:
            print(f"--- Returning cached data for student {student_id} ---")
            return cached_data
        logs_dict = await self.repo.get_logs_by_student(student_id)
        if not logs_dict:
            return self._empty_response(student_id)
        
        df = pd.DataFrame(logs_dict)
        
        if 'time_spent_sec' in df.columns:
            df['time_spent_on_question'] = df['time_spent_sec']
            df['time_spent_on_material'] = df['time_spent_sec']
        
        if 'correct' in df.columns:
            df['avg_correctness'] = df['correct'].astype(int)
        
        required_columns = {
            'time_spent_on_question': 'time_spent_sec',
            'time_spent_on_material': 'time_spent_sec',
            'selected_distractor_freq': 0,
            'study_time_preference': 0,
            'material_id': 0,      
            'attempts': 1,
            'selected_distractor': 'none'
        }

        for col, default in required_columns.items():
            if col not in df.columns:
                if isinstance(default, str) and default in df.columns:
                    df[col] = df[default]
                else:
                    df[col] = default

        eng = EngagementAnalyzer(df)
        eng.calculate_all()
        summary = eng.get_summary(student_id)

        mat = MaterialEffectivenessAnalyzer(df)
        mat_stats = mat.analyze()

        cluster = self.clusterer.predict_cluster(df)

        analysis_result = StudentAnalytics(
            student_id=student_id,
            cluster_group=cluster,
            engagement_score=self._calculate_score(summary),
            success_rate=float(summary.get('learning_patterns', {}).get('avg_correctness', 0)),
            topic_efficiency={k: v['success_rate'] for k, v in mat_stats.items()},
            recommendations=self._build_recs(cluster, summary)
        )
        
        analysis_dict = {
            "student_id": analysis_result.student_id,
            "cluster_group": analysis_result.cluster_group,
            "engagement_score": analysis_result.engagement_score,
            "success_rate": analysis_result.success_rate,
            "recommendations": analysis_result.recommendations,
            "avg_time": 0.0
        }
        await self.repo.save_analytics(student_id, analysis_dict)
        print(f"--- [POSTGRES] Analysis saved for student {student_id} ---")
        
        await self.cache.set_analytics(student_id, analysis_dict)
        print(f"--- Saved analysis to Redis for student {student_id} ---")
        
        print(f"DEBUG: DF columns: {df.columns.tolist()}") 
        print(f"DEBUG: DF head: {df.head().to_dict()}")

        return analysis_result

    def _calculate_score(self, summary):
        # Процент успеха * (1 - уровень пассивности)
        lp = summary.get('learning_patterns', {})
        score = lp.get('avg_correctness', 0) * (1 - lp.get('passive_score', 0))
        return int(score * 100)

    def _build_recs(self, cluster, summary):
        recs = []
        if summary.get('risk', {}).get('risk_flag') == -1:
            recs.append("Risk detected: abnormal learning pattern")
        if cluster == "struggling":
            recs.append("Focus on basic materials")
        else:
            recs.append("Maintain current pace")
        return recs
    
    async def process_new_log(self, data: dict):
        student_id = data.get('student_id')
        if not student_id:
            return
        
        await self.cache.delete_analytics(student_id) 
        
        analysis = await self.get_student_analysis(student_id)
        
        print(f"--- [PYTHON] Анализ пересчитан для студента {student_id} ---")
        
    def _empty_response(self, student_id: int):
        """Возвращает дефолтную структуру, если данных в БД нет"""
        from app.domain.models import StudentAnalytics 
        return StudentAnalytics(
            student_id=student_id,
            cluster_group="unknown",
            engagement_score=0,
            success_rate=0.0,
            topic_efficiency={},
            recommendations=["No data available yet"]
        )