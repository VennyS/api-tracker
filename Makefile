# Загружаем переменные из .env
include .env
export $(shell sed 's/=.*//' .env)

# Имя бинарника
BINARY_NAME=api-tracker

.PHONY: up down build run stop logs clickhouse-console

# Собрать go-приложение
build:
	go build -o $(BINARY_NAME) ./cmd

# Запустить только go-приложение
run: build
	./$(BINARY_NAME)

# Запустить docker-сервисы (без go-приложения)
up:
	docker-compose up -d clickhouse prometheus grafana

# Остановить docker-сервисы
down:
	docker-compose down

# Логи всех контейнеров
logs:
	docker-compose logs -f

# Зайти в ClickHouse
clickhouse-console:
	docker exec -it $$(docker ps -qf "ancestor=clickhouse/clickhouse-server") clickhouse-client

# Очистить скомпилированный бинарник
clean:
	rm -f $(BINARY_NAME)

# Полная перезапуск (пересборка всех)
restart: down clean build up run
