# Users Service

Микросервис для управления пользователями.

## Установка и настройка

### Предварительные требования

1. **Go 1.24+**
2. **PostgreSQL**
3. **migrate** (для миграций БД)
4. **oapi-codegen** (для генерации кода из OpenAPI)

### Установка зависимостей

```bash
# Установка oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Установка migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Настройка проекта

1. **Клонируйте репозиторий**
2. **Перейдите в директорию сервиса:**
   ```bash
   cd internal/users-service
   ```

3. **Автоматическая настройка (рекомендуется):**
   ```bash
   make setup
   ```
   Эта команда автоматически:
   - Устанавливает необходимые инструменты (oapi-codegen, migrate)
   - Генерирует код из `openapi.yaml`
   - Добавляет все необходимые зависимости
   - Проверяет корректность установки

4. **Или ручная настройка:**
   ```bash
   make gen-users
   ```
   Эта команда:
   - Генерирует код из `openapi.yaml`
   - Добавляет необходимые зависимости в `go.mod`
   - Выполняет `go mod tidy`

4. **Настройте базу данных:**
   ```bash
   # Примените миграции
   make migrate
   ```

5. **Запустите сервис:**
   ```bash
   make run
   ```

## Структура проекта

```
users-service/
├── cmd/
│   └── server/
│       └── main.go          # Точка входа приложения
├── internal/
│   ├── database/
│   │   └── db.go           # Подключение к БД
│   ├── transport/
│   │   └── grpc/           # gRPC обработчики
│   ├── user/
│   │   ├── model.go        # Модели пользователей
│   │   ├── repository.go   # Репозиторий для работы с БД
│   │   └── service.go      # Бизнес-логика
│   └── web/
│       └── users/
│           └── api.gen.go  # Сгенерированный HTTP API код
├── migrations/             # SQL миграции
├── go.mod                 # Зависимости Go
├── go.sum                 # Хеши зависимостей
├── makefile               # Команды для сборки и развертывания
└── README.md              # Этот файл
```

## API Endpoints

После запуска сервис предоставляет следующие HTTP endpoints:

- `GET /users` - Получить всех пользователей
- `POST /users` - Создать нового пользователя
- `PUT /users/{id}` - Обновить пользователя
- `DELETE /users/{id}` - Удалить пользователя
- `GET /users/{id}/tasks` - Получить задачи пользователя

## Переменные окружения

Настройте следующие переменные окружения:

```bash
DB_DSN="postgres://username:password@localhost:5432/database?sslmode=disable"
```

## Устранение неполадок

### Ошибка "could not import github.com/labstack/echo/v4"

Если вы видите эту ошибку, выполните:

```bash
make gen-users
```

Эта команда автоматически добавит все необходимые зависимости.

### Проблемы с миграциями

```bash
# Проверить статус миграций
migrate -path ./migrations -database "$DB_DSN" version

# Принудительно установить версию
make migrate-force v=1
``` 