from dataclasses import dataclass, field
from typing import List, Dict
from datetime import datetime

@dataclass
class StudentAnalytics:
    student_id: int
    cluster_group: str
    engagement_score: int
    success_rate: float
    topic_efficiency: Dict[str, float]
    recommendations: List[str]
    analyzed_at: datetime = field(default_factory=datetime.now)