// Конфигурация API
const API_BASE_URL = 'http://localhost/api';

// DOM элементы
let currentStudentId = null;
let topicEfficiencyChart = null;
let activityChart = null;

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    // Загрузка студентов
    loadStudents();
    
    // Обновление времени
    updateCurrentTime();
    setInterval(updateCurrentTime, 1000);
    
    // Обработчики событий
    document.getElementById('triggerAnalysis').addEventListener('click', triggerAnalysis);
    document.getElementById('refreshData').addEventListener('click', loadStudents);
    document.getElementById('logForm').addEventListener('submit', submitLogForm);
    document.getElementById('studentSearch').addEventListener('input', filterStudents);
});

// Загрузка списка студентов
async function loadStudents() {
    try {
        const response = await fetch(`${API_BASE_URL}/students`);
        const data = await response.json();
        
        const studentList = document.getElementById('studentList');
        studentList.innerHTML = '';
        
        data.students.forEach(student => {
            const li = document.createElement('li');
            li.className = 'student-item';
            li.dataset.id = student.id;
            li.innerHTML = `
                <div class="student-name">${student.name}</div>
                <div class="student-email">${student.email}</div>
                <div class="student-role">${student.role}</div>
            `;
            
            li.addEventListener('click', () => selectStudent(student.id));
            studentList.appendChild(li);
        });
        
        // Выбираем первого студента по умолчанию
        if (data.students.length > 0) {
            selectStudent(data.students[0].id);
        }
    } catch (error) {
        console.error('Error loading students:', error);
        showError('Не удалось загрузить список студентов');
    }
}

// Выбор студента
async function selectStudent(studentId) {
    currentStudentId = studentId;
    
    // Обновляем активный элемент
    document.querySelectorAll('.student-item').forEach(item => {
        item.classList.remove('active');
        if (item.dataset.id == studentId) {
            item.classList.add('active');
        }
    });
    
    // Обновляем ID студента
    document.getElementById('currentStudentId').textContent = studentId;
    
    // Загружаем аналитику
    await loadAnalytics(studentId);
    
    // Загружаем логи
    await loadStudentLogs(studentId);
}

// Загрузка аналитики студента
async function loadAnalytics(studentId) {
    try {
        const response = await fetch(`${API_BASE_URL}/analytics/${studentId}`);
        const data = await response.json();
        
        // Обновляем информацию о студенте
        updateStudentInfo(data.analytics);
        
        // Обновляем графики
        updateCharts(data.analytics);
        
        // Обновляем рекомендации
        updateRecommendations(data.analytics);
        
    } catch (error) {
        console.error('Error loading analytics:', error);
        showError('Не удалось загрузить аналитику студента');
    }
}

// Обновление информации о студенте
function updateStudentInfo(analytics) {
    if (!analytics) return;
    
    // Кластер
    document.getElementById('clusterInfo').textContent = 
        analytics.cluster_group || 'Не определен';
    
    // Вовлеченность
    const engagement = analytics.engagement_score || 0;
    document.getElementById('engagementBar').style.width = `${engagement}%`;
    document.getElementById('engagementScore').textContent = `${engagement}/100`;
    
    // Успешность
    const successRate = analytics.success_rate || 0;
    document.getElementById('successRate').textContent = 
        `${(successRate * 100).toFixed(1)}%`;
    
    // Среднее время
    const avgTime = analytics.avg_time_per_task || 0;
    document.getElementById('avgTime').textContent = 
        `${Math.round(avgTime)} сек`;
}

// Обновление графиков
function updateCharts(analytics) {
    // График эффективности по темам
    const topicCtx = document.getElementById('topicEfficiencyChart').getContext('2d');
    
    if (topicEfficiencyChart) {
        topicEfficiencyChart.destroy();
    }
    
    // Пример данных (замени на реальные)
    const topicData = {
        labels: ['Математика', 'Физика', 'Программирование', 'Английский'],
        datasets: [{
            label: 'Эффективность (%)',
            data: [85, 72, 90, 65],
            backgroundColor: [
                '#667eea',
                '#764ba2',
                '#f093fb',
                '#4c51bf'
            ],
            borderWidth: 1
        }]
    };
    
    topicEfficiencyChart = new Chart(topicCtx, {
        type: 'bar',
        data: topicData,
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    max: 100
                }
            }
        }
    });
    
    // График активности по дням
    const activityCtx = document.getElementById('activityChart').getContext('2d');
    
    if (activityChart) {
        activityChart.destroy();
    }
    
    const activityData = {
        labels: ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'],
        datasets: [{
            label: 'Активность',
            data: [12, 19, 8, 15, 22, 10, 5],
            borderColor: '#48bb78',
            backgroundColor: 'rgba(72, 187, 120, 0.2)',
            fill: true,
            tension: 0.4
        }]
    };
    
    activityChart = new Chart(activityCtx, {
        type: 'line',
        data: activityData,
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: false
                }
            }
        }
    });
}

