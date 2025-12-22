import grpc
from proto import analytics_pb2, analytics_pb2_grpc
from app.application.analytics_service import AnalyticsService

class AnalyticsGRPCHandler(analytics_pb2_grpc.AnalyticsServiceServicer):
    def __init__(self, service: AnalyticsService):
        self.service = service

    async def AnalyzeStudent(self, request, context):
        try:
            # Вызываем бизнес-логику
            result = await self.service.get_student_analysis(request.student_id)
            
            # Мапим в Protobuf Response
            return analytics_pb2.AnalyzeStudentResponse(
                student_id=result.student_id,
                cluster=result.cluster_group,
                engagement_score=result.engagement_score,
                success_rate=result.success_rate,
                recommendations=result.recommendations,
                analyzed_at=result.analyzed_at.isoformat()
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return analytics_pb2.AnalyzeStudentResponse()