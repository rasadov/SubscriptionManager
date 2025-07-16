# SubscriptionManager

REST API сервис для агрегации данных об онлайн-подписках пользователей.

## Описание

Этот проект представляет собой REST API для управления подписками пользователей, разработанный на Go с использованием Gin framework, GORM ORM и PostgreSQL в качестве базы данных.

## Функциональность

### API Endpoints

- **POST** `/api/v1/subscriptions` - Создание новой подписки
- **GET** `/api/v1/subscriptions` - Получение списка подписок с фильтрацией и пагинацией
- **GET** `/api/v1/subscriptions/{id}` - Получение подписки по ID
- **PUT** `/api/v1/subscriptions/{id}` - Обновление подписки
- **DELETE** `/api/v1/subscriptions/{id}` - Удаление подписки
- **GET** `/api/v1/subscriptions/total-cost` - Подсчет суммарной стоимости подписок за период

### Дополнительные endpoints

- **GET** `/health` - Проверка состояния сервиса
- **GET** `/swagger/index.html` - Swagger документация

## Структура данных подписки

```json
{
  "id": 1,
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025",
  "created_at": "2025-01-01T12:00:00Z",
  "updated_at": "2025-01-01T12:00:00Z"
}
```

## Технологический стек

- **Go 1.23.4**
- **Gin** - HTTP web framework
- **GORM** - ORM библиотека
- **PostgreSQL** - База данных
- **Swagger** - API документация
- **Docker & Docker Compose** - Контейнеризация
- **Testify** - Тестирование

## Архитектура проекта

```
.
├── cmd/server/           # Точка входа приложения
├── internal/
│   ├── config/          # Конфигурация
│   ├── dto/             # Data Transfer Objects
│   ├── handlers/        # HTTP обработчики
│   ├── models/          # Модели данных
│   ├── repository/      # Слой доступа к данным
│   └── service/         # Бизнес-логика
├── pkg/
│   ├── database/        # Подключение к БД
│   └── logger/          # Логирование
├── tests/               # Тесты
├── docs/                # Swagger документация
├── migrations/          # SQL миграции
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.23+ (для разработки)

### Запуск с Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd subscription-manager
```

2. Создайте файл `.env` на основе `.env.example`:
```bash
cp .env.example .env
```

3. Запустите сервисы:
```bash
docker-compose up -d
```

4. Проверьте состояние сервиса:
```bash
curl http://localhost:8080/health
```

### Локальная разработка

1. Установите зависимости:
```bash
go mod download
```

2. Запустите PostgreSQL:
```bash
docker run --name postgres-dev -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14-alpine
```

3. Настройте переменные окружения:
```bash
export POSTGRES_HOST=localhost
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=postgres
```

4. Запустите приложение:
```bash
go run cmd/server/main.go
```

## API Документация

После запуска сервиса, Swagger документация будет доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

## Примеры использования

### Создание подписки

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

### Получение списка подписок

```bash
curl "http://localhost:8080/api/v1/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&page=1&limit=10"
```

### Подсчет общей стоимости

```bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&start_date=01-2025&end_date=12-2025"
```

## Фильтрация и сортировка

API поддерживает следующие параметры для фильтрации:

- `user_id` - ID пользователя
- `service_name` - Название сервиса
- `start_date_from` - Дата начала подписки (от)
- `end_date_from` - Дата окончания подписки (от)
- `end_date_to` - Дата окончания подписки (до)
- `sort_by` - Поле для сортировки
- `sort_order` - Порядок сортировки (asc/desc)
- `page` - Номер страницы
- `limit` - Количество элементов на странице

## Тестирование

### Запуск всех тестов

```bash
go test ./tests/... -v
```

### Запуск конкретного теста

```bash
go test ./tests -run TestCreateSubscription_Success -v
```

### Покрытие кода

```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Конфигурация

Приложение настраивается через переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `SERVER_PORT` | Порт сервера | `8080` |
| `SERVER_HOST` | Хост сервера | `localhost` |
| `POSTGRES_HOST` | Хост PostgreSQL | `localhost` |
| `POSTGRES_PORT` | Порт PostgreSQL | `5432` |
| `POSTGRES_USER` | Пользователь БД | `postgres` |
| `POSTGRES_PASSWORD` | Пароль БД | `password` |
| `POSTGRES_DB` | Имя БД | `subscriptions` |
| `POSTGRES_SSLMODE` | SSL режим | `disable` |
| `LOG_LEVEL` | Уровень логов | `info` |
| `GIN_MODE` | Режим Gin | `release` |

## База данных

### Миграции

GORM автоматически применяет миграции при запуске приложения. Миграции также доступны в папке `migrations/` для ручного применения.

### Схема таблицы subscriptions

```sql
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date VARCHAR(7) NOT NULL,
    end_date VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Логирование

Приложение использует структурированное логирование с помощью стандартной библиотеки `slog`. Логи включают:

- Информацию о запросах и ответах
- Ошибки базы данных
- Системные события (старт/остановка сервиса)

## Мониторинг

### Health Check

```bash
curl http://localhost:8080/health
```

Ответ:
```json
{
  "status": "healthy",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

## Производительность

- Подключение к БД через пул соединений
- Индексы на часто используемые поля (user_id, service_name, start_date)
- Пагинация для больших результатов
- Graceful shutdown с таймаутом

## Безопасность

- Валидация входных данных
- Использование параметризованных запросов
- Запуск от непривилегированного пользователя в Docker
- SSL отключен только для разработки

## Лицензия

MIT License

## Контакты

Для вопросов и предложений создавайте Issues в репозитории проекта.