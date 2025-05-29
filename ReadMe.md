# Мини-сервис **“Цитатник”**


### в рамках тестового задания для компании
Версия golang: 1.24
---
## Задача и требования
**Задание:**
- Реализуйте REST API-сервис на Go для хранения и управления цитатами.
- Разместите решение в публичном Github репозитории.

**Технические требования:**
- Хранить данные можно в памяти.
- Использовать только стандартные библиотеки Go (максимум gorilla/mux)
- Обязательно: README.md с инструкцией запуска
- Желательно: unit-тесты.

## Описание

Мини-сервис “Цитатник” предоставляет API для работы с цитатами. Сервис поддерживает следующие операции:

1. Добавление новой цитаты (POST /quotes)
2. Получение всех цитат (GET /quotes)
3. Получение случайной цитаты (GET /quotes/random)
4. Фильтрация по автору (GET /quotes?author=Confucius)
5. Удаление цитаты по ID (DELETE /quotes/{id})


### Проверочные команды (curl):
Добавление новой цитаты
```bash
curl -X POST http://localhost:8080/quotes \ -H "Content-Type: application/json" \ -d
'{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```
Получение всех цитат
```bash
curl http://localhost:8080/quotes
```
пример вывода
```json
[
  {
    "id": 1,
    "quote": "Life is simple, but we insist on making it complicated.",
    "author": "Confucius",
    "created_at": "2025-05-29T19:53:01.101063527Z"
  },
  {
    "id": 2,
    "quote": "Life is simple, but we insist on making it complicated.",
    "author": "Confucius",
    "created_at": "2025-05-29T19:53:01.101063527Z"
  }
]
```
Получение случайной цитаты
```bash
curl http://localhost:8080/quotes/random
```
пример вывода:
```json
{
    "id": 1,
    "quote": "Life is simple, but we insist on making it complicated.",
    "author": "Confucius",
    "created_at": "2025-05-29T19:53:01.101063527Z"
}
```

Фильтрация по автору
```bash
curl http://localhost:8080/quotes?author=Confucius
```
пример вывода
```json
[
  {
    "id": 1,
    "quote": "Life is simple, but we insist on making it complicated.",
    "author": "Confucius",
    "created_at": "2025-05-29T19:53:01.101063527Z"
  },
  {
    "id": 2,
    "quote": "Life is simple, but we insist on making it complicated.",
    "author": "Confucius",
    "created_at": "2025-05-29T19:53:01.101063527Z"
  }
]
```
Удаление цитаты по ID
```bash
curl -X DELETE http://localhost:8080/quotes/1
```


---
## Cтруктура проекта

```
quotes_service/
├── cmd/
│   └── main.go
├── internal/
│   ├── handlers/
│   │   ├── handlers_test.go
│   │   └── handlers.go
│   ├── models/
│   │   ├── errs/
│   │   │   └── inmemory_error.go
│   │   └── inmemory_models.go
│   ├── repository/
│   │   └── inmemory/
│   │       ├── quotes_bench_test.go
│   │       ├── quotes_test.go
│   │       └── quotes.go
│   └── server/
│       └── server.go
├── pkg/
│   └── logger/
│       └── logger.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── go.sum
├── Makefile
└── README.md
```

---
## Сборка и запуск

### Сборка и запуск с помощью Makefile
Вызов справки по доступным командам

`make`
```bash
make 
```
или `make help`
```bash
make help
```

### Список доступных команд
```
Доступные команды
Использование   make <команда>
 - build           Собрать проект
 - clean           Очистить проект
 - docker_build    Собрать проект в докере
 - docker_run      Запустить проект в докере
 - docker_stop     Остановить проект в докере
 - help            Показать справку
 - run_bench_storage Запустить бенчмарки для хранилища
 - run_benchs      Запустить все бенчмарки
 - run             Собрать и запустить проект
 - run_test_handler Запустить тесты для хендлера
 - run_test_storage Запустить тесты для хранилища
 - run_tests       Запустить все тесты
```


### Сборка и запуск с помощью Docker-compose
сборка проекта в контейнер
```bash
docker-compose build
```

запуск контейнера
```bash
docker-compose up
```

остановка контейнера
```bash
docker-compose down
```

### Сборка и запуск средствами Golang
сборка проекта
```bash
go build -o . ./cmd/main.go
```

запуск проекта
```bash
./main
```

### Запуск тестов
тесты для хранилища
```bash
go test -v ./internal/repository/inmemory
```

тесты для хендлера
```bash
go test -v ./internal/handlers
```
Запуск бенчмарка хранилища
```bash
go test -bench=. ./internal/repository/inmemory
```

