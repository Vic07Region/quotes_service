.DEFAULT_GOAL := help
BASE_DIR := $(shell pwd)
PROJECT_NAME:= $(shell go list -m)

help: ## Показать справку
	@echo -e "\033[1mДоступные команды\033[0m"
	@echo -e "Использование 	\033[1mmake <команда>\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf " - \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Собрать проект
	@echo "Сборка проекта"
	@CGO_ENABLED=0 go build -o ./$(PROJECT_NAME) ./cmd/main.go
	@echo "Сборка завершена для запуска make run или "./$(PROJECT_NAME)""


run:build ## Собрать и запустить проект
	@echo "Запуск проекта"
	@./$(PROJECT_NAME)

docker_build: ## Собрать проект в докере
	@echo "Сборка проекта в докере"
	@docker-compose build

docker_run: ## Запустить проект в докере
	@echo "Запуск проекта в докере"
	@docker-compose up -d

docker_stop: ## Остановить проект в докере
	@echo "Остановка проекта в докере"
	@docker-compose down

run_test_storage: ## Запустить тесты для хранилища
	@echo "Запуск тестов для хранилища"
	@go test -v ./internal/repository/inmemory

run_bench_storage: ## Запустить бенчмарки для хранилища
	@echo "Запуск бенчмарков для хранилища"
	@go test -bench=. ./internal/repository/inmemory

run_test_handler: ## Запустить тесты для хендлера
	@echo "Запуск тестов для хендлера"
	@go test -v ./internal/handlers

run_tests:: ## Запустить все тесты
	@echo "Запуск всех тестов"
	@make run_test_storage
	@make run_test_handler

run_benchs:: ## Запустить все бенчмарки
	@echo "Запуск всех бенчмарков"
	@make run_bench_storage

clean: ## Очистить проект
	@echo "Очистка проекта"
	@go clean
	@rm -rf $(BASE_DIR)/$(PROJECT_NAME)