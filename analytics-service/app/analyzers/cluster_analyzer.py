import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler
from sklearn.cluster import DBSCAN

class StudentClusterer:
    def __init__(self, eps=1.5, min_samples=3):
        self.eps = eps
        self.min_samples = min_samples
        self.scaler = StandardScaler()

    def predict_cluster(self, df: pd.DataFrame) -> str:
        """Определяет кластер для студента на основе его логов"""
        if df.empty or len(df) < 3:
            return "new_student"

        features = df[['attempts', 'correctness', 'time_spent_on_question', 
                      'time_spent_on_material', 'selected_distractor_freq', 
                      'study_time_preference']]
        
        try:
            scaled = self.scaler.fit_transform(features)
            dbscan = DBSCAN(eps=self.eps, min_samples=self.min_samples)
            clusters = dbscan.fit_predict(scaled)
            
            last_cluster = clusters[-1]
            
            mapping = {
                -1: "outlier",
                0: "high_performer",
                1: "average",
                2: "struggling",
                3: "passive_learner"
            }
            return mapping.get(last_cluster, "unknown")
        except Exception:
            return "analysis_failed"