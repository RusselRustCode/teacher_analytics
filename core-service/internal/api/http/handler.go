//ДЛЯ ФРОНТЕНДА, ПОКА ЧТО НЕ РАБОТАЕТ!

package http

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/gin-gonic/gin"
    
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/domain"
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
)

type HTTPHandler struct {
    service interfaces.AnalyticsService
}

func NewHTTPHandler(service interfaces.AnalyticsService) *HTTPHandler {
    return &HTTPHandler{
        service: service,
    }
}

func (h *HTTPHandler) SendLog(c *gin.Context) {
    var log domain.StudentLog
    
    if err := c.BindJSON(&log); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request body",
            "details": err.Error(),
        })
        return
    }
    
    if err := h.service.SendLog(c.Request.Context(), &log); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to process log",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Log processed successfully",
    })
}

func (h *HTTPHandler) GetAnalytics(c *gin.Context) {
    studentIDStr := c.Param("student_id")
    studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid student ID",
        })
        return
    }
    
    analytics, err := h.service.GetAnalytics(c.Request.Context(), studentID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to get analytics",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, analytics)
}

func (h *HTTPHandler) GetStudentLogs(c *gin.Context) {
    studentIDStr := c.Param("student_id")
    studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
        return
    }
    
    fromStr := c.DefaultQuery("from", "")
    toStr := c.DefaultQuery("to", "")
    
    var from, to time.Time
    
    if fromStr != "" {
        from, err = time.Parse(time.RFC3339, fromStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format"})
            return
        }
    } else {
        from = time.Now().AddDate(0, -1, 0) // Месяц назад по умолчанию сделаем
    }
    
    if toStr != "" {
        to, err = time.Parse(time.RFC3339, toStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format"})
            return
        }
    } else {
        to = time.Now()
    }
    
    logs, err := h.service.GetStudentLogs(c.Request.Context(), studentID, from, to)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to get logs",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "student_id": studentID,
        "logs":       logs,
        "count":      len(logs),
    })
}

func (h *HTTPHandler) GetStudents(c *gin.Context) {
    students, err := h.service.GetStudents(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to get students",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "students": students,
        "count":    len(students),
    })
}

func (h *HTTPHandler) TriggerAnalysis(c *gin.Context) {
    var request struct {
        StudentID uint64 `json:"student_id"`
    }
    
    if err := c.BindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    if err := h.service.TriggerAnalysis(c.Request.Context(), request.StudentID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to trigger analysis",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success":    true,
        "student_id": request.StudentID,
        "message":    "Analysis triggered successfully",
    })
}