import asyncio
import json
import logging
from datetime import datetime
from typing import Dict, Any

import grpc
import pandas as pd
from kafka import KafkaConsumer
from sqlalchemy import create_engine, text

from app.analyzers.engagement_analyzer import EngagementAnalyzer
from app.analyzers.cluster_analyzer import ClusterAnalyzer
from app.analyzers.material_analyzer import MaterialEffectivenessAnalyzer
from app.proto import analytics_pb2, analytics_pb2_grpc
from app.config import Config
from app.database import Database
from app.redis_client import RedisClient

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class AnalyticsService(analytics_pb2_grpc.AnalyticsServiceServicer):
    def __init__(self, db: Database, redis_client: RedisClient):
        self.db = db
        self.redis = redis_client
        self.is_analyzing = False
        
    def AnalyzeStudent(self, request, context):
        """Основной метод анализа студента"""
        try:
            student_id = request.student_id
            
            # Проверка кэша
            cache_key = f"analytics:{student_id}"
            cached = self.redis.get_json(cache_key)
            if cached:
                return analytics_pb2.AnalysisResponse(**cached)
            
            # Получение данных из БД
            logs_df = self.get_student_logs(student_id)
            if logs_df.empty:
                return self._create_empty_response(student_id)
            
            # Анализ вовлеченности
            engagement_analyzer = EngagementAnalyzer(logs_df)
            engagement_analyzer.calculate_engagement_metrics()
            engagement = engagement_analyzer.get_student_engagement_summary(str(student_id))
            
            # Кластерный анализ
            cluster_analyzer = ClusterAnalyzer(logs_df)
            cluster = cluster_analyzer.analyze_clusters()
            
            # Анализ материалов
            material_analyzer = MaterialEffectivenessAnalyzer(logs_df)
            material_stats = material_analyzer.calculate_material_metrics()
            
            # Формирование результата
            result = self._format_analysis_result(
                student_id, engagement, cluster, material_stats
            )
            
            # Кэширование на 1 час
            self.redis.set_json(cache_key, result, ex=3600)
            
            # Сохранение в БД
            self.save_analytics_to_db(student_id, result)
            
            return analytics_pb2.AnalysisResponse(**result)
            
        except Exception as e:
            logger.error(f"Error analyzing student {request.student_id}: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return analytics_pb2.AnalysisResponse()
    
    def get_student_logs(self, student_id: int) -> pd.DataFrame:
        """Получение логов студента из БД"""
        query = """
        SELECT 
            student_id,
            action_type,
            material_id,
            correct,
            time_spent_sec,
            difficulty,
            time_spent_on_mat,
            time_spent_on_question,
            attempts,
            correctness as is_correct,
            selected_distractor,
            timestamp
        FROM student_logs 
        WHERE student_id = :student_id
        AND timestamp > NOW() - INTERVAL '30 days'
        """
        
        df = pd.read_sql(
            text(query), 
            self.db.engine,
            params={"student_id": student_id}
        )
        
        if not df.empty:
            df['timestamp'] = pd.to_datetime(df['timestamp'])
        
        return df
    
    def _format_analysis_result(self, student_id: int, engagement: Dict, 
                               cluster: Dict, material_stats: Dict) -> Dict:
        """Форматирование результатов анализа"""
        return {
            "student_id": student_id,
            "cluster": cluster.get("cluster", "unknown"),
            "engagement_score": engagement.get("engagement_score", 50),
            "recommendations": self._generate_recommendations(engagement, cluster),
            "topic_efficiency": self._calculate_topic_efficiency(material_stats),
            "analyzed_at": datetime.now().isoformat()
        }
    
    def _generate_recommendations(self, engagement: Dict, cluster: Dict) -> list:
        """Генерация рекомендаций"""
        recommendations = []
        
        if engagement.get("engagement_score", 0) < 50:
            recommendations.append("Увеличьте время занятий")
        
        if cluster.get("cluster") == "needs_help":
            recommendations.append("Рекомендуется консультация с преподавателем")
            recommendations.append("Больше практических заданий")
        
        return recommendations
    
    def _calculate_topic_efficiency(self, material_stats: Dict) -> Dict:
        """Расчет эффективности по темам"""
        efficiency = {}
        for material_id, stats in material_stats.items():
            if isinstance(stats, dict):
                efficiency[material_id] = stats.get("success_rate", 0.5)
        return efficiency
    
    def _create_empty_response(self, student_id: int) -> analytics_pb2.AnalysisResponse:
        """Создание пустого ответа"""
        return analytics_pb2.AnalysisResponse(
            student_id=student_id,
            cluster="unknown",
            engagement_score=0,
            recommendations=["Нет данных для анализа"],
            topic_efficiency={},
            analyzed_at=datetime.now().isoformat()
        )
    
    def save_analytics_to_db(self, student_id: int, result: Dict):
        """Сохранение результатов анализа в БД"""
        query = """
        INSERT INTO student_analytics 
        (student_id, cluster_group, engagement_score, recommendation, analyzed_at)
        VALUES (:student_id, :cluster, :engagement_score, :recommendation, NOW())
        ON CONFLICT (student_id) DO UPDATE SET
            cluster_group = EXCLUDED.cluster_group,
            engagement_score = EXCLUDED.engagement_score,
            recommendation = EXCLUDED.recommendation,
            analyzed_at = NOW()
        """
        
        recommendation = "\n".join(result.get("recommendations", []))
        
        self.db.execute(
            query,
            student_id=student_id,
            cluster=result.get("cluster", "unknown"),
            engagement_score=result.get("engagement_score", 0),
            recommendation=recommendation
        )


class KafkaConsumerService:
    def __init__(self, db: Database, redis_client: RedisClient):
        self.db = db
        self.redis = redis_client
        self.consumer = KafkaConsumer(
            Config.KAFKA_TOPIC,
            bootstrap_servers=Config.KAFKA_BOOTSTRAP_SERVERS,
            group_id='analytics-group',
            auto_offset_reset='latest',
            value_deserializer=lambda x: json.loads(x.decode('utf-8'))
        )
    
    async def consume(self):
        """Потребление сообщений из Kafka"""
        logger.info(f"Starting Kafka consumer for topic: {Config.KAFKA_TOPIC}")
        
        for message in self.consumer:
            try:
                log_data = message.value
                student_id = log_data.get("student_id")
                
                # Сохранение лога в БД
                self.save_log_to_db(log_data)
                
                # Инвалидация кэша для этого студента
                cache_key = f"analytics:{student_id}"
                self.redis.delete(cache_key)
                
                logger.info(f"Processed log for student {student_id}")
                
            except Exception as e:
                logger.error(f"Error processing Kafka message: {e}")
    
    def save_log_to_db(self, log_data: Dict):
        """Сохранение лога в БД"""
        query = """
        INSERT INTO student_logs 
        (student_id, action_type, material_id, correct, time_spent_sec, 
         difficulty, time_spent_on_mat, time_spent_on_question, 
         attempts, correctness, selected_distractor)
        VALUES (:student_id, :action_type, :material_id, :correct, 
                :time_spent_sec, :difficulty, :time_spent_on_mat, 
                :time_spent_on_question, :attempts, :correctness, 
                :selected_distractor)
        """
        
        self.db.execute(query, **log_data)


async def serve():
    """Запуск gRPC сервера и Kafka consumer"""
    # Инициализация компонентов
    db = Database()
    redis_client = RedisClient()
    
    # Создание сервиса
    analytics_service = AnalyticsService(db, redis_client)
    
    # Запуск gRPC сервера
    server = grpc.aio.server()
    analytics_pb2_grpc.add_AnalyticsServiceServicer_to_server(
        analytics_service, server
    )
    
    server.add_insecure_port(f'[::]:{Config.GRPC_PORT}')
    
    # Запуск Kafka consumer
    kafka_service = KafkaConsumerService(db, redis_client)
    
    logger.info(f"Starting gRPC server on port {Config.GRPC_PORT}")
    logger.info(f"Starting Kafka consumer for topic: {Config.KAFKA_TOPIC}")
    
    await server.start()
    
    # Запуск consumer в фоне
    asyncio.create_task(kafka_service.consume())
    
    try:
        await server.wait_for_termination()
    except KeyboardInterrupt:
        await server.stop(0)


if __name__ == '__main__':
    asyncio.run(serve())