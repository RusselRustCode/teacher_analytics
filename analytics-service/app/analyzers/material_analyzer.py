import pandas as pd
import numpy as np
from typing import Dict, Any

class MaterialEffectivenessAnalyzer:
    def __init__(self, student_logs_df: pd.DataFrame):
        self.df = student_logs_df.copy()
        self.df["timestamp"] = pd.to_datetime(self.df["timestamp"])
        test_only = self.df[self.df["artifact_type"] == "test_question"]
        self.course_avg_success = test_only["correctness"].mean() if not test_only.empty else 0.0

    def analyze(self) -> Dict[str, Any]:
        test_logs = self.df[self.df["artifact_type"] == "test_question"].copy()
        if test_logs.empty:
            return {}

        material_stats = {}
        for material_id in test_logs["material_id"].unique():
            mat_data = test_logs[test_logs["material_id"] == material_id]
            
            success_rate = mat_data["correctness"].mean()
            
            # Дистракторы
            wrong = mat_data[mat_data["correctness"] < 1.0]
            top_distractors = []
            if not wrong.empty and "selected_distractor" in wrong.columns:
                top_distractors = wrong["selected_distractor"].value_counts().head(5).index.tolist()

            # Кривая обучения
            learning_curve = 0.0
            if len(mat_data) >= 3:
                try:
                    y = mat_data.sort_values("timestamp")["correctness"].values
                    x = np.arange(len(y))
                    slope = np.polyfit(x, y, 1)[0]
                    learning_curve = float(np.clip(slope * 10, -1.0, 1.0))
                except:
                    learning_curve = 0.0

            material_stats[str(material_id)] = {
                "success_rate": float(success_rate),
                "difficulty_index": float(1 - success_rate),
                "learning_curve": learning_curve,
                "top_distractors": top_distractors,
                "unique_students": int(mat_data["student_id"].nunique())
            }
        return material_stats