// Обновление рекомендаций
function updateRecommendations(analytics) {
    const recommendationsList = document.getElementById('recommendationsList');
    
    if (!analytics.recommendation) {
        recommendationsList.innerHTML = '<li>Нет рекомендаций для отображения</li>';
        return;
    }
    
    const recommendations = analytics.recommendation.split('\n').filter(r => r.trim());
    
    if (recommendations.length === 0) {
        recommendationsList.innerHTML = '<li>Нет рекомендаций для отображения</li>';
        return;
    }
    
    recommendationsList.innerHTML = '';
    recommendations.forEach(rec => {
        const li = document.createElement('li');
        li.textContent = rec;
        recommendationsList.appendChild(li);
    });
}

// Загрузка логов студента
async function loadStudentLogs(studentId) {
    try {
        const response = await fetch(`${API_BASE_URL}/students/${studentId}/logs`);
        const data = await response.json();
        
        // Здесь можно отобразить логи, если нужно
        console.log('Student logs:', data);
        
    } catch (error) {
        console.error('Error loading logs:', error);
    }
}

// Запуск анализа
async function triggerAnalysis() {
    if (!currentStudentId) {
        showError('Выберите студента для анализа');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/analyze`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                student_id: currentStudentId
            })
        });
        
        const data = await response.json();
        
        if (data.success) {
            showSuccess('Анализ запущен успешно!');
            
            // Ждем 2 секунды и обновляем данные
            setTimeout(() => {
                loadAnalytics(currentStudentId);
            }, 2000);
            
        } else {
            showError('Ошибка при запуске анализа');
        }
        
    } catch (error) {
        console.error('Error triggering analysis:', error);
        showError('Не удалось запустить анализ');
    }
}

// Отправка лога
async function submitLogForm(event) {
    event.preventDefault();
    
    const formData = {
        student_id: parseInt(document.getElementById('studentId').value),
        action_type: document.getElementById('actionType').value,
        material_id: document.getElementById('materialId').value,
        correct: document.getElementById('correct').value === 'true',
        time_spent_sec: parseInt(document.getElementById('timeSpent').value),
        difficulty: parseInt(document.getElementById('difficulty').value),
        time_spent_on_mat: 0,
        time_spent_on_question: 0,
        attempts: 1,
        correctness: document.getElementById('correct').value === 'true',
        selected_distractor: null
    };
    
    // Устанавливаем время в зависимости от типа действия
    if (formData.action_type === 'view_material') {
        formData.time_spent_on_mat = formData.time_spent_sec;
    } else {
        formData.time_spent_on_question = formData.time_spent_sec;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/log`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ log: formData })
        });
        
        const data = await response.json();
        
        if (data.success) {
            showSuccess('Лог успешно отправлен!');
            document.getElementById('logForm').reset();
            
            // Обновляем аналитику для текущего студента
            if (currentStudentId === formData.student_id) {
                loadAnalytics(currentStudentId);
            }
            
        } else {
            showError('Ошибка при отправке лога: ' + data.message);
        }
        
    } catch (error) {
        console.error('Error submitting log:', error);
        showError('Не удалось отправить лог');
    }
}

// Фильтрация студентов
function filterStudents() {
    const searchTerm = document.getElementById('studentSearch').value.toLowerCase();
    const studentItems = document.querySelectorAll('.student-item');
    
    studentItems.forEach(item => {
        const text = item.textContent.toLowerCase();
        item.style.display = text.includes(searchTerm) ? 'flex' : 'none';
    });
}

// Обновление времени
function updateCurrentTime() {
    const now = new Date();
    const timeString = now.toLocaleTimeString('ru-RU');
    document.getElementById('currentTime').textContent = timeString;
}

// Показать уведомление об успехе
function showSuccess(message) {
    showNotification(message, 'success');
}

// Показать уведомление об ошибке
function showError(message) {
    showNotification(message, 'error');
}

// Показать уведомление
function showNotification(message, type) {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 25px;
        border-radius: 8px;
        color: white;
        font-weight: 600;
        z-index: 1000;
        animation: slideIn 0.3s ease-out;
    `;
    
    if (type === 'success') {
        notification.style.backgroundColor = '#48bb78';
    } else {
        notification.style.backgroundColor = '#f56565';
    }
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Добавляем стили для анимации
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);