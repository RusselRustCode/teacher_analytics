import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler
from sklearn.ensemble import IsolationForest

class EngagementAnalyzer:
    def __init__(self, student_logs_df: pd.DataFrame):
        self.df = student_logs_df.copy()
        self.df['timestamp'] = pd.to_datetime(self.df['timestamp'])
        self.engagement_metrics = {}

    def calculate_all(self):
        """Запускает полный цикл анализа"""
        self.calculate_active_metrics()
        self.calculate_learning_patterns()
        self.calculate_temp_patterns()
        self.calculate_risk_scores()
        return self.engagement_metrics

    def calculate_active_metrics(self):
        activity = self.df.groupby('student_id').agg(
            first_activity=('timestamp', 'min'),
            last_activity=('timestamp', 'max'),
            total_events=('timestamp', 'count'),
            total_time_sec=('time_spent_sec', 'sum')
        )
        activity['activity_duration_days'] = (
            (activity['last_activity'] - activity['first_activity']).dt.days + 1
        )
        activity['events_per_day'] = activity['total_events'] / activity['activity_duration_days']
        self.engagement_metrics['activity'] = activity

    def calculate_learning_patterns(self):
        patterns = {}
        for student_id in self.df['student_id'].unique():
            student_data = self.df[self.df['student_id'] == student_id]
            test_data = student_data[student_data['artifact_type'] == 'test_question']

            total_time = student_data['time_spent_sec'].sum()
            avg_session_time = student_data['time_spent_sec'].mean()

            if not test_data.empty:
                retry_rate = (test_data['attempts'] > 1).mean()
                avg_attempts = test_data['attempts'].mean()
                avg_correctness = test_data['correctness'].mean()
                passive = self._passive_score(
                    time_on_material=total_time - test_data['time_spent_sec'].sum(),
                    correctness=avg_correctness
                )
            else:
                retry_rate = avg_attempts = avg_correctness = 0.0
                passive = 1.0

            patterns[student_id] = {
                'retry_rate': float(retry_rate),
                'avg_attempts': float(avg_attempts),
                'avg_correctness': float(avg_correctness),
                'total_time_sec': float(total_time),
                'avg_session_time': float(avg_session_time),
                'passive_score': float(passive)
            }
        self.engagement_metrics['learning_patterns'] = pd.DataFrame.from_dict(patterns, orient='index')

    def _passive_score(self, time_on_material, correctness):
        if time_on_material > 600 and correctness < 0.5: return 0.9
        if time_on_material > 300 and correctness < 0.7: return 0.6
        return 0.1

    def calculate_temp_patterns(self):
        temp = {}
        for sid in self.df['student_id'].unique():
            d = self.df[self.df['student_id'] == sid].copy()
            d['hour'] = d['timestamp'].dt.hour
            d['dow'] = d['timestamp'].dt.dayofweek
            temp[sid] = {
                'preferred_hour': int(d['hour'].mode().iloc[0]) if not d['hour'].mode().empty else 12,
                'weekend_ratio': float(d[d['dow'].isin([5,6])].shape[0] / len(d)) if len(d) > 0 else 0,
                'activity_regularity': float(1 - d['hour'].std() / 24) if len(d) > 1 else 0.5
            }
        self.engagement_metrics['temp_patterns'] = pd.DataFrame.from_dict(temp, orient='index')

    def calculate_risk_scores(self):
        if 'activity' not in self.engagement_metrics or 'learning_patterns' not in self.engagement_metrics:
            return

        act = self.engagement_metrics['activity'][['events_per_day']]
        pat = self.engagement_metrics['learning_patterns'][['avg_correctness', 'passive_score', 'avg_attempts']]
        key = pd.concat([act, pat], axis=1).fillna(0)
        
        if len(key) < 3:
            return

        scaled = StandardScaler().fit_transform(key)
        iso = IsolationForest(contamination=0.2, random_state=42)
        key['risk_flag'] = iso.fit_predict(scaled)
        self.engagement_metrics['risk'] = key
        
    def get_summary(self, student_id):
        """Безопасно извлекает данные для конкретного ID"""
        summary = {}
        for k, df in self.engagement_metrics.items():
            if student_id in df.index:
                summary[k] = df.loc[student_id].to_dict()
        return summary