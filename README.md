# Chat API

REST API для чатов и сообщений на Go с PostgreSQL.

## Быстрый старт

### Запуск приложения
```bash
docker-compose up --build
```
API будет доступно на `http://localhost:8080`

### Интеграционные тесты
1. **Запустите только базу данных:**
```bash
docker-compose up -d postgres
```

2. **Запустите интеграционные тесты:**
```bash
go test -v ./tests/integration/...
```

## Структура проекта

```
.
├── Dockerfile                      # Конфигурация Docker образа приложения
├── cmd/
│   └── api/
│       └── main.go                 # Основная точка входа приложения
├── config.yaml                     # Файл конфигурации приложения
├── docker-compose.yml              # Конфигурация Docker Compose для развертывания
├── go.mod                          # Файл модуля Go с зависимостями
├── go.sum                          # Контрольные суммы зависимостей Go
├── internal/                       # Внутренние пакеты приложения
│   ├── apperrors/
│   │   └── apperrors.go            # Определение пользовательских ошибок приложения
│   ├── config/
│   │   ├── config.go               # Структуры данных конфигурации
│   │   ├── helper.go               # Вспомогательные функции для конфигурации
│   │   └── loader.go               # Загрузчик конфигурации из файла
│   ├── logger/
│   │   └── logger.go               # Настройка логгера приложения
│   ├── models/
│   │   └── models.go               # Модели данных (Chat, Message, ChatWithMessages)
│   ├── repository/
│   │   └── postgres/               # Репозитории для работы с PostgreSQL
│   │       ├── chat.go             # Репозиторий для операций с чатами
│   │       ├── database.go         # Инициализация подключения к базе данных
│   │       └── message.go          # Репозиторий для операций с сообщениями
│   ├── service/                    # Слой бизнес-логики
│   │   ├── chat.go                 # Сервис для работы с чатами
│   │   ├── message.go              # Сервис для работы с сообщениями
│   │   └── service.go              # Интерфейсы сервисов
│   └── transport/
│       └── http/                   # HTTP транспорт
│           ├── dto/                # Data Transfer Objects
│           │   ├── request.go      # DTO для входящих запросов
│           │   └── response.go     # DTO для исходящих ответов
│           ├── handlers/           # HTTP обработчики (хендлеры)
│           │   ├── chat.go         # Обработчик для операций с чатами
│           │   ├── chat_test.go    # Юнит-тесты для обработчика чатов
│           │   ├── handler.go      # Базовый обработчик с зависимостями
│           │   ├── helper.go       # Вспомогательные функции для HTTP
│           │   ├── message.go      # Обработчик для операций с сообщениями
│           │   └── service_mock_test.go # Мок-объекты для тестирования
│           ├── router/
│           │   └── router.go       # Конфигурация маршрутизации HTTP
│           └── server/
│               └── server.go       # Конфигурация и запуск HTTP сервера
├── migrations/                     # Миграции базы данных
│   └── postgres/                   # Миграции для PostgreSQL
│       ├── 0001_chats.go          # Создание таблицы chats
│       └── 0002_messages.go       # Создание таблицы messages
├── postman/                       # Файлы для тестирования API
│   └── chat-api.postman_collection.json # Postman коллекция для тестирования
└── tests/                         # Тесты
    └── integration/               # Интеграционные тесты
        └── chat_integration_test.go # Интеграционный тест жизненного цикла чата
```

## Postman коллекция

Файл `postman/chat-api.postman_collection.json` содержит готовые запросы для тестирования всех endpoints API.
