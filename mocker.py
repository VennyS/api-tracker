import random
import json
from datetime import datetime, timedelta
from faker import Faker
import requests

# Инициализация Faker для генерации случайных данных
fake = Faker()

# Списки сервисов и путей
SERVICES = ["AuthService", "PaymentService", "UserService", "OrderService", "NotificationService"]
PATHS = ["/api/v1/login", "/api/v1/register", "/api/v1/payments", "/api/v1/orders", "/api/v1/notifications"]

# Функция для генерации одного API-запроса
def generate_log_entry():
    # Генерация случайного смещения от -10 до +10 дней
    delta_days = random.randint(-30, 0)
    delta_seconds = random.randint(0, 86400)  # произвольное время в течение суток
    random_timestamp = datetime.utcnow() + timedelta(days=delta_days, seconds=delta_seconds)

    # Форматирование в ISO 8601 с 'Z'
    timestamp = random_timestamp.strftime('%Y-%m-%dT%H:%M:%SZ')

    return {
        "timestamp": timestamp,
        "method": random.choice(["GET", "POST", "PUT", "DELETE"]),
        "path": random.choice(PATHS),
        "status_code": random.choice([200, 404, 500, 301, 403]),
        "latency_ms": round(random.uniform(20, 1000), 2),
        "ip": fake.ipv4(),
        "user_agent": fake.user_agent(),
        "service_name": random.choice(SERVICES)
    }

# Функция для отправки логов на сервер
def send_log_to_server(log_entry):
    url = "http://localhost:8080/"  # Адрес сервера
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, json=log_entry, headers=headers)
    if response.status_code == 200:
        print(f"Лог успешно отправлен: {log_entry}")
    else:
        print(f"Ошибка при отправке лога: {response.status_code}")

# Генерация и отправка логов
def generate_and_send_logs(num_entries):
    for _ in range(num_entries):
        log_entry = generate_log_entry()
        send_log_to_server(log_entry)

# Главная функция
if __name__ == "__main__":
    num_logs = 10000  # Количество логов, которое нужно отправить
    generate_and_send_logs(num_logs)
